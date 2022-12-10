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

func (p *Parser) parseHeartbeatData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, error) {

	time := dat["Utc"].(string)
	value, err := parseTime(time)
	if err != nil {
		p.ParseTimeError(connection.HeartbeatFile, timestamp, "Utc", err)
	} else {
		p.eventState.Timestamp = value
	}

	// TODO - ignore _kf or flag when no heartbeat recieved?

	p.eventState.Heartbeat = true

	return p.eventState, nil
}
