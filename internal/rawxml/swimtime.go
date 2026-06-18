package rawxml

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// SwimTime is the Lenex "st" data type: "HH:MM:SS.ss" or "NT".
type SwimTime struct {
	Duration time.Duration
	NT       bool
}

func (s *SwimTime) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "NT" {
		s.NT = true
		return nil
	}
	parts := strings.Split(attr.Value, ":")
	if len(parts) != 3 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		return fmt.Errorf("swimtime: invalid format %q", attr.Value)
	}
	var hh, mm, ss, cs int
	if _, err := fmt.Sscanf(parts[0], "%d", &hh); err != nil {
		return fmt.Errorf("swimtime: invalid hours in %q", attr.Value)
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &mm); err != nil {
		return fmt.Errorf("swimtime: invalid minutes in %q", attr.Value)
	}
	secParts := strings.Split(parts[2], ".")
	if len(secParts) != 2 || len(secParts[0]) != 2 || len(secParts[1]) != 2 {
		return fmt.Errorf("swimtime: invalid seconds in %q", attr.Value)
	}
	if _, err := fmt.Sscanf(secParts[0], "%d", &ss); err != nil {
		return fmt.Errorf("swimtime: invalid seconds in %q", attr.Value)
	}
	if _, err := fmt.Sscanf(secParts[1], "%d", &cs); err != nil {
		return fmt.Errorf("swimtime: invalid centiseconds in %q", attr.Value)
	}
	s.Duration = time.Duration(hh)*time.Hour +
		time.Duration(mm)*time.Minute +
		time.Duration(ss)*time.Second +
		time.Duration(cs)*10*time.Millisecond
	return nil
}

func (s SwimTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var val string
	if s.NT {
		val = "NT"
	} else {
		d := s.Duration
		hh := int(d.Hours())
		d -= time.Duration(hh) * time.Hour
		mm := int(d.Minutes())
		d -= time.Duration(mm) * time.Minute
		ss := int(d.Seconds())
		d -= time.Duration(ss) * time.Second
		cs := int(d.Milliseconds() / 10)
		val = fmt.Sprintf("%02d:%02d:%02d.%02d", hh, mm, ss, cs)
	}
	return xml.Attr{Name: name, Value: val}, nil
}
