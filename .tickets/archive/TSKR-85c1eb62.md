---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-9bd2a097
id: TSKR-85c1eb62
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Ticket inspection (taskr show)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr show <id>` — displays full ticket info including frontmatter, body, and compact graph context.

Partial ID resolution: prefix-match filenames in both directories. If multiple matches, print all matching IDs and exit with error asking user to be more specific.

## Acceptance criteria

- [x] `taskr show <full-id>` displays frontmatter + body + graph context
- [x] Partial ID resolves to full ID via filename prefix match
- [x] Searches both `.tickets/` and `.tickets/archive/`
- [x] Ambiguous partial ID prints all matches and exits non-zero
- [x] Graph context shows first-degree deps with their current status
- [x] Graph context shows tickets that list this ticket as a dependency
- [x] Shows "depends on: none" / "required by: none" when empty