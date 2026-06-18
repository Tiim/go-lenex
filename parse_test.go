package lenex

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
	"time"
)

const integrationLenex = `<?xml version="1.0" encoding="UTF-8"?>
<LENEX version="3.0">
  <CONSTRUCTOR name="TestApp" version="1.0">
    <CONTACT email="test@example.com" name="Test Org"/>
  </CONSTRUCTOR>
  <MEETS>
    <MEET name="Test Meet" city="Berlin" nation="GER" course="LCM">
      <SESSIONS>
        <SESSION number="1" date="2024-06-01" course="LCM">
          <EVENTS>
            <EVENT eventid="1" number="1" gender="M">
              <SWIMSTYLE distance="100" stroke="FREE" relaycount="1"/>
              <AGEGROUPS>
                <AGEGROUP agegroupid="1" agemin="-1" agemax="-1">
                  <RANKINGS>
                    <RANKING place="1" resultid="101"/>
                  </RANKINGS>
                </AGEGROUP>
              </AGEGROUPS>
              <HEATS>
                <HEAT heatid="H1" number="1"/>
              </HEATS>
            </EVENT>
          </EVENTS>
        </SESSION>
      </SESSIONS>
      <CLUBS>
        <CLUB clubid="C1" code="SCB" name="Swim Club Berlin" nation="GER">
          <OFFICIALS>
            <OFFICIAL officialid="O1" firstname="Hans" lastname="Richter" gender="M"/>
          </OFFICIALS>
          <ATHLETES>
            <ATHLETE athleteid="A1" firstname="Max" lastname="Mustermann" gender="M" birthdate="2000-01-01">
              <ENTRIES>
                <ENTRY eventid="1" heatid="H1" lane="4" entrytime="00:00:51.00"/>
              </ENTRIES>
              <RESULTS>
                <RESULT resultid="101" eventid="1" heatid="H1" lane="4" swimtime="00:00:50.12" reactiontime="+18">
                  <SPLITS>
                    <SPLIT distance="50" swimtime="00:00:24.00"/>
                    <SPLIT distance="100" swimtime="00:00:50.12"/>
                  </SPLITS>
                </RESULT>
              </RESULTS>
            </ATHLETE>
          </ATHLETES>
        </CLUB>
      </CLUBS>
    </MEET>
  </MEETS>
</LENEX>`

func TestParseIntegration(t *testing.T) {
	doc, warnings, err := Parse(strings.NewReader(integrationLenex))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(warnings) != 0 {
		t.Logf("warnings: %v", warnings)
	}

	// Constructor
	if doc.Constructor.Name != "TestApp" {
		t.Errorf("Constructor.Name = %q", doc.Constructor.Name)
	}
	if doc.Constructor.Contact.Email != "test@example.com" {
		t.Errorf("Constructor.Contact.Email = %q", doc.Constructor.Contact.Email)
	}

	// Meets
	if len(doc.Meets) != 1 {
		t.Fatalf("len(Meets) = %d", len(doc.Meets))
	}
	meet := doc.Meets[0]
	if meet.Name != "Test Meet" || meet.City != "Berlin" {
		t.Errorf("meet = %+v", meet)
	}
	if len(meet.Clubs) != 1 {
		t.Fatalf("len(Clubs) = %d", len(meet.Clubs))
	}
	club := meet.Clubs[0]
	if club.Code != "SCB" {
		t.Errorf("club.Code = %q", club.Code)
	}

	// Index-based reference resolution
	idx := doc.meetsByID[meet.MeetID]
	if idx == nil {
		t.Fatal("meetIndex not found")
	}

	ath := idx.athletesByID["A1"]
	if ath == nil {
		t.Fatal("athlete A1 not found in index")
	}
	if ath.FirstName != "Max" || ath.LastName != "Mustermann" {
		t.Errorf("athlete = %+v", ath)
	}
	if !ath.BirthDate.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("birthdate = %v", ath.BirthDate)
	}

	resolvedClub := ath.Club()
	if resolvedClub == nil {
		t.Fatal("ath.Club() = nil")
	}
	if resolvedClub.Code != "SCB" {
		t.Errorf("ath.Club().Code = %q", resolvedClub.Code)
	}

	// Result
	result := idx.resultsByID["101"]
	if result == nil {
		t.Fatal("result 101 not found")
	}
	want := 50*time.Second + 120*time.Millisecond
	if result.Time.NT || result.Time.Duration != want {
		t.Errorf("result.Time = %+v", result.Time)
	}
	if result.ReactionTime == nil || result.ReactionTime.Hundredths != 18 {
		t.Errorf("result.ReactionTime = %v", result.ReactionTime)
	}
	if len(result.Splits) != 2 {
		t.Errorf("len(Splits) = %d, want 2", len(result.Splits))
	}
	if result.Splits[0].Distance != 50 {
		t.Errorf("split[0].Distance = %d", result.Splits[0].Distance)
	}

	// Date resolution via heatID → session date
	sessionDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if !result.Date.Equal(sessionDate) {
		t.Errorf("result.Date = %v, want %v", result.Date, sessionDate)
	}

	// Ranking resolution
	agIdx := idx.eventsByID["1"]
	if agIdx == nil {
		t.Fatal("event 1 not found")
	}

	// EventKey
	if result.EventKey.Distance != 100 || result.EventKey.Stroke != StrokeFree {
		t.Errorf("EventKey = %+v", result.EventKey)
	}
	if result.EventKey.Course != CourseLCM {
		t.Errorf("EventKey.Course = %q", result.EventKey.Course)
	}

	// Query API
	rs := doc.Results()
	byClub := rs.ByClub("SCB")
	if len(byClub) != 1 {
		t.Errorf("ByClub(SCB) = %d results, want 1", len(byClub))
	}
	byAthlete := rs.ByAthlete("A1")
	if len(byAthlete) != 1 {
		t.Errorf("ByAthlete(A1) = %d results, want 1", len(byAthlete))
	}

	// Ranking → Result lazy resolution
	evt := idx.eventsByID["1"]
	if evt == nil {
		t.Fatal("event 1 not in index")
	}
	if len(evt.AgeGroups) != 1 {
		t.Fatalf("event.AgeGroups = %d, want 1", len(evt.AgeGroups))
	}
	ag := evt.AgeGroups[0]
	if len(ag.Rankings) != 1 {
		t.Fatalf("agegroup.Rankings = %d, want 1", len(ag.Rankings))
	}
	rank := ag.Rankings[0]
	if rank.Place != 1 {
		t.Errorf("rank.Place = %d", rank.Place)
	}
	resolvedResult := rank.Result()
	if resolvedResult == nil {
		t.Fatal("rank.Result() = nil")
	}
	if resolvedResult.ResultID != "101" {
		t.Errorf("rank.Result().ResultID = %q", resolvedResult.ResultID)
	}

	// Official
	off := idx.officialsByID["O1"]
	if off == nil {
		t.Fatal("official O1 not found")
	}
	if off.FirstName != "Hans" || off.LastName != "Richter" {
		t.Errorf("official = %+v", off)
	}
}

func TestParseZip(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("meet.lef")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte(integrationLenex))
	zw.Close()

	doc, _, err := Parse(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("Parse (zip) error: %v", err)
	}
	if len(doc.Meets) != 1 {
		t.Errorf("zip: len(Meets) = %d", len(doc.Meets))
	}
}
