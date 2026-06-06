---
id: "014"
title: "Add --body flag to taskr update and add"
type: feature
mode: afk
status: done
---

## What to build

Extend `taskr update <id>` and `taskr add <title>` with a `--body` flag that sets or replaces the ticket's body (the free-form markdown below the frontmatter fence) non-interactively.

```
taskr update TKT-a3f8bc2d --body "## Description\n\nImplement the thing."
taskr add "My new ticket" --body "## Description\n\nDetails here."
```

This is the first-class way for agents and scripts to write ticket body content without opening the raw `.md` file directly.

## Design decisions (resolved)

- **Body only** — frontmatter is not exposed. `--body` touches only the content below `---`.
- **Replace semantics on `update`** — the new value overwrites the existing body entirely. Empty string (`--body ""`) clears the body.
- **Set semantics on `add`** — the body is written once at creation time. Empty string (or omitting the flag) results in no body section.
- **Extends `update` and `add`** — no separate `taskr edit` command. `update` is already the canonical mutation surface for ticket fields; `add` benefits from inline body creation for scripted workflows.
- **`updated` timestamp is bumped** on every `--body` call to `update`, same as other fields.

## Implementation notes

### `taskr update --body`

`UpdateOptions` in `internal/tickets/update.go` needs a `Body *string` field (nil = no-op, pointer = replace). The `Update()` function already receives `body` as a separate variable from `parseFrontmatter` — set it to `*opts.Body` when non-nil before calling `writeTicket`. Wire the flag in `cmd/update.go` using the same `flags.Changed("body")` pattern as the other flags.

### `taskr add --body`

`AddOptions` in `internal/tickets/add.go` needs a `Body string` field (empty string = no body section). In `Add()`, pass `opts.Body` to `writeTicket` instead of the current `buf`-based write. Wire the flag in `cmd/add.go` with `addCmd.Flags().StringVar(&addBody, "body", "", "initial ticket body (markdown)")`.

## Acceptance criteria

- [ ] `taskr update <id> --body "..."` replaces the ticket body and bumps `updated`
- [ ] `taskr update <id> --body ""` clears the body (results in no body section in the file)
- [ ] Existing body is preserved when `--body` is not passed to `update`
- [ ] `--body` can be combined with any other `update` flags in a single invocation
- [ ] `taskr add <title> --body "..."` creates the ticket with the given body
- [ ] `taskr add <title>` (no `--body`) creates the ticket with no body section
- [ ] Partial ID resolution works the same as for other `update` flags
