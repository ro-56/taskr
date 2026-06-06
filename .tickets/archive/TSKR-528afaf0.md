---
created: "2026-06-06T17:23:35Z"
dependencies:
    - TSKR-9bd2a097
id: TSKR-528afaf0
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Dependency linking (taskr link, taskr unlink, taskr prune)
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Implement directional dependency management between tickets.

`taskr link <dependent-id> <depends-on-id>`: adds `depends-on-id` to `dependent-id`'s `dependencies` list. Cycle-checked — if the link would create a cycle, exit with error and show the cycle path.

`taskr unlink <dependent-id> <depends-on-id>`: removes the dependency link.

`taskr prune <id>`: removes the ticket from all `dependencies` lists that reference it, and clears its own `dependencies` list.

## Acceptance criteria

- [x] `taskr link A B` adds B to A's `dependencies`
- [x] `taskr link` with a cycle-forming pair exits non-zero and prints the cycle path
- [x] `taskr unlink A B` removes B from A's `dependencies`
- [x] `taskr prune A` removes A from all other tickets' `dependencies` and clears A's own `dependencies`
- [x] All commands accept partial IDs
- [x] Linking a ticket to itself exits with error
- [x] Linking already-linked pair is a no-op (or clear message)