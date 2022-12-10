package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"time"
)

func (p *Parser) parseSessionInfoData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, error) {

	info, _ := dat["Meeting"].(map[string]interface{})

	p.eventState.Name = info["Name"].(string)
	// Key
	// OfficialName
	// Location
	// Country
	// Circuit
	// ArchiveStatus

	p.eventState.Heartbeat = true

	switch dat["Name"].(string) {
	case "Race":
		p.eventState.Type = Messages.Race
	case "Qualifying", "Sprint Qualifying":
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
		p.ParseErrorf(connection.SessionInfoFile, timestamp, "Unknown type: ", dat["Type"].(string))
	}

	// TODO handle: StartDate, GmtOffset, ArchiveStatus, Key, Name, EndDate, Path
	// TODO handle: Meeting: Key, OfficialName, Location, Country, Circuit

	p.eventState.Timestamp = timestamp

	return p.eventState, nil
}
