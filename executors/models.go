package executors

// Those models are sent to us from the external API.

type geoInformationModel struct {
	Results []struct {
		Country string `json:"country"`
		Name    string `json:"name"`

		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	} `json:"results"`
}

type weatherInformationModel struct {
	Hourly struct {
		Time []int `json:"time"`

		Temperature []float32 `json:"temperature_2m"`
		Rain        []float32 `json:"rain"`
		Windspeed   []float32 `json:"windspeed_10m"`
	} `json:"hourly"`
}
