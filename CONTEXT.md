# taskr — Domain Glossary

## Implementation Language

Go. Single static binary, built-in cross-compilation for Linux and Windows (`GOOS=linux`/`GOOS=windows`). Key dependencies: `cobra` for the multi-subcommand CLI, `gopkg.in/yaml.v3` for frontmatter parsing.

## Ticket

A markdown file representing a unit of work. Stored in `.tickets/<id>.md` (active) or `.tickets/archive/<id>.md` (closed).

Canonical fields: `id`, `title`, `status`, `type`, `priority`, `mode`, `created`, `updated`, `dependencies`, `links`, `tags`.

## Ticket ID

Format: `<PREFIX>-<n>` where `n` is a monotonically increasing integer with no zero-padding (e.g. `TKT-1`, `TKT-42`). PREFIX set at `taskr init`, stored in `.tickets/config.json`. The next integer to assign is stored as `next_id` in `config.json` and incremented on every `taskr add`.

Resolution: accepts either the full ID (`TKT-42`) or a bare number (`42`). Bare numbers are expanded to `PREFIX-<n>` before lookup. Exact match only — no prefix-matching.

## PREFIX

A short uppercase slug identifying the project. Set by the user at `taskr init` (e.g. `--prefix TKT`). Defaults to an uppercase slug of the current directory name. Stored in `.tickets/config.json`.

## Status

One of: `open`, `in_progress`, `closed`. Terminal status: `closed`. Workflow: `open → in_progress → closed`.

## Priority

Integer 0–3. `0` = critical (highest urgency), `3` = low. Default: `2` (medium). `taskr ready` sorts ascending by priority; ties broken by `updated` descending (most recently edited first).

## Mode

`afk` = agent-runnable without human involvement. `hitl` = requires a human in the loop. Default: `hitl`.

## Type

One of: `bug`, `feature`, `task`, `epic`, `chore`. Default: `task`.

## Dependency

A directional relationship: ticket A **depends on** ticket B means A cannot proceed until B is closed. Stored in A's `dependencies` frontmatter list. Cycle-checked on creation — cycles not allowed.

A ticket is **blocked** if it has at least one dependency with status `open` or `in_progress`.

## Link (related)

A symmetric, non-directional relationship between two tickets. Stored in both tickets' `links` frontmatter list. Created via `taskr relate`, which accepts two or more IDs and links all pairs in a single call. No blocking semantics. `taskr prune` does **not** remove related links — only `dependencies` are pruned.

## Archive

Closed tickets move from `.tickets/` to `.tickets/archive/` on `taskr close`. Active listing excludes archived tickets; partial ID resolution includes them.

## Body

The free-form markdown content below the frontmatter fence in a ticket file. Owned by the user — the tool never rewrites it except when explicitly instructed (e.g. `taskr update --body`). Frontmatter is tool-owned; body is user-owned. An empty body is valid.

## Violation

A broken invariant detected by `taskr check` in a ticket file. Each violation is reported as a single line prefixed by the ticket ID, e.g. `TKT-a3f8bc2d: dependency TKT-deadbeef not found`. Output is machine-readable so an LLM or script can act on it.

## Dangling Reference

A `dependencies` or `links` entry whose ID does not resolve to any file in `.tickets/` or `.tickets/archive/`. Detected by `taskr check`; indicates a ticket was deleted or renamed outside of the CLI.

## Graph Context

Compact dependency summary shown in `taskr show`:
```
depends on:  TKT-a3f8bc2 (open), TKT-1234567 (closed)
required by: TKT-9999abc (open)
```
First-degree only. Full tree via `taskr dep-tree`.

## Ready

A ticket is **ready** if: status is non-terminal (`open` or `in_progress`) AND it has no blocking dependencies (all deps closed). `taskr ready` lists these sorted by priority ascending, then updated descending.

## Config

`.tickets/config.json` — project-level configuration. Contains at minimum the PREFIX. `taskr init` creates this file. Re-running `taskr init` on an already-initialized project is a no-op with a message.

## Tags

Free-form string labels on a ticket. Stored as a list in the `tags` frontmatter field. Set at creation with `taskr add --tags` or updated with `taskr update --tags` (replaces the entire list). `taskr list --tags` prints all unique tags across active tickets.

## Update

`taskr update <id>` modifies one or more frontmatter fields of an existing ticket without changing its status. Only flags that are explicitly supplied are written; omitted flags leave the corresponding field unchanged. Updatable fields: `title`, `priority`, `mode`, `tags`, `body`.

## List

`taskr list` (alias `taskr ls`) prints all active tickets. Flags: `--status <status>` filters by status (use `closed` to include archived tickets); `--tags` prints all unique tags; `--count` prints per-status ticket counts.
