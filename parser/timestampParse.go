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
