---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-9bd2a097
id: TSKR-acb4c6d4
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Workflow transitions (taskr start, taskr close)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr start <id>` and `taskr close <id>` — moves tickets through the workflow and archives on close.

`taskr start`: sets status `open → in_progress`, updates `updated` timestamp.

`taskr close`: sets status to `closed`, updates `updated` timestamp, moves file from `.tickets/<id>.md` to `.tickets/archive/<id>.md`. If `--summary <text>` provided, appends to the `## Notes` section of the ticket body before archiving.

## Acceptance criteria

- [x] `taskr start <id>` sets status to `in_progress`, updates timestamp
- [x] `taskr close <id>` sets status to `closed`, moves file to `.tickets/archive/`
- [x] `taskr close <id> --summary "text"` appends note to `## Notes` before archiving
- [x] `--summary` is optional; close works without it
- [x] Both commands accept partial IDs
- [x] Closing an already-closed ticket exits with a clear error
- [x] Starting a closed ticket exits with a clear error