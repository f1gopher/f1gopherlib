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
	"encoding/json"
	"fmt"
	"github.com/f1gopher/f1gopherlib"
	"github.com/f1gopher/f1gopherlib/Messages"
	"github.com/zsefvlol/timezonemapper"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCreateHistory(t *testing.T) {
	data := buildHistory()

	output, err := os.Create("../sessionHistoryData.go")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer output.Close()

	output.WriteString(`package f1gopherlib

import (
	"f1gopherlib/f1gopherlib/Messages"
	"time"
)

var sessionHistory = [...]RaceEvent{`)

	for _, session := range data {
		output.WriteString(fmt.Sprintf(`	{
		Country:   "%s",
		RaceTime:  time.Date(%d, %d, %d, %d, %d, %d, %d, time.%s),
		EventTime: time.Date(%d, %d, %d, %d, %d, %d, %d, time.%s),
		Type:      Messages.%sSession,
		Name:      "%s",
		timezone:  "%s",
		TrackName: "%s",
		urlName:   "%s",
	},
`,
			session.Country,
			session.RaceTime.Year(),
			session.RaceTime.Month(),
			session.RaceTime.Day(),
			session.RaceTime.Hour(),
			session.RaceTime.Minute(),
			session.RaceTime.Second(),
			session.RaceTime.Nanosecond(),
			session.RaceTime.Location().String(),

			session.EventTime.Year(),
			session.EventTime.Month(),
			session.EventTime.Day(),
			session.EventTime.Hour(),
			session.EventTime.Minute(),
			session.EventTime.Second(),
			session.EventTime.Nanosecond(),
			session.EventTime.Location().String(),

			strings.Replace(session.Type.String(), "_", "", -1),
			session.Name,
			session.Timezone().String(),
			session.TrackName,
			session.Url()))
	}

	output.WriteString("}")
}

type RaceTimetable struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				Circuit struct {
					Location struct {
						Country  string `json:"country"`
						Lat      string `json:"lat"`
						Locality string `json:"locality"`
						Long     string `json:"long"`
					} `json:"Location"`
					CircuitID   string `json:"circuitId"`
					CircuitName string `json:"circuitName"`
					URL         string `json:"url"`
				} `json:"Circuit"`
				FirstPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"FirstPractice"`
				Qualifying struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"Qualifying"`
				SecondPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"SecondPractice"`
				Sprint struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"Sprint"`
				ThirdPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"ThirdPractice"`
				Date     string `json:"date"`
				RaceName string `json:"raceName"`
				Round    string `json:"round"`
				Season   string `json:"season"`
				Time     string `json:"time"`
				URL      string `json:"url"`
			} `json:"Races"`
			Season string `json:"season"`
		} `json:"RaceTable"`
		Limit  string `json:"limit"`
		Offset string `json:"offset"`
		Series string `json:"series"`
		Total  string `json:"total"`
		URL    string `json:"url"`
		Xmlns  string `json:"xmlns"`
	} `json:"MRData"`
}

func buildHistory() []f1gopherlib.RaceEvent {
	result := make([]f1gopherlib.RaceEvent, 0)

	for x := time.Now().Year() + 1; x >= 2018; x-- {
		races := racesForYear(x)

		for y := len(races.MRData.RaceTable.Races) - 1; y >= 0; y-- {

			race := races.MRData.RaceTable.Races[y]

			raceDateTime, err := time.Parse("2006-01-02T15:04:05Z", race.Date+"T"+race.Time)
			if err != nil {
				fmt.Println("")
			}
			raceDate, err := time.Parse("2006-01-02", race.Date)
			if err != nil {
				fmt.Println("")
			}

			country := strings.Replace(race.RaceName, " Grand Prix", "", 1)
			if country == "Brazilian" {
				country = race.Circuit.Location.Locality
			}
			country = strings.Replace(country, " ", "_", -1)

			practice1Time, err := time.Parse("2006-01-02T15:04:05Z", race.FirstPractice.Date+"T"+race.FirstPractice.Time)
			practice2Time, err := time.Parse("2006-01-02T15:04:05Z", race.SecondPractice.Date+"T"+race.SecondPractice.Time)
			practice3Time, err := time.Parse("2006-01-02T15:04:05Z", race.ThirdPractice.Date+"T"+race.ThirdPractice.Time)
			qualifyingTime, err := time.Parse("2006-01-02T15:04:05Z", race.Qualifying.Date+"T"+race.Qualifying.Time)
			sprintTime, err := time.Parse("2006-01-02T15:04:05Z", race.Sprint.Date+"T"+race.Sprint.Time)

			// Some events only have race times
			if len(race.Sprint.Date) != 0 {
				if practice1Time.Year() == 1 {
					practice1Time = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-2, 0, 0, 0, 0, time.UTC)
				}
				if qualifyingTime.Year() == 1 {
					qualifyingTime = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-2, 0, 0, 0, 0, time.UTC)
				}
				if practice2Time.Year() == 1 {
					practice2Time = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-1, 0, 0, 0, 0, time.UTC)
				}
				if sprintTime.Year() == 1 {
					sprintTime = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-1, 0, 0, 0, 0, time.UTC)
				}
			} else {
				if practice1Time.Year() == 1 {
					practice1Time = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-2, 0, 0, 0, 0, time.UTC)
				}
				if practice2Time.Year() == 1 {
					practice2Time = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-2, 0, 0, 0, 0, time.UTC)
				}
				if practice3Time.Year() == 1 {
					practice3Time = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-1, 0, 0, 0, 0, time.UTC)
				}
				if qualifyingTime.Year() == 1 {
					qualifyingTime = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-1, 0, 0, 0, 0, time.UTC)
				}
				if sprintTime.Year() == 1 {
					sprintTime = time.Date(raceDate.Year(), raceDate.Month(), raceDate.Day()-1, 0, 0, 0, 0, time.UTC)
				}
			}

			timezone := timezoneForCountry(race.Circuit.Location.Long, race.Circuit.Location.Lat)

			var sessions []f1gopherlib.RaceEvent

			if (raceDate.Year() == 2021 || raceDate.Year() == 2019 || raceDate.Year() == 2018) && race.RaceName == "Monaco Grand Prix" {
				sessions = monaco2021_2019History(country, race.Circuit.Location.Country, raceDate, raceDateTime, race.RaceName, timezone, race.Circuit.CircuitName)
			} else if raceDate.Year() == 2020 && race.RaceName == "Emilia Romagna Grand Prix" {
				sessions = emiliaRomagna2020History(country, race.Circuit.Location.Country, raceDate, raceDateTime, race.RaceName, timezone, race.Circuit.CircuitName)
			} else if raceDate.Year() == 2020 && race.RaceName == "German Grand Prix" {
				sessions = german2020History(country, race.Circuit.Location.Country, raceDate, raceDateTime, race.RaceName, timezone, race.Circuit.CircuitName)
			} else if raceDate.Year() == 2019 && race.RaceName == "Japanese Grand Prix" {
				sessions = japan2019History(country, race.Circuit.Location.Country, raceDate, raceDateTime, race.RaceName, timezone, race.Circuit.CircuitName)
			} else {
				if len(race.Sprint.Date) != 0 {
					sessions = sprintHistory(
						country,
						race.Circuit.Location.Country,
						race.RaceName,
						race.Circuit.CircuitName,
						timezone,
						raceDateTime,
						practice1Time,
						practice2Time,
						qualifyingTime,
						sprintTime)
				} else {
					sessions = defaultHistory(
						country,
						race.Circuit.Location.Country,
						race.RaceName,
						race.Circuit.CircuitName,
						timezone,
						raceDateTime,
						practice1Time,
						practice2Time,
						practice3Time,
						qualifyingTime)
				}
			}

			result = append(result, sessions...)
		}
	}

	return result
}

func defaultHistory(
	urlCountry string,
	country string,
	name string,
	circuitName string,
	timezone *time.Location,
	raceDateTime time.Time,
	practice1Time time.Time,
	practice2Time time.Time,
	practice3Time time.Time,
	qualifyingTime time.Time) []f1gopherlib.RaceEvent {
	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		qualifyingTime,
		Messages.QualifyingSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		practice3Time,
		Messages.Practice3Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		practice2Time,
		Messages.Practice2Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		practice1Time,
		Messages.Practice1Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	return result
}

func sprintHistory(
	urlCountry string,
	country string,
	name string,
	circuitName string,
	timezone *time.Location,
	raceDateTime time.Time,
	practice1Time time.Time,
	practice2Time time.Time,
	qualifyingTime time.Time,
	sprintTime time.Time) []f1gopherlib.RaceEvent {
	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		sprintTime,
		Messages.SprintSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		practice2Time,
		Messages.Practice2Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		qualifyingTime,
		Messages.QualifyingSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		practice1Time,
		Messages.Practice1Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	return result
}

func monaco2021_2019History(
	urlCountry string,
	country string,
	raceDate time.Time,
	raceDateTime time.Time,
	name string,
	timezone *time.Location,
	circuitName string) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*24),
		Messages.QualifyingSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*24),
		Messages.Practice3Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*72),
		Messages.Practice2Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*72),
		Messages.Practice1Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	return result
}

func emiliaRomagna2020History(
	urlCountry string,
	country string,
	raceDate time.Time,
	raceDateTime time.Time,
	name string,
	timezone *time.Location,
	circuitName string) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*24),
		Messages.QualifyingSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*24),
		Messages.Practice1Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	return result
}

func german2020History(
	urlCountry string,
	country string,
	raceDate time.Time,
	raceDateTime time.Time,
	name string,
	timezone *time.Location,
	circuitName string) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*24),
		Messages.QualifyingSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*24),
		Messages.Practice3Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	return result
}

func japan2019History(
	urlCountry string,
	country string,
	raceDate time.Time,
	raceDateTime time.Time,
	name string,
	timezone *time.Location,
	circuitName string) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate,
		Messages.QualifyingSession,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*48),
		Messages.Practice2Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDate.Add(-time.Hour*48),
		Messages.Practice1Session,
		name,
		circuitName,
		urlCountry,
		timezone.String(),
	))

	return result
}

func racesForYear(year int) RaceTimetable {
	resp, err := http.Get(fmt.Sprintf("https://ergast.com/api/f1/%d.json?limit=100", year))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	res := RaceTimetable{}
	json.Unmarshal(body, &res)

	return res
}

func timezoneForCountry(longitude string, latitude string) *time.Location {

	long, err := strconv.ParseFloat(longitude, 8)
	if err != nil {
		panic(err)
	}
	lat, err := strconv.ParseFloat(latitude, 8)
	if err != nil {
		panic(err)
	}

	timezone := timezonemapper.LatLngToTimezoneString(lat, long)
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(err)
	}

	return loc
}
