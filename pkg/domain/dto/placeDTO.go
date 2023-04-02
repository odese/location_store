package dto

import (
	"location_store/pkg/domain/models"
)

type ListOfPlaces struct {
	TotalCount int            `json:"total_count"`
	Places     []models.Place `json:"places"`
}
