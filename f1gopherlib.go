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

package f1gopherlib

import (
	"context"
	"errors"
	"fmt"
	"github.com/f1gopher/f1gopherlib/Messages"
	"github.com/f1gopher/f1gopherlib/connection"
	"github.com/f1gopher/f1gopherlib/f1log"
	"github.com/f1gopher/f1gopherlib/flowControl"
	"github.com/f1gopher/f1gopherlib/parser"
	"io"
	"path/filepath"
	"sync"
	"time"
)

type F1GopherLib interface {
	Name() string
	Session() Messages.SessionType
	CircuitTimezone() *time.Location
	SessionStart() time.Time
	Track() string
	TrackYear() int

	Weather() <-chan Messages.Weather
	RaceControlMessages() <-chan Messages.RaceControlMessage
	Timing() <-chan Messages.Timing
	Event() <-chan Messages.Event
	Telemetry() <-chan Messages.Telemetry
	Location() <-chan Messages.Location
	Time() <-chan Messages.EventTime
	Radio() <-chan Messages.Radio

	IncrementLap()
	IncrementTime(duration time.Duration)
	SkipToSessionStart()
	TogglePause()
	IsPaused() bool

	Close()
}

type f1gopherlib struct {
	archive string

	session      Messages.SessionType
	name         string
	timezone     *time.Location
	sessionStart time.Time
	track        string
	trackYear    int

	connection   connection.Connection
	dataHandler  *parser.Parser
	replayTiming flowControl.Flow

	weather             chan Messages.Weather
	raceControlMessages chan Messages.RaceControlMessage
	timing              chan Messages.Timing
	event               chan Messages.Event
	telemetry           chan Messages.Telemetry
	location            chan Messages.Location
	eventTime           chan Messages.EventTime
	radio               chan Messages.Radio

	ctxShutdown context.CancelFunc
	ctx         context.Context
	wg          sync.WaitGroup
}

const channelSize = 100000

var f1Log = f1log.CreateLog()

func SetLogOutput(w io.Writer) {
	f1Log.SetLogOutput(w)
}

func CreateRaceEvent(
	country string,
	raceTime time.Time,
	eventTime time.Time,
	sessionType Messages.SessionType,
	name string,
	track string,
	trackYearCreated int,
	urlName string,
	timezone string) *RaceEvent {

	var sessionName string
	switch sessionType {
	case Messages.Practice1Session:
		sessionName = "Practice_1"
	case Messages.Practice2Session:
		sessionName = "Practice_2"
	case Messages.Practice3Session:
		sessionName = "Practice_3"
	case Messages.QualifyingSession:
		sessionName = "Qualifying"
	case Messages.SprintSession:
		sessionName = "Sprint"
	case Messages.RaceSession:
		sessionName = "Race"
	case Messages.PreSeasonSession:
		sessionName = "Test"
	default:
		panic("Unhandled session type: " + sessionType.String())
	}

	urlName = fmt.Sprintf(
		"https://livetiming.formula1.com/static/%d/%d-%02d-%02d_%s_Grand_Prix/%d-%02d-%02d_%s/",
		raceTime.Year(),
		raceTime.Year(),
		raceTime.Month(),
		raceTime.Day(),
		urlName,
		eventTime.Year(),
		eventTime.Month(),
		eventTime.Day(),
		sessionName)

	return &RaceEvent{
		Country:          country,
		RaceTime:         raceTime,
		EventTime:        eventTime,
		Type:             sessionType,
		Name:             name,
		timezone:         timezone,
		TrackName:        track,
		TrackYearCreated: trackYearCreated,
		urlName:          urlName,
	}
}

type RaceEvent struct {
	Country          string
	RaceTime         time.Time
	EventTime        time.Time
	Type             Messages.SessionType
	Name             string
	timezone         string
	TrackName        string
	TrackYearCreated int

	// TODO - add duration

	urlName string
}

func (r *RaceEvent) Timezone() *time.Location {
	tz, _ := time.LoadLocation(r.timezone)
	return tz
}

func (r *RaceEvent) string() string {
	return fmt.Sprintf("%s - %s", r.Name, r.Type.String())
}

func (r *RaceEvent) Url() string {
	return r.urlName
}

func CreateLive(requestedData parser.DataSource, archive string, cache string) (F1GopherLib, error) {

	// TODO - validate path
	// TODO - create archive folder

	currentEvent, exists := liveEvent()

	// No event happening or about to happen so nothing we can do
	if !exists {
		return nil, errors.New("No live event currently happening")
	}

	f1Log.Infof("Creating live session for: %v", currentEvent.string())

	data := f1gopherlib{
		weather:             make(chan Messages.Weather, channelSize),
		raceControlMessages: make(chan Messages.RaceControlMessage, channelSize),
		timing:              make(chan Messages.Timing, channelSize),
		event:               make(chan Messages.Event, channelSize),
		telemetry:           make(chan Messages.Telemetry, channelSize),
		location:            make(chan Messages.Location, channelSize),
		eventTime:           make(chan Messages.EventTime, channelSize),
		radio:               make(chan Messages.Radio, channelSize),

		archive:      archive,
		session:      currentEvent.Type,
		name:         currentEvent.Name,
		timezone:     currentEvent.Timezone(),
		sessionStart: currentEvent.EventTime,
		track:        currentEvent.TrackName,
		trackYear:    currentEvent.TrackYearCreated,
	}
	data.ctx, data.ctxShutdown = context.WithCancel(context.Background())

	err := data.connectLive(requestedData, archive, currentEvent, cache)
	// Always start live paused
	data.TogglePause()
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func CreateDebugReplay(
	requestedData parser.DataSource,
	replayFile string,
	dataFlow flowControl.FlowType) (F1GopherLib, error) {

	// TODO - read event info from debug file
	event := RaceEvent{}

	f1Log.Infof("Creating live replay session for: %v", event.string())

	data := f1gopherlib{
		weather:             make(chan Messages.Weather, channelSize),
		raceControlMessages: make(chan Messages.RaceControlMessage, channelSize),
		timing:              make(chan Messages.Timing, channelSize),
		event:               make(chan Messages.Event, channelSize),
		telemetry:           make(chan Messages.Telemetry, channelSize),
		location:            make(chan Messages.Location, channelSize),
		eventTime:           make(chan Messages.EventTime, channelSize),
		radio:               make(chan Messages.Radio, channelSize),
		session:             event.Type,
		name:                event.Name,
		timezone:            event.Timezone(),
		sessionStart:        event.EventTime,
		track:               event.TrackName,
		trackYear:           event.TrackYearCreated,
	}
	data.ctx, data.ctxShutdown = context.WithCancel(context.Background())

	err := data.connectDebugReplay(requestedData, replayFile, event, dataFlow)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func CreateReplay(
	requestedData parser.DataSource,
	event RaceEvent,
	cache string,
	dataFlow flowControl.FlowType) (F1GopherLib, error) {

	f1Log.Infof("Creating replay session for: %v", event.string())

	data := f1gopherlib{
		weather:             make(chan Messages.Weather, channelSize),
		raceControlMessages: make(chan Messages.RaceControlMessage, channelSize),
		timing:              make(chan Messages.Timing, channelSize),
		event:               make(chan Messages.Event, channelSize),
		telemetry:           make(chan Messages.Telemetry, channelSize),
		location:            make(chan Messages.Location, channelSize),
		eventTime:           make(chan Messages.EventTime, channelSize),
		radio:               make(chan Messages.Radio, channelSize),
		session:             event.Type,
		name:                event.Name,
		timezone:            event.Timezone(),
		sessionStart:        event.EventTime,
		track:               event.TrackName,
		trackYear:           event.TrackYearCreated,
	}
	data.ctx, data.ctxShutdown = context.WithCancel(context.Background())

	err := data.connectReplay(requestedData, event, cache, dataFlow)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (f *f1gopherlib) connectLive(requestedData parser.DataSource, archiveFile string, event RaceEvent, cache string) error {

	cache = f.cachePath(cache, event)

	if len(archiveFile) == 0 {
		f.connection = connection.CreateLive(f.ctx, &f.wg, f1Log)
	} else {
		var connErr error
		f.connection, connErr = connection.CreateArchivingLive(f.ctx, archiveFile)
		if connErr != nil {
			return connErr
		}
	}

	err, dataChannel := f.connection.Connect()
	if err != nil {
		return err
	}

	f.replayTiming = flowControl.CreateFlowControl(
		f.ctx,
		&f.wg,
		flowControl.Realtime,
		f.weather,
		f.raceControlMessages,
		f.timing,
		f.event,
		f.telemetry,
		f.location,
		f.eventTime,
		f.radio)

	assetStore := connection.CreateAssetStore(event.Url(), cache, f1Log)

	dataHandler := parser.Create(
		f.ctx,
		&f.wg,
		requestedData,
		dataChannel,
		f.replayTiming,
		assetStore,
		Messages.RaceSession,
		f1Log)

	go dataHandler.Process()
	go f.replayTiming.Run()

	return nil
}

func (f *f1gopherlib) connectDebugReplay(
	requestedData parser.DataSource,
	replayFile string,
	event RaceEvent,
	dataFlow flowControl.FlowType) error {

	f.connection = connection.CreateArchivedLive(f.ctx, &f.wg, f1Log, replayFile)
	err, dataChannel := f.connection.Connect()

	if err != nil {
		return err
	}

	f.replayTiming = flowControl.CreateFlowControl(
		f.ctx,
		&f.wg,
		dataFlow,
		f.weather,
		f.raceControlMessages,
		f.timing,
		f.event,
		f.telemetry,
		f.location,
		f.eventTime,
		f.radio)

	// Don't use a cache for debug replays because we don't always know the event yet to give it a useful folder name
	assetStore := connection.CreateAssetStore(event.Url(), "", f1Log)

	f.dataHandler = parser.Create(
		f.ctx,
		&f.wg,
		requestedData,
		dataChannel,
		f.replayTiming,
		assetStore,
		event.Type,
		f1Log)

	go f.dataHandler.Process()
	go f.replayTiming.Run()

	return nil
}

func (f *f1gopherlib) connectReplay(
	requestedData parser.DataSource,
	event RaceEvent,
	cache string,
	dataFlow flowControl.FlowType) error {

	url := event.Url()
	cache = f.cachePath(cache, event)

	f.connection = connection.CreateReplay(f.ctx, &f.wg, f1Log, url, event.Type, event.RaceTime.Year(), cache)
	err, dataChannel := f.connection.Connect()

	if err != nil {
		return err
	}

	f.replayTiming = flowControl.CreateFlowControl(
		f.ctx,
		&f.wg,
		dataFlow,
		f.weather,
		f.raceControlMessages,
		f.timing,
		f.event,
		f.telemetry,
		f.location,
		f.eventTime,
		f.radio)

	assetStore := connection.CreateAssetStore(event.Url(), cache, f1Log)

	f.dataHandler = parser.Create(
		f.ctx,
		&f.wg,
		requestedData,
		dataChannel,
		f.replayTiming,
		assetStore,
		event.Type,
		f1Log)

	go f.dataHandler.Process()
	go f.replayTiming.Run()

	return nil
}

func (f *f1gopherlib) cachePath(cache string, event RaceEvent) string {
	return filepath.Join(cache, fmt.Sprintf("%d", event.RaceTime.Year()), fmt.Sprintf("%s_%s", event.RaceTime.Format("2006-01-02"), event.Name), event.Type.String())
}

func (f *f1gopherlib) Session() Messages.SessionType {
	return f.session
}

func (f *f1gopherlib) Name() string {
	return f.name
}

func (f *f1gopherlib) CircuitTimezone() *time.Location {
	return f.timezone
}

func (f *f1gopherlib) SessionStart() time.Time {
	return f.sessionStart
}

func (f *f1gopherlib) Track() string {
	return f.track
}

func (f *f1gopherlib) TrackYear() int {
	return f.trackYear
}

func (f *f1gopherlib) Weather() <-chan Messages.Weather {
	return f.weather
}

func (f *f1gopherlib) RaceControlMessages() <-chan Messages.RaceControlMessage {
	return f.raceControlMessages
}

func (f *f1gopherlib) Timing() <-chan Messages.Timing {
	return f.timing
}

func (f *f1gopherlib) Event() <-chan Messages.Event {
	return f.event
}

func (f *f1gopherlib) Telemetry() <-chan Messages.Telemetry {
	return f.telemetry
}

func (f *f1gopherlib) Location() <-chan Messages.Location {
	return f.location
}

func (f *f1gopherlib) Time() <-chan Messages.EventTime {
	return f.eventTime
}

func (f *f1gopherlib) Radio() <-chan Messages.Radio {
	return f.radio
}

func (f *f1gopherlib) IncrementLap() {
	// Only makes sense for races
	if f.session == Messages.RaceSession || f.session == Messages.SprintSession {
		f.replayTiming.IncrementLap()
	}
}

func (f *f1gopherlib) IncrementTime(duration time.Duration) {
	f.replayTiming.IncrementTime(duration)
}

func (f *f1gopherlib) SkipToSessionStart() {
	f.replayTiming.SkipToSessionStart()
}

func (f *f1gopherlib) TogglePause() {
	f.replayTiming.TogglePause()
}

func (f *f1gopherlib) IsPaused() bool {
	return f.replayTiming.IsPaused()
}

func (f *f1gopherlib) Close() {
	f.name = ""
	f.track = ""

	f.ctxShutdown()
	f.wg.Wait()

	f.connection = nil
	f.dataHandler = nil

	close(f.weather)
	close(f.raceControlMessages)
	close(f.timing)
	close(f.event)
	close(f.telemetry)
	close(f.location)
	close(f.eventTime)
	close(f.radio)
}
