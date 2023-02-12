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
	"github.com/f1gopher/f1gopherlib/Messages"
	"time"
)

type straightThrough struct {
	outputWeather             chan<- Messages.Weather
	outputRaceControlMessages chan<- Messages.RaceControlMessage
	outputTimingMessages      chan<- Messages.Timing
	outputEvent               chan<- Messages.Event
	outputTelemetry           chan<- Messages.Telemetry
	outputLocation            chan<- Messages.Location
	outputEventTime           chan<- Messages.EventTime
	outputRadio               chan<- Messages.Radio
	outputDrivers             chan<- Messages.Drivers

	isPaused bool
}

func (f *straightThrough) Run() {
	// Don't need to do anything
}

func (f *straightThrough) AddWeather(weather Messages.Weather) {
	f.outputWeather <- weather
}

func (f *straightThrough) AddRaceControlMessage(raceControlMessage Messages.RaceControlMessage) {
	f.outputRaceControlMessages <- raceControlMessage
}

func (f *straightThrough) AddTiming(timing Messages.Timing) {
	f.outputTimingMessages <- timing
}

func (f *straightThrough) AddEvent(event Messages.Event) {
	f.outputEvent <- event

	f.outputEventTime <- Messages.EventTime{Timestamp: event.Timestamp}
}

func (f *straightThrough) AddTelemetry(telemetry Messages.Telemetry) {
	f.outputTelemetry <- telemetry
}

func (f *straightThrough) AddLocation(location Messages.Location) {
	f.outputLocation <- location
}

func (f *straightThrough) AddRadio(radio Messages.Radio) {
	f.outputRadio <- radio
}

func (f *straightThrough) AddDrivers(drivers Messages.Drivers) {
	f.outputDrivers <- drivers
}

func (f *straightThrough) IncrementLap() {}

func (f *straightThrough) IncrementTime(duration time.Duration) {}

func (f *straightThrough) SkipToSessionStart(start time.Time) {}

func (f *straightThrough) TogglePause() {
	f.isPaused = !f.isPaused
}

func (f *straightThrough) IsPaused() bool {
	return f.isPaused
}

func (f *straightThrough) IncrementDelay(delay time.Duration) {}

func (f *straightThrough) DecrementDelay(delay time.Duration) {}

func (f *straightThrough) Delay() time.Duration {
	return 0
}
