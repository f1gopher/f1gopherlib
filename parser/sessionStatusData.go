package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
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
