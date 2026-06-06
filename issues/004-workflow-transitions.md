---
id: "004"
title: "Workflow transitions (taskr start, taskr close)"
type: task
mode: afk
status: done
---

## What to build

Implement `taskr start <id>` and `taskr close <id>` — moves tickets through the workflow and archives on close.

`taskr start`: sets status `open → in_progress`, updates `updated` timestamp.

`taskr close`: sets status to `closed`, updates `updated` timestamp, moves file from `.tickets/<id>.md` to `.tickets/archive/<id>.md`. If `--summary <text>` provided, appends to the `## Notes` section of the ticket body before archiving.

Both commands accept partial IDs.

## Acceptance criteria

- [ ] `taskr start <id>` sets status to `in_progress`, updates timestamp
- [ ] `taskr close <id>` sets status to `closed`, moves file to `.tickets/archive/`
- [ ] `taskr close <id> --summary "text"` appends note to `## Notes` before archiving
- [ ] `--summary` is optional; close works without it
- [ ] Both commands accept partial IDs
- [ ] Closing an already-closed ticket exits with a clear error
- [ ] Starting a closed ticket exits with a clear error

## Blocked by

- 002-ticket-creation.md
