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

package flowControl

import (
	"context"
	"github.com/f1gopher/f1gopherlib/Messages"
	"sync"
	"time"
)

type Flow interface {
	Run()

	AddWeather(weather Messages.Weather)
	AddRaceControlMessage(raceControl Messages.RaceControlMessage)
	AddTiming(timing Messages.Timing)
	AddEvent(timing Messages.Event)
	AddTelemetry(timing Messages.Telemetry)
	AddLocation(timing Messages.Location)
	AddRadio(timing Messages.Radio)
	AddDrivers(driver Messages.Drivers)

	IncrementLap()
	IncrementTime(duration time.Duration)
	SkipToSessionStart()
	TogglePause()
	IsPaused() bool
}

type FlowType int

const (
	Realtime FlowType = iota
	StraightThrough
)

func CreateFlowControl(
	ctx context.Context,
	wg *sync.WaitGroup,
	flowType FlowType,
	outputWeather chan<- Messages.Weather,
	outputRaceControlMessages chan<- Messages.RaceControlMessage,
	outputTimingMessages chan<- Messages.Timing,
	outputEvent chan<- Messages.Event,
	outputTelemetry chan<- Messages.Telemetry,
	outputLocation chan<- Messages.Location,
	outputEventTime chan<- Messages.EventTime,
	outputRadio chan<- Messages.Radio,
	outputDrivers chan<- Messages.Drivers) Flow {

	switch flowType {
	case Realtime:
		return &realtime{
			ctx:                       ctx,
			wg:                        wg,
			outputWeather:             outputWeather,
			outputRaceControlMessages: outputRaceControlMessages,
			outputTimingMessages:      outputTimingMessages,
			outputEvent:               outputEvent,
			outputTelemetry:           outputTelemetry,
			outputLocation:            outputLocation,
			outputEventTime:           outputEventTime,
			outputRadio:               outputRadio,
			outputDrivers:             outputDrivers,
		}

	case StraightThrough:
		return &straightThrough{
			outputWeather:             outputWeather,
			outputRaceControlMessages: outputRaceControlMessages,
			outputTimingMessages:      outputTimingMessages,
			outputEvent:               outputEvent,
			outputTelemetry:           outputTelemetry,
			outputLocation:            outputLocation,
			outputEventTime:           outputEventTime,
			outputRadio:               outputRadio,
			outputDrivers:             outputDrivers,
		}

	default:
		panic("Unhandled flow control type")
	}
}
