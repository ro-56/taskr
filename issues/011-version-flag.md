---
id: "011"
title: "Add --version flag"
type: task
mode: afk
status: done
---

## What to build

Add a `--version` flag to the root `taskr` command so users can query the binary's version.

Version string is injected at build time via `ldflags`: `go build -ldflags "-X main.version=1.2.3"`. Falls back to `"dev"` when the flag is not passed (local builds).

Cobra's built-in `Version` field handles the flag wiring. The `-v` shorthand must be stripped — `--version` only.

## Acceptance criteria

- [ ] `taskr --version` prints `taskr version <version>` and exits 0
- [ ] `-v` is NOT a valid shorthand (returns unknown flag error)
- [ ] Local build without ldflags prints `taskr version dev`
- [ ] `go build -ldflags "-X main.version=1.2.3" && ./taskr --version` prints `taskr version 1.2.3`

## Blocked by

None — can start immediately.
