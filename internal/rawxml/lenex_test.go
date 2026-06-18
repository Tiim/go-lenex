package rawxml

import (
	"encoding/xml"
	"testing"
	"time"
)

const minimalLenex = `<?xml version="1.0" encoding="UTF-8"?>
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
                <HEAT heatid="1" number="1"/>
              </HEATS>
            </EVENT>
          </EVENTS>
        </SESSION>
      </SESSIONS>
      <CLUBS>
        <CLUB clubid="C1" code="SCB" name="Swim Club Berlin" nation="GER">
          <ATHLETES>
            <ATHLETE athleteid="A1" firstname="Max" lastname="Mustermann" gender="M" birthdate="2000-01-01">
              <RESULTS>
                <RESULT resultid="101" eventid="1" heatid="1" lane="4" swimtime="00:00:50.12" reactiontime="+18"/>
              </RESULTS>
            </ATHLETE>
          </ATHLETES>
        </CLUB>
      </CLUBS>
    </MEET>
  </MEETS>
</LENEX>`

func TestLenexUnmarshal(t *testing.T) {
	var doc LENEX
	if err := xml.Unmarshal([]byte(minimalLenex), &doc); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if doc.Version != "3.0" {
		t.Errorf("Version = %q, want %q", doc.Version, "3.0")
	}
	if doc.Constructor.Name != "TestApp" {
		t.Errorf("Constructor.Name = %q, want %q", doc.Constructor.Name, "TestApp")
	}
	if doc.Constructor.Contact.Email != "test@example.com" {
		t.Errorf("Constructor.Contact.Email = %q", doc.Constructor.Contact.Email)
	}

	if len(doc.Meets) != 1 {
		t.Fatalf("len(Meets) = %d, want 1", len(doc.Meets))
	}
	meet := doc.Meets[0]
	if meet.Name != "Test Meet" || meet.City != "Berlin" || meet.Nation != "GER" {
		t.Errorf("meet = %+v", meet)
	}

	if len(meet.Sessions) != 1 {
		t.Fatalf("len(Sessions) = %d, want 1", len(meet.Sessions))
	}
	sess := meet.Sessions[0]
	if sess.Number != 1 {
		t.Errorf("session.Number = %d, want 1", sess.Number)
	}
	wantDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if !sess.Date.Equal(wantDate) {
		t.Errorf("session.Date = %v, want %v", sess.Date, wantDate)
	}

	if len(sess.Events) != 1 {
		t.Fatalf("len(Events) = %d, want 1", len(sess.Events))
	}
	evt := sess.Events[0]
	if evt.EventID != "1" || evt.Number != 1 || evt.Gender != "M" {
		t.Errorf("event = %+v", evt)
	}
	if evt.SwimStyle.Distance != 100 || evt.SwimStyle.Stroke != "FREE" {
		t.Errorf("swimstyle = %+v", evt.SwimStyle)
	}
	if len(evt.AgeGroups) != 1 || len(evt.AgeGroups[0].Rankings) != 1 {
		t.Errorf("agegroups/rankings not parsed correctly")
	}
	if evt.AgeGroups[0].Rankings[0].Place != 1 || evt.AgeGroups[0].Rankings[0].ResultID != "101" {
		t.Errorf("ranking = %+v", evt.AgeGroups[0].Rankings[0])
	}

	if len(meet.Clubs) != 1 {
		t.Fatalf("len(Clubs) = %d, want 1", len(meet.Clubs))
	}
	club := meet.Clubs[0]
	if club.ClubID != "C1" || club.Code != "SCB" {
		t.Errorf("club = %+v", club)
	}

	if len(club.Athletes) != 1 {
		t.Fatalf("len(Athletes) = %d, want 1", len(club.Athletes))
	}
	ath := club.Athletes[0]
	if ath.AthleteID != "A1" || ath.FirstName != "Max" || ath.LastName != "Mustermann" {
		t.Errorf("athlete = %+v", ath)
	}
	if ath.BirthDate == nil || ath.BirthDate.Year() != 2000 {
		t.Errorf("birthdate = %v", ath.BirthDate)
	}

	if len(ath.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(ath.Results))
	}
	res := ath.Results[0]
	if res.ResultID != "101" || res.EventID != "1" {
		t.Errorf("result ids = %+v", res)
	}
	if res.SwimTime.NT || res.SwimTime.Duration != 50*time.Second+120*time.Millisecond {
		t.Errorf("result swimtime = %+v", res.SwimTime)
	}
	if res.ReactionTime == nil || res.ReactionTime.Hundredths != 18 {
		t.Errorf("reactiontime = %v", res.ReactionTime)
	}
}
