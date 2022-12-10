package Messages

import (
	"time"
)

type Weather struct {
	Timestamp time.Time

	AirTemp       float64
	Humidity      float64
	AirPressure   float64
	Rainfall      bool
	TrackTemp     float64
	WindDirection float64
	WindSpeed     float64
}
