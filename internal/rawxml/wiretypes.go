package rawxml

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// ReactionTime is the Lenex "rt" data type: signed hundredths of a second.
// "+14" = +14 hundredths, "-05" = -5 hundredths, "0" = zero.
// A nil *ReactionTime in a struct field means the attribute was absent.
type ReactionTime struct {
	Hundredths int
}

func (r *ReactionTime) UnmarshalXMLAttr(attr xml.Attr) error {
	v := attr.Value
	if v == "" {
		return nil
	}
	if v == "0" {
		r.Hundredths = 0
		return nil
	}
	if len(v) < 2 {
		return fmt.Errorf("reactiontime: invalid format %q", v)
	}
	sign := 1
	switch v[0] {
	case '+':
		sign = 1
	case '-':
		sign = -1
	default:
		return fmt.Errorf("reactiontime: expected leading sign in %q", v)
	}
	var n int
	if _, err := fmt.Sscanf(v[1:], "%d", &n); err != nil {
		return fmt.Errorf("reactiontime: invalid value %q: %w", v, err)
	}
	r.Hundredths = sign * n
	return nil
}

// Date is the Lenex "d" data type: "YYYY-MM-DD".
type Date struct{ time.Time }

func (d *Date) UnmarshalXMLAttr(attr xml.Attr) error {
	t, err := time.Parse("2006-01-02", attr.Value)
	if err != nil {
		return fmt.Errorf("date: invalid format %q: %w", attr.Value, err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: d.Format("2006-01-02")}, nil
}

// Daytime is the Lenex "t" data type: "HH:MM".
type Daytime struct{ Hour, Minute int }

func (d *Daytime) UnmarshalXMLAttr(attr xml.Attr) error {
	var h, m int
	if _, err := fmt.Sscanf(attr.Value, "%d:%d", &h, &m); err != nil {
		return fmt.Errorf("daytime: invalid format %q: %w", attr.Value, err)
	}
	d.Hour = h
	d.Minute = m
	return nil
}

// Timestamp is the Lenex "ts" data type: "YYYY-MM-DDThh:mm.ss" (UTC).
// Note: the spec uses a dot before seconds, not a colon.
type Timestamp struct{ time.Time }

func (ts *Timestamp) UnmarshalXMLAttr(attr xml.Attr) error {
	// The spec format is "YYYY-MM-DDThh:mm.ss" — try both dot and colon before seconds.
	for _, layout := range []string{"2006-01-02T15:04.05", "2006-01-02T15:04:05"} {
		t, err := time.ParseInLocation(layout, attr.Value, time.UTC)
		if err == nil {
			ts.Time = t
			return nil
		}
	}
	return fmt.Errorf("timestamp: invalid format %q", attr.Value)
}

// UID is the Lenex "uid" data type: letter + number, separators ignored when comparing.
type UID string

// Equal compares two UIDs ignoring separator characters (space, dash, dot).
func (a UID) Equal(b UID) bool {
	return stripSeparators(string(a)) == stripSeparators(string(b))
}

func stripSeparators(s string) string {
	var b strings.Builder
	for _, r := range s {
		if !unicode.IsSpace(r) && r != '-' && r != '.' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
