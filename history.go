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

package f1gopherlib

import (
	"github.com/f1gopher/f1gopherlib/Messages"
	"time"
)

func RaceHistory() []RaceEvent {
	result := make([]RaceEvent, 0)

	for _, session := range sessionHistory {
		sessionEnd := session.EventTime
		switch session.Type {
		case Messages.Practice1Session, Messages.Practice2Session, Messages.Practice3Session:
			sessionEnd = sessionEnd.Add(time.Hour * 1)

		case Messages.QualifyingSession:
			sessionEnd = sessionEnd.Add(time.Hour * 1)

		case Messages.SprintSession:
			sessionEnd = sessionEnd.Add(time.Hour * 1)

		case Messages.RaceSession:
			sessionEnd = sessionEnd.Add(time.Hour * 3)
		}

		if sessionEnd.Before(time.Now()) {
			result = append(result, session)
		}
	}

	return result
}

func NextSession() (RaceEvent, bool) {
	all := sessionHistory

	for x := 0; x < len(all); x++ {

		if all[x].EventTime.After(time.Now().In(time.UTC)) {
			continue
		}

		// No next session
		if x == 0 {
			return RaceEvent{}, false
		}

		return all[x-1], true
	}

	return RaceEvent{}, false
}

func liveEvent() (event RaceEvent, exists bool) {
	all := sessionHistory

	// TODO - handle timezones
	for x := 0; x < len(all); x++ {
		if all[x].EventTime.Year() == time.Now().Year() &&
			all[x].EventTime.Month() == time.Now().Month() &&
			all[x].EventTime.Day() == time.Now().Day() {

			// Start up to 45mins before thet start  of the event
			if time.Now().After(all[x].EventTime.Add(-time.Minute * 45)) {

				duringEvent := false
				switch all[x].Type {
				case Messages.Practice1Session, Messages.Practice2Session, Messages.Practice3Session:
					// Usually 60 mins but tire tests are 90 so cover both since it won't overlap with anything else
					duringEvent = time.Now().Before(all[x].EventTime.Add(time.Hour * 2))

				case Messages.QualifyingSession, Messages.SprintSession, Messages.RaceSession:
					// Last events in the day so just assume it's that event
					duringEvent = true

				default:
					panic("History: Unhandled session type: " + all[x].Type.String())
				}

				if duringEvent {
					return all[x], true
				}
			}
		}
	}

	return RaceEvent{}, false
}
