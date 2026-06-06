---
id: "012"
title: "Validation command (taskr check)"
type: task
mode: afk
status: done
---

## What to build

Implement `taskr check` — a read-only command that scans all tickets in `.tickets/` and `.tickets/archive/` and reports every broken invariant as a violation.

Each violation is printed as one line to stdout, prefixed by the full ticket ID:

```
TKT-a3f8bc2d: dependency TKT-deadbeef not found in .tickets/ or .tickets/archive/
TKT-a3f8bc2d: status "wip" is not valid; must be one of: open, in_progress, closed
TKT-99999999: link TKT-aaaaaaaa is not symmetric (TKT-aaaaaaaa does not link back)
```

Exits non-zero if any violations are found, zero if clean. Output is intended to be machine-readable (LLM or script can parse and act on it).

## Invariants checked

**Field values:**
- `status` must be one of `open`, `in_progress`, `closed`
- `type` must be one of `bug`, `feature`, `task`, `epic`, `chore`
- `priority` must be an integer 0–3
- `mode` must be one of `afk`, `hitl`

**Required fields:** `id`, `title`, `status`, `type`, `priority`, `mode`, `created`, `updated` must all be present and non-empty.

**ID integrity:**
- `id` in frontmatter must match the filename stem (e.g. frontmatter `id: TKT-a3f8bc2d` → file must be `TKT-a3f8bc2d.md`)
- `id` must match the format `<PREFIX>-[0-9a-f]{8}` where PREFIX is read from `.tickets/config.json`

**File location:**
- A ticket with `status: closed` must live in `.tickets/archive/`, not `.tickets/`
- A ticket with `status` other than `closed` must live in `.tickets/`, not `.tickets/archive/`

**Dependencies:**
- Every ID in `dependencies` must resolve to a file in `.tickets/` or `.tickets/archive/`
- No dependency cycles (same BFS algorithm used by `taskr link`)

**Links:**
- Every ID in `links` must resolve to a file in `.tickets/` or `.tickets/archive/`
- Links must be symmetric: if A lists B in `links`, B must list A in `links`

## Acceptance criteria

- [ ] `taskr check` exits 0 and prints nothing when all tickets are valid
- [ ] `taskr check` exits non-zero and prints one line per violation when issues are found
- [ ] Each violation line is prefixed with the full ticket ID
- [ ] Violation messages are specific enough for an LLM to determine the correct fix without additional context
- [ ] All 13 invariant categories above are checked
- [ ] Both `.tickets/` and `.tickets/archive/` are scanned
- [ ] Running on an uninitialized directory exits with an appropriate error

## Blocked by

- 005-dependency-linking.md
- 006-symmetric-links.md
