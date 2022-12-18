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
	"fmt"
	"github.com/f1gopher/f1gopherlib/Messages"
	"github.com/f1gopher/f1gopherlib/connection"
	"image/color"
	"strconv"
	"time"
)

func (p *Parser) parseDriverList(dat map[string]interface{}, timestamp time.Time) {

	for driverNum, info := range dat {
		if driverNum == "_kf" {
			continue
		}

		record := info.(map[string]interface{})

		current, exists := p.driverTimes[driverNum]

		if !exists {
			number, _ := strconv.Atoi(driverNum)

			line := 0
			rawLine, exists := record["Line"]
			if exists {
				line = int(rawLine.(float64))
			}

			fullName, _ := record["FullName"].(string)
			shortName, _ := record["Tla"].(string)
			teamName, _ := record["TeamName"].(string)
			teamHexColour, _ := record["TeamColour"].(string)
			teamColor := color.RGBA{A: 0xFF}
			_, err := fmt.Sscanf(teamHexColour, "#%02x%02x%02x", &teamColor.R, &teamColor.G, &teamColor.B)
			if err != nil {
				p.ParseErrorf(connection.DriverListFile, timestamp, "Unable to parse team color: '%s', %v", teamColor, err)

				// Fallback to default colors
				teamColor.R = 0xFF
				teamColor.G = 0xFF
				teamColor.B = 0xFF
			}

			current = Messages.Timing{
				Number:    number,
				Position:  line,
				Name:      fullName,
				ShortName: shortName,
				Team:      teamName,
				HexColor:  "#" + teamHexColour,
				Color:     teamColor,
			}
		}

		p.driverTimes[driverNum] = current
	}
}
