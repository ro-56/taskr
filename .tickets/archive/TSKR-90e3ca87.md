---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-acb4c6d4
    - TSKR-528afaf0
id: TSKR-90e3ca87
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Work queue (taskr ready, taskr blocked)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr ready` and `taskr blocked` — query commands for finding actionable work.

`taskr ready`: lists all non-terminal tickets that have no blocking dependencies. Sorted by priority ascending (0 first), ties broken by `updated` descending. `--mode afk` filters to only `mode=afk` tickets.

`taskr blocked`: lists all non-terminal tickets that have at least one dependency with status `open` or `in_progress`.

## Acceptance criteria

- [x] `taskr ready` lists non-terminal, non-blocked tickets
- [x] Results sorted by priority asc, then updated desc on tie
- [x] `taskr ready --mode afk` filters to afk-mode tickets only
- [x] `taskr blocked` lists non-terminal tickets with at least one non-closed dep
- [x] Both commands print ticket ID, title, priority, and status per line
- [x] Both return empty output (not error) when no matching tickets