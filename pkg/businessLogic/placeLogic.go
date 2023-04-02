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
func CreatePlaces(places []models.Place) (err error) {
	for i := 0; i < len(places); i++ {
		place := places[i]
		err = CreatePlace(place)
		if err != nil {
			if utils.ErrDuplicateKey(err) {
				continue
			} else {
				return err
			}
		}
	}

	return nil
}

// CreatePlace, creates a place in the database.
func CreatePlace(place models.Place) (err error) {
	err = postgreRepo.InsertPlace(place)
	if err != nil {
		if utils.ErrDuplicateKey(err) {
			log.Warn().Err(err).Str("Key", place.PlaceID).Msg("Duplicate key")
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
		log.Error().Err(err).Msg("Error counting places from DB")
		return response, err
	}

	return response, err
}

// GetAndSaveNearbyPlacesByList, gets nearby places from Yelp and saves them in the database.
func GetAndSaveNearbyPlacesByList(places []models.Place) (err error) {
	for i := 0; i < len(places); i++ {
		place := places[i]
		GetAndSaveNearbyPlaces(place)
	}
	return err
}

// GetAndSaveNearbyPlaces, gets nearby places from Yelp.
func GetAndSaveNearbyPlaces(place models.Place) (err error) {
	businessList, err := CollectAllNearbyYelpBusinesses(place)
	if err != nil {
		log.Error().Err(err).Msg("Error on collecting businesses")
		return err
	}

	if len(businessList) == 0 {
		log.Info().Str("Place key", place.PlaceID).Msg("Has no nearby suggestions")
		return err
	}

	placeList := convertBusinessesToPlaces(businessList) 

	err = CreatePlaces(placeList)
	if err != nil {
		log.Error().Err(err).Msg("Error on creating places")
	}


	return err
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

			// TODO: In case of requests are too fast, or any other error, 
			// we should implement some retry mechanism until it get success.
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
	place.Latitude = fmt.Sprintln(business.Coordinates.Latitude)  
	place.Longitude = fmt.Sprintln(business.Coordinates.Longitude)
	place.Name = business.Name

	return place
}