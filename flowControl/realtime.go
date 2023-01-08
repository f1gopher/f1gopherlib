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

type realtime struct {
	outputWeather             chan<- Messages.Weather
	outputRaceControlMessages chan<- Messages.RaceControlMessage
	outputTimingMessages      chan<- Messages.Timing
	outputEvent               chan<- Messages.Event
	outputTelemetry           chan<- Messages.Telemetry
	outputLocation            chan<- Messages.Location
	outputEventTime           chan<- Messages.EventTime
	outputRadio               chan<- Messages.Radio

	weatherLock     sync.Mutex
	weather         []Messages.Weather
	raceControlLock sync.Mutex
	raceControl     []Messages.RaceControlMessage
	timingLock      sync.Mutex
	timing          []Messages.Timing
	eventLock       sync.Mutex
	event           []Messages.Event
	telemetryLock   sync.Mutex
	telemetry       []Messages.Telemetry
	locationLock    sync.Mutex
	location        []Messages.Location
	radioLock       sync.Mutex
	radio           []Messages.Radio

	currentTime   time.Time
	currentLap    int
	currentStatus Messages.SessionState
	remainingTime time.Duration
	clockStopped  bool

	incrementLapCount int
	incrementTime     time.Duration
	skipToStart       bool
	isPaused          bool

	sessionStart  time.Time
	sessionLength time.Duration

	ctx context.Context
	wg  *sync.WaitGroup
}

func (f *realtime) Run() {
	f.wg.Add(1)
	defer f.wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond)
	counter := 2

	for {
		select {
		case <-f.ctx.Done():
			ticker.Stop()
			return

		case <-ticker.C:
			if f.isPaused {
				continue
			}

			if f.skipToStart {
				f.eventLock.Lock()
				if len(f.event) > 0 {
					needToStopAndStart := f.event[0].Status == Messages.Started

					for x := range f.event {
						if f.event[x].Timestamp.IsZero() {
							continue
						}

						if f.event[x].Status == Messages.Started && !needToStopAndStart {
							f.currentTime = f.event[x].Timestamp
							break
						}

						if f.event[x].Status != Messages.Started {
							needToStopAndStart = false
						}
					}

					// We want to skip any radio messages when we jump forward in time
					f.radioLock.Lock()
					for len(f.radio) > 0 && (f.radio[0].Timestamp.Before(f.currentTime) || f.radio[0].Timestamp.Equal(f.currentTime)) {
						f.radio = f.radio[1:]
					}
					f.radioLock.Unlock()
				}
				f.eventLock.Unlock()

				f.skipToStart = false
			}

			if counter == 3 {
				counter = 0

				f.eventLock.Lock()
				if len(f.event) > 0 {

					if f.currentTime.IsZero() && !f.event[0].Timestamp.IsZero() {
						f.currentTime = f.event[0].Timestamp
						f.clockStopped = f.event[0].ClockStopped
					}

					increment := f.incrementLapCount

					if increment > 0 {
						// TODO - thread safe
						f.incrementLapCount = f.incrementLapCount - increment
						targetLap := f.currentLap + increment
						var incrementTime time.Time

						for len(f.event) > 0 && f.currentLap < targetLap {
							select {
							case f.outputEvent <- f.event[0]:
								f.currentLap = f.event[0].CurrentLap
								f.currentStatus = f.event[0].Status
								incrementTime = f.event[0].Timestamp

								f.sessionStart = f.event[0].SessionStartTime
								f.sessionLength = f.event[0].RemainingTime
								f.clockStopped = f.event[0].ClockStopped

							default:
								// Data loss
							}

							f.event = f.event[1:]
						}

						f.currentTime = incrementTime

						// We want to skip any radio messages when we jump forward in time
						f.radioLock.Lock()
						for len(f.radio) > 0 && (f.radio[0].Timestamp.Before(f.currentTime) || f.radio[0].Timestamp.Equal(f.currentTime)) {
							f.radio = f.radio[1:]
						}
						f.radioLock.Unlock()

					} else {
						for len(f.event) > 0 && (f.event[0].Timestamp.Before(f.currentTime) || f.event[0].Timestamp.Equal(f.currentTime)) {
							select {
							case f.outputEvent <- f.event[0]:
								f.currentLap = f.event[0].CurrentLap
								f.currentStatus = f.event[0].Status

								f.sessionStart = f.event[0].SessionStartTime
								f.sessionLength = f.event[0].RemainingTime
								f.clockStopped = f.event[0].ClockStopped

							default:
								// Data loss
							}

							f.event = f.event[1:]
						}
					}
				}
				f.eventLock.Unlock()

				f.raceControlLock.Lock()
				if len(f.raceControl) > 0 {
					for len(f.raceControl) > 0 && (f.raceControl[0].Timestamp.Before(f.currentTime) || f.raceControl[0].Timestamp.Equal(f.currentTime)) {
						select {
						case f.outputRaceControlMessages <- f.raceControl[0]:
						default:
							// Data loss
						}

						f.raceControl = f.raceControl[1:]
					}
				}
				f.raceControlLock.Unlock()

				f.weatherLock.Lock()
				if len(f.weather) > 0 {
					for len(f.weather) > 0 && (f.weather[0].Timestamp.Before(f.currentTime) || f.weather[0].Timestamp.Equal(f.currentTime)) {
						select {
						case f.outputWeather <- f.weather[0]:
						default:
							// Data loss
						}

						f.weather = f.weather[1:]
					}
				}
				f.weatherLock.Unlock()

				f.timingLock.Lock()
				if len(f.timing) > 0 {
					for len(f.timing) > 0 && (f.timing[0].Timestamp.Before(f.currentTime) || f.timing[0].Timestamp.Equal(f.currentTime)) {
						select {
						case f.outputTimingMessages <- f.timing[0]:
						default:
							// Data loss
						}

						f.timing = f.timing[1:]
					}
				}
				f.timingLock.Unlock()

				f.telemetryLock.Lock()
				if len(f.telemetry) > 0 {
					for len(f.telemetry) > 0 && (f.telemetry[0].Timestamp.Before(f.currentTime) || f.telemetry[0].Timestamp.Equal(f.currentTime)) {
						select {
						case f.outputTelemetry <- f.telemetry[0]:
						default:
							// Data loss
						}

						f.telemetry = f.telemetry[1:]
					}
				}
				f.telemetryLock.Unlock()

				f.radioLock.Lock()
				if len(f.radio) > 0 {
					for len(f.radio) > 0 && (f.radio[0].Timestamp.Before(f.currentTime) || f.radio[0].Timestamp.Equal(f.currentTime)) {
						select {
						case f.outputRadio <- f.radio[0]:
						default:
							// Data loss
						}

						f.radio = f.radio[1:]
					}
				}
				f.radioLock.Unlock()
			} else {
				counter++
			}

			f.locationLock.Lock()
			if len(f.location) > 0 {
				for len(f.location) > 0 && (f.location[0].Timestamp.Before(f.currentTime) || f.location[0].Timestamp.Equal(f.currentTime)) {
					select {
					case f.outputLocation <- f.location[0]:
					default:
						// Data loss
					}

					f.location = f.location[1:]
				}
			}
			f.locationLock.Unlock()

			if !f.currentTime.IsZero() {
				increment := f.incrementTime
				if increment > 0 {
					f.currentTime = f.currentTime.Add(increment)

					// TODO - do thread safe
					f.incrementTime = f.incrementTime - increment

					// We want to skip any radio messages when we jump forward in time
					f.radioLock.Lock()
					for len(f.radio) > 0 && (f.radio[0].Timestamp.Before(f.currentTime) || f.radio[0].Timestamp.Equal(f.currentTime)) {
						f.radio = f.radio[1:]
					}
					f.radioLock.Unlock()
				}

				if !f.sessionStart.IsZero() && !f.clockStopped {
					f.remainingTime = f.sessionLength - f.currentTime.Sub(f.sessionStart)

					// Things keep happening after the time has run out so just stop at 0
					if f.remainingTime < 0 {
						f.remainingTime = 0
					}
				}

				f.outputEventTime <- Messages.EventTime{Timestamp: f.currentTime, Remaining: f.remainingTime}

				f.currentTime = f.currentTime.Add(time.Millisecond * 500)
			}
		}
	}
}

func (f *realtime) AddWeather(weather Messages.Weather) {
	f.weatherLock.Lock()
	defer f.weatherLock.Unlock()
	f.weather = append(f.weather, weather)
}

func (f *realtime) AddRaceControlMessage(raceControlMessage Messages.RaceControlMessage) {
	f.raceControlLock.Lock()
	defer f.raceControlLock.Unlock()
	f.raceControl = append(f.raceControl, raceControlMessage)
}

func (f *realtime) AddTiming(timing Messages.Timing) {
	f.timingLock.Lock()
	defer f.timingLock.Unlock()
	f.timing = append(f.timing, timing)
}

func (f *realtime) AddEvent(event Messages.Event) {
	f.eventLock.Lock()
	defer f.eventLock.Unlock()
	f.event = append(f.event, event)
}

func (f *realtime) AddTelemetry(telemetry Messages.Telemetry) {
	f.telemetryLock.Lock()
	defer f.telemetryLock.Unlock()
	f.telemetry = append(f.telemetry, telemetry)
}

func (f *realtime) AddLocation(location Messages.Location) {
	f.locationLock.Lock()
	defer f.locationLock.Unlock()
	f.location = append(f.location, location)
}

func (f *realtime) AddRadio(radio Messages.Radio) {
	f.radioLock.Lock()
	defer f.radioLock.Unlock()
	f.radio = append(f.radio, radio)
}

func (f *realtime) IncrementLap() {
	f.incrementLapCount++
}

func (f *realtime) IncrementTime(duration time.Duration) {
	f.incrementTime += duration
}

func (f *realtime) SkipToSessionStart() {
	f.skipToStart = true
}

func (f *realtime) TogglePause() {
	f.isPaused = !f.isPaused
}

func (f *realtime) IsPaused() bool {
	return f.isPaused
}

func (f *realtime) IncrementDelay(delay time.Duration) {}

func (f *realtime) DecrementDelay(delay time.Duration) {}

func (f *realtime) Delay() time.Duration {
	return 0
}
