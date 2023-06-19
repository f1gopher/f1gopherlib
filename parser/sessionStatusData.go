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

func (p *Parser) parseSessionStatusData(dat map[string]interface{}, timestamp time.Time) (Messages.Event, error) {

	status := dat["Status"].(string)

	switch status {
	case "Inactive":
		p.eventState.Status = Messages.Inactive
	case "Started":
		p.eventState.Status = Messages.Started
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

	return p.eventState, nil
}
