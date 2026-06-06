---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-9bd2a097
id: TSKR-e5faf175
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Symmetric links (taskr relate)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr relate <id_1> <id_2>` — creates a bidirectional, non-blocking "related" link between two tickets.

Both tickets' `links` frontmatter lists are updated atomically. No cycle checking needed — symmetric links carry no ordering semantics.

## Acceptance criteria

- [x] `taskr relate A B` adds B to A's `links` and A to B's `links`
- [x] Both files updated atomically
- [x] Accepts partial IDs
- [x] Relating already-related pair is a no-op (or clear message)
- [x] Relating a ticket to itself exits with error