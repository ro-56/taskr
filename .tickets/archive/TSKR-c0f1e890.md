---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-9bd2a097
id: TSKR-c0f1e890
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Ticket annotation (taskr update)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr update <id>` — non-interactive frontmatter updates via flags.

Any combination of flags can be passed in one invocation. Only the specified fields are changed; all other frontmatter is preserved. Updates `updated` timestamp on write.

`--tags` replaces the entire tags list (not appended).

## Acceptance criteria

- [x] `--title` updates the `title` field
- [x] `--priority` updates the `priority` field (validates 0–3)
- [x] `--mode` updates the `mode` field (validates `afk`/`hitl`)
- [x] `--tags` replaces tags list (comma-separated input)
- [x] Multiple flags work in one invocation
- [x] Unspecified fields are unchanged
- [x] `updated` timestamp updated on write
- [x] Accepts partial IDs