#!/usr/bin/env bash
# Applies every patches/*.patch to vendor/ (paths relative to vendor/, -p1).
# Run by `make vendor` after `go mod vendor`; fails if a patch no longer applies.
set -euo pipefail

cd "$(dirname "$0")/.."

if [ ! -f vendor/modules.txt ]; then
  echo "vendor/ is missing: run 'make vendor'" >&2
  exit 1
fi

for p in patches/*.patch; do
  if patch -p1 -d vendor -N -s --dry-run <"$p" >/dev/null; then
    patch -p1 -d vendor -N -s <"$p"
    echo "  patched: $p"
  elif patch -p1 -d vendor -R -s --dry-run <"$p" >/dev/null; then
    echo "  already applied: $p"
  else
    echo "cannot apply $p to vendor/: rebase it onto the vendored module or drop it" >&2
    exit 1
  fi
done
