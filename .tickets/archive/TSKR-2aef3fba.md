---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-9bd2a097
id: TSKR-2aef3fba
links: []
mode: afk
priority: 2
status: closed
tags: []
title: List commands (taskr list, taskr ls)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr list` and its alias `taskr ls` — various views over the ticket collection.

- `taskr list --tags`: prints all unique tags across all tickets (active only), sorted alphabetically
- `taskr list --count`: prints per-status counts
- `taskr list --status <status>`: prints tickets matching that status

## Acceptance criteria

- [x] `taskr list` with no flags prints all active tickets with ID, title, status
- [x] `taskr list --tags` prints sorted unique tags from active tickets
- [x] `taskr list --count` prints count per status (including 0 counts)
- [x] `taskr list --status open` prints only open tickets
- [x] `taskr ls` behaves identically to `taskr list`
- [x] `--status` accepts all valid statuses: `open`, `in_progress`, `closed`
- [x] `--status closed` includes archived tickets in count/list