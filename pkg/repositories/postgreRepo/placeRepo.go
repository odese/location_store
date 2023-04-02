package postgreRepo

import (
	"location_store/pkg/domain/models"
	"location_store/pkg/infrastructure/db/postgre"
)

// InsertPlace, inserts a new place into the database.
func InsertPlace(place models.Place) (err error) {
	query := `INSERT INTO public."Places" (id, latitude, longitude, name) VALUES ($1, $2, $3, $4)`
	_, err = postgre.InsertQuery(query, place.PlaceID, place.Latitude, place.Longitude, place.Name)
	return err
}

// GetAllPlaces, gets all places from the database.
func GetAllPlaces() (places []models.Place, err error) {
	places = make([]models.Place, 0)

	query := `SELECT id, latitude, longitude, name FROM public."Places"`
	rows, err := postgre.SelectQuery(query)
	if err != nil {
		return places, err
	}

	for rows.Next() {
		var place models.Place
		err = rows.Scan(&place.PlaceID, &place.Latitude, &place.Longitude, &place.Name)
		if err != nil {
			return places, err
		}
		places = append(places, place)
	}

	return places, err
}

// GetPlacesByLimit, gets places by chosen offset and limit from the database.
func GetPlacesByLimit(offset, limit int) (places []models.Place, err error) {
	places = make([]models.Place, 0)

	query := `SELECT id, latitude, longitude, name FROM public."Places" LIMIT $1 OFFSET $2`
	rows, err := postgre.SelectQuery(query, limit, offset)
	if err != nil {
		return places, err
	}

	for rows.Next() {
		var place models.Place
		err = rows.Scan(&place.PlaceID, &place.Latitude, &place.Longitude, &place.Name)
		if err != nil {
			return places, err
		}
		places = append(places, place)
	}

	return places, err
}

func CountPlaces() (count int, err error) {
	query := `SELECT COUNT(*) FROM public."Places"`
	rows, err := postgre.SelectQuery(query)
	if err != nil {
		return count, err
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}
	
	return count, err
}