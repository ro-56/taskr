---
id: "008"
title: "Dependency tree (taskr dep-tree)"
type: task
mode: afk
status: open
---

## What to build

Implement `taskr dep-tree <id>` — renders an ASCII tree of a ticket's dependencies.

Default: first-degree only (direct deps). `--full` expands recursively to leaf nodes. Tree is downward only (what this ticket depends on, not who depends on it).

Example output:
```
TKT-a3f8bc2d (in_progress)
├── TKT-1234567a (closed)
└── TKT-9999abcd (open)
    └── TKT-bbbbbbbb (open)
```

Cycles cannot exist (enforced by `taskr link`), so no infinite loop risk.

## Acceptance criteria

- [ ] `taskr dep-tree <id>` prints first-degree deps as ASCII tree
- [ ] `--full` recursively expands all deps to leaves
- [ ] Each node shows ticket ID and current status
- [ ] Accepts partial IDs
- [ ] Ticket with no deps prints just the root node
- [ ] Works for tickets in both `.tickets/` and `.tickets/archive/`

## Blocked by

- 005-dependency-linking.md
