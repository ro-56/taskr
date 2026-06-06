#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${HOME}/.local/bin"
BINARY="${INSTALL_DIR}/taskr"

if [ -f "${BINARY}" ]; then
  rm "${BINARY}"
  echo "Removed: ${BINARY}"
else
  echo "taskr not found at ${BINARY}"
fi

# Remove PATH lines added by the installer from .bashrc and .zshrc
for RC in "${HOME}/.bashrc" "${HOME}/.zshrc"; do
  if [ -f "${RC}" ] && grep -qF "${INSTALL_DIR}" "${RC}"; then
    # Remove the comment line and the export line together
    sed -i.bak '/# Added by taskr installer/{N;d;}' "${RC}"
    rm -f "${RC}.bak"
    echo "Removed PATH entry from ${RC}"
  fi
done
