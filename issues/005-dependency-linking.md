---
id: "005"
title: "Dependency linking (taskr link, taskr unlink, taskr prune)"
type: task
mode: afk
status: done
---

## What to build

Implement directional dependency management between tickets.

`taskr link <dependent-id> <depends-on-id>`: adds `depends-on-id` to `dependent-id`'s `dependencies` list. Cycle-checked — if the link would create a cycle, exit with error and show the cycle path.

`taskr unlink <dependent-id> <depends-on-id>`: removes the dependency link.

`taskr prune <id>`: removes the ticket from all `dependencies` lists that reference it, and clears its own `dependencies` list. Useful before deleting or archiving manually.

All commands accept partial IDs and update `updated` timestamps on modified files.

## Acceptance criteria

- [ ] `taskr link A B` adds B to A's `dependencies`
- [ ] `taskr link` with a cycle-forming pair exits non-zero and prints the cycle path
- [ ] `taskr unlink A B` removes B from A's `dependencies`
- [ ] `taskr prune A` removes A from all other tickets' `dependencies` and clears A's own `dependencies`
- [ ] All commands accept partial IDs
- [ ] Linking a ticket to itself exits with error
- [ ] Linking already-linked pair is a no-op (or clear message)

## Blocked by

- 002-ticket-creation.md
