#!/usr/bin/env bash
# Apply vendor patches for provider-upjet-cloudflare.
# Run after `go mod vendor` to fix upstream issues.
set -e

UPJET_FILES="vendor/github.com/crossplane/upjet/v2/pkg/terraform/files.go"

# Fix: skip setting empty ID in synthetic TF state.
# Without this, upjet creates a state with id="" which causes TF provider v5
# to attempt a Read with no ID, failing before Create can be attempted.
if [ -f "$UPJET_FILES" ]; then
  if grep -q 'if fp.hasTFID {' "$UPJET_FILES"; then
    sed -i.bak 's/if fp.hasTFID {/if fp.hasTFID \&\& tfID != "" {/' "$UPJET_FILES"
    rm -f "${UPJET_FILES}.bak"
    echo "  patched: $UPJET_FILES (empty tfID guard)"
  else
    echo "  skipped: $UPJET_FILES (already patched or changed upstream)"
  fi
fi
