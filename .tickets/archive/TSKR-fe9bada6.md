---
created: "2026-06-06T17:36:19Z"
dependencies: []
id: TSKR-fe9bada6
links: []
mode: hitl
priority: 1
status: closed
tags: []
title: replace random hex IDs with sequential integer IDs
type: feature
updated: "2026-06-07T14:04:56Z"
---

## Design decisions

- **ID format**: `PREFIX-<n>` with no zero-padding (e.g. `TSKR-1`, `TSKR-42`)
- **Counter storage**: `next_id` field in `.tickets/config.json`, starts at 1, incremented on every `add`
- **Resolution**: bare number (`42`) expands to `PREFIX-42`; exact match only — no prefix-matching
- **New projects**: fresh start; no migration command needed
- **This repo**: one-off migration script to rename existing hex-ID tickets and rewrite cross-references
## Notes

Implemented sequential PREFIX-<n> IDs with next_id counter in config.json, exact-match resolution with bare-number expansion, and check.go support for both new and legacy ID formats. No migration script — existing hex-ID tickets remain reachable by full ID; new tickets start at TSKR-1. Updated CLAUDE.md docs accordingly.
