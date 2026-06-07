# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & test

```sh
# Build (version is injected at link time)
go build -ldflags "-X main.version=1.0.0" -o taskr .

# Run all tests
go test ./...

# Run a single test or package
go test ./internal/tickets/... -run TestAdd
go test ./internal/tickets/... -run TestLink

# Run with verbose output
go test -v ./internal/tickets/...
```

There is no linter config; the standard `go vet ./...` is the only static check in use.

## Architecture

The repo has two layers:

**`cmd/`** — Cobra command definitions. Each file registers exactly one subcommand via a package-level `init()` that calls `rootCmd.AddCommand(...)`. Commands parse flags, resolve the working directory (`os.Getwd()`), and delegate to `internal/tickets`.

**`internal/tickets/`** — All business logic. Functions take `dir string` (the project root) as their first argument. No global state.

The **frontmatter is tool-owned**; the **body is user-owned**. Commands only rewrite the body when explicitly asked (e.g., `close --summary`, `update --body`).

### Storage layout

```
.tickets/
  config.json          # { "prefix": "TKT", "next_id": 1 }
  <PREFIX>-<n>.md      # active tickets (open, in_progress)
  archive/
    <PREFIX>-<n>.md    # closed tickets
```

Ticket IDs are `<PREFIX>-<n>`, a monotonically increasing integer with no zero-padding (e.g. `TKT-1`, `TKT-42`). The next integer to assign is stored as `next_id` in `config.json` and incremented on every `add`. ID resolution accepts either the full ID or a bare number (`42` expands to `PREFIX-42`); matching is exact — no prefix-matching.

### Dependency model

Dependencies are stored as a list of full IDs in the dependent ticket's `dependencies` frontmatter field. `link.go` runs a BFS cycle check before writing. `check.go` validates dangling references and symmetry of `links` (related, non-directional) entries. A ticket is **blocked** if any dependency has status `open` or `in_progress`; it is **ready** if non-terminal and unblocked.

### Adding a new command

1. Create `cmd/<name>.go` — define the `cobra.Command`, wire flags, call `internal/tickets`.
2. Create `internal/tickets/<name>.go` — implement the logic; accept `dir string` as first param.
3. Add tests in `internal/tickets/<name>_test.go` using `t.TempDir()` + `tickets.Init(dir, "TKT")` for isolation.

## Domain model (quick reference)

| Concept | Value |
|---|---|
| Statuses | `open`, `in_progress`, `closed` |
| Types | `bug`, `feature`, `task`, `epic`, `chore` |
| Modes | `afk` (agent-runnable), `hitl` (needs human) |
| Priority | `0`=critical … `3`=low; default `2` |
| ID format | `PREFIX-<n>` (sequential integer, no zero-padding) |
| Relationships | `dependencies` (directional, blocking); `links` (symmetric, non-blocking via `relate`) |

### Commands at a glance

| Command | Purpose |
|---|---|
| `init` | Create `.tickets/` scaffold |
| `add` | Create a ticket |
| `show` | Display ticket + dependency graph |
| `update` | Edit fields non-interactively |
| `start` / `close` | Transition status |
| `list` / `ls` | List tickets; `--status`, `--tags`, `--count` |
| `ready` / `blocked` | Work-queue queries |
| `link` / `unlink` / `prune` | Manage `dependencies` |
| `relate` | Add symmetric `links` (2+ IDs) |
| `dep-tree` | ASCII dependency tree |
| `check` | Validate invariants |

See `CONTEXT.md` for the full authoritative glossary.
