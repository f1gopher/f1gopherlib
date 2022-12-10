package Messages

import (
	"time"
)

type EventTime struct {
	Timestamp time.Time

	Remaining time.Duration
}
