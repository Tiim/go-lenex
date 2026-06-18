# go-lenex

A Go library for reading [LENEX 3.0](https://www.swimrankings.net/files/lenex/Lenex_3.0_Technical_Documentation.pdf) files — the international XML data exchange format for swimming meets, entries, startlists, results, and records.

## Installation

```sh
go get github.com/Tiim/go-lenex
```

## Quick start

```go
import "github.com/Tiim/go-lenex"

f, _ := os.Open("meet.lxf")
defer f.Close()

doc, warnings, err := lenex.Parse(f)
if err != nil {
    log.Fatal(err)
}

for _, p := range doc.Results().All() {
    fmt.Printf("%s %s: %s\n", p.Athlete().FirstName, p.Athlete().LastName, p.Time)
}
```

## File formats

`Parse` accepts both `.lef` (plain XML) and `.lxf` (zip-wrapped XML) without any configuration — it detects the format automatically from the input stream.

## Core concepts

### Document

`Parse` returns a `*Document`, which is the root of everything. It is fully built before returning and safe for concurrent reads.

```go
doc.Meets             // []*Meet
doc.RecordLists       // []*RecordList
doc.TimeStandardLists // []*TimeStandardList
doc.Constructor       // Constructor — who produced the file
```

### Performance

`Performance` is the unified type for a single swim. It represents an entry, a result, or a record, distinguished by `Performance.Source`:

```go
const (
    SourceEntry  PerformanceSource = iota // from an ENTRY element
    SourceResult                          // from a RESULT element
    SourceRecord                          // from a RECORD element
)
```

Key fields:

```go
p.Time          // SwimTime — the swim time, or NT
p.Status        // PerformanceStatus — "", DSQ, DNS, DNF, EXH, WDR, SICK, RJC
p.Date          // time.Time — derived from session date or MEETINFO date
p.EventKey      // EventKey — event identity (distance, stroke, course, relay count)
p.ReactionTime  // *ReactionTime — nil if not recorded
p.Splits        // []Split
p.RelayPositions // []RelayPosition — for relay swims

// Resolution methods (lazy, backed by indexes):
p.Athlete()     // *Athlete — nil for relay or target-time performances
p.Club()        // *Club
```

### SwimTime

```go
type SwimTime struct {
    Duration time.Duration
    NT       bool // true means "no time" — Duration is meaningless
}

t.String() // "1:23.45", "59.90", "NT"
```

### EventKey

Events are identified by a stable key. Course is part of the key, so SCM and LCM swims in the "same" event are distinct.

```go
type EventKey struct {
    Distance   int
    Stroke     Stroke    // "FREE", "BACK", "BREAST", "FLY", "MEDLEY", ...
    RelayCount int       // 1 = individual, 4 = 4×relay, etc.
    Course     Course    // "LCM", "SCM", "SCY"
    Code       string    // only compared when Stroke == StrokeUnknown
}
```

## Querying

### Results and entries

```go
doc.Results()  // ResultSet — all SourceResult performances across all meets
doc.Entries()  // ResultSet — all SourceEntry performances across all meets
```

`ResultSet` supports filtering:

```go
rs.All()                  // []*Performance
rs.ByClub("SCB")          // filter by club code
rs.ByAthlete("12345")     // filter by athlete ID
```

### Meet structure

```go
for _, meet := range doc.Meets {
    fmt.Println(meet.Name, meet.City, meet.Nation)

    for _, session := range meet.Sessions {
        fmt.Println(session.Date, session.Course)
    }

    for _, club := range meet.Clubs {
        fmt.Println(club.Code, club.Name)
    }
}
```

### Records and time standards

```go
for _, rl := range doc.RecordLists {
    fmt.Println(rl.Name, rl.Course)
    for _, rec := range rl.Records {
        fmt.Println(rec.EventKey, rec.Time)
        if ath := rec.Athlete(); ath != nil {
            fmt.Println(ath.FirstName, ath.LastName)
        }
    }
}

for _, tsl := range doc.TimeStandardLists {
    for _, ts := range tsl.Standards {
        fmt.Println(ts.EventKey, ts.Time)
    }
}
```

## Reference resolution

Domain types resolve cross-references lazily through the document. You never need to manage IDs manually:

```go
athlete.Club()     // *Club
relay.Club()       // *Club
judge.Official()   // *Official
ranking.Result()   // *Performance
```

## Enumerations

```go
// Strokes
lenex.StrokeFree, StrokeBack, StrokeBreast, StrokeFly
lenex.StrokeMedley, StrokeIMRelay
lenex.StrokeUnknown // custom/non-standard events

// Courses
lenex.CourseLCM, CourseSCM, CourseSCY

// Genders
lenex.GenderMale, GenderFemale, GenderMixed, GenderAll

// Performance status
lenex.StatusOK, StatusDSQ, StatusDNS, StatusDNF
lenex.StatusEXH, StatusWDR, StatusSICK, StatusRJC
```

## Error handling

```go
doc, warnings, err := lenex.Parse(r)
```

`err` is non-nil only for unreadable input or malformed XML. `warnings` carries non-fatal issues (currently unused — the type is reserved for future use as real-world malformed files surface).

## Limitations

- Read-only: there is no serialization path back to LENEX XML.
- Relay entries are indexed internally under a synthetic key; a dedicated relay-entries query API may be added in a future version.
- `Warning` cases are not yet populated — the slice will always be empty.
