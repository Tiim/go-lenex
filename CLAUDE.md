# go-lenex

A Go library for reading and writing LENEX 3.0 files — the international XML
data exchange format for swimming (meets, entries, startlists, results,
records, time standards). The full spec is in
`docs/Lenex_3.0_Technical_Documentation.pdf`, transcribed to searchable
markdown at `docs/LENEX.md` (grep that instead of the PDF when looking up
element/attribute details). The intended library architecture is sketched in
`docs/plan.md`. This repo is currently empty of code (no `go.mod` yet), so
there are no established package conventions to follow yet — check current
repo state before assuming structure described below has been implemented.

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
