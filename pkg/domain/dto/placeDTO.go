package dto

import (
	"location_store/pkg/domain/models"
)

type ListOfPlaces struct {
	TotalCount int            `json:"total_count"`
	Places     []models.Place `json:"places"`
}

type CountsByPlaces struct {
	ActualCount int                    `json:"total_count"`
	Places      []PlaceWithNearbyCount `json:"places"`
}

type PlaceWithNearbyCount struct {
	PlaceName   string `json:"place_name"`
	NearbyCount int    `json:"nearby_count"`
	UniqueCount int    `json:"unique_count"`
}
