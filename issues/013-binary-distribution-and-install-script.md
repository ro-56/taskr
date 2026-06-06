---
id: "013"
title: "Binary distribution + install script"
type: feature
mode: afk
status: open
---

## What to build

Ship pre-built binaries via GitHub Releases and provide one-liner install scripts so end-users can install taskr without having Go installed.

## Design decisions

- **Release hosting**: GitHub Releases, triggered by pushing a semver tag (`git tag vX.Y.Z && git push origin vX.Y.Z`)
- **Build pipeline**: GoReleaser + GitHub Actions (`.github/workflows/release.yml`)
- **Target platforms**: Linux amd64 + arm64, macOS amd64 + arm64, Windows amd64
- **Version injection**: GoReleaser injects the tag into the binary via `-ldflags "-X main.version={{.Version}}"` (hook already exists in `main.go`)
- **Install scripts**: two scripts — `install.sh` (Linux/macOS) and `install.ps1` (Windows), hosted at `raw.githubusercontent.com/ro-56/taskr/main/`
- **Install location**: user-local — `~/.local/bin` (Linux/macOS), `%LOCALAPPDATA%\taskr\` (Windows); no elevated privileges required
- **PATH management**: scripts automatically append to shell config (`.bashrc`/`.zshrc` on Linux/macOS; registry on Windows) and print a confirmation of what changed
- **Checksums**: scripts download GoReleaser's `checksums.txt` and verify SHA256 before installing
- **Version selection**: installs latest by default; `TASKR_VERSION=v1.2.3` env var overrides to a specific release

## Acceptance criteria

- [ ] `.goreleaser.yml` builds all 5 platform targets and produces `.tar.gz` (Unix) / `.zip` (Windows) archives
- [ ] GoReleaser config injects version via ldflags; `taskr --version` on an installed binary prints the release tag
- [ ] `.github/workflows/release.yml` triggers on `push: tags: ['v*']` and runs GoReleaser
- [ ] `install.sh` detects OS (`uname -s`) and arch (`uname -m`), downloads correct archive, verifies SHA256 against `checksums.txt`, extracts binary to `~/.local/bin`, adds to PATH in `.bashrc` and `.zshrc` if not already present, prints confirmation
- [ ] `install.sh` respects `TASKR_VERSION` env var; defaults to latest via GitHub API
- [ ] `install.ps1` mirrors the above for Windows: downloads correct `.zip`, verifies hash, extracts to `%LOCALAPPDATA%\taskr\`, adds to user PATH via registry, prints confirmation
- [ ] `install.ps1` respects `TASKR_VERSION` env var
- [ ] README updated with one-liner install commands for Linux/macOS and Windows
- [ ] Manual smoke test: fresh environment (no Go), run install script, `taskr --version` prints correct tag

## Blocked by

None — can start immediately.
