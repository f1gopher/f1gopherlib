package Messages

import (
	"time"
)

type Radio struct {
	Timestamp time.Time

	Driver string
	Msg    []byte
}
