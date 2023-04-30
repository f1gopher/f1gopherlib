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

	output.WriteString(`// F1GopherLib - Copyright (C) 2023 f1gopher
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
	"github.com/f1gopher/f1gopherlib/Messages"
	"time"
)

var sessionHistory = [...]RaceEvent{
`)

	for _, session := range data {
		output.WriteString(fmt.Sprintf(`	{
		Country:   "%s",
		RaceTime:  time.Date(%d, %d, %d, %d, %d, %d, %d, time.%s),
		EventTime: time.Date(%d, %d, %d, %d, %d, %d, %d, time.%s),
		Type:      Messages.%sSession,
		Name:      "%s",
		timezone:  "%s",
		TrackName: "%s",
		TrackYearCreated: %d,
		TimeLostInPitlane: time.Duration(%d) * time.Millisecond,
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

			session.RaceTime.Year(),
			session.EventTime.Month(),
			session.EventTime.Day(),
			session.EventTime.Hour(),
			session.EventTime.Minute(),
			session.EventTime.Second(),
			session.EventTime.Nanosecond(),
			session.EventTime.Location().String(),

			strings.Replace(strings.Replace(session.Type.String(), "_", "", -1), " ", "", -1),
			session.Name,
			session.Timezone().String(),
			session.TrackName,
			session.TrackYearCreated,
			session.TimeLostInPitlane.Milliseconds(),
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

var pitlaneTimes = map[string]time.Duration{
	"Yas Marina Circuit":                          (time.Duration(19914) * time.Millisecond) - (15 * time.Second),
	"Las Vegas Strip Street Circuit":              (time.Duration(0) * time.Millisecond) - (10 * time.Second), // Estimate
	"Autódromo José Carlos Pace":                  (time.Duration(20865) * time.Millisecond) - (21 * time.Second),
	"Autódromo Hermanos Rodríguez":                (time.Duration(20300) * time.Millisecond) - (15 * time.Second),
	"Circuit of the Americas":                     (time.Duration(21290) * time.Millisecond) - (12 * time.Second),
	"Losail International Circuit":                (time.Duration(22900) * time.Millisecond) - (10 * time.Second), // Estimate
	"Suzuka Circuit":                              (time.Duration(22516) * time.Millisecond) - (10 * time.Second),
	"Marina Bay Street Circuit":                   (time.Duration(27315) * time.Millisecond) - (10 * time.Second),
	"Autodromo Nazionale di Monza":                (time.Duration(21598) * time.Millisecond) - (8 * time.Second),
	"Circuit Park Zandvoort":                      (time.Duration(16745) * time.Millisecond) - (15 * time.Second),
	"Circuit de Spa-Francorchamps":                (time.Duration(20421) * time.Millisecond) - (16 * time.Second),
	"Hungaroring":                                 (time.Duration(18918) * time.Millisecond) - (13 * time.Second),
	"Silverstone Circuit":                         (time.Duration(26303) * time.Millisecond) - (19 * time.Second),
	"Red Bull Ring":                               (time.Duration(29106) * time.Millisecond) - (15 * time.Second),
	"Circuit Gilles Villeneuve":                   (time.Duration(21150) * time.Millisecond) - (14 * time.Second),
	"Circuit de Barcelona-Catalunya":              (time.Duration(19564) * time.Millisecond) - (13 * time.Second),
	"Circuit de Monaco":                           (time.Duration(22109) * time.Millisecond) - (15 * time.Second),
	"Autodromo Enzo e Dino Ferrari":               (time.Duration(27800) * time.Millisecond) - (9 * time.Second),
	"Miami International Autodrome":               (time.Duration(16186) * time.Millisecond) - (13 * time.Second),
	"Baku City Circuit":                           (time.Duration(17831) * time.Millisecond) - (10 * time.Second),
	"Albert Park Grand Prix Circuit":              (time.Duration(15164) * time.Millisecond) - (14 * time.Second),
	"Jeddah Corniche Circuit":                     (time.Duration(17873) * time.Millisecond) - (11 * time.Second),
	"Bahrain International Circuit":               (time.Duration(22336) * time.Millisecond) - (8 * time.Second),
	"Circuit Paul Ricard":                         (time.Duration(32822) * time.Millisecond) - (16 * time.Second),
	"Istanbul Park":                               (time.Duration(20292) * time.Millisecond) - (15 * time.Second),
	"Sochi Autodrom":                              (time.Duration(22119) * time.Millisecond) - (15 * time.Second),
	"Autódromo Internacional do Algarve":          (time.Duration(20898) * time.Millisecond) - (9 * time.Second),
	"Bahrain International Circuit - Outer Track": (time.Duration(21974) * time.Millisecond) - (10 * time.Second), // Estimate
	"Nürburgring":                                 (time.Duration(20206) * time.Millisecond) - (10 * time.Second), // Estimate
	"Autodromo Internazionale del Mugello":        (time.Duration(15571) * time.Millisecond) - (10 * time.Second), // Estimate
	"Hockenheimring":                              (time.Duration(17262) * time.Millisecond) - (10 * time.Second), // Estimate
	"Shanghai International Circuit":              (time.Duration(20115) * time.Millisecond) - (10 * time.Second), // Estimate
}

func buildHistory() []f1gopherlib.RaceEvent {
	result := make([]f1gopherlib.RaceEvent, 0)
	const defaultTrackCreatedYear = 2018

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

			// Correct for different layout at the same track
			if race.RaceName == "Sakhir Grand Prix" && raceDate.Year() == 2020 {
				race.Circuit.CircuitName = "Bahrain International Circuit - Outer Track"
			}

			// Include the year the track was created/last changed so we know which map to use if the track has changed
			// over time. For tracks that have changed use the change data
			trackCreatedYear := defaultTrackCreatedYear
			if race.Circuit.CircuitName == "Albert Park Grand Prix Circuit" && raceDate.Year() >= 2022 {
				trackCreatedYear = 2022
			} else if race.Circuit.CircuitName == "Yas Marina Circuit" && raceDate.Year() >= 2021 {
				trackCreatedYear = 2021
			} else if race.Circuit.CircuitName == "Jeddah Corniche Circuit" && raceDate.Year() == 2021 {
				trackCreatedYear = 2021
			} else if race.Circuit.CircuitName == "Jeddah Corniche Circuit" && raceDate.Year() >= 2022 {
				trackCreatedYear = 2022
			} else if race.Circuit.CircuitName == "Jeddah Corniche Circuit" && raceDate.Year() >= 2023 {
				trackCreatedYear = 2023
			} else if race.Circuit.CircuitName == "Circuit de Spa-Francorchamps" && raceDate.Year() >= 2022 {
				trackCreatedYear = 2022
			} else if race.Circuit.CircuitName == "Autodromo Internazionale del Mugello" {
				trackCreatedYear = 2020
			} else if race.Circuit.CircuitName == "Autódromo Internacional do Algarve" {
				trackCreatedYear = 2020
			} else if race.Circuit.CircuitName == "Autodromo Enzo e Dino Ferrari" {
				trackCreatedYear = 2020
			} else if race.Circuit.CircuitName == "Istanbul Park" {
				trackCreatedYear = 2020
			} else if race.Circuit.CircuitName == "Bahrain International Circuit - Outer Track" && race.RaceName == "Sakhir Grand Prix" {
				trackCreatedYear = 2020
			} else if race.Circuit.CircuitName == "Circuit Park Zandvoort" {
				trackCreatedYear = 2021
			} else if race.Circuit.CircuitName == "Miami International Autodrome" {
				trackCreatedYear = 2022
			} else if race.Circuit.CircuitName == "Las Vegas Strip Street Circuit" {
				trackCreatedYear = 2023
			} else if race.Circuit.CircuitName == "Losail International Circuit" {
				trackCreatedYear = 2021
			} else if race.Circuit.CircuitName == "Circuit de Barcelona-Catalunya" && raceDate.Year() >= 2023 {
				trackCreatedYear = 2023
			}

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

			pitlaneTime, exists := pitlaneTimes[race.Circuit.CircuitName]
			if !exists {
				pitlaneTime = 0
			}

			var sessions []f1gopherlib.RaceEvent

			if (raceDate.Year() == 2021 || raceDate.Year() == 2019 || raceDate.Year() == 2018) && race.RaceName == "Monaco Grand Prix" {
				sessions = monaco2021_2019History(
					country,
					race.Circuit.Location.Country,
					raceDate,
					raceDateTime,
					race.RaceName,
					timezone,
					race.Circuit.CircuitName,
					trackCreatedYear,
					pitlaneTime)
			} else if raceDate.Year() == 2020 && race.RaceName == "Emilia Romagna Grand Prix" {
				sessions = emiliaRomagna2020History(
					country,
					race.Circuit.Location.Country,
					raceDate,
					raceDateTime,
					race.RaceName,
					timezone,
					race.Circuit.CircuitName,
					trackCreatedYear,
					pitlaneTime)
			} else if raceDate.Year() == 2020 && race.RaceName == "German Grand Prix" {
				sessions = german2020History(
					country,
					race.Circuit.Location.Country,
					raceDate,
					raceDateTime,
					race.RaceName,
					timezone,
					race.Circuit.CircuitName,
					trackCreatedYear,
					pitlaneTime)
			} else if raceDate.Year() == 2019 && race.RaceName == "Japanese Grand Prix" {
				sessions = japan2019History(
					country,
					race.Circuit.Location.Country,
					raceDate,
					raceDateTime,
					race.RaceName,
					timezone,
					race.Circuit.CircuitName,
					trackCreatedYear,
					pitlaneTime)
			} else {
				if len(race.Sprint.Date) != 0 {
					sessions = sprintHistory(
						country,
						race.Circuit.Location.Country,
						race.RaceName,
						race.Circuit.CircuitName,
						trackCreatedYear,
						timezone,
						raceDateTime,
						practice1Time,
						practice2Time,
						qualifyingTime,
						sprintTime,
						pitlaneTime)
				} else {
					sessions = defaultHistory(
						country,
						race.Circuit.Location.Country,
						race.RaceName,
						race.Circuit.CircuitName,
						trackCreatedYear,
						timezone,
						raceDateTime,
						practice1Time,
						practice2Time,
						practice3Time,
						qualifyingTime,
						pitlaneTime)
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
	trackYearCreated int,
	timezone *time.Location,
	raceDateTime time.Time,
	practice1Time time.Time,
	practice2Time time.Time,
	practice3Time time.Time,
	qualifyingTime time.Time,
	pitlaneTime time.Duration) []f1gopherlib.RaceEvent {
	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
	trackYearCreated int,
	timezone *time.Location,
	raceDateTime time.Time,
	practice1Time time.Time,
	practice2Time time.Time,
	qualifyingTime time.Time,
	sprintTime time.Time,
	pitlaneTime time.Duration) []f1gopherlib.RaceEvent {
	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
		urlCountry,
		timezone.String(),
	))

	if raceDateTime.Year() < 2023 {
		result = append(result, *f1gopherlib.CreateRaceEvent(
			country,
			raceDateTime,
			practice2Time,
			Messages.Practice2Session,
			name,
			circuitName,
			trackYearCreated,
			pitlaneTime,
			urlCountry,
			timezone.String(),
		))
	} else {
		result = append(result, *f1gopherlib.CreateRaceEvent(
			country,
			raceDateTime,
			practice2Time,
			Messages.QualifyingSession,
			name,
			circuitName,
			trackYearCreated,
			pitlaneTime,
			urlCountry,
			timezone.String(),
		))
	}

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		qualifyingTime,
		Messages.QualifyingSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
	circuitName string,
	trackYearCreated int,
	pitlaneTime time.Duration) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
	circuitName string,
	trackYearCreated int,
	pitlaneTime time.Duration) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
	circuitName string,
	trackYearCreated int,
	pitlaneTime time.Duration) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
	circuitName string,
	trackYearCreated int,
	pitlaneTime time.Duration) []f1gopherlib.RaceEvent {

	result := make([]f1gopherlib.RaceEvent, 0)

	result = append(result, *f1gopherlib.CreateRaceEvent(
		country,
		raceDateTime,
		raceDateTime,
		Messages.RaceSession,
		name,
		circuitName,
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
		trackYearCreated,
		pitlaneTime,
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
