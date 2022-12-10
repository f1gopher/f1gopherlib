package Messages

import (
	"time"
)

type Telemetry struct {
	Timestamp    time.Time
	DriverNumber int

	RPM      float64
	Speed    float64
	Gear     float64
	Throttle float64
	Brake    float64
	DRS      bool
}
