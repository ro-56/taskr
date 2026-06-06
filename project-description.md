# Summary

CLI file-based issue ticket tracker to maintain independently-grabbable issues, defined in markdown files.

# Requirements

- Project name: `taskr`
- `taskr --help`: tool documentation
- `taskr init`: set up a project (creates .tickets/)

## Create a ticket
- `taskr add`: add new issue ticket. Returns the id of the created issue ticket.
- defaults: type=task, priority=2, mode=hitl

## Inspect one ticket
- `taskr show <id>`: return information for ticket with id \<id\>
- show frontmatter, body, computed graph context
- partial ids work: 0123abcd → unique full id

## Move work through the workflow
- `taskr start <id>`: set status of issue ticket with id \<id\> to `in_progress` (open → in_progress → closed)
- `taskr close <id> --summary "shipped in 0123abcd"`: close issue ticket with id \<id\>, append note to ticket, and archive it

## Find what to pick up
- `taskr ready`: non-terminal + non-blocked, sorted by priority
- `taskr ready --mode afk`:only agent-runnable
- `taskr blocked`: tickets with at least one open dep

## Relationships
- `taskr link <dependent-id> <depends-on-id>`: link tickets, cycle-checked, not allow cycles
- `taskr unlink <dependent-id> <depends-on-id>`: unlink tickets
- `taskr prune <id>`: unlink ticket from all dependents and all dependencies
- `taskr dep-tree <id>`: ASCII tree, first degree dependencies, `--full` to expand
- `taskr link --all <id_1> <id_2> ... <id_3>`: link tickets as related, not dependencies

## Annotation
- `taskr update <id> --priority 0 --tags p0,auth --mode hitl --title "New title"`: updates ticked with id \<id\>, non-interactive frontmatter write

## Lists
- `taskr list --tags`: Returns all tags in all tickets
- `taskr ls`: alias for `taskr ls`
- `taskr list --count`: Returns the count all issue tickets by status
- `taskr list --status <status>`: Returns the count issue tickets with status \<status\>


# Schema
:statuses ["open", "in_progress", "closed"]
:types ["bug", "feature", "task", "epic", "chore"]
:modes ["afk", "hitl"] - afk = agent-runnable; hitl = needs a human.

# Template

<ticket-template>
---
id: ID1234
title: "Ticket title"
status: open
type: task
priority: 3
mode: hitl
created: '2026-01-01T00:00:00.000000000Z'
updated: '2026-01-01T00:00:00.000000000Z'
dependencies:
- ID1233
links:
- ID1235
- ID1236
tags:
- tag_01
- tag_02
---

## Description

A concise description of this issue.

## Implementation notes

...

## Acceptance criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Open questions

1. **Question one topic** Question one

2. **Question two topic** Question two

3. **Question three topic** Question three

## Notes

Added notes (only to be added by `taskr`).

</ticket-template>

## Possible future features

- `taskr update <id> --description "New desc."`: updates ticked with id \<id\>,replace ## Description in place
- `taskr update <id> --body "Plain body.`: updates ticked with id \<id\>,destructive whole-body replace
