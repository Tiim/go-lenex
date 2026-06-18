package rawxml

import (
	"encoding/xml"
	"testing"
	"time"
)

type swimTimeHolder struct {
	XMLName xml.Name `xml:"holder"`
	Time    SwimTime `xml:"time,attr"`
}

func TestSwimTimeUnmarshal(t *testing.T) {
	cases := []struct {
		in       string
		wantNT   bool
		wantDur  time.Duration
		wantErr  bool
	}{
		{in: "NT", wantNT: true},
		{in: "00:01:23.45", wantDur: 1*time.Minute + 23*time.Second + 450*time.Millisecond},
		{in: "00:00:00.00", wantDur: 0},
		{in: "01:59:59.99", wantDur: 1*time.Hour + 59*time.Minute + 59*time.Second + 990*time.Millisecond},
		{in: "99:59:59.99", wantDur: 99*time.Hour + 59*time.Minute + 59*time.Second + 990*time.Millisecond},
		{in: "bogus", wantErr: true},
		{in: "1:23.45", wantErr: true}, // wrong format, must be HH:MM:SS.ss
	}

	for _, c := range cases {
		data := `<holder time="` + c.in + `"></holder>`
		var h swimTimeHolder
		err := xml.Unmarshal([]byte(data), &h)
		if c.wantErr {
			if err == nil {
				t.Errorf("input %q: expected error, got none", c.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", c.in, err)
		}
		if h.Time.NT != c.wantNT {
			t.Errorf("input %q: NT = %v, want %v", c.in, h.Time.NT, c.wantNT)
		}
		if !c.wantNT && h.Time.Duration != c.wantDur {
			t.Errorf("input %q: Duration = %v, want %v", c.in, h.Time.Duration, c.wantDur)
		}
	}
}

func TestSwimTimeMarshal(t *testing.T) {
	cases := []struct {
		in   SwimTime
		want string
	}{
		{in: SwimTime{NT: true}, want: "NT"},
		{in: SwimTime{Duration: 1*time.Minute + 23*time.Second + 450*time.Millisecond}, want: "00:01:23.45"},
		{in: SwimTime{Duration: 0}, want: "00:00:00.00"},
		{in: SwimTime{Duration: 99*time.Hour + 59*time.Minute + 59*time.Second + 990*time.Millisecond}, want: "99:59:59.99"},
	}

	for _, c := range cases {
		h := swimTimeHolder{Time: c.in}
		out, err := xml.Marshal(h)
		if err != nil {
			t.Fatalf("marshal %+v: unexpected error: %v", c.in, err)
		}
		want := `<holder time="` + c.want + `"></holder>`
		if string(out) != want {
			t.Errorf("marshal %+v = %q, want %q", c.in, out, want)
		}
	}
}
