package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"time"
)

func (p *Parser) parseCurrentLapData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, error) {

	value, exists := dat["CurrentLap"]
	if exists {
		p.eventState.CurrentLap = int(value.(float64))
	}

	value, exists = dat["TotalLaps"]
	if exists {
		p.eventState.TotalLaps = int(value.(float64))
	}

	p.eventState.Timestamp = timestamp

	return p.eventState, nil
}
