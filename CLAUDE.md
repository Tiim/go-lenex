# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# go-lenex

A Go library for reading LENEX 3.0 files — the international XML data exchange format for swimming meets, entries, startlists, results, records, and time standards. The full spec is in `docs/Lenex_3.0_Technical_Documentation.pdf`, transcribed to searchable markdown at `docs/LENEX.md` (grep that instead of the PDF when looking up element/attribute details). The intended library architecture is detailed in `docs/plan.md`.

## Commands

```sh
# Run all tests
go test ./...

# Run a single test
go test -run TestName ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./internal/rawxml/...
```

No build step required — this is a library.

## Commit guidelines

- Use semantic commit messages: `type: short description` (e.g. `feat:`, `fix:`, `ci:`, `docs:`, `refactor:`, `test:`)
- Do not add a `Co-Authored-By` trailer

## Architecture

The library has two layers:

**`internal/rawxml/`** — Raw XML decode layer (not exported)
- `wiretypes.go` — Custom attribute types: `SwimTime`, `ReactionTime`, `Date`, `Daytime`, `Timestamp`, `UID`. Each implements `xml.UnmarshalerAttr`/`xml.MarshalerAttr` for LENEX's non-standard string formats.
- `lenex.go` — Structs mirroring the full LENEX 3.0 XML tree, used only as a decode target.
- `swimtime.go` — SwimTime parsing logic (format: `HH:MM:SS.ss`, or literal `"NT"`).

**Root `lenex` package** — Public API
- `parse.go` — `Parse(io.Reader)` entry point + full normalization pipeline. Transparently handles `.lef` (plain XML) and `.lxf` (zip-wrapped XML) inputs.
- `container.go` — Zip magic-byte detection and extraction.
- `document.go` — `Document`, `meetIndex`, `ResultSet`, and query methods (`Results()`, `Entries()`, `ByClub()`, `ByAthlete()`).
- `types.go` — All exported domain types: `Performance`, `Athlete`, `Club`, `Event`, `Heat`, `Session`, `Relay`, `Official`, `Judge`, `AgeGroup`, `Ranking`, `TimeStandard`, `RecordList`, `TimeStandardList`, enumerations, and value types.

### Key design decisions

**Per-meet ID scoping.** The spec defines athlete/event/heat/result IDs as unique only within one meet, not across a document. `Document.meetsByID` maps an internal `meetID` (the meet's slice index, stringified) to a `meetIndex` holding separate maps for each ID type. Every domain type whose IDs are meet-scoped carries private `doc *Document` and `meetID string` fields; lazy resolution methods (`Athlete.Club()`, `Ranking.Result()`, etc.) use these to pick the right index.

**Parse-then-normalize pipeline.** `Parse` does: container detection → `xml.Unmarshal` into `rawxml.*` structs → `normalize()` which builds indexes and domain objects in one pass. The returned `Document` is fully built — no background initialization, safe for concurrent reads.

**`Performance` as the unified swim type.** `ENTRY`, `RESULT`, and `RECORD` all normalize to `*Performance`. `Performance.Source` (`SourceEntry`/`SourceResult`/`SourceRecord`) distinguishes origin. `Performance.Date` is derived contextually: for meet results/entries from `SESSION.date` (via heat→event→session index); for records/entries-with-meetinfo from `MEETINFO.date`.

**`EventKey` as event identity.** `EventKey{Distance, Stroke, RelayCount, Course, Code}` uniquely identifies an event. `Course` is resolved contextually from the session/meet (not from `SWIMSTYLE`). `Code` is only compared when `Stroke == StrokeUnknown`.

**Relay entries stopgap.** Relay entries are indexed under `entriesByAthleteID["relay:<relayID>"]`; a proper `entriesByRelayID` map may replace this later.

## Known gaps (from `docs/progress.md`)

- Relay entries are indexed under `entriesByAthleteID["relay:<relayID>"]`; a proper `entriesByRelayID` map may replace this later.
- `Warning` is a placeholder type; no concrete warning cases implemented yet.

## Tree shape (for modeling Go structs)

```
LENEX
  CONSTRUCTOR (CONTACT)
  MEETS
    MEET
      SESSIONS
        SESSION (POOL)
          EVENTS
            EVENT
              AGEGROUPS -> AGEGROUP (RANKINGS -> RANKING)
              HEATS -> HEAT
              SWIMSTYLE
              TIMESTANDARDREFS -> TIMESTANDARDREF
          JUDGES -> JUDGE
      CLUBS
        CLUB
          ATHLETES -> ATHLETE
            ENTRIES -> ENTRY (MEETINFO)
            RESULTS -> RESULT (SPLITS)
          COACHES -> COACH
          RELAYS -> RELAY
            ENTRIES -> ENTRY (RELAYPOSITIONS, MEETINFO)
            RESULTS -> RESULT (RELAYPOSITIONS, SPLITS)
          OFFICIALS -> OFFICIAL
  RECORDLISTS
    RECORDLIST (AGEGROUP)
      RECORDS -> RECORD (SWIMSTYLE, ATHLETE | RELAY, MEETINFO, SPLITS)
  TIMESTANDARDLISTS
    TIMESTANDARDLIST (AGEGROUP)
      TIMESTANDARDS -> TIMESTANDARD (SWIMSTYLE)
```
