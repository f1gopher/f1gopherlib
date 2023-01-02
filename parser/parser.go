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

package parser

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/f1gopher/f1gopherlib/Messages"
	"github.com/f1gopher/f1gopherlib/connection"
	"github.com/f1gopher/f1gopherlib/f1log"
	"github.com/f1gopher/f1gopherlib/flowControl"
	"io"
	"strings"
	"sync"
	"time"
)

type DataSource int

const (
	EventTime DataSource = 1 << iota
	Event
	RaceControl
	Weather
	Timing
	Telemetry
	Location
	TeamRadio
)

type Parser struct {
	requestedData DataSource

	incoming <-chan connection.Payload
	output   flowControl.Flow

	driverTimes map[string]Messages.Timing
	eventState  Messages.Event

	assets connection.AssetStore

	session Messages.SessionType

	log *f1log.F1GopherLibLog

	ctx context.Context
	wg  *sync.WaitGroup
}

func Create(
	ctx context.Context,
	wg *sync.WaitGroup,
	requestedData DataSource,
	incoming <-chan connection.Payload,
	output flowControl.Flow,
	assets connection.AssetStore,
	session Messages.SessionType,
	log *f1log.F1GopherLibLog) *Parser {

	abc := Parser{
		ctx:           ctx,
		wg:            wg,
		requestedData: requestedData,
		incoming:      incoming,
		output:        output,
		driverTimes:   make(map[string]Messages.Timing),
		assets:        assets,
		session:       session,
		log:           log,
	}

	return &abc
}

func (p *Parser) ParseErrorf(file string, timestamp time.Time, msg string, a ...any) {
	p.log.Errorf("%s - %v: %s", file, timestamp, fmt.Sprintf(msg, a))
}

func (p *Parser) ParseTimeError(file string, timestamp time.Time, field string, err error) {
	p.log.Errorf("%s - %v: Unable to parse time for '%s': %v", file, timestamp, field, err)
}

func (p *Parser) Process() {
	p.wg.Add(1)
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return

		case msg := <-p.incoming:
			switch msg.Name {
			case connection.EndOfDataFile:
				return

			case connection.CatchupFile:
				var dat map[string]interface{}
				if err := json.Unmarshal([]byte(msg.Data), &dat); err != nil {
					p.log.Errorf("Catchup data parse error: '%v' for data: %s", err, msg.Data)
					continue
				}

				zeroTimestamp := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

				for _, fileName := range connection.OrderedFiles {
					if fileName == connection.TeamRadioFile ||
						fileName == connection.ContentStreamsFile ||
						fileName == connection.AudioStreamsFile {
						continue
					}

					fileData, exists := dat[fileName]
					if exists {
						if strings.HasSuffix(fileName, ".z") {
							abc, err := p.decompressData([]byte(fileData.(string)))
							if err != nil {
								p.log.Errorf("Decompressing data for file '%s': %v with data: %s", msg.Name, err, msg.Data)
								continue
							}

							p.handleMessage(fileName, abc, zeroTimestamp)
						} else {
							p.handleMessage(fileName, fileData.(map[string]interface{}), zeroTimestamp)
						}
					}
				}

			default:
				var dat map[string]interface{}
				var err error
				if strings.HasSuffix(msg.Name, ".z") {
					dat, err = p.decompressData(msg.Data)
					if err != nil {
						p.log.Errorf("Decompressing data for file '%s': %v with data: %s", msg.Name, err, msg.Data)
						continue
					}

				} else {
					if err := json.Unmarshal([]byte(msg.Data), &dat); err != nil {
						p.log.Errorf("Data parse error for file '%s': '%v' for data: %s", msg.Name, err, msg.Data)
						continue
					}
				}

				dataTime, err := parseTime(msg.Timestamp)
				if err != nil {
					p.log.Errorf("Parsing file timestamp for '%s' with value '%s': %v", msg.Name, msg.Timestamp, err)
				}

				p.handleMessage(msg.Name, dat, dataTime)
			}
		}
	}
}

func (p *Parser) handleMessage(name string, dat map[string]interface{}, timestamp time.Time) {
	switch name {
	case connection.WeatherDataFile:
		if p.requestedData&Weather == Weather {
			outgoing, err := parseWeatherData(dat, timestamp)
			if err == nil {
				p.output.AddWeather(outgoing)
			}
		}

	case connection.SessionDataFile:
		if p.requestedData&Event == Event {
			outgoing, err := p.parseSessionDataData(dat, timestamp)
			if err == nil {
				for _, rcMsg := range outgoing {
					p.output.AddEvent(rcMsg)
				}
			}
		}

	case connection.TimingDataFile:
		if p.requestedData&Timing == Timing {
			outgoing, err := p.parseTimingData(dat, timestamp)
			if err == nil {
				for _, rcMsg := range outgoing {
					p.output.AddTiming(rcMsg)
				}
			}
		}

	case connection.TimingAppDataFile:
		if p.requestedData&Timing == Timing {
			outgoing, err := p.parseTimingAppData(dat, timestamp)
			if err == nil {
				for _, rcMsg := range outgoing {
					p.output.AddTiming(rcMsg)
				}
			}
		}

	case connection.HeartbeatFile:
		if p.requestedData&Event == Event {
			outgoing, err := p.parseHeartbeatData(dat, timestamp)
			if err == nil {
				p.output.AddEvent(outgoing)
			}
		}

	case connection.CarDataFile:
		if p.requestedData&Telemetry == Telemetry || p.requestedData&Timing == Timing {
			outgoing, timingOutgoing, err := p.parseCarData(dat, timestamp)
			if err == nil {
				if p.requestedData&Telemetry == Telemetry {
					for _, rcMsg := range outgoing {
						p.output.AddTelemetry(rcMsg)
					}
				}

				if p.requestedData&Timing == Timing {
					for _, rcMsg := range timingOutgoing {
						p.output.AddTiming(rcMsg)
					}
				}
			}
		}

	case connection.PositionFile:
		if p.requestedData&Location == Location {
			outgoing, err := p.parsePositionData(dat, timestamp)
			if err == nil {
				for _, rcMsg := range outgoing {
					p.output.AddLocation(rcMsg)
				}
			}
		}

	case connection.SessionInfoFile:
		if p.requestedData&Event == Event {
			outgoing, err := p.parseSessionInfoData(dat, timestamp)
			if err == nil {
				p.output.AddEvent(outgoing)
			}
		}

	case connection.LapCountFile:
		if p.requestedData&Event == Event {
			outgoing, err := p.parseCurrentLapData(dat, timestamp)
			if err == nil {
				p.output.AddEvent(outgoing)
			}
		}

	case connection.RaceControlMessagesFile:
		if p.requestedData&RaceControl == RaceControl || p.requestedData&Event == Event || p.requestedData&Timing == Timing {
			outgoingRcm, outgoingEvent, outgoingTiming, err := p.parseRaceControlMessagesData(dat, timestamp)
			if err == nil {

				if p.requestedData&RaceControl == RaceControl {
					for _, rcMsg := range outgoingRcm {
						p.output.AddRaceControlMessage(rcMsg)
					}
				}

				if p.requestedData&Event == Event {
					for _, rcMsg := range outgoingEvent {
						p.output.AddEvent(rcMsg)
					}
				}

				if p.requestedData&Timing == Timing {
					for _, timingMsg := range outgoingTiming {
						p.output.AddTiming(timingMsg)
					}
				}
			}
		}

	case connection.SessionStatusFile:
		outgoing, timingOutgoing, err := p.parseSessionStatusData(dat, timestamp)
		if p.requestedData&Event == Event && err == nil {
			p.output.AddEvent(outgoing)
		}

		if p.requestedData&Timing == Timing {
			for _, rcMsg := range timingOutgoing {
				p.output.AddTiming(rcMsg)
			}
		}

	case connection.TeamRadioFile:
		if p.requestedData&TeamRadio == TeamRadio {
			outgoing, err := p.parseTeamRadioData(dat, timestamp)
			if err == nil {
				for _, radioMsg := range outgoing {
					p.output.AddRadio(radioMsg)
				}
			}
		}

	case connection.DriverListFile:
		p.parseDriverList(dat, timestamp)

	case connection.ExtrapolatedClockFile:
		outgoing, err := p.parseExtrapolatedClockData(dat, timestamp)
		if p.requestedData&Event == Event && err == nil {
			p.output.AddEvent(outgoing)
		}

	case connection.TrackStatusFile:
	case connection.TopThreeFile:
	case connection.TimingStatsFile:
	case connection.AudioStreamsFile:
	case connection.ContentStreamsFile:

	default:

	}
}

func (p *Parser) decompressData(data []byte) (map[string]interface{}, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write([]byte("Welcome to CodeBeautify"))
	zw.Close()

	b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))

	stuff, err := io.ReadAll(b64)
	stuff = append(buf.Bytes()[:10], stuff...)

	zr, err := gzip.NewReader(bytes.NewBuffer(stuff))
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	uncompressed, err := io.ReadAll(zr)

	return p.deserializeData(uncompressed)
}

func (p *Parser) deserializeData(data []byte) (map[string]interface{}, error) {
	var dat map[string]interface{}
	err := json.Unmarshal(data, &dat)
	return dat, err
}
