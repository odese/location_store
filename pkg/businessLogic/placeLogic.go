package businessLogic

import (
	"fmt"
	"io"
	"location_store/pkg/domain/dto"
	"location_store/pkg/domain/models"
	log "location_store/pkg/infrastructure/logger"
	"location_store/pkg/repositories/postgreRepo"
	"location_store/pkg/repositories/yelpRepo"
	"location_store/pkg/utils"
	"strconv"
	"sync"
)

// ParseCsvIntoPlaceModel parses the CSV into a Place model
func ParseCsvIntoPlaceModel(payload io.ReadCloser) (places []models.Place, err error) {
	places = make([]models.Place, 0)

	records, err := utils.ReadCsvFromPayload(payload)
	if err != nil {
		log.Error().Err(err).Msg("Error reading CSV from payload")
		return places, err
	}

	for i := 1; i < len(records); i++ {
		row := records[i]
		var place models.Place
		place.PlaceID = row[0]
		place.Latitude = row[1]
		place.Longitude = row[2]
		place.Name = row[3]
		places = append(places, place)
	}

	return places, err
}

// CreatePlaces, creates a list of places in the database.
func CreatePlaces(places []models.Place) (successAmount int, err error) {
	for i := 0; i < len(places); i++ {
		place := places[i]
		err = CreatePlace(place)
		if err != nil {
			if utils.PostgreDuplicateKeyErr(err) {
				continue
			} else {
				return successAmount, err
			}
		} else {
			successAmount++
		}
	}

	return successAmount, nil
}

// CreatePlace, creates a place in the database.
// XXX: How to act if there is an error like too mant connections to DB?
func CreatePlace(place models.Place) (err error) {
	err = postgreRepo.InsertPlace(place)
	if err != nil {
		if utils.PostgreDuplicateKeyErr(err) {
			// log.Debug().Err(err).Str("Key", place.PlaceID).Msg("Duplicate key")
			return err
		} else if utils.PostgreConnBusyErr(err) {
			log.Warn().Err(err).Str("Key", place.PlaceID).Msg("Connection busy")
			err = CreatePlace(place)
			return err
		} else {
			log.Error().Err(err).Msg("Error on inserting place")
			return err
		}
	}
	return err
}

// ListPlaces, lists places from the database.
func ListPlaces(offset, limit int) (response dto.ListOfPlaces, err error) {
	if limit == 0 {
		response.Places, err = postgreRepo.GetAllPlaces()
	} else {
		response.Places, err = postgreRepo.GetPlacesByLimit(offset, limit)
	}
	if err != nil {
		log.Error().Err(err).Msg("Error getting places from DB")
		return response, err
	}

	response.TotalCount, err = postgreRepo.CountPlaces()
	if err != nil {
		log.Error().Err(err).Msg("Error counting places on DB")
		return response, err
	}

	return response, err
}

func CountPlaces() (count int, err error) {
	count, err = postgreRepo.CountPlaces()
	if err != nil {
		log.Error().Err(err).Msg("Error counting places on DB")
		return count, err
	}

	return count, err
}

// GetAndSaveNearbyPlacesByList, gets nearby places from Yelp and saves them in the database.
func GetAndSaveNearbyPlacesByList(places []models.Place) (err error) {
	for i := 0; i < len(places); i++ {
		place := places[i]
		GetAndSaveNearbyPlaces(place)
	}
	return err
}

// GetAndSaveNearbyPlacesByListConcurently, gets nearby places from Yelp and saves them in the database.
func GetAndSaveNearbyPlacesByListConcurently(places []models.Place) (response dto.CountsByPlaces, err error) {
	// Job is defined as searching & collecting & saving nearby businesses for each individual place in the csv file.
	// So, the total number of jobs is the number of places in the csv file.
	totalJobs := places
	totalJobCount := len(totalJobs)
	// Create a channel to assign job via (a.k.a) place.
	jobs := make(chan models.Place, totalJobCount)

	var worker sync.WaitGroup

	// XXX: What if the number of places in the csv file is less than the maxWorker?
	maxWorker := 5 // In other words max concurrent jobs
	for i := 0; i < maxWorker; i++ {
		worker.Add(1)
		go work(jobs, &worker, &response)
	}

	// Assign jobs to workers
	for jobIndex := 0; jobIndex < totalJobCount; jobIndex++ {
		job := totalJobs[jobIndex]

		var nearby dto.PlaceWithNearbyCount
		nearby.PlaceName = job.Name
		response.Places = append(response.Places, nearby)

		jobs <- job
	}
	close(jobs)

	worker.Wait()

	return response, err
}

func work(jobs <-chan models.Place, worker *sync.WaitGroup, response *dto.CountsByPlaces) {
	defer worker.Done()

	for place := range jobs {
		nearbyCount, uniqueCount, err := GetAndSaveNearbyPlaces(place)
		if err != nil {
			// XXX: What if there is an error on one of the jobs?
			// Should we stop the whole process?
			// Should we continue the process?
			// Should we retry the failed job?
			// Should we retry the whole process?
			log.Error().Err(err).Str("Place key", place.PlaceID).Msg("Error on getting and saving nearby places")
		}

		// Update the response
		for i := 0; i < len(response.Places); i++ {
			if response.Places[i].PlaceName == place.Name {
				response.Places[i].NearbyCount = nearbyCount
				response.Places[i].UniqueCount = uniqueCount
			}
		}

	}
}

// GetAndSaveNearbyPlaces, gets nearby places from Yelp.
func GetAndSaveNearbyPlaces(place models.Place) (nearbyCount, uniqueCount int, err error) {
	businessList, err := CollectAllNearbyYelpBusinesses(place)
	nearbyCount = len(businessList)
	if err != nil {
		log.Error().Err(err).Msg("Error on collecting businesses")
		return nearbyCount, uniqueCount, err
	}

	// if len(businessList) == 0 {
	// 	log.Info().Str("Place key", place.PlaceID).Msg("Has no nearby suggestions")
	// 	return nearbyCount, uniqueCount, err
	// }

	placeList := convertBusinessesToPlaces(businessList)

	uniqueCount, err = CreatePlaces(placeList)
	if err != nil {
		log.Error().Err(err).Msg("Error on creating places")
		return nearbyCount, uniqueCount, err
	}

	return nearbyCount, uniqueCount, err
}

func CollectAllNearbyYelpBusinesses(place models.Place) (businessList []models.YelpBusiness, err error) {
	businessList = make([]models.YelpBusiness, 0)

	searchResponse, err := yelpRepo.SearchBusinesses("", place.Latitude, place.Longitude, "0")
	if err != nil {
		log.Error().Err(err).Msg("Error on searching businesses")
		return businessList, err
	}

	businessList = append(businessList, searchResponse.Businesses...)

	if searchResponse.Total > utils.LimitForYelpInt {
		iterationRequired := (searchResponse.Total / utils.LimitForYelpInt) + 1

		for i := 1; i <= iterationRequired; i++ {

			offset := i * utils.LimitForYelpInt
			offsetStr := strconv.Itoa(offset)

			searchResponse, err = yelpRepo.SearchBusinesses("", place.Latitude, place.Longitude, offsetStr)
			if err != nil {
				log.Warn().Err(err).Msg("Error on searching businesses")
				continue
			} else {
				businessList = append(businessList, searchResponse.Businesses...)
			}

		}

	}

	return businessList, nil
}

func convertBusinessesToPlaces(businessList []models.YelpBusiness) (places []models.Place) {
	places = make([]models.Place, 0)

	for i := 0; i < len(businessList); i++ {
		business := businessList[i]
		place := convertBusinessToPlace(business)
		places = append(places, place)
	}

	return places
}

func convertBusinessToPlace(business models.YelpBusiness) (place models.Place) {
	place.PlaceID = business.ID
	place.Latitude = fmt.Sprintf("%f", business.Coordinates.Latitude)
	place.Longitude = fmt.Sprintf("%f", business.Coordinates.Longitude)
	place.Name = business.Name

	return place
}
