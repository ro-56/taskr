---
id: "014"
title: "Add --body flag to taskr update"
type: feature
mode: afk
status: open
---

## What to build

Extend `taskr update <id>` with a `--body` flag that replaces the ticket's body (the free-form markdown below the frontmatter fence) non-interactively.

```
taskr update TKT-a3f8bc2d --body "## Description\n\nImplement the thing."
```

This is the first-class way for agents and scripts to write ticket body content without opening the raw `.md` file directly.

## Design decisions (resolved)

- **Body only** — frontmatter is not exposed. `--body` touches only the content below `---`.
- **Replace semantics** — the new value overwrites the existing body entirely. Empty string (`--body ""`) clears the body.
- **Extends `update`** — no separate `taskr edit` command. `update` is already the canonical mutation surface for ticket fields.
- **`updated` timestamp is bumped** on every `--body` call, same as other fields.

## Implementation notes

`UpdateOptions` in `internal/tickets/update.go` needs a `Body *string` field (nil = no-op, pointer = replace). The `Update()` function already receives `body` as a separate variable from `parseFrontmatter` — set it to `*opts.Body` when non-nil before calling `writeTicket`. Wire the flag in `cmd/update.go` using the same `flags.Changed("body")` pattern as the other flags.

## Acceptance criteria

- [ ] `taskr update <id> --body "..."` replaces the ticket body and bumps `updated`
- [ ] `taskr update <id> --body ""` clears the body (results in no body section in the file)
- [ ] Existing body is preserved when `--body` is not passed
- [ ] `--body` can be combined with any other update flags in a single invocation
- [ ] Partial ID resolution works the same as for other `update` flags
</content>
</invoke>