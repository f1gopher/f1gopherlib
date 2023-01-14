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
	"github.com/f1gopher/f1gopherlib/Messages"
	"github.com/f1gopher/f1gopherlib/connection"
	"strconv"
	"time"
)

func (p *Parser) parseCarData(dat map[string]interface{}, timestamp time.Time) ([]Messages.Telemetry, []Messages.Timing, error) {

	result := make([]Messages.Telemetry, 0)
	timingResult := make([]Messages.Timing, 0)

	entries := dat["Entries"].([]interface{})
	for _, record := range entries {

		timestampStr := record.(map[string]interface{})["Utc"].(string)
		utcTimestamp, err := parseTime(timestampStr)
		if err != nil {
			p.ParseTimeError(connection.CarDataFile, timestamp, "Utc", err)
		}
		localTimestamp := utcTimestamp.In(p.timezone)

		for driverId, car := range record.(map[string]interface{})["Cars"].(map[string]interface{}) {
			driverNum, _ := strconv.Atoi(driverId)

			t := Messages.Telemetry{
				Timestamp:    localTimestamp,
				DriverNumber: driverNum,
			}

			for id, channel := range car.(map[string]interface{})["Channels"].(map[string]interface{}) {
				switch id {
				case "0": // RPM
					t.RPM = int16(channel.(float64))
				case "2": // Speed
					t.Speed = float32(channel.(float64))
				case "3": // Gear
					t.Gear = byte(channel.(float64))
				case "4": // Throttle
					t.Throttle = float32(channel.(float64))
				case "5": // Brake
					t.Brake = float32(channel.(float64))
				case "45": // DRS
					driverInfo, _ := p.driverTimes[driverId]

					drsOpen := false
					drsValue := int(channel.(float64))
					if drsValue == 10 || drsValue == 12 || drsValue == 14 {
						t.DRS = true
						drsOpen = true
					} else {
						t.DRS = false
						drsOpen = false
					}

					if drsOpen != driverInfo.DRSOpen {
						driverInfo.DRSOpen = drsOpen
						p.driverTimes[driverId] = driverInfo
						timingResult = append(timingResult, driverInfo)
					}

				default:
					p.ParseErrorf(connection.CarDataFile, timestamp, "Unhandled channel id '%s'", id)
				}
			}

			result = append(result, t)
		}
	}

	return result, timingResult, nil
}
