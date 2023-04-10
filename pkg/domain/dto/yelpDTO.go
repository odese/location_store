package dto

import "location_store/pkg/domain/models"

type YelpSearchBusinessesResponse struct {
	Total  int `json:"total"`
	Region struct {
		Center models.YelpCoordinates `json:"center"`
	} `json:"region"`
	Businesses []models.YelpBusiness `json:"businesses"`
}

type YelpErrorResponse struct {
	Error struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"error"`
}
