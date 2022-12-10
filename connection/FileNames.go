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

package connection

const DriverListFile = "DriverList"
const SessionInfoFile = "SessionInfo"
const LapCountFile = "LapCount"
const ExtrapolatedClockFile = "ExtrapolatedClock"
const TrackStatusFile = "TrackStatus"
const TimingDataFile = "TimingData"
const TimingAppDataFile = "TimingAppData"
const SessionStatusFile = "SessionStatus"
const SessionDataFile = "SessionData"
const HeartbeatFile = "Heartbeat"
const TimingStatsFile = "TimingStats"
const CarDataFile = "CarData.z"
const PositionFile = "Position.z"
const WeatherDataFile = "WeatherData"
const RaceControlMessagesFile = "RaceControlMessages"
const TopThreeFile = "TopThree"
const AudioStreamsFile = "AudioStreams"
const TeamRadioFile = "TeamRadio"
const ContentStreamsFile = "ContentStreams"

// These are special files that don't come from the raw data but we use internally
const EndOfDataFile = "EndOfData"
const CatchupFile = "Catchup"

var OrderedFiles = [...]string{
	DriverListFile,
	SessionInfoFile,
	LapCountFile,
	ExtrapolatedClockFile,
	TrackStatusFile,
	TimingDataFile,
	TimingAppDataFile,
	SessionStatusFile,
	SessionDataFile,
	HeartbeatFile,
	TimingStatsFile,
	CarDataFile,
	PositionFile,
	WeatherDataFile,
	RaceControlMessagesFile,
	TopThreeFile,
	AudioStreamsFile,
	TeamRadioFile,
	ContentStreamsFile,
}
