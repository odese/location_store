package businessLogic

import (
	"location_store/pkg/domain/models"
	"location_store/pkg/infrastructure/config"
	"location_store/pkg/infrastructure/db/postgre"
	log "location_store/pkg/infrastructure/logger"
	"strconv"
	"testing"
)

func initTestSetup() {
	log.SetupLogger()
	config.InitForTest()
	postgre.Init()
}

func createTestPlaces() (places []models.Place) {
	places = make([]models.Place, 0)

	for i := 0; i < 5; i++ {
		var place models.Place
		place.PlaceID = "test" + strconv.Itoa(i)
		place.Latitude = "test"
		place.Longitude = "test"
		place.Name = "test"

		places = append(places, place)
	}

	return places
}

func TestCreatePlace(t *testing.T) {
	initTestSetup()

	var place models.Place
	place.PlaceID = "test"
	place.Latitude = "test"
	place.Longitude = "test"
	place.Name = "test"

	err := CreatePlace(place)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestGetAndSaveNearbyPlacesByListConcurently(t *testing.T) {
	initTestSetup()

	places := createTestPlaces()
	_, err := GetAndSaveNearbyPlacesByListConcurently(places)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestListPlaces(t *testing.T) {
	initTestSetup()

	_, err := ListPlaces(0, 10)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}