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
