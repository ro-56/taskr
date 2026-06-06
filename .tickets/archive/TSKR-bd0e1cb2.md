---
created: "2026-06-06T17:23:35Z"
dependencies: []
id: TSKR-bd0e1cb2
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Project scaffold + taskr init
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Bootstrap the CLI project and implement `taskr init`. This is the foundation all other commands build on.

Re-running on an already-initialized project prints "already initialized" and exits cleanly.

Config file shape:
```json
{ "prefix": "TKT" }
```

PREFIX defaults to an uppercase slug of the current directory name. Can be set via `taskr init --prefix MYPREFIX`.

## Acceptance criteria

- [x] `taskr --help` prints usage docs for all commands
- [x] `taskr init` creates `.tickets/`, `.tickets/archive/`, `.tickets/config.json`
- [x] PREFIX defaults to uppercased slug of cwd name
- [x] `taskr init --prefix FOO` sets prefix to `FOO`
- [x] Re-running `taskr init` in initialized directory prints "already initialized" and exits 0
- [x] `.tickets/config.json` is valid JSON with a `prefix` field