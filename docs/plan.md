# Go-Lenex Library Plan

The library should treat LENEX as an input format, not as its public model.
The public API should expose a domain-oriented swim results model with resolved references, indexes, and query helpers.

## 1. Core Architecture

Use a single root object representing a loaded document.

It owns:
- Parsed entities
- Indexes for fast lookup
- Reference resolution logic

**ID scoping is per-meet, not per-document.** The spec defines `athleteid`,
`eventid`, `heatid`, `officialid`, and `resultid` as unique only "over all
[X] of a meet" — not over a whole file. `LENEX.MEETS` can legally contain
more than one `MEET`, and two meets are free to reuse the same ids. A flat
`map[string]*Athlete` keyed only by id would silently collide in that case.
Indexes are therefore scoped per meet:

```go
type Document struct {
    Constructor       *Constructor // required root element; holds the file producer's identity/contact info
    Meets             []*Meet
    RecordLists       []*RecordList
    TimeStandardLists []*TimeStandardList

    meetsByID map[string]*meetIndex
}

type Constructor struct {
    Name          string
    Version       string
    Contact       Contact
    RegisteredBy  string // optional, e.g. timing system name
}

type Contact struct {
    Name  string
    Email string
    // ... remaining CONTACT attributes (street, city, nation, phone, fax, internet)
}

// one set of indexes per MEET element
type meetIndex struct {
    clubsByID    map[string]*Club
    athletesByID map[string]*Athlete
    relaysByID   map[string]*Relay
    eventsByID   map[string]*Event
    heatsByID    map[string]*Heat
    officialsByID map[string]*Official
    resultsByID   map[string]*Result // needed for RANKING.resultid lookups

    entriesByAthleteID map[string][]*Entry
    resultsByEvent     map[string][]*Result
    judgesBySessionID  map[string][]*Judge
}
```

Domain objects that resolve relationships (`Athlete.Club()`, etc.) need to
know which meet they belong to, since lookups go through that meet's index,
not a document-wide one. `RecordList`/`TimeStandardList` entities have no
such per-meet scoping concern — their ids (`timestandardlistid`, etc.) are
already document-unique per the spec.

Concretely, every domain type whose id is only unique per-meet
(`Athlete`, `Relay`, `Official`, `Coach`, `Judge`, `Heat`, `Event`, and
`Performance`) carries a private `meetID` field alongside `doc`, used
internally to pick the right `meetIndex` before doing an id lookup. It is
not exported — consumers never need to see or compare meet ids directly,
they just call `athlete.Club()` and the resolution method does the
scoped lookup for them. Types backed by `RecordList`/`TimeStandardList`
(e.g. a `Performance` built from a `RECORD`) leave `meetID` empty, since
those ids are already document-unique.

All relationships are resolved through this object, not stored as raw XML
structure.

`meetID` is assigned during step 5 of the parsing pipeline (§2), while
walking the raw `MEET` subtree: the normalization code processing
`MEET.CLUBS.CLUB.ATHLETES.ATHLETE` (etc.) knows which `MEET` it is currently
inside and stamps that meet's id onto every domain object it builds from
that subtree. It is never inferred or looked up after the fact — by the
time a domain object exists, its `meetID` is already fixed.

Since this library is read-only (parse + query, no write/serialization
path — see §2), `Document` and all domain objects are built once during
`Parse` and never mutated afterward. Concurrent reads from multiple
goroutines are therefore safe without any locking.

`Parse` MUST fully build the `Document` — indexes, resolved domain
objects, everything — before returning. No partially-built `Document` is
ever handed to a caller, and no background goroutine continues
initializing it after `Parse` returns. This is what makes the no-locking
claim above actually hold: there is no window where one goroutine could
observe a `Document` that another goroutine is still mutating.

Note: `Club`/`Athlete`/`Relay` mean different things depending on whether they
sit under a `MEET` subtree or a `RECORDLIST`/`TIMESTANDARDLIST` subtree.
The Document needs to track which context each parsed entity
came from, since ID scoping and field requiredness differ between them.
`Official`, `Coach`, and `Judge` follow the same lazy-resolution pattern as
`Athlete`/`Club` (`Judge.Official()` via `officialsByID`, etc.) and are
first-class domain types, not deferred to a later version.

## 2. Parsing Strategy

This library is read-only for now: parse a Lenex file into the domain model
and query it. There is no write/serialization path (domain model →
XML/Lenex output) in scope — the normalization in §5 (derived `Date`,
resolved `Course`, ids dropped in favor of lazy resolution) is one-way and
deliberately lossy, so round-tripping back to XML is not a goal.

Parsing should be followed by a normalization pipeline:
1. Detect container format and obtain raw XML bytes (see below)
2. XML decode into raw structs
3. Build indexes
4. Resolve references
5. Normalize domain objects

```go
# r is io.Reader
doc := lenex.Parse(r)
```

No manual linking or initialization is required by consumers.

### File format: `.lef` vs `.lxf`

A `.lef` file is plain XML; a `.lxf` file is the same XML wrapped in a zip
archive (the common distribution format). `Parse` must accept either
transparently:

1. Peek the input for the zip local-file-header magic (`PK\x03\x04`).
   - If present, treat it as `.lxf`: read it as a zip archive (via
     `archive/zip`, which needs an `io.ReaderAt` + size — buffer the input
     if `r` isn't already a seekable/sized source) and locate the single
     XML entry inside it.
   - If absent, treat the input as raw `.lef` XML directly.
2. Feed the resulting XML bytes into the XML-decode step.

No public distinction is made between the two at the API level — `Parse`
takes any `io.Reader` and figures out the container itself.


## 3. Reference Resolution Model

Objects store only IDs internally and resolve lazily through the document.

```go
type Athlete struct {
    doc    *Document
    meetID string // private; selects which meetIndex ClubID resolves against

    ClubID string
}

func (a *Athlete) Club() *Club
```

`Relay`, `Official`, `Coach`, `Judge`, `Heat`, and `Event` follow the same
shape — `doc` + private `meetID` + exported id fields for their references.

Example usage:

```go
club := athlete.Club()
```

Do not eager link the full graph linking. Use lazy resolution backed by indexes.

This ID-based model applies consistently to all domain types, including
`Performance` (see section 5) — it stores `AthleteID`/`ClubID`, never a
denormalized field like a club code, so there is exactly one resolution path
through `Document` for every relationship.

## 4. Index Layer

All cross-references are backed by precomputed maps for O(1) access.

Examples:
- athleteByID
- clubByID
- entriesByAthleteID
- resultsByEvent
- heatByID (needed since ENTRY/RESULT.heatid references HEAT, and HEAT carries
  round/status info relevant to filtering, e.g. excluding non-podium rounds)

This enables both traversal and query APIs.

## 5. Domain Model Layer

Expose swim-specific abstractions instead of XML structures.

Key concept: normalized "Performance". It represents one swim — either an
`ENTRY`, a `RESULT`, or a `RECORD` — for an athlete or relay, in or out of a
live meet context.

`TIMESTANDARD` is deliberately **not** modeled as a `Performance`: it has no
athlete, relay, club, or date at all — just a `SWIMSTYLE` and a `swimtime`
("a maximal/minimal time for a meet"). Forcing it into `Performance` would
mean every consumer has to remember that those fields are meaningless for
it. It gets its own type instead (see below).

A `RECORD` can also legally have no `ATHLETE`/`RELAY` — "a record standard
time"/`TARGETTIME`, used to define a qualifying time for a future record
without a current holder. Unlike `TIMESTANDARD`, it still carries a date and
location (`MEETINFO` is required for `RECORD`'s `city`/`date`/`nation`), so
it stays a `Performance` with `AthleteID == "" && RelayID == ""` rather than
moving to `TimeStandard` — consumers checking "is this an individual or
relay swim" must handle both being empty.

```go
type PerformanceStatus string

const (
    StatusOK   PerformanceStatus = ""
    StatusEXH  PerformanceStatus = "EXH"  // ENTRY or RESULT
    StatusRJC  PerformanceStatus = "RJC"  // ENTRY only: rejected entry
    StatusDSQ  PerformanceStatus = "DSQ"  // RESULT only: disqualified
    StatusDNS  PerformanceStatus = "DNS"  // RESULT only: did not start
    StatusDNF  PerformanceStatus = "DNF"  // RESULT only: did not finish
    StatusSICK PerformanceStatus = "SICK" // ENTRY or RESULT
    StatusWDR  PerformanceStatus = "WDR"  // ENTRY or RESULT
)
```

`ENTRY.status` and `RESULT.status` are not the same enumeration — `RJC` only
makes sense for an entry (it was rejected before it became a result), while
`DSQ`/`DNS`/`DNF` only make sense for a result (something observed during
the actual swim). `EXH`/`SICK`/`WDR` are shared. `Performance` uses one
unified type for API simplicity, but a `Performance` built from an `ENTRY`
should never carry `DSQ`/`DNS`/`DNF`, and one built from a `RESULT` should
never carry `RJC` — worth asserting/warning on in the normalization pipeline
rather than trusting the source file.

```go
type Split struct {
    Distance int
    Time     SwimTime
}

// Distinct from PerformanceStatus (not a subtype of it): RELAYPOSITION.status
// only ever has empty/DSQ/DNF in the spec. Giving it its own type makes that
// a compile-time guarantee instead of a comment readers have to trust.
type RelayPositionStatus string

const (
    RelayPositionOK  RelayPositionStatus = ""
    RelayPositionDSQ RelayPositionStatus = "DSQ"
    RelayPositionDNF RelayPositionStatus = "DNF"
)

type RelayPosition struct {
    Number       int // -1 means reserve swimmer
    AthleteID    string
    ReactionTime *ReactionTime
    Status       RelayPositionStatus
}

type Performance struct {
    doc    *Document
    meetID string // private; empty when built from a RECORD/TIMESTANDARD (document-unique ids already)

    AthleteID string // empty for relay performances, and for RECORD "target times"
    RelayID   string // empty for individual performances, and for RECORD "target times"
    ClubID    string
    EventKey  EventKey

    Time         SwimTime
    Status       PerformanceStatus
    ReactionTime *ReactionTime
    Splits       []Split
    RelayPositions []RelayPosition // only for relay performances

    Date time.Time
}

// TimeStandard has no swimmer, club, or date — just a target time for a style.
type TimeStandard struct {
    doc *Document

    EventKey EventKey
    Time     SwimTime
}
```

`Date` is derived, not copied verbatim from one XML field, because the
source differs by context:
- Performance built from a live-meet `RESULT` or `ENTRY`: no date lives on
  either element itself — derive it from the owning `SESSION.date` (via
  heatid → event → session, or directly from session if no heat).
- Performance built from a `RECORD`, or an `ENTRY`/`RELAYPOSITION` carrying
  `MEETINFO` (the qualification swim, not the meet being entered into): use
  `MEETINFO.date` directly.

`Status` matters because `Time` being `NT` is ambiguous on its own — `NT`
plus `Status == StatusDSQ` means "disqualified", `NT` plus `StatusDNS` means
"didn't start", and `NT` with no status on an `ENTRY` just means "no entry
time declared yet". Consumers must be able to tell these apart.

Consumers primarily work with performances, not XML entities.

### Session, Pool, AgeGroup, Ranking

A handful of structural entities are queried directly rather than folded
into `Performance`, since they describe meet structure/context rather than
a swim:

```go
type Session struct {
    doc    *Document
    meetID string

    Number int
    Date   time.Time
    Course Course
    Pool   Pool
}

type Pool struct {
    Lanes  int
    Course Course // pool's nominal course; SESSION.course can still override per-session
}

type AgeGroup struct {
    doc    *Document
    meetID string // empty when owned by a RecordList/TimeStandardList instead of a MEET's EVENT

    AgeMin    int // -1 sentinel means unbounded, per spec
    AgeMax    int
    Gender    Gender
    Rankings  []Ranking // only populated for the meet-context AGEGROUP, via AGEGROUP.RANKINGS
}

// Ranking resolves RANKING.resultid through the owning meet's resultsByID index.
type Ranking struct {
    doc    *Document
    meetID string

    Place     int
    ResultID  string
}

func (r *Ranking) Result() *Performance
```

`AgeGroup` is shared structurally between a meet `EVENT` and a
`RECORDLIST`/`TIMESTANDARDLIST`, same caveat as `Club`/`Athlete`/`Relay` in
§1 — fields like `Rankings` are only ever populated in the meet context.

## 6. Result-Oriented API

Provide direct extraction of useful data:

```go
results := doc.Results()
```

Filtering is first-class:

```go
results.ByClub("SCB")
results.ByAthlete("123")
```

## 7. Time and Event Abstractions

### Wire type layer

Every Lenex attribute type needs a Go type with custom
`encoding/xml` attribute (un)marshaling, since they're string-encoded XML
attributes, not elements. Beyond `SwimTime`/`ReactionTime` below, the spec
also requires:

```go
// Date: "YYYY-MM-DD" (e.g. SESSION.date, MEETINFO.date)
type Date struct {
    time.Time
}

// Daytime: "HH:MM", hours 0-24 (e.g. SESSION.daytime)
type Daytime struct {
    Hour, Minute int
}

// Timestamp: ISO 8601 "YYYY-MM-DDThh:mm.ss" UTC (e.g. LENEX.created, .revisiondate)
type Timestamp struct {
    time.Time
}

// Unique id: letter + number; separators ignored when comparing
// (e.g. athleteid, eventid, heatid, officialid, resultid)
type UID string

func (a UID) Equal(b UID) bool // compares with separators stripped
```

These are internal to the XML decode step (§9) — by the time domain
objects are built, ids are plain `string` and dates/times are `time.Time`,
so consumers never see `Date`/`Daytime`/`Timestamp`/`UID` directly.

Times must round-trip the spec's `"NT"` sentinel and the fixed
`HH:MM:SS.ss` format (HH 0-99). Embedding `time.Duration` directly leaves no
room for an explicit "no time" state distinct from a real zero-length swim,
so it needs an explicit flag:

```go
type SwimTime struct {
    Duration time.Duration
    NT       bool // true means the spec's literal "NT"; Duration is meaningless
}
```

Reaction time has its own distinct wire format (hundredths of a second,
signed, "+14"/"-05", "0.00" written as "0", missing means empty string) and
needs its own type so the sign and "missing vs zero" distinction survive a
round trip:

```go
type ReactionTime struct {
    Hundredths int  // signed
}
// nil *ReactionTime means the attribute was absent, not zero
```

Events should be reduced to stable keys. `SWIMSTYLE.relaycount` (1 =
individual, >1 = relay) must be part of the key or e.g. "100 Free" and
"4x100 Free" collide. `stroke` can be `UNKNOWN`, in which case identity
instead comes from `code`/`swimstyleid` — that needs to be captured too, or
all custom/UNKNOWN-stroke events collapse into one key:

```go
type EventKey struct {
    Distance   int
    Stroke     Stroke
    RelayCount int
    Course     Course

    // Only meaningful (and only compared) when Stroke == UNKNOWN.
    Code string
}
```

Course is part of the key (not a separate dimension on Performance):
SCM/LCM/SCY times for "the same" event are tracked as distinct EventKeys,
not merged for cross-course history.

`Course` is not an attribute of `SWIMSTYLE` — it has to be resolved
contextually, the same way `Performance.Date` is, since the source differs
by where the swim happened:
- A meet `RESULT`/`ENTRY`'s event was actually swum at the course given by
  `SESSION.course`, falling back to `MEET.course` if the session doesn't
  override it.
- A `RECORD`'s course comes from its owning `RECORDLIST.course`.
- A `TIMESTANDARD`'s course comes from its owning `TIMESTANDARDLIST.course`.

Two more course-shaped fields exist but are *not* `EventKey.Course` and must
not be confused with it: `ENTRY.entrycourse` (the course the entry time
itself was swum in, relevant only for cross-course seeding conversion) and
`MEETINFO.course` (the course of the qualification swim referenced by an
entry, separate again from the course of the meet being entered). Both
belong on `Performance` as auxiliary fields, if exposed at all — neither
should feed into `EventKey`.

## 8. Error Handling and Validation

Validation is part of parsing:
- Broken references
- Missing athletes/clubs
- Invalid event data


Optionally expose validation diagnostics.
Some non-fatal validation errors should return warnings instead of errors.
The client should be as fault tolerant as possible.

```go
doc, warnings, err := lenex.Parse(r)
```

No concrete warning cases are known yet — there is no real-world malformed
LENEX file on hand to model one after. `Warning` is therefore introduced
now only as a placeholder type (`[]lenex.Warning`, currently an empty
struct) to lock in the `Parse` signature above; populating it with actual
cases (which references can dangle, which fields tolerate bad data and how)
is future work, once invalid real-world files surface and the library
grows specific workarounds for them.

## 9. XML Strategy

XML structs remain internal. They are not part of the public API.

Post-parse conversion ensures:
- Clean domain model
- Stable API across LENEX versions
- No dependency of consumers on XML schema


## 10. Design Outcome

Consumers should primarily interact with:
- Document
- Performance
- TimeStandard
- SwimTime
- Query methods

```go
doc := lenex.Parse(r)

for _, p := range doc.Results().ByClub("SCB") {
    fmt.Println(p.Time)
}
```

`doc.Results()` (and the rest of the §6 query API) operates across all
meets in the document, not one meet at a time — a document with multiple
`MEET`s (rare, but legal per §1) is queried as a single combined result
set. This is a document-level convenience, not a contradiction of the
per-meet ID scoping in §1: scoping still happens internally per meet, the
query layer just flattens the output across meets for the caller.

There is no document-level "this is a results file" attribute in the
spec — `MEET.type` is for federation-specific meet categories (e.g.
masters), not file purpose — and no `HasResults()` heuristic is being
committed to yet. Distinguishing invitation/registration/startlist/results
files reliably needs real-world examples of each to check assumptions
against (e.g. whether a startlist file can carry `MEET.status ==
OFFICIAL` with no `RESULT` elements, which would break a naive
"status implies results" heuristic). This is deferred until those sample
files are available.

The goal is to make LENEX act as an ingestion format, not a data model exposed to users.
