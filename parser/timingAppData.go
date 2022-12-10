package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"strconv"
	"time"
)

func (p *Parser) parseTimingAppData(dat map[string]interface{}, timestamp time.Time) ([]Messages.Timing, error) {

	result := make([]Messages.Timing, 0)

	for driverStr, line := range dat["Lines"].(map[string]interface{}) {

		currentDriver, exists := p.driverTimes[driverStr]
		if !exists {
			continue
		}
		currentDriver.Timestamp = timestamp

		value, exists := line.(map[string]interface{})["GridPos"]
		if exists {
			value, _ := strconv.ParseInt(value.(string), 10, 8)
			currentDriver.Position = int(value)
		}

		value, exists = line.(map[string]interface{})["Line"]
		if exists {
			currentDriver.Position = int(value.(float64))
		}

		value, exists = line.(map[string]interface{})["Stints"]
		if exists {

			switch value.(type) {
			case map[string]interface{}:
				for _, stintData := range value.(map[string]interface{}) {
					p.readTimingAppData(stintData, &currentDriver, timestamp)
				}

			case []interface{}:
				for _, stintData := range value.([]interface{}) {
					p.readTimingAppData(stintData, &currentDriver, timestamp)
				}

			default:
				p.ParseErrorf(connection.TimingAppDataFile, timestamp, "Unhandled data format: %v", dat)
			}
		}

		p.driverTimes[driverStr] = currentDriver

		result = append(result, currentDriver)
	}

	return result, nil
}

func (p *Parser) readTimingAppData(stintData interface{}, currentDriver *Messages.Timing, timestamp time.Time) {
	tyre, hasTyre := stintData.(map[string]interface{})["Compound"]
	if !hasTyre {
		return
	}

	switch tyre.(string) {
	case "SOFT":
		currentDriver.Tire = Messages.Soft
	case "MEDIUM":
		currentDriver.Tire = Messages.Medium
	case "HARD":
		currentDriver.Tire = Messages.Hard
	case "INTERMEDIATE":
		currentDriver.Tire = Messages.Intermediate
	case "WET":
		currentDriver.Tire = Messages.Wet
	case "UNKNOWN", "C": // Apparently a thing!
		currentDriver.Tire = Messages.Unknown
	case "TEST", "TEST_UNKNOWN":
		currentDriver.Tire = Messages.Test
	case "HYPERSOFT":
		currentDriver.Tire = Messages.HYPERSOFT
	case "SUPERSOFT":
		currentDriver.Tire = Messages.SUPERSOFT
	case "ULTRASOFT":
		currentDriver.Tire = Messages.ULTRASOFT
	default:
		p.ParseErrorf(connection.TimingAppDataFile, timestamp, "Unhandled Compound '%s'", tyre.(string))
	}

	//drivers[driverNumber].PitStops = append(drivers[driverNumber].PitStops, driver.PitStop{
	//	Lap: drivers[driverNumber].Lap,
	//})

	// TODO - Handle: LapFlags, New, TyresNotChanged, TotalLaps, StartLaps

	totalLaps, exists := stintData.(map[string]interface{})["TotalLaps"]
	if exists {
		currentDriver.LapsOnTire = int(totalLaps.(float64))
	}
}
