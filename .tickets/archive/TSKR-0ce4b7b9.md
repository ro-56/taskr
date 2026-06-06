---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-528afaf0
    - TSKR-e5faf175
id: TSKR-0ce4b7b9
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Validation command (taskr check)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr check` — a read-only command that scans all tickets and reports every broken invariant as a violation.

Each violation is printed as one line prefixed by the full ticket ID. Exits non-zero if any violations found, zero if clean.

## Invariants checked

- Field values: `status`, `type`, `priority`, `mode` must be valid enum values
- Required fields: `id`, `title`, `status`, `type`, `priority`, `mode`, `created`, `updated`
- ID integrity: frontmatter `id` must match filename stem and format
- File location: closed tickets must be in archive/, others must not
- Dependencies: all dep IDs must resolve; no cycles
- Links: all link IDs must resolve; links must be symmetric

## Acceptance criteria

- [x] `taskr check` exits 0 and prints nothing when all tickets are valid
- [x] `taskr check` exits non-zero and prints one line per violation
- [x] Each violation line is prefixed with the full ticket ID
- [x] Violation messages are specific enough for an LLM to determine the correct fix
- [x] All 13 invariant categories are checked
- [x] Both `.tickets/` and `.tickets/archive/` are scanned
- [x] Running on an uninitialized directory exits with an appropriate error