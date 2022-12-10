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
