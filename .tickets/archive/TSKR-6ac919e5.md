---
created: "2026-06-06T17:23:35Z"
dependencies: []
id: TSKR-6ac919e5
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Add --body flag to taskr update and add
type: feature
updated: "2026-06-06T17:23:35Z"
---

## What to build

Extend `taskr update <id>` and `taskr add <title>` with a `--body` flag that sets or replaces the ticket's body non-interactively.

- **Body only** — `--body` touches only the content below the frontmatter fence
- **Replace semantics on `update`** — new value overwrites existing body entirely; empty string clears it
- **Set semantics on `add`** — body written once at creation time

## Acceptance criteria

- [x] `taskr update <id> --body "..."` replaces the ticket body and bumps `updated`
- [x] `taskr update <id> --body ""` clears the body
- [x] Existing body is preserved when `--body` is not passed to `update`
- [x] `--body` can be combined with any other `update` flags in a single invocation
- [x] `taskr add <title> --body "..."` creates the ticket with the given body
- [x] `taskr add <title>` (no `--body`) creates the ticket with no body section
- [x] Partial ID resolution works the same as for other `update` flags