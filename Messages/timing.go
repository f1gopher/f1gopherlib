package Messages

import (
	"time"
)

type CarLocation int

// TODO - add garage and grid - need to calculate these based on speed and session type
const (
	NoLocation CarLocation = iota
	Pitlane
	PitOut
	OutLap
	OnTrack
	OutOfRace
	Stopped
)

func (c CarLocation) String() string {
	return [...]string{"Unknown", "Pitlane", "Pit Exit", "Out Lap", "On Track", "Out", "Stopped"}[c]
}

type TireType int

const (
	Unknown TireType = iota
	Soft
	Medium
	Hard
	Intermediate
	Wet
	Test
	HYPERSOFT
	ULTRASOFT
	SUPERSOFT
)

func (t TireType) String() string {
	return [...]string{"", "Soft", "Medium", "Hard", "Inter", "Wet", "Test", "Hyp Soft", "Ult Soft", "Sup Soft"}[t]
}

type SegmentType int

const (
	None SegmentType = iota
	YellowSegment
	GreenSegment
	InvalidSegment // Doesn't get displayed, cut corner/boundaries or invalid segment time?
	PurpleSegment
	RedSegment     // After chequered flag/stopped on track
	PitlaneSegment // In pitlane
	Mystery
	Mystery2 // ??? 2021 - Turkey Practice_2
	Mystery3 // ??? 2020 - Italy Race
)

type Timing struct {
	Timestamp time.Time

	Position int

	Name      string
	ShortName string
	Number    int
	Team      string
	Color     string

	TimeDiffToFastest       time.Duration
	TimeDiffToPositionAhead time.Duration
	GapToLeader             time.Duration

	Segment                [MaxSegments]SegmentType
	Sector1                time.Duration
	Sector1PersonalFastest bool
	Sector1OverallFastest  bool
	Sector2                time.Duration
	Sector2PersonalFastest bool
	Sector2OverallFastest  bool
	Sector3                time.Duration
	Sector3PersonalFastest bool
	Sector3OverallFastest  bool
	LastLap                time.Duration
	LastLapPersonalFastest bool
	LastLapOverallFastest  bool

	FastestLap        time.Duration
	OverallFastestLap bool

	KnockedOutOfQualifying bool
	ChequeredFlag          bool

	Tire       TireType
	LapsOnTire int
	Lap        int

	DRSOpen bool

	Pitstops int

	Location CarLocation

	SpeedTrap                int
	SpeedTrapPersonalFastest bool
	SpeedTrapOverallFastest  bool
}
