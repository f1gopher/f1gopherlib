# F1Gopher Lib

A library for understanding and using the session data from Formula1.com writtein in Go.

There is an example GUI and command line client for the library:

* [F1Gopher Command Line](https://github.com/f1gopher/f1gopher-cmdline)
* [F1Gopher GUI](https://github.com/f1gopher/f1gopher)

## Features

* Supports data for all live sessions (pre-season testing, practice, qualifying, sprint and race)
* Supports replays of all session from 2018 and onward
* Live session can be paused and skipped forward to the live time
* Replay sessions can be paused and skipped through
* Provides data for:
  * Timing
  * Location on track
  * Car telemetry
  * Race control messages
  * Team radio messages (audio)
  * Weather

## Data

### Timing

* Current driver position and starting position
* Team color, name, abbreviated name, team for drivers
* Segment times
* Sector times (is personal or overall fastest)
* Last lap time (is personal or overall fastest)
* Current tire and laps on the tire
* Location (on track, outlap, pitlane, stopped...)
* Safety car status
* Track status (red flag, green flag...)
* Current lap and total number of laps
* Session time remaining
* Is DRS enabled
* Gap to the fastest time and gap to car infront
* Pitstop times
* Speed trap

### Location on Track

* X, Y, Z co-ordinate locations for all cars 
* Includes safety car when active

### Car Telemetry

* Six channels of telemetry for every car:
  * Throttle %
  * Brake %
  * RPM
  * Gear
  * Speed
  * DRS

### Race Control Messages

* Full text and timestamp for all race control messages

### Team Radio

* The mp3 audio for each message and the driver talking

### Weather

* Whether it is raining
* Air temperature
* Track temperature
* Wind speed
* Wind direction
* Air pressure
* Humidity