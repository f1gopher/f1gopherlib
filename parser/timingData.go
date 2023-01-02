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
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (p *Parser) parseTimingData(dat map[string]interface{}, timestamp time.Time) ([]Messages.Timing, error) {

	result := make([]Messages.Timing, 0)

	lines, exists := dat["Lines"]
	if !exists {
		return result, nil
	}

	fastestLapChanged := false
	var currentFastestLap time.Duration
	var err error

	for driverNumber, data := range lines.(map[string]interface{}) {
		record := data.(map[string]interface{})

		currentDriver, exists := p.driverTimes[driverNumber]
		if !exists {
			continue
		}

		currentDriver.Timestamp = timestamp

		value, exists := record["Stopped"]

		intValue, exists := record["NumberOfPitStops"]
		if exists {
			currentDriver.Pitstops = int(intValue.(float64))
		}

		// TODO - handle line hee for display purposes?
		value, exists = record["Position"].(string)
		if exists {
			pos, _ := strconv.Atoi(value.(string))

			//line, _ := record["Line"].(int64)
			//if int(line) != pos {
			//	fmt.Println(fmt.Sprintf("line %d doesn't match position %d", line, pos))
			//}

			currentDriver.Position = pos
		}

		// TODO - do we ever get both values at the same time? Should we just use the value we get as the gap?
		value, exists = record["TimeDiffToFastest"].(string)
		if exists {
			var t time.Duration

			if len(value.(string)) > 0 {
				t, err = parseDuration(value.(string))
				if err != nil {
					p.ParseTimeError(connection.TimingDataFile, timestamp, "TimeDiffToFastest", err)
				}
			}

			currentDriver.TimeDiffToFastest = t
		}

		value, exists = record["TimeDiffToPositionAhead"].(string)
		if exists {
			var t time.Duration

			if len(value.(string)) > 0 {
				t, err = parseDuration(value.(string))
				if err != nil {
					p.ParseTimeError(connection.TimingDataFile, timestamp, "TimeDiffToPositionAhead", err)
				}
			}

			currentDriver.TimeDiffToPositionAhead = t
		}

		value, exists = record["GapToLeader"].(string)
		if exists {
			var t time.Duration

			if len(value.(string)) > 0 &&
				!strings.HasPrefix(value.(string), "LAP") &&
				!strings.HasSuffix(value.(string), "L") {

				t, err = parseDuration(value.(string))
				if err != nil {
					p.ParseTimeError(connection.TimingDataFile, timestamp, "GapToLeader", err)
				}
			}

			// TODO - handle gap being '1 L'

			// TODO -  Leaders lap number - use it somewhere else?

			currentDriver.GapToLeader = t
		}

		value, exists = record["IntervalToPositionAhead"].(map[string]interface{})
		if exists {
			interval, exists := value.(map[string]interface{})["Value"]

			if exists {

				strInterval := interval.(string)

				if len(strInterval) > 0 &&
					!strings.HasPrefix(strInterval, "LAP") &&
					!strings.HasSuffix(strInterval, "L") {

					t, err := parseDuration(strInterval)
					if err != nil {
						p.ParseTimeError(connection.TimingDataFile, timestamp, "IntervalToPositionAhead Value", err)
					}

					currentDriver.TimeDiffToPositionAhead = t
				}

			}
		} else {
			// For races the leader doesn't have there time ahead cleared
			if p.eventState.Type == Messages.Race && currentDriver.Position == 1 {
				currentDriver.TimeDiffToPositionAhead = 0
			}
		}

		stats, exists := record["Stats"].(map[string]interface{})
		if exists {
			// TODO - has per sector (0, 1, 2) data but do we care?
			for _, diff := range stats {
				value, exists = diff.(map[string]interface{})["TimeDiffToPositionAhead:"].(string)
				if exists {
					var t time.Duration

					if len(value.(string)) > 0 {
						t, err = parseDuration(value.(string))
						if err != nil {
							p.ParseTimeError(connection.TimingDataFile, timestamp, "TimeDiffToPositionAhead", err)
						}
					}

					currentDriver.TimeDiffToPositionAhead = t
				}

				value, exists = diff.(map[string]interface{})["TimeDiffToFastest"].(string)
				if exists {
					var t time.Duration

					if len(value.(string)) > 0 {
						t, err = parseDuration(value.(string))
						if err != nil {
							p.ParseTimeError(connection.TimingDataFile, timestamp, "TimeDiffToFastest", err)
						}
					}

					currentDriver.TimeDiffToFastest = t
				}
			}
		}

		// TODO - handle Status = 96, 608, 288, 800, 768, 576

		// Handle NumberOfLaps
		value, exists = record["NumberOfLaps"].(float64)
		if exists {
			currentDriver.Lap = int(value.(float64))
			if currentDriver.Location == Messages.OutLap {
				currentDriver.Location = Messages.OnTrack
			}
		}

		sectors, exists := record["Sectors"]

		// TODO - sectors count from 0?
		if exists {
			// Work out how many segments this track has
			if p.eventState.Sector1Segments == 0 {
				var one, two, three interface{}

				switch sectors.(type) {
				case map[string]interface{}:
					one = sectors.(map[string]interface{})["0"]
					two = sectors.(map[string]interface{})["1"]
					three = sectors.(map[string]interface{})["2"]

				case []interface{}:
					one = sectors.([]interface{})[0]
					two = sectors.([]interface{})[1]
					three = sectors.([]interface{})[2]

				default:
					p.ParseErrorf(connection.TimingDataFile, timestamp, "Unhandled data format: %v", dat)
				}

				// Older data doesn't have this
				if one != nil {
					// Older data doesn't have this
					_, exists = one.(map[string]interface{})["Segments"]
					if exists {

						if reflect.TypeOf(one.(map[string]interface{})["Segments"]).Kind() == reflect.Slice {
							p.eventState.Sector1Segments = len(one.(map[string]interface{})["Segments"].([]interface{}))
							p.eventState.Sector2Segments = len(two.(map[string]interface{})["Segments"].([]interface{}))
							p.eventState.Sector3Segments = len(three.(map[string]interface{})["Segments"].([]interface{}))
							p.eventState.TotalSegments = p.eventState.Sector1Segments + p.eventState.Sector2Segments + p.eventState.Sector3Segments
						}
					}
				}
			}

			switch sectors.(type) {
			case map[string]interface{}:
				for key, value2 := range sectors.(map[string]interface{}) {
					p.processSectorTimes(key, value2, &currentDriver, timestamp)
				}

			case []interface{}:

				for key, value2 := range sectors.([]interface{}) {
					p.processSectorTimes(string(key), value2, &currentDriver, timestamp)
				}

			default:
				p.ParseErrorf(connection.TimingDataFile, timestamp, "Unhandled data format: %v", dat)
			}
		}

		// Override location after reading if from the segments
		value, exists = record["Stopped"]
		if exists {
			if value.(bool) {
				currentDriver.Location = Messages.Stopped
			}
		}

		value, exists = record["Retired"].(string)
		if exists {
			if value.(bool) {
				currentDriver.Location = Messages.OutOfRace
			}
		}

		// We use the segments to work out when we are in the pitlane
		//
		//inPit, exists := record["InPit"]
		//if exists && inPit.(bool) {
		//		currentDriver.Location = Messages.Pitlane
		//}
		//
		//boolValue, exists := record["PitOut"]
		//if exists && boolValue.(bool) {
		//		currentDriver.Location = Messages.PitOut
		//}
		//
		//// TODO - we can do this earlier if we look at the segments and update then
		//// If we were Pit Out but now aren't then out lap
		//if exists && !boolValue.(bool) && currentDriver.Location == Messages.PitOut {
		//		currentDriver.Location = Messages.OutLap
		//}

		bestLapTime, exists := record["BestLapTime"].(map[string]interface{})
		if exists {
			var t time.Duration

			lapTime, exists := bestLapTime["Value"]
			if exists {
				if len(lapTime.(string)) > 0 {
					t, err = parseDuration(lapTime.(string))
					if err != nil {
						p.ParseTimeError(connection.TimingDataFile, timestamp, "BestLapTime Value", err)
					}
				}

				currentDriver.FastestLap = t
			}

			//lapTime, exists := bestLapTime["_deleted"]
			//if exists {
			//	// TODO - handle deleted lap time
			//}
		}

		lastLapTime, exists := record["LastLapTime"].(map[string]interface{})
		if exists {
			value, exists := lastLapTime["Value"]
			if exists && len(value.(string)) > 0 {
				t, err := parseDuration(value.(string))
				if err != nil {
					p.ParseTimeError(connection.TimingDataFile, timestamp, "LastLapTime Value", err)
				}

				currentDriver.LastLap = t
			}

			overallFastest, exists := lastLapTime["OverallFastest"]
			if exists {
				currentDriver.LastLapOverallFastest = overallFastest.(bool)
				if currentDriver.LastLapOverallFastest {
					fastestLapChanged = true
					currentFastestLap = currentDriver.LastLap
					currentDriver.OverallFastestLap = true
				}
			}

			personalFastest, exists := lastLapTime["PersonalFastest"]
			if exists {
				currentDriver.LastLapPersonalFastest = personalFastest.(bool)
			}
		}

		speeds, exists := record["Speeds"]
		if exists {
			// TODO - handle 'I1', 'I2', 'FL'

			speedTrap, exists := speeds.(map[string]interface{})["ST"]
			if exists {
				val, exists := speedTrap.(map[string]interface{})["Value"]
				if exists {
					st, _ := strconv.Atoi(val.(string))
					currentDriver.SpeedTrap = st // KM/hr
				}

				overallFastest, exists := speedTrap.(map[string]interface{})["OverallFastest"]
				if exists {
					currentDriver.SpeedTrapOverallFastest = overallFastest.(bool)
				}

				personalFastest, exists := speedTrap.(map[string]interface{})["PersonalFastest"]
				if exists {
					currentDriver.SpeedTrapPersonalFastest = personalFastest.(bool)
				}
			}
		}

		knockedOut, exists := record["KnockedOut"]
		if exists {
			currentDriver.KnockedOutOfQualifying = knockedOut.(bool)
		}

		p.driverTimes[driverNumber] = currentDriver

		result = append(result, currentDriver)
	}

	// Quali doesn't give us gap times so we have to calculate them when the overall fastest lap changes
	if fastestLapChanged && p.session == Messages.QualifyingSession {
		result = make([]Messages.Timing, 0)

		orderedDrivers := make([]Messages.Timing, 0)

		for _, info := range p.driverTimes {
			orderedDrivers = append(orderedDrivers, info)
		}

		sort.SliceStable(orderedDrivers, func(i, j int) bool {
			return orderedDrivers[i].FastestLap < orderedDrivers[j].FastestLap
		})

		for x := range orderedDrivers {
			// TODO - 2022 british quali - first driver doesn't match currentFastestLap
			if x == 0 { //orderedDrivers[x].FastestLap == currentFastestLap || orderedDrivers[x].FastestLap == 0 {
				orderedDrivers[x].TimeDiffToFastest = 0
				orderedDrivers[x].TimeDiffToPositionAhead = 0
			} else {
				// TODO - this value is sometimes 0 and it shouldn't be
				if orderedDrivers[x].FastestLap > 0 {
					orderedDrivers[x].TimeDiffToPositionAhead = orderedDrivers[x].FastestLap - orderedDrivers[x-1].FastestLap
					orderedDrivers[x].TimeDiffToFastest = orderedDrivers[x].FastestLap - currentFastestLap

					if orderedDrivers[x].TimeDiffToFastest < 0 {
						p.ParseErrorf(connection.TimingDataFile, timestamp, "TimeDiffToFastest < 0 '%v'", orderedDrivers[x].TimeDiffToFastest)
					}
				}
			}

			p.driverTimes[strconv.Itoa(orderedDrivers[x].Number)] = orderedDrivers[x]
			result = append(result, orderedDrivers[x])
		}
	} else if fastestLapChanged && p.session == Messages.RaceSession || p.session == Messages.SprintSession {
		// For races we need to know who has the overall fastest lap
		result = make([]Messages.Timing, 0)
		for x, info := range p.driverTimes {
			info.OverallFastestLap = info.FastestLap == currentFastestLap
			p.driverTimes[strconv.Itoa(p.driverTimes[x].Number)] = info
			result = append(result, p.driverTimes[x])
		}
	}

	return result, nil
}

func (p *Parser) processSectorTimes(key string, value interface{}, driver *Messages.Timing, timestamp time.Time) {

	segments, exists := value.(map[string]interface{})["Segments"]
	if exists {
		segmentState := Messages.None
		currentSegmentIndex := 0
		useSegmentChange := true

		if reflect.TypeOf(segments).Kind() == reflect.Slice {
			for x, info := range segments.([]interface{}) {
				currentSegmentIndex = x
				segmentState, useSegmentChange = p.calcSegment(key, info, timestamp, currentSegmentIndex, driver)

				if useSegmentChange {
					p.updateLocation(driver, segmentState)
				}
			}

		} else if reflect.TypeOf(segments).Kind() == reflect.Map {
			for x, info := range segments.(map[string]interface{}) {
				currentSegmentIndex, _ = strconv.Atoi(x)
				segmentState, useSegmentChange = p.calcSegment(key, info, timestamp, currentSegmentIndex, driver)

				if useSegmentChange {
					p.updateLocation(driver, segmentState)
				}
			}
		}
	}

	previousValue, exists := value.(map[string]interface{})["Value"]
	if exists {
		var sectorTime time.Duration
		var err error

		if len(previousValue.(string)) > 0 {
			sectorTime, err = parseDuration(previousValue.(string))
			if err != nil {
				p.ParseTimeError(connection.TimingDataFile, timestamp, "Sector Value", err)
			}
		}

		switch key {
		case "0":
			driver.Sector1 = sectorTime

		case "1":
			driver.Sector2 = sectorTime

		case "2":
			driver.Sector3 = sectorTime

			if p.eventState.TrackStatus == Messages.ChequeredFlag {
				driver.ChequeredFlag = true
			}

			if sectorTime != 0 {
				driver.LapsOnTire++
			}
		}
	}

	isFastest, exists := value.(map[string]interface{})["OverallFastest"]
	if exists {
		switch key {
		case "0":
			driver.Sector1OverallFastest = isFastest.(bool)

		case "1":
			driver.Sector2OverallFastest = isFastest.(bool)

		case "2":
			driver.Sector3OverallFastest = isFastest.(bool)
		}
	}

	isFastest, exists = value.(map[string]interface{})["PersonalFastest"]
	if exists {
		switch key {
		case "0":
			driver.Sector1PersonalFastest = isFastest.(bool)
		case "1":
			driver.Sector2PersonalFastest = isFastest.(bool)
		case "2":
			driver.Sector3PersonalFastest = isFastest.(bool)
		}
	}
}

func (p *Parser) updateLocation(driver *Messages.Timing, segmentState Messages.SegmentType) {

	// Sometimes it is none when we are on track so leave location as is
	if segmentState == Messages.None && (driver.Location != Messages.OutLap && driver.Location != Messages.OnTrack) {
		driver.Location = Messages.Pitlane
	} else if segmentState == Messages.PitlaneSegment {
		driver.Location = Messages.Pitlane
	} else {
		if driver.Segment[0] == Messages.PitlaneSegment || driver.Segment[1] == Messages.PitlaneSegment {
			driver.Location = Messages.OutLap
		} else {
			driver.Location = Messages.OnTrack
		}
	}
}

func (p *Parser) calcSegment(
	key string,
	info interface{},
	timestamp time.Time,
	currentSegment int,
	driver *Messages.Timing) (Messages.SegmentType, bool) {

	segmentState := Messages.None
	abc := info.(map[string]interface{})
	status := int(abc["Status"].(float64))
	useSegmentChange := false

	segmentIndex := currentSegment
	if key == "1" {
		segmentIndex = p.eventState.Sector1Segments + currentSegment
	} else if key == "2" {
		segmentIndex = p.eventState.Sector1Segments + p.eventState.Sector2Segments + currentSegment
	}

	if status != 0 {
		switch status {
		case 2048:
			segmentState = Messages.YellowSegment
		case 2049:
			segmentState = Messages.GreenSegment
		case 2050:
			segmentState = Messages.InvalidSegment
		case 2051:
			segmentState = Messages.PurpleSegment
		case 2052:
			segmentState = Messages.RedSegment
		case 2064:
			segmentState = Messages.PitlaneSegment
		case 2065:
			segmentState = Messages.Mystery2
		case 2066:
			segmentState = Messages.Mystery3
		case 2068:
			segmentState = Messages.Mystery
		default:
			p.ParseErrorf(connection.TimingDataFile, timestamp, "Unhandled segment state value: %d", status)
		}

		useSegmentChange = segmentIndex >= driver.PreviousSegmentIndex ||
			// if previous was sector 3 and new one is sector 1
			(driver.PreviousSegmentIndex > (p.eventState.Sector1Segments+p.eventState.Sector2Segments) &&
				segmentIndex < p.eventState.Sector1Segments)

		switch key {
		case "0":
			// If the last segment was in the third sector then we have started a new lap so clear everything
			if driver.PreviousSegmentIndex > (p.eventState.Sector1Segments + p.eventState.Sector2Segments) {
				for y := 0; y < len(driver.Segment); y++ {
					driver.Segment[y] = Messages.None
				}
			}

			driver.Segment[segmentIndex] = segmentState
			if segmentIndex > driver.PreviousSegmentIndex ||
				driver.PreviousSegmentIndex > (p.eventState.Sector1Segments+p.eventState.Sector2Segments) {
				driver.PreviousSegmentIndex = segmentIndex
			}

		case "1":
			driver.Segment[segmentIndex] = segmentState

			if segmentIndex > driver.PreviousSegmentIndex {
				driver.PreviousSegmentIndex = segmentIndex
			}
		case "2":
			// If we get late data and we have already started a new lap then ignore
			if !(driver.PreviousSegmentIndex < p.eventState.Sector1Segments) {
				driver.Segment[segmentIndex] = segmentState

				if segmentIndex > driver.PreviousSegmentIndex {
					driver.PreviousSegmentIndex = segmentIndex
				}
			}
		}
	}

	return segmentState, useSegmentChange
}
