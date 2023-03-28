package executors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Those are the friendlier structs we're producing from *Model structures.

type GeoInformation struct {
	Country string
	City    string

	Longitude float32
	Latitude  float32
}

type WeatherInformation struct {
	Temperature float32
	RainLevel   float32
	WindSpeed   float32
}

func FetchGeoInformation(city string) (GeoInformation, error) {
	// Send API request.
	URL := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&language=ru", url.QueryEscape(city))

	response, err := http.Get(URL)
	if err != nil {
		return GeoInformation{},
			errors.New(
				fmt.Sprintf("Geocoding request failed for city %s: %s", city, err))
	}

	// Make sure to close the body reader at the end of this function.
	defer response.Body.Close()

	// Parse the JSON model.
	var responseModel geoInformationModel

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseModel)
	if err != nil {
		return GeoInformation{}, errors.New(
			fmt.Sprintf("Geocode data parsing failed for city %s: %s", city, err))
	}

	// Validate the parsed model.
	if len(responseModel.Results) <= 0 {
		return GeoInformation{}, errors.New(
			fmt.Sprintf("Geocode data not found for city %s.", city))
	}

	// Form a response and return it.
	bestResult := responseModel.Results[0]

	return GeoInformation{
		Country: bestResult.Country,
		City:    bestResult.Name,

		Longitude: bestResult.Longitude,
		Latitude:  bestResult.Latitude,
	}, nil
}

func FetchWeatherInformation(geo GeoInformation) (WeatherInformation, error) {
	// Send API request.
	URL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&hourly=temperature_2m,rain,windspeed_10m&timeformat=unixtime",
		geo.Latitude, geo.Longitude)

	response, err := http.Get(URL)
	if err != nil {
		return WeatherInformation{},
			errors.New(
				fmt.Sprintf(
					"Weather forecast request failed for lat, lon %.2f, %.2f: %s",
					geo.Latitude, geo.Longitude, err))
	}

	// Make sure to close the body reader at the end of this function.
	defer response.Body.Close()

	// Parse the JSON model.
	var responseModel weatherInformationModel

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseModel)
	if err != nil {
		return WeatherInformation{}, errors.New(
			fmt.Sprintf(
				"Weather data parsing failed for lat, lon %.2f, %.2f: %s",
				geo.Latitude, geo.Longitude, err))
	}

	// Find closest record to the current timestamp in the response.
	closestResultIndex := 0
	closestResultTimeDelta := responseModel.Hourly.Time[closestResultIndex]

	for i, timestamp := range responseModel.Hourly.Time {
		timeDelta := absInt(timestamp - int(time.Now().Unix()))

		if timeDelta < closestResultTimeDelta {
			closestResultTimeDelta = timeDelta
			closestResultIndex = i
		}
	}

	// Form a response and return it.
	return WeatherInformation{
		Temperature: responseModel.Hourly.Temperature[closestResultIndex],
		RainLevel:   responseModel.Hourly.Rain[closestResultIndex],
		WindSpeed:   responseModel.Hourly.Windspeed[closestResultIndex],
	}, nil
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
