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

package Messages

import (
	"image/color"
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

type PitStop struct {
	Lap          int
	PitlaneEntry time.Time
	PitlaneExit  time.Time
	PitlaneTime  time.Duration
}

type Timing struct {
	Timestamp time.Time

	Position int

	Name      string
	ShortName string
	Number    int
	Team      string
	HexColor  string
	Color     color.RGBA

	TimeDiffToFastest       time.Duration
	TimeDiffToPositionAhead time.Duration
	GapToLeader             time.Duration

	PreviousSegmentIndex   int
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

	Pitstops     int
	PitStopTimes []PitStop

	Location CarLocation

	SpeedTrap                int
	SpeedTrapPersonalFastest bool
	SpeedTrapOverallFastest  bool
}
