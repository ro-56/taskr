---
id: "007"
title: "Work queue (taskr ready, taskr blocked)"
type: task
mode: afk
status: done
---

## What to build

Implement `taskr ready` and `taskr blocked` — query commands for finding actionable work.

`taskr ready`: lists all non-terminal tickets (`open` or `in_progress`) that have no blocking dependencies (all deps are `closed`). Sorted by priority ascending (0 first), ties broken by `updated` descending (most recently edited first). `--mode afk` filters to only `mode=afk` tickets.

`taskr blocked`: lists all non-terminal tickets that have at least one dependency with status `open` or `in_progress`.

## Acceptance criteria

- [ ] `taskr ready` lists non-terminal, non-blocked tickets
- [ ] Results sorted by priority asc, then updated desc on tie
- [ ] `taskr ready --mode afk` filters to afk-mode tickets only
- [ ] `taskr blocked` lists non-terminal tickets with at least one non-closed dep
- [ ] Both commands print ticket ID, title, priority, and status per line
- [ ] Both return empty output (not error) when no matching tickets

## Blocked by

- 004-workflow-transitions.md
- 005-dependency-linking.md
