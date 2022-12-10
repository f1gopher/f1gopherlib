package Messages

import (
	"time"
)

type Location struct {
	Timestamp time.Time

	DriverNumber int
	X            float64
	Y            float64
	Z            float64
}
