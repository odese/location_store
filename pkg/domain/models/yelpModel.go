package models

type YelpBusiness struct {
	ID           string                 `json:"id"`
	Alias        string                 `json:"alias"`
	Name         string                 `json:"name"`
	ImageURL     string                 `json:"image_url"`
	IsClosed     bool                   `json:"is_closed"`
	URL          string                 `json:"url"`
	ReviewCount  int                    `json:"review_count"`
	Categories   []YelpBusinessCategory `json:"categories"`
	Rating       float64                `json:"rating"`
	Coordinates  YelpCoordinates        `json:"coordinates"`
	Transactions []interface{}          `json:"transactions"`
	Price        string                 `json:"price,omitempty"`
	Location     YelpLocation           `json:"location"`
	Phone        string                 `json:"phone"`
	DisplayPhone string                 `json:"display_phone"`
	Distance     float64                `json:"distance"`
}

type YelpBusinessCategory struct {
	Alias string `json:"alias"`
	Title string `json:"title"`
}

type YelpCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type YelpLocation struct {
	Address1       string   `json:"address1"`
	Address2       string   `json:"address2"`
	Address3       string   `json:"address3"`
	City           string   `json:"city"`
	ZipCode        string   `json:"zip_code"`
	Country        string   `json:"country"`
	State          string   `json:"state"`
	DisplayAddress []string `json:"display_address"`
}
