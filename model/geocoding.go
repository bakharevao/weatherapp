package model

type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

type GeocodingResult struct {
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
}
