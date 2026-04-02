#!/usr/bin/env bash
# build-family.sh — Build xpkg packages for provider family sub-providers.
#
# Usage:
#   ./scripts/build-family.sh <subpackage> <version> <output-dir> <crossplane-cli> <build-registry> <arch> <platform>
#
# Each sub-provider gets:
#   1. Filtered CRDs (only its API group)
#   2. Rendered crossplane.yaml (with dependsOn for non-config/monolith)
#   3. xpkg built with embedded runtime image
#
# CRD naming convention:
#   <apigroup>.upjet-cloudflare.upbound.io_<kind>.yaml   (cluster-scoped)
#   <apigroup>.upjet-cloudflare.m.upbound.io_<kind>.yaml  (namespaced)
#
# The sub-provider name matches the API group prefix (e.g., "dns" matches dns.upjet-cloudflare.*)

set -euo pipefail

SUBPACKAGE="${1:?Usage: build-family.sh <subpackage> <version> <output-dir> <crossplane-cli> <build-registry> <arch> <platform>}"
VERSION="${2:?}"
OUTPUT_DIR="${3:?}"
CROSSPLANE_CLI="${4:?}"
BUILD_REGISTRY="${5:?}"
ARCH="${6:?}"
PLATFORM="${7:?}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
PACKAGE_DIR="${ROOT_DIR}/package"
CRD_DIR="${PACKAGE_DIR}/crds"

# Map subpackage name to xpkg repo name
map_repo_name() {
    case "$1" in
    config) echo "provider-family-cloudflare" ;;
    monolith) echo "provider-cloudflare" ;;
    *) echo "provider-cloudflare-$1" ;;
    esac
}

REPO_NAME=$(map_repo_name "${SUBPACKAGE}")
WORK_DIR=$(mktemp -d)
trap "rm -rf ${WORK_DIR}" EXIT

mkdir -p "${WORK_DIR}/crds"

# Filter CRDs for this subpackage
if [ "${SUBPACKAGE}" = "monolith" ]; then
    # Monolith gets all CRDs
    cp "${CRD_DIR}"/*.yaml "${WORK_DIR}/crds/"
elif [ "${SUBPACKAGE}" = "config" ]; then
    # Config/family provider gets ProviderConfig + StoreConfig CRDs
    # These are in the upjet-cloudflare.upbound.io and upjet-cloudflare.m.upbound.io groups directly
    for f in "${CRD_DIR}"/*.yaml; do
        basename=$(basename "$f")
        if echo "${basename}" | grep -qE '^upjet-cloudflare\.(m\.)?upbound\.io_'; then
            cp "$f" "${WORK_DIR}/crds/"
        fi
    done
else
    # Sub-provider gets only its API group CRDs
    # Match: <subpackage>.upjet-cloudflare.upbound.io_* and <subpackage>.upjet-cloudflare.m.upbound.io_*
    for f in "${CRD_DIR}/${SUBPACKAGE}.upjet-cloudflare."*.yaml; do
        [ -f "$f" ] && cp "$f" "${WORK_DIR}/crds/"
    done
fi

# Count CRDs
CRD_COUNT=$(ls "${WORK_DIR}/crds/"*.yaml 2>/dev/null | wc -l | tr -d ' ')
echo "  ${SUBPACKAGE}: ${CRD_COUNT} CRDs"

if [ "${CRD_COUNT}" -eq 0 ] && [ "${SUBPACKAGE}" != "config" ]; then
    echo "WARNING: No CRDs found for subpackage '${SUBPACKAGE}'" >&2
fi

# Render crossplane.yaml
cat >"${WORK_DIR}/crossplane.yaml" <<YAML
apiVersion: meta.pkg.crossplane.io/v1
kind: Provider
metadata:
  name: ${REPO_NAME}
  annotations:
    meta.crossplane.io/maintainer: "Wildbit"
    meta.crossplane.io/source: "https://github.com/wildbitca/provider-upjet-cloudflare"
    meta.crossplane.io/license: "Apache-2.0"
YAML

# Add family label to all family members (config + sub-providers, not monolith)
# Crossplane RBAC manager uses this label to share RBAC across family members.
# Without it on the config provider, sub-providers can't access ProviderConfig CRDs.
if [ "${SUBPACKAGE}" != "monolith" ]; then
    cat >>"${WORK_DIR}/crossplane.yaml" <<YAML
  labels:
    pkg.crossplane.io/provider-family: provider-family-cloudflare
YAML
fi

# Add dependsOn for sub-providers (not config, not monolith)
if [ "${SUBPACKAGE}" != "config" ] && [ "${SUBPACKAGE}" != "monolith" ]; then
    cat >>"${WORK_DIR}/crossplane.yaml" <<YAML
spec:
  dependsOn:
    - provider: ghcr.io/wildbitca/provider-family-cloudflare
      version: ">= ${VERSION}"
  capabilities:
    - SafeStart
YAML
else
    cat >>"${WORK_DIR}/crossplane.yaml" <<YAML
spec:
  capabilities:
    - SafeStart
YAML
fi

# Build xpkg
IMAGE_NAME="${BUILD_REGISTRY}/${REPO_NAME}-${ARCH}"
mkdir -p "${OUTPUT_DIR}/${PLATFORM}"

${CROSSPLANE_CLI} xpkg build \
    --embed-runtime-image "${IMAGE_NAME}" \
    --package-root "${WORK_DIR}" \
    --package-file "${OUTPUT_DIR}/${PLATFORM}/${REPO_NAME}-${VERSION}.xpkg"

echo "  Built: ${OUTPUT_DIR}/${PLATFORM}/${REPO_NAME}-${VERSION}.xpkg"
