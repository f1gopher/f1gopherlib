package test

import (
	"f1gopherlib/f1gopherlib"
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/connection"
	"f1gopherlib/f1gopherlib/f1log"
	"f1gopherlib/f1gopherlib/parser"
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

		replay := connection.CreateReplay(log, session.Url(), session.Type, session.RaceTime.Year(), "")
		err, payload := replay.Connect()

		if err != nil {
			t.Error(err)
			continue
		}

		assetStore := connection.CreateAssetStore(
			session.Url(),
			filepath.Join("./cache", strings.Replace(session.Url(), "https://livetiming.formula1.com/static/", "", 1)),
			log)

		p := parser.Create(parser.EventTime|parser.Event|parser.RaceControl|parser.Weather|parser.Timing|parser.Telemetry|parser.Location|parser.TeamRadio,
			payload,
			dummy,
			assetStore,
			session.Type,
			log)

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
func (d *dummyFlowControl) IncrementLap()                                                 {}
func (d *dummyFlowControl) IncrementTime(duration time.Duration)                          {}
func (d *dummyFlowControl) SkipToSessionStart()                                           {}
func (d *dummyFlowControl) TogglePause()                                                  {}
func (d *dummyFlowControl) IsPaused() bool                                                { return false }
func (d *dummyFlowControl) IncrementDelay(delay time.Duration)                            {}
func (d *dummyFlowControl) DecrementDelay(delay time.Duration)                            {}
func (d *dummyFlowControl) Delay() time.Duration                                          { return 0 }
