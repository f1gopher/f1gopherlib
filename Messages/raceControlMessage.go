package Messages

import (
	"time"
)

type RaceControlMessage struct {
	Timestamp time.Time

	Msg  string
	Flag FlagState
}
