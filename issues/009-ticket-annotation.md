---
id: "009"
title: "Ticket annotation (taskr update)"
type: task
mode: afk
status: open
---

## What to build

Implement `taskr update <id>` — non-interactive frontmatter updates via flags.

Any combination of flags can be passed in one invocation. Only the specified fields are changed; all other frontmatter is preserved. Updates `updated` timestamp on write.

```
taskr update <id> --priority 0 --tags p0,auth --mode hitl --title "New title"
```

`--tags` replaces the entire tags list (not appended).

## Acceptance criteria

- [ ] `--title` updates the `title` field
- [ ] `--priority` updates the `priority` field (validates 0–3)
- [ ] `--mode` updates the `mode` field (validates `afk`/`hitl`)
- [ ] `--tags` replaces tags list (comma-separated input)
- [ ] Multiple flags work in one invocation
- [ ] Unspecified fields are unchanged
- [ ] `updated` timestamp updated on write
- [ ] Accepts partial IDs

## Blocked by

- 002-ticket-creation.md
