#!/usr/bin/env bash
# Script to checkout or update the DML wiki repository in .docs-cache/dml.wiki
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
WIKI_DIR="${REPO_ROOT}/.docs-cache/dml.wiki"
WIKI_URL="https://github.com/intel/device-modeling-language.wiki.git"

if [ -d "${WIKI_DIR}/.git" ]; then
    echo "DML wiki already checked out at: ${WIKI_DIR}"
    if [ "${1:-}" = "--update" ] || [ "${1:-}" = "-u" ]; then
        echo "Updating wiki repository..."
        git -C "${WIKI_DIR}" pull --ff-only
    fi
else
    echo "Cloning DML wiki repository to: ${WIKI_DIR}..."
    mkdir -p "$(dirname "${WIKI_DIR}")"
    git clone --depth 1 "${WIKI_URL}" "${WIKI_DIR}"
    echo "Successfully cloned DML wiki."
fi
