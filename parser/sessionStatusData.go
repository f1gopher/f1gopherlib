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

func (p *Parser) parseSessionStatusData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, []Messages.Timing, error) {

	status := dat["Status"].(string)

	timingResult := make([]Messages.Timing, 0)

	switch status {
	case "Inactive":
		p.eventState.Status = Messages.Inactive

	case "Started":
		p.eventState.Status = Messages.Started

		// Clear the chequered flag state for all cars
		for driverNum, info := range p.driverTimes {
			info.ChequeredFlag = false
			info.Sector1 = 0
			info.Sector2 = 0
			info.Sector3 = 0
			info.OverallFastestLap = false
			info.FastestLap = 0
			info.TimeDiffToPositionAhead = 0
			info.TimeDiffToFastest = 0
			info.GapToLeader = 0
			info.LastLap = 0
			info.SpeedTrap = 0
			info.SpeedTrapOverallFastest = false
			info.SpeedTrapPersonalFastest = false
			info.Sector1OverallFastest = false
			info.Sector1PersonalFastest = false
			info.Sector2OverallFastest = false
			info.Sector2PersonalFastest = false
			info.Sector3OverallFastest = false
			info.Sector3PersonalFastest = false
			for x := range info.Segment {
				info.Segment[x] = Messages.None
			}
			info.Location = Messages.NoLocation

			p.driverTimes[driverNum] = info

			timingResult = append(timingResult, info)
		}

	case "Aborted":
		p.eventState.Status = Messages.Aborted
	case "Finished":
		p.eventState.Status = Messages.Finished
	case "Finalised":
		p.eventState.Status = Messages.Finalised
	case "Ends":
		p.eventState.Status = Messages.Ended
	default:
		p.ParseErrorf(connection.SessionStatusFile, timestamp, "SessionStatus: Unhandled Status '%s'", status)
	}

	p.eventState.Timestamp = timestamp

	return p.eventState, timingResult, nil
}
