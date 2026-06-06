---
id: "006"
title: "Symmetric links (taskr relate)"
type: task
mode: afk
status: done
---

## What to build

Implement `taskr relate <id_1> <id_2>` — creates a bidirectional, non-blocking "related" link between two tickets.

Both tickets' `links` frontmatter lists are updated atomically (both writes succeed or neither does). No cycle checking needed — symmetric links carry no ordering semantics.

Accepts partial IDs. Updates `updated` timestamp on both tickets.

## Acceptance criteria

- [ ] `taskr relate A B` adds B to A's `links` and A to B's `links`
- [ ] Both files updated atomically
- [ ] Accepts partial IDs
- [ ] Relating already-related pair is a no-op (or clear message)
- [ ] Relating a ticket to itself exits with error

## Blocked by

- 002-ticket-creation.md
