package lenex

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Tiim/go-lenex/internal/rawxml"
)

// Warning is a non-fatal issue found during parsing.
// Currently a placeholder — specific cases will be added as real-world
// malformed files are encountered.
type Warning struct{}

// Parse reads a LENEX file (either .lef plain XML or .lxf zip-wrapped XML)
// and returns a fully built Document. The parse pipeline runs to completion
// before returning; the returned Document is safe for concurrent reads.
func Parse(r io.Reader) (*Document, []Warning, error) {
	xmlBytes, err := readXML(r)
	if err != nil {
		return nil, nil, err
	}

	var raw rawxml.LENEX
	if err := xml.Unmarshal(xmlBytes, &raw); err != nil {
		return nil, nil, fmt.Errorf("xml decode: %w", err)
	}

	doc, warnings := normalize(&raw)
	return doc, warnings, nil
}

// normalize converts the raw XML tree into the domain model.
func normalize(raw *rawxml.LENEX) (*Document, []Warning) {
	doc := &Document{
		Constructor: normalizeConstructor(raw.Constructor),
		meetsByID:   make(map[string]*meetIndex),
	}

	for i := range raw.Meets {
		meetID := strconv.Itoa(i)
		meet, idx := normalizeMeet(&raw.Meets[i], meetID, doc)
		doc.Meets = append(doc.Meets, meet)
		doc.meetsByID[meetID] = idx
	}

	for i := range raw.RecordLists {
		rl := normalizeRecordList(&raw.RecordLists[i], doc)
		doc.RecordLists = append(doc.RecordLists, rl)
	}

	for i := range raw.TimeStandardLists {
		tsl := normalizeTimeStandardList(&raw.TimeStandardLists[i])
		doc.TimeStandardLists = append(doc.TimeStandardLists, tsl)
	}

	return doc, nil
}

func normalizeConstructor(rc rawxml.Constructor) Constructor {
	return Constructor{
		Name:         rc.Name,
		Version:      rc.Version,
		Registration: rc.Registration,
		Contact:      normalizeContact(rc.Contact),
	}
}

func normalizeContact(rc rawxml.Contact) Contact {
	return Contact{
		Name:     rc.Name,
		Email:    rc.Email,
		City:     rc.City,
		Country:  rc.Country,
		Fax:      rc.Fax,
		Internet: rc.Internet,
		Mobile:   rc.Mobile,
		Phone:    rc.Phone,
		State:    rc.State,
		Street:   rc.Street,
		Street2:  rc.Street2,
		Zip:      rc.Zip,
	}
}

func normalizeMeet(rm *rawxml.Meet, meetID string, doc *Document) (*Meet, *meetIndex) {
	idx := &meetIndex{
		clubsByID:          make(map[string]*Club),
		athletesByID:       make(map[string]*Athlete),
		relaysByID:         make(map[string]*Relay),
		eventsByID:         make(map[string]*Event),
		heatsByID:          make(map[string]*Heat),
		officialsByID:      make(map[string]*Official),
		resultsByID:        make(map[string]*Performance),
		entriesByAthleteID: make(map[string][]*Performance),
		resultsByEvent:     make(map[string][]*Performance),
		judgesBySessionID:  make(map[string][]*Judge),
		dateByEventID:      make(map[string]time.Time),
		dateByHeatID:       make(map[string]time.Time),
	}

	meet := &Meet{
		doc:    doc,
		MeetID: meetID,
		Name:   rm.Name,
		City:   rm.City,
		Nation: rm.Nation,
		Course: Course(rm.Course),
		Status: rm.Status,
		Type:   rm.Type,
	}

	// Resolve session course: sessions inherit meet course if not overridden.
	for i := range rm.Sessions {
		rs := &rm.Sessions[i]
		sessionCourse := Course(rs.Course)
		if sessionCourse == "" {
			sessionCourse = meet.Course
		}

		sess := &Session{
			doc:    doc,
			meetID: meetID,
			Number: rs.Number,
			Date:   rs.Date.Time,
			Course: sessionCourse,
			Name:   rs.Name,
			Status: rs.Status,
		}
		if rs.Pool != nil {
			sess.Pool = normalizePool(rs.Pool)
		}
		meet.Sessions = append(meet.Sessions, sess)

		// Index events and heats within this session.
		for j := range rs.Events {
			re := &rs.Events[j]
			evt := normalizeEvent(re, meetID, doc, sessionCourse)
			idx.eventsByID[evt.EventID] = evt
			idx.dateByEventID[evt.EventID] = sess.Date

			for k := range re.Heats {
				rh := &re.Heats[k]
				heat := &Heat{
					doc:    doc,
					meetID: meetID,
					HeatID: rh.HeatID,
					Number: rh.Number,
					Final:  rh.Final,
					Status: rh.Status,
				}
				idx.heatsByID[heat.HeatID] = heat
				idx.dateByHeatID[heat.HeatID] = sess.Date
			}
		}

		// Index judges, keyed by session number (string) for lookups.
		sessKey := strconv.Itoa(rs.Number)
		for k := range rs.Judges {
			rj := &rs.Judges[k]
			judge := &Judge{
				doc:        doc,
				meetID:     meetID,
				OfficialID: rj.OfficialID,
				Number:     rj.Number,
				Role:       rj.Role,
				Status:     rj.Status,
				Remarks:    rj.Remarks,
			}
			idx.judgesBySessionID[sessKey] = append(idx.judgesBySessionID[sessKey], judge)
		}
	}

	// Process clubs, athletes, officials, relays.
	for i := range rm.Clubs {
		rc := &rm.Clubs[i]
		clubID := rc.ClubID
		if clubID == "" {
			clubID = rc.Code
		}
		club := &Club{
			doc:       doc,
			meetID:    meetID,
			ClubID:    clubID,
			Code:      rc.Code,
			Name:      rc.Name,
			NameEn:    rc.NameEn,
			Nation:    rc.Nation,
			Number:    rc.Number,
			Region:    rc.Region,
			ShortName: rc.ShortName,
			Type:      rc.Type,
		}
		idx.clubsByID[clubID] = club
		meet.Clubs = append(meet.Clubs, club)

		for j := range rc.Officials {
			ro := &rc.Officials[j]
			off := &Official{
				doc:        doc,
				meetID:     meetID,
				OfficialID: ro.OfficialID,
				FirstName:  ro.FirstName,
				LastName:   ro.LastName,
				Gender:     Gender(ro.Gender),
				Grade:      ro.Grade,
				License:    ro.License,
				Nation:     ro.Nation,
			}
			idx.officialsByID[off.OfficialID] = off
			meet.Officials = append(meet.Officials, off)
		}

		for j := range rc.Athletes {
			ra := &rc.Athletes[j]
			ath := &Athlete{
				doc:       doc,
				meetID:    meetID,
				AthleteID: ra.AthleteID,
				ClubID:    clubID,
				FirstName: ra.FirstName,
				LastName:  ra.LastName,
				Gender:    Gender(ra.Gender),
				License:   ra.License,
				Nation:    ra.Nation,
			}
			if ra.BirthDate != nil {
				ath.BirthDate = ra.BirthDate.Time
			}
			idx.athletesByID[ath.AthleteID] = ath

			// Results
			for k := range ra.Results {
				rr := &ra.Results[k]
				evt := idx.eventsByID[rr.EventID]
				p := normalizeResult(rr, ath.AthleteID, clubID, meetID, doc, evt, idx)
				idx.resultsByID[rr.ResultID] = p
				idx.resultsByEvent[rr.EventID] = append(idx.resultsByEvent[rr.EventID], p)
			}

			// Entries
			for k := range ra.Entries {
				re := &ra.Entries[k]
				evt := idx.eventsByID[re.EventID]
				p := normalizeEntry(re, ath.AthleteID, clubID, meetID, doc, evt, idx)
				idx.entriesByAthleteID[ath.AthleteID] = append(idx.entriesByAthleteID[ath.AthleteID], p)
			}
		}

		for j := range rc.Relays {
			rr := &rc.Relays[j]
			relayID := rr.RelayID
			if relayID == "" {
				relayID = fmt.Sprintf("%s-relay-%d", clubID, j)
			}
			relay := &Relay{
				doc:     doc,
				meetID:  meetID,
				RelayID: relayID,
				ClubID:  clubID,
				Name:    rr.Name,
				Number:  rr.Number,
				Gender:  Gender(rr.Gender),
			}
			idx.relaysByID[relayID] = relay

			for k := range rr.Results {
				res := &rr.Results[k]
				evt := idx.eventsByID[res.EventID]
				p := normalizeResult(res, "", clubID, meetID, doc, evt, idx)
				p.RelayID = relayID
				idx.resultsByID[res.ResultID] = p
				idx.resultsByEvent[res.EventID] = append(idx.resultsByEvent[res.EventID], p)
			}

			for k := range rr.Entries {
				ent := &rr.Entries[k]
				evt := idx.eventsByID[ent.EventID]
				p := normalizeEntry(ent, "", clubID, meetID, doc, evt, idx)
				p.RelayID = relayID
				idx.entriesByAthleteID["relay:"+relayID] = append(idx.entriesByAthleteID["relay:"+relayID], p)
			}
		}
	}

	return meet, idx
}

func normalizePool(rp *rawxml.Pool) *Pool {
	p := &Pool{Type: rp.Type}
	if rp.LaneMin != "" {
		p.LaneMin, _ = strconv.Atoi(rp.LaneMin)
	}
	if rp.LaneMax != "" {
		p.LaneMax, _ = strconv.Atoi(rp.LaneMax)
	}
	return p
}

func normalizeEvent(re *rawxml.Event, meetID string, doc *Document, course Course) *Event {
	return &Event{
		doc:     doc,
		meetID:  meetID,
		EventID: re.EventID,
		Number:  re.Number,
		Gender:  Gender(re.Gender),
		Round:   re.Round,
		Status:  re.Status,
		Key: EventKey{
			Distance:   re.SwimStyle.Distance,
			Stroke:     Stroke(re.SwimStyle.Stroke),
			RelayCount: re.SwimStyle.RelayCount,
			Course:     course,
			Code:       re.SwimStyle.Code,
		},
	}
}

func normalizeResult(
	rr *rawxml.Result,
	athleteID, clubID, meetID string,
	doc *Document,
	evt *Event,
	idx *meetIndex,
) *Performance {
	p := &Performance{
		doc:      doc,
		meetID:   meetID,
		AthleteID: athleteID,
		ClubID:   clubID,
		ResultID: rr.ResultID,
		HeatID:   rr.HeatID,
		Status:   PerformanceStatus(rr.Status),
		Source:   SourceResult,
		Time:     SwimTime(rr.SwimTime),
	}
	if evt != nil {
		p.EventKey = evt.Key
	}
	if rr.ReactionTime != nil {
		p.ReactionTime = &ReactionTime{Hundredths: rr.ReactionTime.Hundredths}
	}
	if rr.Lane != "" {
		p.Lane, _ = strconv.Atoi(rr.Lane)
	}
	if rr.Points != "" {
		p.Points, _ = strconv.Atoi(rr.Points)
	}

	// Date comes from the heat's session (via heatID → session date index).
	if rr.HeatID != "" {
		p.Date = idx.dateByHeatID[rr.HeatID]
	} else if evt != nil {
		p.Date = idx.dateByEventID[evt.EventID]
	}

	for _, rs := range rr.Splits {
		p.Splits = append(p.Splits, Split{
			Distance: rs.Distance,
			Time:     SwimTime(rs.SwimTime),
		})
	}
	for _, rp := range rr.RelayPositions {
		p.RelayPositions = append(p.RelayPositions, RelayPosition{
			Number:    rp.Number,
			AthleteID: rp.AthleteID,
			Status:    RelayPositionStatus(rp.Status),
		})
		if rp.ReactionTime != nil {
			p.RelayPositions[len(p.RelayPositions)-1].ReactionTime = &ReactionTime{Hundredths: rp.ReactionTime.Hundredths}
		}
	}
	return p
}

func normalizeEntry(
	re *rawxml.Entry,
	athleteID, clubID, meetID string,
	doc *Document,
	evt *Event,
	idx *meetIndex,
) *Performance {
	p := &Performance{
		doc:      doc,
		meetID:   meetID,
		AthleteID: athleteID,
		ClubID:   clubID,
		HeatID:   re.HeatID,
		Status:   PerformanceStatus(re.Status),
		Source:   SourceEntry,
	}
	if evt != nil {
		p.EventKey = evt.Key
	}
	if re.EntryTime != nil {
		p.Time = SwimTime(*re.EntryTime)
	} else {
		p.Time = SwimTime{NT: true}
	}

	// Date for entries: from MEETINFO if present, otherwise from heat/event session.
	if re.MeetInfo != nil && re.MeetInfo.Date != nil {
		p.Date = re.MeetInfo.Date.Time
	} else if re.HeatID != "" {
		p.Date = idx.dateByHeatID[re.HeatID]
	} else if evt != nil {
		p.Date = idx.dateByEventID[evt.EventID]
	}

	for _, rp := range re.RelayPositions {
		p.RelayPositions = append(p.RelayPositions, RelayPosition{
			Number:    rp.Number,
			AthleteID: rp.AthleteID,
			Status:    RelayPositionStatus(rp.Status),
		})
	}
	return p
}

func normalizeRecordList(rl *rawxml.RecordList, doc *Document) *RecordList {
	out := &RecordList{
		Name:   rl.Name,
		Type:   rl.Type,
		Gender: Gender(rl.Gender),
		Course: Course(rl.Course),
		Nation: rl.Nation,
		Region: rl.Region,
	}
	if rl.AgeGroup != nil {
		ag := normalizeAgeGroup(rl.AgeGroup, "", doc)
		out.AgeGroup = ag
	}
	for i := range rl.Records {
		rr := &rl.Records[i]
		p := normalizeRecord(rr, out.Course, doc)
		out.Records = append(out.Records, p)
	}
	return out
}

func normalizeRecord(rr *rawxml.Record, course Course, doc *Document) *Performance {
	p := &Performance{
		doc:    doc,
		meetID: "", // document-unique ids; no per-meet scoping needed
		Source: SourceRecord,
		Time:   SwimTime(rr.SwimTime),
		Status: PerformanceStatus(rr.Status),
		EventKey: EventKey{
			Distance:   rr.SwimStyle.Distance,
			Stroke:     Stroke(rr.SwimStyle.Stroke),
			RelayCount: rr.SwimStyle.RelayCount,
			Course:     course,
			Code:       rr.SwimStyle.Code,
		},
	}
	if rr.Athlete != nil {
		p.AthleteID = rr.Athlete.AthleteID
		if rr.Athlete.Club != nil {
			p.ClubID = rr.Athlete.Club.ClubID
			if p.ClubID == "" {
				p.ClubID = rr.Athlete.Club.Code
			}
		}
	}
	if rr.Relay != nil {
		p.RelayID = rr.Relay.RelayID
		if rr.Relay.Club != nil {
			p.ClubID = rr.Relay.Club.ClubID
			if p.ClubID == "" {
				p.ClubID = rr.Relay.Club.Code
			}
		}
	}
	if rr.MeetInfo != nil && rr.MeetInfo.Date != nil {
		p.Date = rr.MeetInfo.Date.Time
	}
	for _, rs := range rr.Splits {
		p.Splits = append(p.Splits, Split{
			Distance: rs.Distance,
			Time:     SwimTime(rs.SwimTime),
		})
	}
	return p
}

func normalizeAgeGroup(ra *rawxml.AgeGroup, meetID string, doc *Document) *AgeGroup {
	ag := &AgeGroup{
		doc:        doc,
		meetID:     meetID,
		AgeGroupID: ra.AgeGroupID,
		AgeMin:     ra.AgeMin,
		AgeMax:     ra.AgeMax,
		Gender:     Gender(ra.Gender),
		Name:       ra.Name,
	}
	return ag
}

func normalizeTimeStandardList(tsl *rawxml.TimeStandardList) *TimeStandardList {
	out := &TimeStandardList{
		TimeStandardListID: tsl.TimeStandardListID,
		Name:               tsl.Name,
		Gender:             Gender(tsl.Gender),
		Course:             Course(tsl.Course),
		Type:               tsl.Type,
	}
	for i := range tsl.TimeStandards {
		rts := &tsl.TimeStandards[i]
		ts := &TimeStandard{
			Time: SwimTime(rts.SwimTime),
			EventKey: EventKey{
				Distance:   rts.SwimStyle.Distance,
				Stroke:     Stroke(rts.SwimStyle.Stroke),
				RelayCount: rts.SwimStyle.RelayCount,
				Course:     out.Course,
				Code:       rts.SwimStyle.Code,
			},
		}
		out.Standards = append(out.Standards, ts)
	}
	return out
}
