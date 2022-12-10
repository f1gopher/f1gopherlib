package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"strconv"
	"time"
)

func parseWeatherData(dat map[string]interface{}, timestamp time.Time) (Messages.Weather, error) {

	airTemp, _ := strconv.ParseFloat(dat["AirTemp"].(string), 8)
	humidity, _ := strconv.ParseFloat(dat["Humidity"].(string), 8)
	pressure, _ := strconv.ParseFloat(dat["Pressure"].(string), 8)
	rainfall, _ := strconv.ParseBool(dat["Rainfall"].(string))
	trackTemp, _ := strconv.ParseFloat(dat["TrackTemp"].(string), 8)
	windDirection, _ := strconv.ParseFloat(dat["WindDirection"].(string), 8)
	windSpeed, _ := strconv.ParseFloat(dat["WindSpeed"].(string), 8)

	return Messages.Weather{
		Timestamp:     timestamp,
		AirTemp:       airTemp,
		Humidity:      humidity,
		AirPressure:   pressure,
		Rainfall:      rainfall,
		TrackTemp:     trackTemp,
		WindDirection: windDirection,
		WindSpeed:     windSpeed,
	}, nil
}
