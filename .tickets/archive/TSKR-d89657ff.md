---
created: "2026-06-06T17:23:35Z"
dependencies: []
id: TSKR-d89657ff
links: []
mode: afk
priority: 2
status: closed
tags: []
title: Add --version flag
type: task
updated: "2026-06-06T17:23:35Z"
---

## What to build

Add a `--version` flag to the root `taskr` command. Version string is injected at build time via `ldflags`. Falls back to `"dev"` when not passed.

The `-v` shorthand must be stripped — `--version` only.

## Acceptance criteria

- [x] `taskr --version` prints `taskr version <version>` and exits 0
- [x] `-v` is NOT a valid shorthand (returns unknown flag error)
- [x] Local build without ldflags prints `taskr version dev`
- [x] `go build -ldflags "-X main.version=1.2.3" && ./taskr --version` prints `taskr version 1.2.3`