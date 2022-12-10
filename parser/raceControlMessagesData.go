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
	"time"
)

func (p *Parser) parseRaceControlMessagesData(dat map[string]interface{}, timestamp time.Time) ([]Messages.RaceControlMessage, []Messages.Event, []Messages.Timing, error) {

	result := make([]Messages.RaceControlMessage, 0)
	eventResult := make([]Messages.Event, 0)
	timingResult := make([]Messages.Timing, 0)

	if reflect.TypeOf(dat["Messages"]).Kind() == reflect.Slice {
		for _, msg := range dat["Messages"].([]interface{}) {
			p.readRaceControlMessage(msg, timestamp, &result, &eventResult, &timingResult)
		}
	} else if reflect.TypeOf(dat["Messages"]).Kind() == reflect.Map {
		for _, msg := range dat["Messages"].(map[string]interface{}) {
			p.readRaceControlMessage(msg, timestamp, &result, &eventResult, &timingResult)
		}
	} else {
		p.ParseErrorf(connection.RaceControlMessagesFile, timestamp, "Unhandled data format: %v", dat)
	}

	return result, eventResult, timingResult, nil
}

func (p *Parser) readRaceControlMessage(
	msg interface{},
	timestamp time.Time,
	result *[]Messages.RaceControlMessage,
	eventResult *[]Messages.Event,
	timingResult *[]Messages.Timing) {

	time, err := parseTime(msg.(map[string]interface{})["Utc"].(string))
	if err != nil {
		p.ParseTimeError(connection.RaceControlMessagesFile, timestamp, "Utc", err)
		return
	}

	// Lap
	// Category
	status := msg.(map[string]interface{})["Message"].(string)

	//status := msg.(map[string]interface{})["Message"].(string)

	flagTxt, exists := msg.(map[string]interface{})["Flag"].(string)
	flag := Messages.NoFlag
	if exists {
		switch flagTxt {
		case "BLUE":
			flag = Messages.BlueFlag
		case "YELLOW":
			flag = Messages.YellowFlag
		case "DOUBLE YELLOW":
			flag = Messages.DoubleYellowFlag
		case "CHEQUERED":
			flag = Messages.ChequeredFlag
		case "CLEAR", "GREEN":
			flag = Messages.GreenFlag
		case "RED":
			flag = Messages.RedFlag
		case "BLACK AND WHITE":
			flag = Messages.BlackAndWhite
		}
	}

	*result = append(*result, Messages.RaceControlMessage{
		Timestamp: time,
		Msg:       status,
		Flag:      flag,
	})

	switch status {
	case "GREEN LIGHT - PIT EXIT OPEN":
		p.eventState.PitExitOpen = true
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "RED LIGHT - PIT EXIT CLOSED":
		p.eventState.PitExitOpen = false
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "VIRTUAL SAFETY CAR DEPLOYED":
		p.eventState.SafetyCar = Messages.VirtualSafetyCar
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "VIRTUAL SAFETY CAR ENDING":
		p.eventState.SafetyCar = Messages.VirtualSafetyCarEnding
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "SAFETY CAR DEPLOYED":
		p.eventState.SafetyCar = Messages.SafetyCar
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "SAFETY CAR IN THIS LAP":
		p.eventState.SafetyCar = Messages.SafetyCarEnding
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "DRS ENABLED":
		p.eventState.DRSEnabled = Messages.DRSEnabled
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

	case "DRS DISABLED":
		p.eventState.DRSEnabled = Messages.DRSDisabled
		p.eventState.Timestamp = time
		*eventResult = append(*eventResult, p.eventState)

		//default:
		//	fmt.Println("Unhandled RC: " + status)
	}

	if exists {
		scope, _ := msg.(map[string]interface{})["Scope"].(string)
		sectorNum := 0
		if scope == "Sector" {
			sectorNum = int(msg.(map[string]interface{})["Sector"].(float64))
			sectorNum -= 1 // 0 indexing
			// TODO - 2021 - Saudi Arabia Qualifying uses sector 0 so 1 isn't the first sector?
		}

		switch flagTxt {
		case "RED":
			if scope == "Track" {
				p.eventState.TrackStatus = Messages.RedFlag
			}
			if scope == "Sector" {
				p.eventState.SegmentFlags[sectorNum] = Messages.RedFlag
			}
			p.eventState.Timestamp = time
			*eventResult = append(*eventResult, p.eventState)

		case "YELLOW":
			if scope == "Track" {
				p.eventState.TrackStatus = Messages.YellowFlag
			}
			if scope == "Sector" && sectorNum >= 0 {
				p.eventState.SegmentFlags[sectorNum] = Messages.YellowFlag
			}
			p.eventState.Timestamp = time
			*eventResult = append(*eventResult, p.eventState)

		case "DOUBLE YELLOW":
			if scope == "Track" {
				p.eventState.TrackStatus = Messages.DoubleYellowFlag
			}
			if scope == "Sector" && sectorNum >= 0 {
				p.eventState.SegmentFlags[sectorNum] = Messages.DoubleYellowFlag
			}
			p.eventState.Timestamp = time
			*eventResult = append(*eventResult, p.eventState)

		case "GREEN", "CLEAR":
			if scope == "Track" {
				p.eventState.TrackStatus = Messages.GreenFlag
				p.eventState.SafetyCar = Messages.Clear
			}
			if scope == "Sector" {
				if sectorNum >= 0 {
					p.eventState.SegmentFlags[sectorNum] = Messages.GreenFlag
				}
			} else {
				for x := range p.eventState.SegmentFlags {
					p.eventState.SegmentFlags[x] = Messages.GreenFlag
				}
			}
			p.eventState.Timestamp = time
			*eventResult = append(*eventResult, p.eventState)

		case "CHEQUERED":
			p.eventState.TrackStatus = Messages.ChequeredFlag
			p.eventState.Timestamp = time
			*eventResult = append(*eventResult, p.eventState)

			// Anyone in the pitlane when chequered flag waves for quali will be out
			if p.eventState.Type == Messages.Qualifying1 ||
				p.eventState.Type == Messages.Qualifying2 ||
				p.eventState.Type == Messages.Qualifying3 {

				for x, driver := range p.driverTimes {
					if driver.Location == Messages.Pitlane || driver.Location == Messages.PitOut {
						driver.ChequeredFlag = true
						p.driverTimes[x] = driver

						*timingResult = append(*timingResult, driver)
					}
				}
			}
		}
	}

	// TODO - DRS enabled/disabled
	// TODO - yellow flags per sector
}
