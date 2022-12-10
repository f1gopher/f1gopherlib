package connection

import (
	"bufio"
	"encoding/json"
	"errors"
	"f1gopherlib/f1gopherlib/Messages"
	"f1gopherlib/f1gopherlib/f1log"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fileInfo struct {
	name         string
	data         *bufio.Scanner
	nextLine     string
	nextLineTime time.Time
}

type replay struct {
	log      *f1log.F1GopherLibLog
	cache    string
	dataFeed chan Payload

	eventUrl  string
	session   Messages.SessionType
	eventYear int

	dataFiles []fileInfo
}

const NotFoundResponse = "<?xml version='1.0' encoding='UTF-8'?><Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message></Error>"

func CreateReplay(log *f1log.F1GopherLibLog, url string, session Messages.SessionType, eventYear int, cache string) *replay {
	return &replay{
		log:       log,
		dataFeed:  make(chan Payload, 100000),
		eventUrl:  url,
		session:   session,
		eventYear: eventYear,
		cache:     cache,
	}
}

func (r *replay) Connect() (error, <-chan Payload) {

	r.dataFiles = make([]fileInfo, 0)

	for _, name := range OrderedFiles {

		if (name == PositionFile || name == ContentStreamsFile) && r.eventYear <= 2018 {
			continue
		}

		if name == LapCountFile && !(r.session == Messages.RaceSession || r.session == Messages.SprintSession) {
			continue
		}

		r.dataFiles = append(r.dataFiles, fileInfo{
			name:         name,
			data:         r.get(r.eventUrl + name + ".jsonStream"),
			nextLine:     "",
			nextLineTime: time.Time{},
		})
	}

	go r.readEntries()

	return nil, r.dataFeed
}

func (r *replay) readEntries() {

	sessionStartTime, err := r.findSessionStartTime()
	if err != nil {
		r.dataFeed <- Payload{
			Name: EndOfDataFile,
		}
		return
	}

	currentTime := sessionStartTime

	hasData := true

	for hasData {
		hasData = false

		for x := range r.dataFiles {

			if strings.HasSuffix(r.dataFiles[x].name, ".z") {
				hasData = r.sim(
					r.dataFiles[x].data,
					currentTime,
					&r.dataFiles[x].nextLineTime,
					&r.dataFiles[x].nextLine,
					sessionStartTime,
					r.compressedDataTime,
					r.dataFiles[x].name) || hasData
			} else {
				hasData = r.sim(
					r.dataFiles[x].data,
					currentTime,
					&r.dataFiles[x].nextLineTime,
					&r.dataFiles[x].nextLine,
					sessionStartTime,
					r.uncompressedDataTime,
					r.dataFiles[x].name) || hasData
			}
		}

		currentTime = currentTime.Add(time.Second)
	}

	r.dataFeed <- Payload{
		Name: EndOfDataFile,
	}
}

func (r *replay) sim(
	dataBuffer *bufio.Scanner,
	currentRaceTime time.Time,
	nextTime *time.Time,
	nextData *string,
	sessionStartTime time.Time,
	splitData func(data string, sessionStart time.Time) (timestamp time.Time, payload string, err error),
	name string) bool {

	// If no data then skip
	if dataBuffer == nil {
		return false
	}

	if *nextData != "" {
		if nextTime.After(currentRaceTime) {
			return true
		}

		r.dataFeed <- Payload{
			Name:      name,
			Data:      []byte(*nextData),
			Timestamp: nextTime.Format("2006-01-02T15:04:05.999Z"),
		}
	}

	var err error
	for dataBuffer.Scan() {
		line := dataBuffer.Text()

		if line == NotFoundResponse {
			r.log.Errorf("Replay file not found '%s'", name)
			return false
		}

		*nextTime, *nextData, err = splitData(line, sessionStartTime)
		if err != nil {
			continue
		}

		if nextTime.After(currentRaceTime) {
			return true
		}

		r.dataFeed <- Payload{
			Name:      name,
			Data:      []byte(*nextData),
			Timestamp: nextTime.Format("2006-01-02T15:04:05.999Z"),
		}
	}

	if !dataBuffer.Scan() {
		*nextData = ""
		return false
	}

	return true
}

func (r *replay) findSessionStartTime() (time.Time, error) {
	dataBuffer := r.get(r.eventUrl + ExtrapolatedClockFile + ".jsonStream")

	if dataBuffer == nil {
		r.log.Errorf("Unable to find session start time because file doesn't exist")
		return time.Time{}, errors.New("No file for session start time")
	}

	dataBuffer.Scan()
	line := dataBuffer.Text()

	if line == NotFoundResponse {
		r.log.Errorf("Session start time not found because %s file was not found.", ExtrapolatedClockFile)
		return time.Time{}, errors.New("session start time not found")
	}

	timeEnd := strings.Index(line, "{")
	data := line[timeEnd:]
	timestamp := line[timeEnd-12 : timeEnd]

	abc := fmt.Sprintf("%sh%sm%ss%sms", timestamp[:2], timestamp[3:5], timestamp[6:8], timestamp[9:12])

	offsetFromStart, err := time.ParseDuration(abc)

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		r.log.Errorf("Session start date file was invalid: %s", err)
		return time.Time{}, err
	}

	timestampStr := dat["Utc"].(string)

	sessionUtc, err := time.Parse("2006-01-02T15:04:05.9999999Z", timestampStr)
	if err != nil {
		r.log.Errorf("Session start timestamp was invalid: %s", err)
		return time.Time{}, err
	}

	return sessionUtc.Add(-offsetFromStart), nil
}

func (r *replay) uncompressedDataTime(data string, sessionStart time.Time) (timestamp time.Time, payload string, err error) {
	timeEnd := strings.Index(data, "{")

	timestamp, err = r.raceTime(data[timeEnd-12:timeEnd], sessionStart)
	if err != nil {
		return time.Time{}, "", err
	}

	return timestamp, data[timeEnd:], nil
}

func (r *replay) compressedDataTime(data string, sessionStart time.Time) (timestamp time.Time, payload string, err error) {
	timeEnd := strings.Index(data, "\"")

	timestamp, err = r.raceTime(data[timeEnd-12:timeEnd], sessionStart)
	if err != nil {
		return time.Time{}, "", err
	}

	return timestamp, data[timeEnd+1 : len(data)-1], nil
}

func (r *replay) raceTime(value string, sessionStart time.Time) (time.Time, error) {

	// There is some weird characters at the start of the string you can't see
	abc := fmt.Sprintf("%sh%sm%ss%sms", value[:2], value[3:5], value[6:8], value[9:12])

	timestamp, err := time.ParseDuration(abc)
	if err != nil {
		r.log.Errorf("Replay error parsing time '%s': %s", abc, err)
		return time.Time{}, err
	}

	return sessionStart.Add(timestamp), nil
}

func (r *replay) get(url string) *bufio.Scanner {

	if len(r.cache) > 0 {
		dataPath := strings.Replace(url, "https://livetiming.formula1.com/static/", "", 1)

		// If file matching url doesn't exist then retrieve
		cachedFile := filepath.Join(r.cache, dataPath)
		cachedFile, _ = filepath.Abs(cachedFile)
		f, err := os.Open(cachedFile)

		if os.IsNotExist(err) {
			f.Close()

			var resp *http.Response
			resp, err = http.Get(url)
			if err != nil {
				r.log.Errorf("Replay url error for '%s': %s", url, err)
				return nil
			}
			defer resp.Body.Close()

			if resp.ContentLength == int64(len(NotFoundResponse)) {
				content, _ := io.ReadAll(resp.Body)
				if string(content) == NotFoundResponse {
					r.log.Errorf("Replay url not found '%s'", url)
					return nil
				}
			}

			scanner := bufio.NewScanner(resp.Body)

			err = os.MkdirAll(filepath.Dir(cachedFile), 0755)

			// Write body to file - using url as name
			var newFile *os.File
			newFile, err = os.Create(cachedFile)
			defer newFile.Close()
			for scanner.Scan() {
				_, err = newFile.Write(scanner.Bytes())

				// need newline for scanner to split
				newFile.WriteString("\n")
			}
			f, err = os.Open(cachedFile)
		}

		return bufio.NewScanner(f)
	}

	var resp *http.Response
	resp, err := http.Get(url)
	if err != nil {
		r.log.Errorf("Replay get url '%s': %s", err)
		return nil
	}
	// TODO - probably need to tidy this up but if we have no cache then we can't close it here or no data
	//defer resp.Body.Close()

	if resp.ContentLength == int64(len(NotFoundResponse)) {
		content, _ := io.ReadAll(resp.Body)
		if string(content) == NotFoundResponse {
			r.log.Errorf("Replay url not found '%s'", url)
			return nil
		}
	}

	return bufio.NewScanner(resp.Body)
}
