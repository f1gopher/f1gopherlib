package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
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

		for driverId, car := range record.(map[string]interface{})["Cars"].(map[string]interface{}) {
			driverNum, _ := strconv.Atoi(driverId)

			t := Messages.Telemetry{
				Timestamp:    utcTimestamp,
				DriverNumber: driverNum,
			}

			for id, channel := range car.(map[string]interface{})["Channels"].(map[string]interface{}) {
				switch id {
				case "0": // RPM
					t.RPM = channel.(float64)
				case "2": // Speed
					t.Speed = channel.(float64)
				case "3": // Gear
					t.Gear = channel.(float64)
				case "4": // Throttle
					t.Throttle = channel.(float64)
				case "5": // Brake
					t.Brake = channel.(float64)
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
