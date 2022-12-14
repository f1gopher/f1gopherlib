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

func HappeningSessions() (liveSession RaceEvent, nextSession RaceEvent, hasLiveSession bool, hasNextSession bool) {
	all := sessionHistory
	utcNow := time.Now().UTC()

	for x := 0; x < len(all); x++ {
		// If we are the same day as a session see if it is live
		if all[x].EventTime.Year() == utcNow.Year() &&
			all[x].EventTime.Month() == utcNow.Month() &&
			all[x].EventTime.Day() == utcNow.Day() {

			// Start up to 45mins before the start  of the event
			if utcNow.After(all[x].EventTime.Add(-time.Minute * 45)) {

				duringEvent := false
				switch all[x].Type {
				case Messages.Practice1Session, Messages.Practice2Session, Messages.Practice3Session:
					// Usually 60 mins but tire tests are 90 so cover both since it won't overlap with anything else
					duringEvent = utcNow.Before(all[x].EventTime.Add(time.Hour * 2))

				case Messages.QualifyingSession, Messages.SprintSession, Messages.RaceSession:
					// Last events in the day so just assume it's that event
					duringEvent = true

				default:
					panic("History: Unhandled session type: " + all[x].Type.String())
				}

				if duringEvent {
					if x == 0 {
						return all[x], RaceEvent{}, true, false
					}

					return all[x], all[x-1], true, true
				}
			}
		} else if utcNow.Before(all[x].EventTime) {
			// If this is the first session that is after now then nothing is live and this is the next session
			return RaceEvent{}, all[x], false, true
		} else {
			// No live or upcoming sessions
			break
		}
	}

	return RaceEvent{}, RaceEvent{}, false, false
}

func liveEvent() (event RaceEvent, exists bool) {
	live, _, exists, _ := HappeningSessions()
	return live, exists
}
