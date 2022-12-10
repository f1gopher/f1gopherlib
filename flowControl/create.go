package flowControl

import (
	"f1gopherlib/f1gopherlib/Messages"
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
	flowType FlowType,
	outputWeather chan<- Messages.Weather,
	outputRaceControlMessages chan<- Messages.RaceControlMessage,
	outputTimingMessages chan<- Messages.Timing,
	outputEvent chan<- Messages.Event,
	outputTelemetry chan<- Messages.Telemetry,
	outputLocation chan<- Messages.Location,
	outputEventTime chan<- Messages.EventTime,
	outputRadio chan<- Messages.Radio) Flow {

	switch flowType {
	case Realtime:
		return &realtime{
			outputWeather:             outputWeather,
			outputRaceControlMessages: outputRaceControlMessages,
			outputTimingMessages:      outputTimingMessages,
			outputEvent:               outputEvent,
			outputTelemetry:           outputTelemetry,
			outputLocation:            outputLocation,
			outputEventTime:           outputEventTime,
			outputRadio:               outputRadio,
			weather:                   make([]Messages.Weather, 0),
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
			weather:                   make([]Messages.Weather, 0),
		}

	default:
		panic("Unhandled flow control type")
	}
}
