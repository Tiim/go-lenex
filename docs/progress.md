# Implementation Progress

Tracks status of implementing `docs/plan.md`. Updated as work proceeds.

## Status

- [x] 1. Init go module (`github.com/Tiim/go-lenex`) and project layout
      (`internal/rawxml` for XML structs, root `lenex` package for public API)
- [ ] 2. Wire types: SwimTime, ReactionTime, Date, Daytime, Timestamp, UID
- [ ] 3. Raw XML structs for full LENEX tree
- [ ] 4. Container format detection (.lef vs .lxf)
- [ ] 5. Document/meetIndex skeleton and id-scoped indexes
- [ ] 6. Domain model: Club, Athlete, Relay, Official, Coach, Judge, Heat, Event
- [ ] 7. EventKey and Course/Date resolution
- [ ] 8. Performance and TimeStandard normalization
- [ ] 9. Session, Pool, AgeGroup, Ranking domain types
- [ ] 10. Result-oriented query API
- [ ] 11. Parse() entrypoint wiring full pipeline + Warning type

## Notes

- Using TDD throughout: tests written before implementation for each unit.
- Task list (numbers above) mirrors the TaskCreate/TaskUpdate tracker used
  during this session — see that for live status if this file lags behind.
