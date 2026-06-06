#!/usr/bin/env bash
set -euo pipefail

REPO="ro-56/taskr"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS
case "$(uname -s)" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)      echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

# Detect arch
case "$(uname -m)" in
  x86_64)           ARCH="amd64" ;;
  aarch64 | arm64)  ARCH="arm64" ;;
  *)                echo "Unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac

# Determine version
if [ -z "${TASKR_VERSION:-}" ]; then
  TASKR_VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
fi

VERSION_NUM="${TASKR_VERSION#v}"
ARCHIVE="taskr_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/${TASKR_VERSION}"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "Downloading taskr ${TASKR_VERSION} (${OS}/${ARCH})..."
curl -fsSL "${BASE_URL}/${ARCHIVE}" -o "${TMP_DIR}/${ARCHIVE}"
curl -fsSL "${BASE_URL}/checksums.txt" -o "${TMP_DIR}/checksums.txt"

# Verify SHA256
EXPECTED=$(grep " ${ARCHIVE}" "${TMP_DIR}/checksums.txt" | awk '{print $1}')
if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "${TMP_DIR}/${ARCHIVE}" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "${TMP_DIR}/${ARCHIVE}" | awk '{print $1}')
else
  echo "Warning: no sha256sum or shasum found; skipping checksum verification" >&2
  ACTUAL="${EXPECTED}"
fi

if [ "${ACTUAL}" != "${EXPECTED}" ]; then
  echo "SHA256 mismatch for ${ARCHIVE}" >&2
  echo "  expected: ${EXPECTED}" >&2
  echo "  got:      ${ACTUAL}" >&2
  exit 1
fi
echo "Checksum verified."

# Extract and install
tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "${TMP_DIR}"
mkdir -p "${INSTALL_DIR}"
install -m 755 "${TMP_DIR}/taskr" "${INSTALL_DIR}/taskr"
echo "Installed: ${INSTALL_DIR}/taskr"

# PATH management — add to .bashrc always, .zshrc only if it exists
PATH_LINE="export PATH=\"\${PATH}:${INSTALL_DIR}\""
ADDED_TO=()

for RC in "${HOME}/.bashrc" "${HOME}/.zshrc"; do
  if [ -f "${RC}" ] && grep -qF "${INSTALL_DIR}" "${RC}"; then
    continue
  fi
  if [ "${RC}" = "${HOME}/.bashrc" ] || [ -f "${RC}" ]; then
    printf '\n# Added by taskr installer\n%s\n' "${PATH_LINE}" >> "${RC}"
    ADDED_TO+=("${RC}")
  fi
done

if [ "${#ADDED_TO[@]}" -gt 0 ]; then
  echo "Added ${INSTALL_DIR} to PATH in: ${ADDED_TO[*]}"
  echo "Run 'source ~/.bashrc' or restart your shell to apply."
else
  echo "${INSTALL_DIR} is already in your shell PATH configuration."
fi
