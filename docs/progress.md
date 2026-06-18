# Implementation Progress

Tracks status of implementing `docs/plan.md`. Updated as work proceeds.

## Status

- [x] 1. Init go module (`github.com/Tiim/go-lenex`) and project layout
      (`internal/rawxml` for XML structs, root `lenex` package for public API)
- [x] 2. Wire types: SwimTime, ReactionTime, Date, Daytime, Timestamp, UID
- [x] 3. Raw XML structs for full LENEX tree
- [x] 4. Container format detection (.lef vs .lxf)
- [x] 5. Document/meetIndex skeleton and id-scoped indexes
- [x] 6. Domain model: Club, Athlete, Relay, Official, Coach, Judge, Heat, Event
- [x] 7. EventKey and Course/Date resolution
- [x] 8. Performance and TimeStandard normalization
- [x] 9. Session, Pool, AgeGroup, Ranking domain types
- [x] 10. Result-oriented query API
- [x] 11. Parse() entrypoint wiring full pipeline + Warning type

## Notes

- Wire types in `internal/rawxml` implement `xml.UnmarshalerAttr`/`xml.MarshalerAttr`
  for Lenex's non-standard attribute formats (SwimTime, ReactionTime, Date, Daytime,
  Timestamp, UID).
- Raw XML structs in `internal/rawxml/lenex.go` mirror the full LENEX 3.0 tree;
  all types are package-internal.
- `Parse()` transparently handles both `.lef` (plain XML) and `.lxf` (zip-wrapped).
- meetIndex provides O(1) id-scoped lookup per meet; dateByHeatID/dateByEventID
  allow deriving session dates for results and entries without back-references.
- All cross-references resolve lazily (e.g., `Athlete.Club()`, `Ranking.Result()`).
- AgeGroup Rankings are normalized with `doc`/`meetID` stamps so `Ranking.Result()`
  resolves after the result index is fully built.
- `Document.Results()` / `Entries()` + `ResultSet.ByClub()` / `ByAthlete()` form
  the query API (§6 of plan).
- `Warning` is a placeholder type; concrete cases deferred until malformed files surface.

## Remaining work / known gaps

- Session-level Pool is parsed but not yet attached to Session in the domain model pool resolution.
- `Meet.AgeGroups` is declared but never populated (events own their age groups).
- Relay entries and results are indexed under `entriesByAthleteID["relay:<id>"]` as
  a stopgap; a proper `entriesByRelayID` index may be cleaner.
- Time-standard list AGEGROUP not yet normalized into a domain AgeGroup.
- No federation-specific extensions (GER/SUI) handled — treated as unknown values.
- Full validation / Warning population deferred until real malformed files surface.
