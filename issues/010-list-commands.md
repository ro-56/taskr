---
id: "010"
title: "List commands (taskr list, taskr ls)"
type: task
mode: afk
status: done
---

## What to build

Implement `taskr list` and its alias `taskr ls` — various views over the ticket collection.

- `taskr list --tags`: prints all unique tags across all tickets (active only), sorted alphabetically
- `taskr list --count`: prints per-status counts (open, in_progress, closed)
- `taskr list --status <status>`: prints tickets matching that status, one per line with ID and title
- `taskr ls`: identical to `taskr list` (all flags apply)

With no flags, `taskr list` prints all active tickets (excludes archive).

## Acceptance criteria

- [ ] `taskr list` with no flags prints all active tickets with ID, title, status
- [ ] `taskr list --tags` prints sorted unique tags from active tickets
- [ ] `taskr list --count` prints count per status (including 0 counts)
- [ ] `taskr list --status open` prints only open tickets
- [ ] `taskr ls` behaves identically to `taskr list`
- [ ] `--status` accepts all valid statuses: `open`, `in_progress`, `closed`
- [ ] `--status closed` includes archived tickets in count/list

## Blocked by

- 002-ticket-creation.md
