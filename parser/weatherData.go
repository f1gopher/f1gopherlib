// F1GopherLib - Copyright (C) 2022 f1gopher
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
