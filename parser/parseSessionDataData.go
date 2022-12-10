package parser

import (
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"fmt"
	"reflect"
	"time"
)

func (p *Parser) parseSessionDataData(dat map[string]interface{}, timestamp time.Time) ([]Messages.Event, error) {

	var result []Messages.Event

	if dat["Series"] != nil && reflect.TypeOf(dat["Series"]).Kind() == reflect.Map {
		for _, series := range dat["Series"].(map[string]interface{}) {

			var text string
			msg, exists := series.(map[string]interface{})["SessionStatus"]
			if exists {
				text = msg.(string)

			} else {
				msg, exists = series.(map[string]interface{})["QualifyingPart"]
				if exists {
					text = fmt.Sprintf("%g", msg.(float64))

					switch text {
					case "0":
						p.eventState.Type = Messages.Qualifying0
					case "1":
						p.eventState.Type = Messages.Qualifying1
					case "2":
						p.eventState.Type = Messages.Qualifying2
					case "3":
						p.eventState.Type = Messages.Qualifying3
					default:
						p.ParseErrorf(connection.SessionDataFile, timestamp, "SessionData: Unhandled value for QualifyingPart '%s'", text)
					}
				}
			}

			time := series.(map[string]interface{})["Utc"].(string)
			value, err := parseTime(time)
			if err != nil {
				p.ParseTimeError(connection.SessionDataFile, timestamp, "Utc", err)
			} else {
				p.eventState.Timestamp = value
			}

			result = append(result, p.eventState)
		}
	} else {
		if reflect.TypeOf(dat["StatusSeries"]).Kind() == reflect.Map {
			for _, series := range dat["StatusSeries"].(map[string]interface{}) {
				time := series.(map[string]interface{})["Utc"].(string)

				value, err := parseTime(time)
				if err != nil {
					p.ParseTimeError(connection.SessionDataFile, timestamp, "Utc", err)
				} else {
					p.eventState.Timestamp = value
				}

				result = append(result, p.eventState)
			}
		}
	}

	return result, nil
}
