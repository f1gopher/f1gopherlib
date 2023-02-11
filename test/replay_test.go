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

package test

import (
	"github.com/f1gopher/f1gopherlib"
	"github.com/f1gopher/f1gopherlib/Messages"
	"github.com/f1gopher/f1gopherlib/connection"
	"github.com/f1gopher/f1gopherlib/f1log"
	"github.com/f1gopher/f1gopherlib/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReplay(t *testing.T) {
	data := f1gopherlib.RaceHistory()
	dummy := &dummyFlowControl{}

	log := f1log.CreateLog()
	log.SetLogOutput(os.Stdout)

	for x := 0; x < len(data); x++ {
		session := data[x]

		t.Logf("Testing: %d %d - %s %s...", x, session.RaceTime.Year(), session.Country, session.Type.String())

		replay := connection.CreateReplay(nil, nil, log, session.Url(), session.Type, session.RaceTime.Year(), "")
		err, payload := replay.Connect()

		if err != nil {
			t.Error(err)
			continue
		}

		assetStore := connection.CreateAssetStore(
			session.Url(),
			filepath.Join("./cache", strings.Replace(session.Url(), "https://livetiming.formula1.com/static/", "", 1)),
			log)

		p := parser.Create(nil, nil, parser.EventTime|parser.Event|parser.RaceControl|parser.Weather|parser.Timing|parser.Telemetry|parser.Location|parser.TeamRadio,
			payload,
			dummy,
			assetStore,
			session.Type,
			log,
			time.UTC)

		p.Process()
	}
}

type dummyFlowControl struct {
}

func (d *dummyFlowControl) Run()                                                          {}
func (d *dummyFlowControl) AddWeather(weather Messages.Weather)                           {}
func (d *dummyFlowControl) AddRaceControlMessage(raceControl Messages.RaceControlMessage) {}
func (d *dummyFlowControl) AddTiming(timing Messages.Timing)                              {}
func (d *dummyFlowControl) AddEvent(timing Messages.Event)                                {}
func (d *dummyFlowControl) AddTelemetry(timing Messages.Telemetry)                        {}
func (d *dummyFlowControl) AddLocation(timing Messages.Location)                          {}
func (d *dummyFlowControl) AddRadio(timing Messages.Radio)                                {}
func (d *dummyFlowControl) AddDrivers(driver Messages.Drivers)                            {}
func (d *dummyFlowControl) IncrementLap()                                                 {}
func (d *dummyFlowControl) IncrementTime(duration time.Duration)                          {}
func (d *dummyFlowControl) SkipToSessionStart()                                           {}
func (d *dummyFlowControl) TogglePause()                                                  {}
func (d *dummyFlowControl) IsPaused() bool                                                { return false }
func (d *dummyFlowControl) IncrementDelay(delay time.Duration)                            {}
func (d *dummyFlowControl) DecrementDelay(delay time.Duration)                            {}
func (d *dummyFlowControl) Delay() time.Duration                                          { return 0 }
