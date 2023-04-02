package models

type Place struct {
	PlaceID   string `json:"place_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
}
