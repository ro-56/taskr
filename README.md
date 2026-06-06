# taskr

A file-based issue tracker for the command line. Tickets are plain markdown files — version-controlled, greppable, and independently grabbable by agents or humans.

## How it works

`taskr init` creates a `.tickets/` directory in your project. Each ticket is a markdown file with YAML frontmatter (`TKT-a3f8bc2d.md`). Closed tickets are archived to `.tickets/archive/`. Everything lives alongside your code.

## Install

**Linux / macOS:**

```sh
curl -fsSL https://raw.githubusercontent.com/ro-56/taskr/main/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/ro-56/taskr/main/install.ps1 | iex
```

To pin a specific version, set `TASKR_VERSION` before running:

```sh
TASKR_VERSION=v1.2.3 curl -fsSL https://raw.githubusercontent.com/ro-56/taskr/main/install.sh | bash
```

Pre-built binaries for all platforms are available on the [releases page](https://github.com/ro-56/taskr/releases).

**Build from source** (requires Go):

```sh
go install github.com/ro-56/taskr@latest
```

Or:

```sh
git clone https://github.com/ro-56/taskr
cd taskr
go build -ldflags "-X main.version=1.0.0" -o taskr .
```

## Uninstall

**Linux / macOS:**

```sh
curl -fsSL https://raw.githubusercontent.com/ro-56/taskr/main/uninstall.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/ro-56/taskr/main/uninstall.ps1 | iex
```

## Quick start

```sh
taskr init                        # set up .tickets/ in current directory
taskr add "Fix login timeout"     # create a ticket, prints its ID
taskr start TKT-a3f8bc2d          # mark in_progress
taskr close TKT-a3f8bc2d --summary "shipped in abc1234"
```

## Commands

### Setup

| Command | Description |
|---------|-------------|
| `taskr --version` / `taskr -v` | Print the binary version and exit. |
| `taskr init` | Create `.tickets/` and `.tickets/config.json`. Re-running on an initialized project is a no-op. |
| `taskr init --prefix FOO` | Set a custom prefix (default: uppercased directory name). |

### Tickets

| Command | Description |
|---------|-------------|
| `taskr add "title"` | Create a ticket. Returns its ID. |
| `taskr show <id>` | Show frontmatter, body, and dependency graph. Accepts partial IDs. |
| `taskr update <id>` | Update one or more ticket fields non-interactively. |
| `taskr start <id>` | Set status to `in_progress`. |
| `taskr close <id>` | Close and archive a ticket. |

`taskr add` flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--type` | `task` | Ticket type: `bug`, `feature`, `task`, `epic`, `chore` |
| `--priority` | `2` | Priority `0`=critical … `3`=low |
| `--mode` | `hitl` | Mode: `afk` or `hitl` |
| `--tags` | _(none)_ | Comma-separated tags |
| `--body` | _(empty)_ | Initial ticket body (markdown) |

`taskr update` flags (only supplied flags are changed):

| Flag | Description |
|------|-------------|
| `--title` | New title |
| `--priority` | New priority (`0`–`3`) |
| `--mode` | New mode (`afk` or `hitl`) |
| `--tags` | Comma-separated tags (replaces existing) |
| `--body` | New body content (replaces existing; empty string clears body) |

`taskr close` flags: `--summary` (optional note appended to the ticket before archiving)

### Finding work

| Command | Description |
|---------|-------------|
| `taskr ready` | Non-terminal, unblocked tickets sorted by priority, then by most recently updated. |
| `taskr ready --mode afk` | Only agent-runnable tickets. |
| `taskr blocked` | Tickets with at least one open or in-progress dependency. |
| `taskr list` / `taskr ls` | List all active tickets with ID, status, and title. |
| `taskr list --status <status>` | Filter by status (`open`, `in_progress`, `closed`). |
| `taskr list --tags` | Print all unique tags across active tickets. |
| `taskr list --count` | Print per-status ticket counts. |

### Dependencies

| Command | Description |
|---------|-------------|
| `taskr link <dependent> <depends-on>` | Add a dependency (cycle-checked). |
| `taskr unlink <dependent> <depends-on>` | Remove a dependency. |
| `taskr prune <id>` | Clear a ticket's own dependencies and remove it from all other tickets' dependency lists. |
| `taskr dep-tree <id>` | ASCII tree of first-degree dependencies. |
| `taskr dep-tree <id> --full` | Fully recursive tree. |

### Related tickets

| Command | Description |
|---------|-------------|
| `taskr relate <id> <id> [id ...]` | Create bidirectional related links between all provided tickets. Accepts two or more IDs. |

### Maintenance

| Command | Description |
|---------|-------------|
| `taskr check` | Validate all tickets and report invariant violations. Exits non-zero if any are found. |

## Ticket schema

```yaml
---
id: TKT-a3f8bc2d
title: "Fix login timeout"
status: open           # open | in_progress | closed
type: task             # bug | feature | task | epic | chore
priority: 2            # 0 = critical, 1 = high, 2 = medium, 3 = low
mode: hitl             # afk (agent-runnable) | hitl (needs a human)
created: '2026-01-01T00:00:00Z'
updated: '2026-01-01T00:00:00Z'
dependencies: []
links: []
tags: []
---
```

**Priority:** `0` = critical (highest), `3` = low. Default: `2`.

**Mode:** `afk` tickets can be picked up by an agent without human involvement. `hitl` tickets require a human in the loop.

**Blocking:** a ticket is blocked if any of its `dependencies` has status `open` or `in_progress`.

## Partial IDs

Any command that takes an ID also accepts a prefix. `taskr show a3f8` resolves to the unique matching ticket. If multiple tickets match, taskr lists them and exits with an error.

## License

MIT
