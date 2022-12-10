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
