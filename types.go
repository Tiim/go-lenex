package lenex

import (
	"fmt"
	"time"
)

// ── Enumerations ─────────────────────────────────────────────────────────────

type Course string

const (
	CourseLCM   Course = "LCM"
	CourseSCM   Course = "SCM"
	CourseSCY   Course = "SCY"
	CourseSCM16 Course = "SCM16"
	CourseSCM20 Course = "SCM20"
	CourseSCM33 Course = "SCM33"
	CourseSCY20 Course = "SCY20"
	CourseSCY27 Course = "SCY27"
	CourseSCY33 Course = "SCY33"
	CourseSCY36 Course = "SCY36"
	CourseOpen  Course = "OPEN"
)

type Gender string

const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
	GenderMixed  Gender = "X"
	GenderAll    Gender = "A"
)

type Stroke string

const (
	StrokeBack      Stroke = "BACK"
	StrokeBreast    Stroke = "BREAST"
	StrokeFly       Stroke = "FLY"
	StrokeFree      Stroke = "FREE"
	StrokeIMRelay   Stroke = "IMRELAY"
	StrokeMedley    Stroke = "MEDLEY"
	StrokeUnknown   Stroke = "UNKNOWN"
	// Fin swimming strokes
	StrokeApnea         Stroke = "APNEA"
	StrokeBifins        Stroke = "BIFINS"
	StrokeDynamic       Stroke = "DYNAMIC"
	StrokeDynamicBifins Stroke = "DYNAMIC_BIFINS"
	StrokeDynamicNofins Stroke = "DYNAMIC_NOFINS"
	StrokeImmersion     Stroke = "IMMERSION"
	StrokeMultiple      Stroke = "MULTIPLE"
	StrokeSpeedApnea    Stroke = "SPEED_APNEA"
	StrokeSpeedEndurance Stroke = "SPEED_ENDURANCE"
	StrokeStatic        Stroke = "STATIC"
	StrokeSurface       Stroke = "SURFACE"
)

type PerformanceStatus string

const (
	StatusOK   PerformanceStatus = ""
	StatusEXH  PerformanceStatus = "EXH"
	StatusRJC  PerformanceStatus = "RJC"
	StatusDSQ  PerformanceStatus = "DSQ"
	StatusDNS  PerformanceStatus = "DNS"
	StatusDNF  PerformanceStatus = "DNF"
	StatusSICK PerformanceStatus = "SICK"
	StatusWDR  PerformanceStatus = "WDR"
)

type RelayPositionStatus string

const (
	RelayPositionOK  RelayPositionStatus = ""
	RelayPositionDSQ RelayPositionStatus = "DSQ"
	RelayPositionDNF RelayPositionStatus = "DNF"
)

// ── Value types ──────────────────────────────────────────────────────────────

// SwimTime represents a swim time or NT (no time).
type SwimTime struct {
	Duration time.Duration
	NT       bool
}

func (s SwimTime) String() string {
	if s.NT {
		return "NT"
	}
	d := s.Duration
	hh := int(d.Hours())
	d -= time.Duration(hh) * time.Hour
	mm := int(d.Minutes())
	d -= time.Duration(mm) * time.Minute
	ss := int(d.Seconds())
	d -= time.Duration(ss) * time.Second
	cs := int(d.Milliseconds() / 10)
	if hh > 0 {
		return fmt.Sprintf("%d:%02d:%02d.%02d", hh, mm, ss, cs)
	}
	return fmt.Sprintf("%d:%02d.%02d", mm, ss, cs)
}

// Daytime is a time of day (local time of the meet), e.g. a session start.
type Daytime struct {
	Hour   int
	Minute int
}

// Duration is the time since midnight.
func (d Daytime) Duration() time.Duration {
	return time.Duration(d.Hour)*time.Hour + time.Duration(d.Minute)*time.Minute
}

func (d Daytime) String() string { return fmt.Sprintf("%02d:%02d", d.Hour, d.Minute) }

// ReactionTime represents a reaction time in hundredths of a second.
// A nil *ReactionTime means the reaction time was not recorded.
type ReactionTime struct {
	Hundredths int
}

// EventKey uniquely identifies an event across courses and relay types.
type EventKey struct {
	Distance   int
	Stroke     Stroke
	RelayCount int
	Course     Course
	// Only compared when Stroke == StrokeUnknown.
	Code string
}

// Split is a single cumulative split time at a given distance.
type Split struct {
	Distance int
	Time     SwimTime
}

// RelayPosition is one swimmer's leg in a relay.
type RelayPosition struct {
	Number       int
	AthleteID    string
	ReactionTime *ReactionTime
	Status       RelayPositionStatus
}

// ── Domain types ─────────────────────────────────────────────────────────────

// Contact holds contact address information.
type Contact struct {
	Name     string
	Email    string
	City     string
	Country  string
	Fax      string
	Internet string
	Mobile   string
	Phone    string
	State    string
	Street   string
	Street2  string
	Zip      string
}

// Constructor identifies the software that created the LENEX file.
type Constructor struct {
	Name         string
	Version      string
	Registration string
	Contact      Contact
}

// Club is a swimming club or national team.
type Club struct {
	doc    *Document
	meetID string

	// ClubID is unique within the meet: CLUB.clubid, else the club code if
	// no other club of the meet shares it, else a generated "club-<index>".
	ClubID    string
	Code      string
	Name      string
	NameEn    string
	Nation    string
	Number    string
	Region    string
	ShortName string
	Type      string
}

// Athlete is a swimmer in the context of a meet or a record.
type Athlete struct {
	doc    *Document
	meetID string

	AthleteID string
	ClubID    string // set from parent CLUB in meet context; from inline CLUB in record context
	FirstName string
	LastName  string
	Gender    Gender
	BirthDate time.Time
	License   string
	Nation    string
}

// Club returns the athlete's club via the meet index.
func (a *Athlete) Club() *Club {
	if a.doc == nil || a.meetID == "" {
		return nil
	}
	idx := a.doc.meetIndexFor(a.meetID)
	if idx == nil {
		return nil
	}
	return idx.clubsByID[a.ClubID]
}

// Official is a meet official.
type Official struct {
	doc    *Document
	meetID string

	OfficialID string
	FirstName  string
	LastName   string
	Gender     Gender
	Grade      string
	License    string
	Nation     string
}

// Coach is a club coach or staff member.
type Coach struct {
	FirstName string
	LastName  string
	Gender    Gender
	Nation    string
	Type      string
}

// Judge links an official to a session with a specific role.
type Judge struct {
	doc    *Document
	meetID string

	OfficialID string
	Number     string
	Role       string
	Status     string
	Remarks    string
}

// Official returns the official referenced by this judge assignment.
func (j *Judge) Official() *Official {
	if j.doc == nil || j.meetID == "" {
		return nil
	}
	idx := j.doc.meetIndexFor(j.meetID)
	if idx == nil {
		return nil
	}
	return idx.officialsByID[j.OfficialID]
}

// Heat contains heat-level details within an event.
type Heat struct {
	doc    *Document
	meetID string

	HeatID  string
	Number  int
	Final   string
	Status  string
}

// Event is one event (or round) within a session.
type Event struct {
	doc    *Document
	meetID string

	EventID   string
	Number    int
	Gender    Gender
	Round     string
	Status    string
	Key       EventKey
	AgeGroups []*AgeGroup
}

// Session is one session of a meet.
type Session struct {
	doc    *Document
	meetID string

	Number int
	Date   time.Time
	Course Course
	Name   string
	Status string
	Pool   *Pool

	// Times of day; nil when not given in the file.
	DayTime     *Daytime // start of the session
	EndTime     *Daytime // estimated end of the session
	WarmupFrom  *Daytime
	WarmupUntil *Daytime
}

// Start is the start of the session (date and daytime), or ok=false when the
// file gives no start time.
func (s *Session) Start() (time.Time, bool) {
	if s.DayTime == nil {
		return time.Time{}, false
	}
	return s.Date.Add(s.DayTime.Duration()), true
}

// End is the estimated end of the session, or ok=false when the file gives
// no end time. An end time before the start is taken as the next day.
func (s *Session) End() (time.Time, bool) {
	if s.EndTime == nil {
		return time.Time{}, false
	}
	end := s.Date.Add(s.EndTime.Duration())
	if start, ok := s.Start(); ok && end.Before(start) {
		end = end.AddDate(0, 0, 1)
	}
	return end, true
}

// Pool describes a swimming pool.
type Pool struct {
	LaneMin int
	LaneMax int
	Type    string
}

// AgeGroup is an age range used in events and record/time-standard lists.
type AgeGroup struct {
	doc    *Document
	meetID string // empty when owned by a RecordList/TimeStandardList

	AgeGroupID string
	AgeMin     int
	AgeMax     int
	Gender     Gender
	Name       string
	Rankings   []Ranking
}

// Ranking references a result and its place within an age group.
type Ranking struct {
	doc    *Document
	meetID string

	Place    int
	ResultID string
	Medal    string
}

// Result returns the Performance referenced by this Ranking.
func (r *Ranking) Result() *Performance {
	if r.doc == nil || r.meetID == "" {
		return nil
	}
	idx := r.doc.meetIndexFor(r.meetID)
	if idx == nil {
		return nil
	}
	return idx.resultsByID[r.ResultID]
}

// Relay is a relay team.
type Relay struct {
	doc    *Document
	meetID string

	RelayID string
	ClubID  string
	Name    string
	Number  string
	Gender  Gender
}

// Club returns the relay's club.
func (r *Relay) Club() *Club {
	if r.doc == nil || r.meetID == "" {
		return nil
	}
	idx := r.doc.meetIndexFor(r.meetID)
	if idx == nil {
		return nil
	}
	return idx.clubsByID[r.ClubID]
}

// Performance represents a single swim — an Entry, Result, or Record.
type Performance struct {
	doc    *Document
	meetID string // empty for record-context performances

	AthleteID string
	RelayID   string
	ClubID    string
	EventKey  EventKey

	Time         SwimTime
	Status       PerformanceStatus
	ReactionTime *ReactionTime
	Splits       []Split
	RelayPositions []RelayPosition

	Date time.Time

	// Source distinguishes the origin of this performance.
	Source PerformanceSource

	// For result-context performances:
	ResultID string
	HeatID   string
	Lane     int
	Points   int
}

// PerformanceSource describes where a Performance came from.
type PerformanceSource int

const (
	SourceEntry  PerformanceSource = iota // from ENTRY element
	SourceResult                          // from RESULT element
	SourceRecord                          // from RECORD element
)

// Athlete returns the athlete who swam this performance (nil for relays and target times).
func (p *Performance) Athlete() *Athlete {
	if p.AthleteID == "" || p.doc == nil {
		return nil
	}
	if p.meetID == "" {
		return nil
	}
	idx := p.doc.meetIndexFor(p.meetID)
	if idx == nil {
		return nil
	}
	return idx.athletesByID[p.AthleteID]
}

// Club returns the club for this performance.
func (p *Performance) Club() *Club {
	if p.ClubID == "" || p.doc == nil {
		return nil
	}
	if p.meetID == "" {
		return nil
	}
	idx := p.doc.meetIndexFor(p.meetID)
	if idx == nil {
		return nil
	}
	return idx.clubsByID[p.ClubID]
}

// TimeStandard is a target time for a specific event style.
type TimeStandard struct {
	EventKey EventKey
	Time     SwimTime
}

// RecordList is a named collection of records.
type RecordList struct {
	Name     string
	Type     string
	Gender   Gender
	Course   Course
	Nation   string
	Region   string
	AgeGroup *AgeGroup
	Records  []*Performance
}

// TimeStandardList is a named collection of time standards.
type TimeStandardList struct {
	TimeStandardListID string
	Name               string
	Gender             Gender
	Course             Course
	Type               string
	AgeGroup           *AgeGroup
	Standards          []*TimeStandard
}

// Meet is one meet within a LENEX document.
type Meet struct {
	doc *Document

	MeetID    string // internal id (index into doc.Meets)
	Name      string
	NameEn    string
	City      string
	CityEn    string
	HostClub  string
	Organizer string
	Nation  string
	Course  Course
	Status  string
	Type    string
	Pool    *Pool

	Sessions           []*Session
	Clubs              []*Club
	Officials          []*Official
	AgeGroups          []*AgeGroup
}

// Results returns the performances of type SourceResult of this meet.
func (m *Meet) Results() ResultSet {
	var rs ResultSet
	if m.doc == nil {
		return rs
	}
	if idx := m.doc.meetsByID[m.MeetID]; idx != nil {
		for _, p := range idx.resultsByID {
			rs.all = append(rs.all, p)
		}
	}
	return rs
}
