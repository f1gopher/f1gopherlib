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

		// Don't use this to update the driver position because it results in multiple drivers
		// with the same position.
		//
		//value, exists = line.(map[string]interface{})["Line"]
		//if exists {
		//	currentDriver.Position = int(value.(float64))
		//}

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
