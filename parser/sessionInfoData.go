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
	"time"
)

func (p *Parser) parseSessionInfoData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, []Messages.Timing, error) {

	info, _ := dat["Meeting"].(map[string]interface{})
	timingResult := make([]Messages.Timing, 0)

	p.eventState.Name = info["Name"].(string)
	// Key
	// OfficialName
	// Location
	// Country
	// Circuit
	// ArchiveStatus

	p.eventState.Heartbeat = true
	previousType := p.eventState.Type

	switch dat["Name"].(string) {
	case "Race":
		p.eventState.Type = Messages.Race
	case "Qualifying", "Sprint Qualifying", "Sprint Shootout":
		p.eventState.Type = Messages.Qualifying1
	case "Sprint":
		p.eventState.Type = Messages.Sprint
	case "Practice 1":
		p.eventState.Type = Messages.Practice1
	case "Practice 2":
		p.eventState.Type = Messages.Practice2
	case "Practice 3":
		p.eventState.Type = Messages.Practice3
	default:
		p.ParseErrorf(connection.SessionInfoFile, timestamp, "Unknown type: %s", dat["Name"].(string))
	}

	if previousType != p.eventState.Type {
		// Clear the chequered flag state for all cars
		for driverNum, driverInfo := range p.driverTimes {
			driverInfo.ChequeredFlag = false
			driverInfo.Sector1 = 0
			driverInfo.Sector2 = 0
			driverInfo.Sector3 = 0
			driverInfo.OverallFastestLap = false
			driverInfo.FastestLap = 0
			driverInfo.TimeDiffToPositionAhead = 0
			driverInfo.TimeDiffToFastest = 0
			driverInfo.GapToLeader = 0
			driverInfo.LastLap = 0
			driverInfo.SpeedTrap = 0
			driverInfo.SpeedTrapOverallFastest = false
			driverInfo.SpeedTrapPersonalFastest = false
			driverInfo.Sector1OverallFastest = false
			driverInfo.Sector1PersonalFastest = false
			driverInfo.Sector2OverallFastest = false
			driverInfo.Sector2PersonalFastest = false
			driverInfo.Sector3OverallFastest = false
			driverInfo.Sector3PersonalFastest = false
			for x := range driverInfo.Segment {
				driverInfo.Segment[x] = Messages.None
			}
			driverInfo.Location = Messages.NoLocation

			p.driverTimes[driverNum] = driverInfo

			timingResult = append(timingResult, driverInfo)
		}
	}

	// TODO handle: StartDate, GmtOffset, ArchiveStatus, Key, Name, EndDate, Path
	// TODO handle: Meeting: Key, OfficialName, Location, Country, Circuit

	p.eventState.Timestamp = timestamp

	return p.eventState, timingResult, nil
}
