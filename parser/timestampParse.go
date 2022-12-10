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
	"strings"
	"time"
)

func parseTime(timestamp string) (time.Time, error) {

	dataTime, err := time.Parse("2006-01-02T15:04:05.999Z", timestamp)
	if err != nil {
		dataTime, err = time.Parse("2006-01-02T15:04:05", timestamp)
		if err != nil {
			return time.Parse("2006-01-02T15:04:05.9999999Z", timestamp)
		}
	}

	return dataTime, err
}

func parseDuration(duration string) (time.Duration, error) {
	previousValueStr := strings.Replace(duration, "+", "", 1)
	previousValueStr = strings.Replace(previousValueStr, ":", "m", 1)
	previousValueStr = strings.Replace(previousValueStr, ".", "s", 1)
	previousValueStr = previousValueStr + "ms"
	return time.ParseDuration(previousValueStr)
}
