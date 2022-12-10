package flowControl

import (
	"f1gopherlib/f1gopherlib/Messages"
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

	weather []Messages.Weather

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

func (f *straightThrough) IncrementLap() {}

func (f *straightThrough) IncrementTime(duration time.Duration) {}

func (f *straightThrough) SkipToSessionStart() {}

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
