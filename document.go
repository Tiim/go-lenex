package lenex

import "time"

// Document is the root of a parsed LENEX file. All relationships are resolved
// through this object. Safe for concurrent reads after Parse returns.
type Document struct {
	Constructor       Constructor
	Meets             []*Meet
	RecordLists       []*RecordList
	TimeStandardLists []*TimeStandardList

	meetsByID map[string]*meetIndex
}

// meetIndex holds id-scoped indexes for one MEET element.
type meetIndex struct {
	clubsByID     map[string]*Club
	athletesByID  map[string]*Athlete
	relaysByID    map[string]*Relay
	eventsByID    map[string]*Event
	heatsByID     map[string]*Heat
	officialsByID map[string]*Official
	resultsByID   map[string]*Performance

	entriesByAthleteID map[string][]*Performance
	resultsByEvent     map[string][]*Performance
	judgesBySessionID  map[string][]*Judge

	// Date resolution: heatID/eventID → session date
	dateByEventID map[string]time.Time
	dateByHeatID  map[string]time.Time
}

func (d *Document) meetIndexFor(meetID string) *meetIndex {
	return d.meetsByID[meetID]
}

// Results returns all performances of type SourceResult across all meets.
func (d *Document) Results() ResultSet {
	var rs ResultSet
	for _, m := range d.Meets {
		idx := d.meetsByID[m.MeetID]
		if idx == nil {
			continue
		}
		for _, p := range idx.resultsByID {
			rs.all = append(rs.all, p)
		}
	}
	return rs
}

// Entries returns all performances of type SourceEntry across all meets.
func (d *Document) Entries() ResultSet {
	var rs ResultSet
	for _, m := range d.Meets {
		idx := d.meetsByID[m.MeetID]
		if idx == nil {
			continue
		}
		for _, entries := range idx.entriesByAthleteID {
			rs.all = append(rs.all, entries...)
		}
	}
	return rs
}

// ResultSet is a filterable set of performances.
type ResultSet struct {
	all []*Performance
}

// All returns all performances in this set.
func (rs ResultSet) All() []*Performance {
	return rs.all
}

// ByClub filters performances to those belonging to the given club code.
func (rs ResultSet) ByClub(code string) []*Performance {
	var out []*Performance
	for _, p := range rs.all {
		club := p.Club()
		if club != nil && club.Code == code {
			out = append(out, p)
		}
	}
	return out
}

// ByClubID filters performances to those belonging to the given club ID
// (Club.ClubID). Unlike the club code, every club has an ID.
func (rs ResultSet) ByClubID(clubID string) []*Performance {
	var out []*Performance
	for _, p := range rs.all {
		if p.ClubID == clubID {
			out = append(out, p)
		}
	}
	return out
}

// ByAthlete filters performances to those for the given athlete ID.
func (rs ResultSet) ByAthlete(athleteID string) []*Performance {
	var out []*Performance
	for _, p := range rs.all {
		if p.AthleteID == athleteID {
			out = append(out, p)
		}
	}
	return out
}
