---
id: "003"
title: "Ticket inspection (taskr show)"
type: task
mode: afk
status: done
---

## What to build

Implement `taskr show <id>` — displays full ticket info including frontmatter, body, and compact graph context.

End-to-end: `taskr show TKT-a3f` resolves the partial ID against both `.tickets/` and `.tickets/archive/`, then prints the ticket's frontmatter fields, markdown body, and a two-line graph context summary:

```
depends on:  TKT-a3f8bc2d (open), TKT-1234567a (closed)
required by: TKT-9999abcd (open)
```

Partial ID resolution: prefix-match filenames in both directories. If multiple matches, print all matching IDs and exit with error asking user to be more specific.

## Acceptance criteria

- [ ] `taskr show <full-id>` displays frontmatter + body + graph context
- [ ] Partial ID resolves to full ID via filename prefix match
- [ ] Searches both `.tickets/` and `.tickets/archive/`
- [ ] Ambiguous partial ID prints all matches and exits non-zero
- [ ] Graph context shows first-degree deps with their current status
- [ ] Graph context shows tickets that list this ticket as a dependency
- [ ] Shows "depends on: none" / "required by: none" when empty

## Blocked by

- 002-ticket-creation.md
