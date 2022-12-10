package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"strings"
	"time"
)

func (p *Parser) parseExtrapolatedClockData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, error) {

	var err error
	remaining := dat["Remaining"].(string)

	remaining = strings.Replace(remaining, ":", "h", 1)
	remaining = strings.Replace(remaining, ":", "m", 1)
	remaining = remaining + "s"
	p.eventState.RemainingTime, err = time.ParseDuration(remaining)
	if err != nil {
		p.ParseTimeError(connection.ExtrapolatedClockFile, timestamp, "Remaining", err)
	}

	extrapolating, exists := dat["Extrapolating"].(bool)
	if exists && extrapolating {
		abc, err := parseTime(dat["Utc"].(string))
		if err != nil {
			p.ParseTimeError(connection.ExtrapolatedClockFile, timestamp, "Utc", err)
		} else {
			p.eventState.SessionStartTime = abc
		}
	}

	if exists {
		p.eventState.ClockStopped = !extrapolating
	}

	p.eventState.Timestamp = timestamp

	return p.eventState, nil
}
