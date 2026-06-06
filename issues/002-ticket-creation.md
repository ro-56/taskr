---
id: "002"
title: "Ticket creation (taskr add)"
type: task
mode: afk
status: open
---

## What to build

Implement `taskr add` — creates a new ticket file and prints the generated ID.

End-to-end: `taskr add "Fix login bug" --type bug --priority 1 --mode afk` generates an ID like `TKT-a3f8bc2d`, writes `.tickets/TKT-a3f8bc2d.md` using the canonical template, and prints the ID to stdout.

ID format: `<PREFIX>-<8-char lowercase hex>` where the hex is randomly generated.

Defaults: `type=task`, `priority=2`, `mode=hitl`. `created` and `updated` timestamps set to current UTC time.

## Acceptance criteria

- [ ] `taskr add "<title>"` creates a ticket with defaults and prints ID
- [ ] `--type`, `--priority`, `--mode`, `--tags` flags override defaults
- [ ] Generated ID is `PREFIX-xxxxxxxx` where x is random hex
- [ ] Ticket file written to `.tickets/<id>.md` with correct frontmatter
- [ ] `created` and `updated` set to current UTC ISO timestamp
- [ ] `dependencies` and `links` initialized as empty lists
- [ ] Exits with error if run outside an initialized `.tickets/` directory

## Blocked by

- 001-project-scaffold-and-init.md
