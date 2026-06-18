package lenex_test

import (
	"os"
	"testing"
	"time"

	lenex "github.com/Tiim/go-lenex"
)

func openFixture(t *testing.T, name string) *lenex.Document {
	t.Helper()
	f, err := os.Open("testdata/" + name)
	if err != nil {
		t.Fatalf("open %s: %v", name, err)
	}
	t.Cleanup(func() { f.Close() })
	doc, warnings, err := lenex.Parse(f)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	if len(warnings) > 0 {
		t.Logf("%s: %d warnings: %v", name, len(warnings), warnings)
	}
	return doc
}

// ── Invitation file ────────────────────────────────────────────────────────────

func TestInvitation_ParseSucceeds(t *testing.T) {
	doc := openFixture(t, "2026-IBC-invitation.lxf")

	if len(doc.Meets) != 1 {
		t.Fatalf("len(Meets) = %d, want 1", len(doc.Meets))
	}
	m := doc.Meets[0]
	if m.Name != "Internationaler Basler Cup 2026" {
		t.Errorf("meet Name = %q", m.Name)
	}
	if m.City != "Basel" {
		t.Errorf("meet City = %q", m.City)
	}
	if m.Nation != "SUI" {
		t.Errorf("meet Nation = %q", m.Nation)
	}
	if m.Course != lenex.CourseLCM {
		t.Errorf("meet Course = %q", m.Course)
	}
}

func TestInvitation_Sessions(t *testing.T) {
	doc := openFixture(t, "2026-IBC-invitation.lxf")
	m := doc.Meets[0]

	if len(m.Sessions) != 6 {
		t.Fatalf("len(Sessions) = %d, want 6", len(m.Sessions))
	}

	wantDates := []time.Time{
		time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
	}
	for i, s := range m.Sessions {
		if s.Number != i+1 {
			t.Errorf("Session[%d].Number = %d, want %d", i, s.Number, i+1)
		}
		if !s.Date.Equal(wantDates[i]) {
			t.Errorf("Session[%d].Date = %v, want %v", i, s.Date, wantDates[i])
		}
		if s.Course != lenex.CourseLCM {
			t.Errorf("Session[%d].Course = %q, want LCM", i, s.Course)
		}
	}
}

func TestInvitation_NoResultsOrEntries(t *testing.T) {
	doc := openFixture(t, "2026-IBC-invitation.lxf")

	if n := len(doc.Results().All()); n != 0 {
		t.Errorf("Results count = %d, want 0", n)
	}
	if n := len(doc.Entries().All()); n != 0 {
		t.Errorf("Entries count = %d, want 0", n)
	}
	if n := len(doc.Meets[0].Clubs); n != 0 {
		t.Errorf("Clubs count = %d, want 0", n)
	}
}

// ── Results file ───────────────────────────────────────────────────────────────

func TestResults_ParseSucceeds(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")

	if len(doc.Meets) != 1 {
		t.Fatalf("len(Meets) = %d, want 1", len(doc.Meets))
	}
	m := doc.Meets[0]
	if m.Name != "Internationaler Basler Cup" {
		t.Errorf("meet Name = %q", m.Name)
	}
	if m.City != "Basel" {
		t.Errorf("meet City = %q", m.City)
	}
	if m.Nation != "SUI" {
		t.Errorf("meet Nation = %q", m.Nation)
	}
	if m.Course != lenex.CourseLCM {
		t.Errorf("meet Course = %q", m.Course)
	}
}

func TestResults_Sessions(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")
	m := doc.Meets[0]

	if len(m.Sessions) != 6 {
		t.Fatalf("len(Sessions) = %d, want 6", len(m.Sessions))
	}
	wantDates := []time.Time{
		time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC),
	}
	for i, s := range m.Sessions {
		if !s.Date.Equal(wantDates[i]) {
			t.Errorf("Session[%d].Date = %v, want %v", i, s.Date, wantDates[i])
		}
		if s.Course != lenex.CourseLCM {
			t.Errorf("Session[%d].Course = %q, want LCM", i, s.Course)
		}
	}
}

func TestResults_ClubCount(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")
	if n := len(doc.Meets[0].Clubs); n != 45 {
		t.Errorf("len(Clubs) = %d, want 45", n)
	}
}

func TestResults_TotalResultCount(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")
	if n := len(doc.Results().All()); n != 3416 {
		t.Errorf("len(Results) = %d, want 3416", n)
	}
}

func TestResults_ByClub(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")
	aare := doc.Results().ByClub("AARE")
	if len(aare) != 138 {
		t.Errorf("ByClub(AARE) = %d, want 138", len(aare))
	}
}

func TestResults_AthleteResolution(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")

	var found *lenex.Performance
	for _, p := range doc.Results().All() {
		if p.ResultID == "246621143" {
			found = p
			break
		}
	}
	if found == nil {
		t.Fatal("result 246621143 not found")
	}

	// Athlete
	ath := found.Athlete()
	if ath == nil {
		t.Fatal("Athlete() = nil")
	}
	if ath.FirstName != "Alessia" || ath.LastName != "Suter" {
		t.Errorf("athlete = %q %q", ath.FirstName, ath.LastName)
	}
	wantBirth := time.Date(2012, 8, 23, 0, 0, 0, 0, time.UTC)
	if !ath.BirthDate.Equal(wantBirth) {
		t.Errorf("BirthDate = %v, want %v", ath.BirthDate, wantBirth)
	}
	if ath.Nation != "SUI" {
		t.Errorf("Nation = %q, want SUI", ath.Nation)
	}

	// Club via athlete
	club := ath.Club()
	if club == nil {
		t.Fatal("ath.Club() = nil")
	}
	if club.Code != "AARE" {
		t.Errorf("club.Code = %q, want AARE", club.Code)
	}
	if club.Name != "Schwimmclub Aarefisch" {
		t.Errorf("club.Name = %q", club.Name)
	}
	if club.Nation != "SUI" {
		t.Errorf("club.Nation = %q", club.Nation)
	}

	// Club via performance
	pClub := found.Club()
	if pClub == nil {
		t.Fatal("p.Club() = nil")
	}
	if pClub.Code != "AARE" {
		t.Errorf("p.Club().Code = %q", pClub.Code)
	}
}

func TestResults_PerformanceValues(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")

	var found *lenex.Performance
	for _, p := range doc.Results().All() {
		if p.ResultID == "246621143" {
			found = p
			break
		}
	}
	if found == nil {
		t.Fatal("result 246621143 not found")
	}

	// Swim time: 1:08.99
	wantTime := 1*time.Minute + 8*time.Second + 990*time.Millisecond
	if found.Time.NT || found.Time.Duration != wantTime {
		t.Errorf("Time = %v (NT=%v), want 1:08.99", found.Time.Duration, found.Time.NT)
	}
	if found.Status != lenex.StatusOK {
		t.Errorf("Status = %q, want empty", found.Status)
	}

	// EventKey: 100m FREE LCM
	if found.EventKey.Distance != 100 {
		t.Errorf("EventKey.Distance = %d, want 100", found.EventKey.Distance)
	}
	if found.EventKey.Stroke != lenex.StrokeFree {
		t.Errorf("EventKey.Stroke = %q, want FREE", found.EventKey.Stroke)
	}
	if found.EventKey.Course != lenex.CourseLCM {
		t.Errorf("EventKey.Course = %q, want LCM", found.EventKey.Course)
	}
	if found.EventKey.RelayCount != 1 {
		t.Errorf("EventKey.RelayCount = %d, want 1", found.EventKey.RelayCount)
	}

	// Date from session
	wantDate := time.Date(2026, 5, 29, 0, 0, 0, 0, time.UTC)
	if !found.Date.Equal(wantDate) {
		t.Errorf("Date = %v, want %v", found.Date, wantDate)
	}

	// Lane
	if found.Lane != 2 {
		t.Errorf("Lane = %d, want 2", found.Lane)
	}

	// Reaction time: +0.78s = 78 hundredths
	if found.ReactionTime == nil {
		t.Fatal("ReactionTime = nil")
	}
	if found.ReactionTime.Hundredths != 78 {
		t.Errorf("ReactionTime.Hundredths = %d, want 78", found.ReactionTime.Hundredths)
	}

	// Splits: 50m at 0:32.84
	if len(found.Splits) != 1 {
		t.Fatalf("len(Splits) = %d, want 1", len(found.Splits))
	}
	if found.Splits[0].Distance != 50 {
		t.Errorf("Splits[0].Distance = %d, want 50", found.Splits[0].Distance)
	}
	wantSplit := 32*time.Second + 840*time.Millisecond
	if found.Splits[0].Time.Duration != wantSplit {
		t.Errorf("Splits[0].Time = %v, want 0:32.84", found.Splits[0].Time)
	}
}

func TestResults_PerformanceWithMultipleSplits(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")

	// Result 246621145: Alessia Suter 200m FREE, splits at 50/100/150
	var found *lenex.Performance
	for _, p := range doc.Results().All() {
		if p.ResultID == "246621145" {
			found = p
			break
		}
	}
	if found == nil {
		t.Fatal("result 246621145 not found")
	}
	if len(found.Splits) != 3 {
		t.Fatalf("len(Splits) = %d, want 3", len(found.Splits))
	}
	wantDistances := []int{50, 100, 150}
	for i, s := range found.Splits {
		if s.Distance != wantDistances[i] {
			t.Errorf("Splits[%d].Distance = %d, want %d", i, s.Distance, wantDistances[i])
		}
	}
	// Total time: 2:29.16
	wantTime := 2*time.Minute + 29*time.Second + 160*time.Millisecond
	if found.Time.Duration != wantTime {
		t.Errorf("Time = %v, want 2:29.16", found.Time)
	}
}

func TestResults_ByAthleteID(t *testing.T) {
	doc := openFixture(t, "2026-IBC-results.lxf")

	// Alessia Suter has 10 results (verified from exploration)
	var athID string
	for _, p := range doc.Results().All() {
		if p.ResultID == "246621143" {
			athID = p.AthleteID
			break
		}
	}
	if athID == "" {
		t.Fatal("could not find athlete ID from result 246621143")
	}

	byAth := doc.Results().ByAthlete(athID)
	if len(byAth) != 10 {
		t.Errorf("ByAthlete(%s) = %d, want 10", athID, len(byAth))
	}
}
