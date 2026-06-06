---
id: "001"
title: "Project scaffold + taskr init"
type: task
mode: afk
status: done
---

## What to build

Bootstrap the CLI project and implement `taskr init`. This is the foundation all other commands build on.

End-to-end: running `taskr init` in a directory creates `.tickets/`, `.tickets/archive/`, and `.tickets/config.json` with the project PREFIX. Re-running on an already-initialized project prints "already initialized" and exits cleanly. `taskr --help` prints tool documentation.

Config file shape:
```json
{ "prefix": "TKT" }
```

PREFIX defaults to an uppercase slug of the current directory name. Can be set via `taskr init --prefix MYPREFIX`.

## Acceptance criteria

- [ ] `taskr --help` prints usage docs for all commands
- [ ] `taskr init` creates `.tickets/`, `.tickets/archive/`, `.tickets/config.json`
- [ ] PREFIX defaults to uppercased slug of cwd name
- [ ] `taskr init --prefix FOO` sets prefix to `FOO`
- [ ] Re-running `taskr init` in initialized directory prints "already initialized" and exits 0
- [ ] `.tickets/config.json` is valid JSON with a `prefix` field

## Blocked by

None — can start immediately.
