package rawxml

import (
	"encoding/xml"
	"testing"
	"time"
)

type reactionTimeHolder struct {
	RT *ReactionTime `xml:"rt,attr"`
}

func TestReactionTimeUnmarshal(t *testing.T) {
	cases := []struct {
		in         string
		wantNil    bool
		wantHundredths int
		wantErr    bool
	}{
		{in: `<h></h>`, wantNil: true},
		{in: `<h rt="0"></h>`, wantHundredths: 0},
		{in: `<h rt="+14"></h>`, wantHundredths: 14},
		{in: `<h rt="-05"></h>`, wantHundredths: -5},
		{in: `<h rt="bogus"></h>`, wantErr: true},
	}
	for _, c := range cases {
		var h reactionTimeHolder
		err := xml.Unmarshal([]byte(c.in), &h)
		if c.wantErr {
			if err == nil {
				t.Errorf("input %q: expected error, got none", c.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", c.in, err)
		}
		if c.wantNil {
			if h.RT != nil {
				t.Errorf("input %q: expected nil, got %v", c.in, h.RT)
			}
			continue
		}
		if h.RT == nil {
			t.Errorf("input %q: expected non-nil", c.in)
			continue
		}
		if h.RT.Hundredths != c.wantHundredths {
			t.Errorf("input %q: Hundredths = %d, want %d", c.in, h.RT.Hundredths, c.wantHundredths)
		}
	}
}

type dateHolder struct {
	D Date `xml:"d,attr"`
}

func TestDateUnmarshal(t *testing.T) {
	h := dateHolder{}
	err := xml.Unmarshal([]byte(`<h d="2004-03-09"></h>`), &h)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2004, 3, 9, 0, 0, 0, 0, time.UTC)
	if !h.D.Equal(want) {
		t.Errorf("got %v, want %v", h.D.Time, want)
	}
}

type daytimeHolder struct {
	DT Daytime `xml:"dt,attr"`
}

func TestDaytimeUnmarshal(t *testing.T) {
	h := daytimeHolder{}
	err := xml.Unmarshal([]byte(`<h dt="14:30"></h>`), &h)
	if err != nil {
		t.Fatal(err)
	}
	if h.DT.Hour != 14 || h.DT.Minute != 30 {
		t.Errorf("got %+v, want {Hour:14 Minute:30}", h.DT)
	}
}

type tsHolder struct {
	TS Timestamp `xml:"ts,attr"`
}

func TestTimestampUnmarshal(t *testing.T) {
	h := tsHolder{}
	err := xml.Unmarshal([]byte(`<h ts="2004-03-09T14:30.45"></h>`), &h)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2004, 3, 9, 14, 30, 45, 0, time.UTC)
	if !h.TS.Equal(want) {
		t.Errorf("got %v, want %v", h.TS.Time, want)
	}
}

func TestUIDEqual(t *testing.T) {
	cases := [][2]UID{
		{"A1", "A1"},
		{"A 1", "A1"},
		{"A-1", "A1"},
		{"A.1", "A1"},
		{"A 1-2.3", "A123"},
	}
	for _, c := range cases {
		if !c[0].Equal(c[1]) {
			t.Errorf("%q.Equal(%q) = false, want true", c[0], c[1])
		}
	}
	notEqual := [][2]UID{
		{"A1", "B1"},
		{"A1", "A2"},
	}
	for _, c := range notEqual {
		if c[0].Equal(c[1]) {
			t.Errorf("%q.Equal(%q) = true, want false", c[0], c[1])
		}
	}
}
