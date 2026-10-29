# Migration guide — Cloudflare TF provider 5.27.0 → 5.28.1

This release upgrades the upstream Cloudflare Terraform provider from `v5.27.0`
to `v5.28.1`. Edit existing managed resources **in place** (never delete and
recreate) so the `crossplane.io/external-name` annotation keeps tracking the
same Cloudflare object.

| Kind (group) | Change | Action |
|--------------|--------|--------|
| `Domain` (`registrar`) | Removed upstream, no replacement. The CRDs and the `registrar` sub-provider are gone. | Set `managementPolicies` to `["Observe"]` or remove finalizers/orphan before upgrading; manage domains outside Crossplane. |
| `ShieldDiscoveryOperation` (`api`) | `operationId` renamed to `discoveryId`; `zoneId` is now required; computed `id` removed. | Rename the field and set `zoneId`. |
| `SecurityTrustedDomains` (`email`) | `isRecent`, `isRegex`, `isSimilarity` are now required. | Set all three explicitly (`false` to keep protections enabled). |
| `TrustOrganization` (`zero`) | `mfaSshPivKeyRequirements` replaced by `mfaPivKeyRequirements`. | Rename the field; nested attributes are unchanged. |
| DLP entries (`zero`) | Computed `wordList` is now a list of strings. | Update consumers that parsed the JSON string. |

New: `Rules` (`monetization` group, `cloudflare_monetization_rules`).
Notes from upstream: Pages build image version `1` is retired; use `3`.
