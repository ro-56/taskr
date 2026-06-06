---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-bd0e1cb2
id: TSKR-9bd2a097
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Ticket creation (taskr add)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement `taskr add` — creates a new ticket file and prints the generated ID.

ID format: `<PREFIX>-<8-char lowercase hex>` where the hex is randomly generated.

Defaults: `type=task`, `priority=2`, `mode=hitl`. `created` and `updated` timestamps set to current UTC time.

## Acceptance criteria

- [x] `taskr add "<title>"` creates a ticket with defaults and prints ID
- [x] `--type`, `--priority`, `--mode`, `--tags` flags override defaults
- [x] Generated ID is `PREFIX-xxxxxxxx` where x is random hex
- [x] Ticket file written to `.tickets/<id>.md` with correct frontmatter
- [x] `created` and `updated` set to current UTC ISO timestamp
- [x] `dependencies` and `links` initialized as empty lists
- [x] Exits with error if run outside an initialized `.tickets/` directory