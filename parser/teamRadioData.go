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
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"reflect"
	"time"
)

func (p *Parser) parseTeamRadioData(dat map[string]interface{}, timestamp time.Time) ([]Messages.Radio, error) {

	result := make([]Messages.Radio, 0)

	if reflect.TypeOf(dat["Captures"]).Kind() == reflect.Map {
		for _, data := range dat["Captures"].(map[string]interface{}) {
			p.readTeamRadio(data, timestamp, &result)
		}

	} else if reflect.TypeOf(dat["Captures"]).Kind() == reflect.Slice {
		for _, data := range dat["Captures"].([]interface{}) {
			p.readTeamRadio(data, timestamp, &result)
		}
	} else {
		p.ParseErrorf(connection.TeamRadioFile, timestamp, "Unhandled data format: %v", dat)
	}

	return result, nil
}

func (p *Parser) readTeamRadio(data interface{}, timestamp time.Time, result *[]Messages.Radio) {
	record := data.(map[string]interface{})

	time := record["Utc"].(string)
	driverNumber := record["RacingNumber"].(string)
	path := record["Path"].(string)

	radio, err := p.assets.TeamRadio(path)

	if err == nil {

		msgTime, err := parseTime(time)
		if err != nil {
			p.ParseTimeError(connection.TeamRadioFile, timestamp, "Utc", err)
			return
		}

		msg := Messages.Radio{
			Timestamp: msgTime,
			Driver:    p.driverTimes[driverNumber].Name,
			Msg:       radio,
		}

		*result = append(*result, msg)
	}
}
