package Messages

import (
	"image/color"
	"time"
)

type DriverInfo struct {
	StartPosition int
	Name          string
	ShortName     string
	Number        int
	Team          string
	HexColor      string
	Color         color.RGBA
}

type Drivers struct {
	Timestamp time.Time

	Drivers []DriverInfo
}
