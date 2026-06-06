---
created: "2026-06-06T17:36:19Z"
dependencies: []
id: TSKR-fe9bada6
links: []
mode: hitl
priority: 1
status: open
tags: []
title: replace random hex IDs with sequential integer IDs
type: feature
updated: "2026-06-06T17:36:19Z"
---

## Design decisions

- **ID format**: `PREFIX-<n>` with no zero-padding (e.g. `TSKR-1`, `TSKR-42`)
- **Counter storage**: `next_id` field in `.tickets/config.json`, starts at 1, incremented on every `add`
- **Resolution**: bare number (`42`) expands to `PREFIX-42`; exact match only — no prefix-matching
- **New projects**: fresh start; no migration command needed
- **This repo**: one-off migration script to rename existing hex-ID tickets and rewrite cross-references