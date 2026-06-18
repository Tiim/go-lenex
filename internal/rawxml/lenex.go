// Package rawxml contains unexported XML structs that mirror the LENEX 3.0 schema.
// These are used only by the parse pipeline; consumers never see them.
package rawxml

import "encoding/xml"

type LENEX struct {
	XMLName           xml.Name           `xml:"LENEX"`
	Version           string             `xml:"version,attr"`
	Created           *Timestamp         `xml:"created,attr"`
	RevisionDate      *Date              `xml:"revisiondate,attr"`
	Constructor       Constructor        `xml:"CONSTRUCTOR"`
	Meets             []Meet             `xml:"MEETS>MEET"`
	RecordLists       []RecordList       `xml:"RECORDLISTS>RECORDLIST"`
	TimeStandardLists []TimeStandardList `xml:"TIMESTANDARDLISTS>TIMESTANDARDLIST"`
}

// ── CONSTRUCTOR / CONTACT ────────────────────────────────────────────────────

type Constructor struct {
	Name         string  `xml:"name,attr"`
	Version      string  `xml:"version,attr"`
	Registration string  `xml:"registration,attr"`
	Contact      Contact `xml:"CONTACT"`
}

type Contact struct {
	Name     string `xml:"name,attr"`
	Email    string `xml:"email,attr"`
	City     string `xml:"city,attr"`
	Country  string `xml:"country,attr"`
	Fax      string `xml:"fax,attr"`
	Internet string `xml:"internet,attr"`
	Mobile   string `xml:"mobile,attr"`
	Phone    string `xml:"phone,attr"`
	State    string `xml:"state,attr"`
	Street   string `xml:"street,attr"`
	Street2  string `xml:"street2,attr"`
	Zip      string `xml:"zip,attr"`
}

// ── MEET ─────────────────────────────────────────────────────────────────────

type Meet struct {
	SwrID              string      `xml:"swrid,attr"`
	Name               string      `xml:"name,attr"`
	NameEn             string      `xml:"name.en,attr"`
	City               string      `xml:"city,attr"`
	CityEn             string      `xml:"city.en,attr"`
	Nation             string      `xml:"nation,attr"`
	Course             string      `xml:"course,attr"`
	Altitude           string      `xml:"altitude,attr"`
	Deadline           *Date       `xml:"deadline,attr"`
	DeadlineTime       *Daytime    `xml:"deadlinetime,attr"`
	EntryStartDate     *Date       `xml:"entrystartdate,attr"`
	EntryType          string      `xml:"entrytype,attr"`
	HostClub           string      `xml:"hostclub,attr"`
	HostClubURL        string      `xml:"hostclub.url,attr"`
	MaxEntriesAthlete  string      `xml:"maxentriesathlete,attr"`
	MaxEntriesRelay    string      `xml:"maxentriesrelay,attr"`
	Number             string      `xml:"number,attr"`
	Organizer          string      `xml:"organizer,attr"`
	OrganizerURL       string      `xml:"organizer.url,attr"`
	ReserveCount       string      `xml:"reservecount,attr"`
	ResultURL          string      `xml:"result.url,attr"`
	StartMethod        string      `xml:"startmethod,attr"`
	Status             string      `xml:"status,attr"`
	Timing             string      `xml:"timing,attr"`
	TouchpadMode       string      `xml:"touchpadmode,attr"`
	Type               string      `xml:"type,attr"`
	WithdrawUntil      *Date       `xml:"withdrawuntil,attr"`
	AgeDate            *AgeDate    `xml:"AGEDATE"`
	Bank               *Bank       `xml:"BANK"`
	Contact            *Contact    `xml:"CONTACT"`
	Facility           *Facility   `xml:"FACILITY"`
	Fees               []Fee       `xml:"FEES>FEE"`
	Pool               *Pool       `xml:"POOL"`
	PointTable         *PointTable `xml:"POINTTABLE"`
	Qualify            *Qualify    `xml:"QUALIFY"`
	Sessions           []Session   `xml:"SESSIONS>SESSION"`
	Clubs              []Club      `xml:"CLUBS>CLUB"`
}

type AgeDate struct {
	Type  string `xml:"type,attr"`
	Value Date   `xml:"value,attr"`
}

type Bank struct {
	AccountHolder string `xml:"accountholder,attr"`
	BIC           string `xml:"bic,attr"`
	IBAN          string `xml:"iban,attr"`
	Name          string `xml:"name,attr"`
	Note          string `xml:"note,attr"`
}

type Facility struct {
	City    string `xml:"city,attr"`
	Nation  string `xml:"nation,attr"`
	Name    string `xml:"name,attr"`
	State   string `xml:"state,attr"`
	Street  string `xml:"street,attr"`
	Street2 string `xml:"street2,attr"`
	Zip     string `xml:"zip,attr"`
}

type Fee struct {
	Currency string `xml:"currency,attr"`
	Type     string `xml:"type,attr"`
	Value    string `xml:"value,attr"`
}

type Pool struct {
	LaneMax     string `xml:"lanemax,attr"`
	LaneMin     string `xml:"lanemin,attr"`
	Temperature string `xml:"temperature,attr"`
	Type        string `xml:"type,attr"`
}

type PointTable struct {
	Name         string `xml:"name,attr"`
	PointTableID string `xml:"pointtableid,attr"`
	Version      string `xml:"version,attr"`
}

type Qualify struct {
	Conversion string `xml:"conversion,attr"`
	From       Date   `xml:"from,attr"`
	Percent    string `xml:"percent,attr"`
	Until      *Date  `xml:"until,attr"`
}

// ── SESSION ──────────────────────────────────────────────────────────────────

type Session struct {
	Number              int      `xml:"number,attr"`
	Date                Date     `xml:"date,attr"`
	Course              string   `xml:"course,attr"`
	DayTime             *Daytime `xml:"daytime,attr"`
	EndTime             *Daytime `xml:"endtime,attr"`
	MaxEntriesAthlete   string   `xml:"maxentriesathlete,attr"`
	MaxEntriesRelay     string   `xml:"maxentriesrelay,attr"`
	Name                string   `xml:"name,attr"`
	OfficialMeeting     *Daytime `xml:"officialmeeting,attr"`
	RemarksJudge        string   `xml:"remarksjudge,attr"`
	Status              string   `xml:"status,attr"`
	TeamLeaderMeeting   *Daytime `xml:"teamleadermeeting,attr"`
	Timing              string   `xml:"timing,attr"`
	TouchpadMode        string   `xml:"touchpadmode,attr"`
	WarmupFrom          *Daytime `xml:"warmupfrom,attr"`
	WarmupUntil         *Daytime `xml:"warmupuntil,attr"`
	Pool                *Pool    `xml:"POOL"`
	Events              []Event  `xml:"EVENTS>EVENT"`
	Fees                []Fee    `xml:"FEES>FEE"`
	Judges              []Judge  `xml:"JUDGES>JUDGE"`
}

// ── EVENT ────────────────────────────────────────────────────────────────────

type Event struct {
	EventID     string          `xml:"eventid,attr"`
	Number      int             `xml:"number,attr"`
	Gender      string          `xml:"gender,attr"`
	DayTime     *Daytime        `xml:"daytime,attr"`
	MaxEntries  string          `xml:"maxentries,attr"`
	Order       string          `xml:"order,attr"`
	PrevEventID string          `xml:"preveventid,attr"`
	Round       string          `xml:"round,attr"`
	Run         string          `xml:"run,attr"`
	Status      string          `xml:"status,attr"`
	Timing      string          `xml:"timing,attr"`
	Type        string          `xml:"type,attr"`
	SwimStyle   SwimStyle       `xml:"SWIMSTYLE"`
	AgeGroups   []AgeGroup      `xml:"AGEGROUPS>AGEGROUP"`
	Heats       []Heat          `xml:"HEATS>HEAT"`
	TSRefs      []TimeStdRef    `xml:"TIMESTANDARDREFS>TIMESTANDARDREF"`
	Fee         *Fee            `xml:"FEE"`
}

type SwimStyle struct {
	Distance    int    `xml:"distance,attr"`
	Stroke      string `xml:"stroke,attr"`
	RelayCount  int    `xml:"relaycount,attr"`
	Code        string `xml:"code,attr"`
	Name        string `xml:"name,attr"`
	SwimStyleID string `xml:"swimstyleid,attr"`
	Technique   string `xml:"technique,attr"`
}

type AgeGroup struct {
	AgeGroupID string    `xml:"agegroupid,attr"`
	AgeMin     int       `xml:"agemin,attr"`
	AgeMax     int       `xml:"agemax,attr"`
	Gender     string    `xml:"gender,attr"`
	Calculate  string    `xml:"calculate,attr"`
	Handicap   string    `xml:"handicap,attr"`
	LevelMax   string    `xml:"levelmax,attr"`
	LevelMin   string    `xml:"levelmin,attr"`
	Levels     string    `xml:"levels,attr"`
	Name       string    `xml:"name,attr"`
	Rankings   []Ranking `xml:"RANKINGS>RANKING"`
}

type Heat struct {
	HeatID     string   `xml:"heatid,attr"`
	Number     int      `xml:"number,attr"`
	AgeGroupID string   `xml:"agegroupid,attr"`
	DayTime    *Daytime `xml:"daytime,attr"`
	Final      string   `xml:"final,attr"`
	Order      string   `xml:"order,attr"`
	Status     string   `xml:"status,attr"`
}

type Ranking struct {
	Place    int    `xml:"place,attr"`
	ResultID string `xml:"resultid,attr"`
	Medal    string `xml:"medal,attr"`
	Order    string `xml:"order,attr"`
}

type TimeStdRef struct {
	TimeStandardListID string `xml:"timestandardlistid,attr"`
	Marker             string `xml:"marker,attr"`
	Fee                *Fee   `xml:"FEE"`
}

// ── JUDGE ────────────────────────────────────────────────────────────────────

type Judge struct {
	OfficialID string `xml:"officialid,attr"`
	Number     string `xml:"number,attr"`
	Remarks    string `xml:"remarks,attr"`
	Role       string `xml:"role,attr"`
	Status     string `xml:"status,attr"`
}

// ── CLUB ─────────────────────────────────────────────────────────────────────

type Club struct {
	ClubID    string     `xml:"clubid,attr"`
	Code      string     `xml:"code,attr"`
	Name      string     `xml:"name,attr"`
	NameEn    string     `xml:"name.en,attr"`
	Nation    string     `xml:"nation,attr"`
	Number    string     `xml:"number,attr"`
	Region    string     `xml:"region,attr"`
	ShortName string     `xml:"shortname,attr"`
	ShortNameEn string   `xml:"shortname.en,attr"`
	SwrID     string     `xml:"swrid,attr"`
	Type      string     `xml:"type,attr"`
	Contact   *Contact   `xml:"CONTACT"`
	Athletes  []Athlete  `xml:"ATHLETES>ATHLETE"`
	Coaches   []Coach    `xml:"COACHES>COACH"`
	Officials []Official `xml:"OFFICIALS>OFFICIAL"`
	Relays    []Relay    `xml:"RELAYS>RELAY"`
}

// ── ATHLETE ──────────────────────────────────────────────────────────────────

type Athlete struct {
	AthleteID   string    `xml:"athleteid,attr"`
	FirstName   string    `xml:"firstname,attr"`
	FirstNameEn string    `xml:"firstname.en,attr"`
	LastName    string    `xml:"lastname,attr"`
	LastNameEn  string    `xml:"lastname.en,attr"`
	NamePrefix  string    `xml:"nameprefix,attr"`
	Gender      string    `xml:"gender,attr"`
	BirthDate   *Date     `xml:"birthdate,attr"`
	License     string    `xml:"license,attr"`
	LicenseIPC  string    `xml:"license_ipc,attr"`
	Level       string    `xml:"level,attr"`
	Nation      string    `xml:"nation,attr"`
	Passport    string    `xml:"passport,attr"`
	Status      string    `xml:"status,attr"`
	SwrID       string    `xml:"swrid,attr"`
	Handicap    *Handicap `xml:"HANDICAP"`
	// In meet context:
	Entries []Entry  `xml:"ENTRIES>ENTRY"`
	Results []Result `xml:"RESULTS>RESULT"`
	// In record context: inline CLUB element
	Club *Club `xml:"CLUB"`
}

type Handicap struct {
	Breast       string `xml:"breast,attr"`
	BreastStatus string `xml:"breaststatus,attr"`
	Exception    string `xml:"exception,attr"`
	Free         string `xml:"free,attr"`
	FreeStatus   string `xml:"freestatus,attr"`
	Medley       string `xml:"medley,attr"`
	MedleyStatus string `xml:"medleystatus,attr"`
}

// ── COACH ────────────────────────────────────────────────────────────────────

type Coach struct {
	FirstName  string   `xml:"firstname,attr"`
	LastName   string   `xml:"lastname,attr"`
	NamePrefix string   `xml:"nameprefix,attr"`
	Gender     string   `xml:"gender,attr"`
	Nation     string   `xml:"nation,attr"`
	Passport   string   `xml:"passport,attr"`
	Type       string   `xml:"type,attr"`
	Contact    *Contact `xml:"CONTACT"`
}

// ── OFFICIAL ─────────────────────────────────────────────────────────────────

type Official struct {
	OfficialID string   `xml:"officialid,attr"`
	FirstName  string   `xml:"firstname,attr"`
	LastName   string   `xml:"lastname,attr"`
	NamePrefix string   `xml:"nameprefix,attr"`
	Gender     string   `xml:"gender,attr"`
	Grade      string   `xml:"grade,attr"`
	License    string   `xml:"license,attr"`
	Nation     string   `xml:"nation,attr"`
	Passport   string   `xml:"passport,attr"`
	Contact    *Contact `xml:"CONTACT"`
}

// ── RELAY ────────────────────────────────────────────────────────────────────

type Relay struct {
	RelayID      string   `xml:"relayid,attr"`
	Name         string   `xml:"name,attr"`
	Number       string   `xml:"number,attr"`
	Gender       string   `xml:"gender,attr"`
	AgeMin       string   `xml:"agemin,attr"`
	AgeMax       string   `xml:"agemax,attr"`
	AgeTotalMin  string   `xml:"agetotalmin,attr"`
	AgeTotalMax  string   `xml:"agetotalmax,attr"`
	Handicap     string   `xml:"handicap,attr"`
	Entries      []Entry  `xml:"ENTRIES>ENTRY"`
	Results      []Result `xml:"RESULTS>RESULT"`
	// In record context:
	Club           *Club           `xml:"CLUB"`
	RelayPositions []RelayPosition `xml:"RELAYPOSITIONS>RELAYPOSITION"`
}

// ── ENTRY ────────────────────────────────────────────────────────────────────

type Entry struct {
	EventID        string          `xml:"eventid,attr"`
	HeatID         string          `xml:"heatid,attr"`
	Lane           string          `xml:"lane,attr"`
	AgeGroupID     string          `xml:"agegroupid,attr"`
	EntryCourse    string          `xml:"entrycourse,attr"`
	EntryDistance  string          `xml:"entrydistance,attr"`
	EntryTime      *SwimTime       `xml:"entrytime,attr"`
	Handicap       string          `xml:"handicap,attr"`
	Status         string          `xml:"status,attr"`
	MeetInfo       *MeetInfo       `xml:"MEETINFO"`
	RelayPositions []RelayPosition `xml:"RELAYPOSITIONS>RELAYPOSITION"`
}

// ── RESULT ───────────────────────────────────────────────────────────────────

type Result struct {
	ResultID       string          `xml:"resultid,attr"`
	EventID        string          `xml:"eventid,attr"`
	HeatID         string          `xml:"heatid,attr"`
	Lane           string          `xml:"lane,attr"`
	SwimTime       SwimTime        `xml:"swimtime,attr"`
	Comment        string          `xml:"comment,attr"`
	Handicap       string          `xml:"handicap,attr"`
	Points         string          `xml:"points,attr"`
	ReactionTime   *ReactionTime   `xml:"reactiontime,attr"`
	Status         string          `xml:"status,attr"`
	SwimDistance   string          `xml:"swimdistance,attr"`
	RelayPositions []RelayPosition `xml:"RELAYPOSITIONS>RELAYPOSITION"`
	Splits         []Split         `xml:"SPLITS>SPLIT"`
}

// ── RELAYPOSITION ────────────────────────────────────────────────────────────

type RelayPosition struct {
	Number       int           `xml:"number,attr"`
	AthleteID    string        `xml:"athleteid,attr"`
	ReactionTime *ReactionTime `xml:"reactiontime,attr"`
	Status       string        `xml:"status,attr"`
	// In record context: inline ATHLETE
	Athlete  *Athlete  `xml:"ATHLETE"`
	MeetInfo *MeetInfo `xml:"MEETINFO"`
}

// ── MEETINFO ─────────────────────────────────────────────────────────────────

type MeetInfo struct {
	City              string    `xml:"city,attr"`
	Nation            string    `xml:"nation,attr"`
	State             string    `xml:"state,attr"`
	Name              string    `xml:"name,attr"`
	Date              *Date     `xml:"date,attr"`
	DayTime           *Daytime  `xml:"daytime,attr"`
	Course            string    `xml:"course,attr"`
	Timing            string    `xml:"timing,attr"`
	Approved          string    `xml:"approved,attr"`
	QualificationTime *SwimTime `xml:"qualificationtime,attr"`
	Pool              *Pool     `xml:"POOL"`
}

// ── SPLIT ────────────────────────────────────────────────────────────────────

type Split struct {
	Distance int      `xml:"distance,attr"`
	SwimTime SwimTime `xml:"swimtime,attr"`
}

// ── RECORDLIST ───────────────────────────────────────────────────────────────

type RecordList struct {
	RecordListID string     `xml:"recordlistid,attr"`
	Name         string     `xml:"name,attr"`
	Type         string     `xml:"type,attr"`
	Gender       string     `xml:"gender,attr"`
	Course       string     `xml:"course,attr"`
	Handicap     string     `xml:"handicap,attr"`
	Nation       string     `xml:"nation,attr"`
	Order        string     `xml:"order,attr"`
	Region       string     `xml:"region,attr"`
	Updated      *Date      `xml:"updated,attr"`
	AgeGroup     *AgeGroup  `xml:"AGEGROUP"`
	Records      []Record   `xml:"RECORDS>RECORD"`
}

type Record struct {
	SwimTime SwimTime  `xml:"swimtime,attr"`
	Comment  string    `xml:"comment,attr"`
	Status   string    `xml:"status,attr"`
	SwimStyle SwimStyle `xml:"SWIMSTYLE"`
	Athlete  *Athlete  `xml:"ATHLETE"`
	Relay    *Relay    `xml:"RELAY"`
	MeetInfo *MeetInfo `xml:"MEETINFO"`
	Splits   []Split   `xml:"SPLITS>SPLIT"`
}

// ── TIMESTANDARDLIST ─────────────────────────────────────────────────────────

type TimeStandardList struct {
	TimeStandardListID string        `xml:"timestandardlistid,attr"`
	Name               string        `xml:"name,attr"`
	Gender             string        `xml:"gender,attr"`
	Course             string        `xml:"course,attr"`
	Handicap           string        `xml:"handicap,attr"`
	Type               string        `xml:"type,attr"`
	AgeGroup           *AgeGroup     `xml:"AGEGROUP"`
	TimeStandards      []TimeStandard `xml:"TIMESTANDARDS>TIMESTANDARD"`
}

type TimeStandard struct {
	SwimTime  SwimTime  `xml:"swimtime,attr"`
	SwimStyle SwimStyle `xml:"SWIMSTYLE"`
}
