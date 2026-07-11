# Changelog

All notable changes to this provider family are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.12] - 2026-07-10

Exposes two managed resources that the upstream Cloudflare Terraform provider
`v5.22.0` made available. Both are folded into existing sub-providers, so no new
family package or cluster wiring is introduced. Purely additive.

### Added

- `MOQRelay` (from `cloudflare_moq_relay`) — a Media over QUIC (MoQ) relay, new in upstream `v5.22.0`. It is homed in the existing `stream` sub-provider (`xpkg.upbound.io/wildbitca/provider-cloudflare-stream`) rather than a standalone group, so no new package/wiring is required. Account-scoped; `accountId` supports cross-resource references/selectors (`accountIdRef`/`accountIdSelector`). Available in both the cluster-scoped (`stream.upjet-cloudflare.upbound.io`) and namespaced (`stream.upjet-cloudflare.m.upbound.io`) API groups.
- `OAuthClient` (from `cloudflare_oauth_client`) — an account-scoped OAuth application. The upstream resource gained full CRUD in `v5.22.0` (it was previously incomplete and therefore not exposed), so it is now included. It is homed in the existing `account` sub-provider (`xpkg.upbound.io/wildbitca/provider-cloudflare-account`); `accountId` supports references/selectors. Available in both `account.upjet-cloudflare.upbound.io` and `account.upjet-cloudflare.m.upbound.io` groups.
- Managed resource count increased from 213 to 215 per scope. No sub-providers were added or removed.

## [0.2.11] - 2026-07-10

Upgrades the upstream Cloudflare Terraform provider from `v5.21.1` to `v5.22.0`
and regenerates the provider schema, CRDs, and controllers. This is a minor
upstream release and is almost entirely additive (many resources gained new
optional fields). One managed resource has a **breaking** schema change — see
below.

### Changed

- Upgraded the upstream Cloudflare Terraform provider from `v5.21.1` to `v5.22.0`; provider schema, CRDs, and controllers were regenerated. No sub-providers were added or removed (213 managed resources per scope, unchanged).
- **BREAKING** — `SearchInstance` (`ai` group, from `cloudflare_ai_search_instance`): the `crawlOptions` block (`spec.forProvider.crawlOptions` with `depth`, `includeExternalLinks`, `includeSubdomains`, `maxAge`, `source`) was removed upstream and replaced by a new optional `customDomains` (`custom_domains`) string list. Any resource that set `crawlOptions` must migrate to `customDomains`.
- Relaxed from required to optional (non-breaking): `Transforms` (`managed` group) — `managedRequestHeaders` and `managedResponseHeaders` are no longer required.
- Additive new optional fields across many resources, notably `workers` `Script` and `worker` `Version` (e.g. `cacheOptions`, `packageDependencies`), `stream` `LiveInput` (`keysRotatedAt`), and several `zero` trust resources.

## [0.2.10] - 2026-07-09

Upgrades the upstream Cloudflare Terraform provider from `v5.18.0` to `v5.21.1`
and adds a new Secrets Store sub-provider. Despite the patch-level tag, this
release contains **breaking changes** to several existing managed resources
(driven by the upstream provider upgrade). See
[`docs/migration/cloudflare-5.21-breaking-changes.md`](docs/migration/cloudflare-5.21-breaking-changes.md)
and the migration notes below **before** upgrading sub-providers that manage
live resources.

### Added

- New `secretsstore` sub-provider (`xpkg.upbound.io/wildbitca/provider-cloudflare-secretsstore`) with two managed resources:
  - `Store` (from `cloudflare_secrets_store`) — an account-scoped secrets store.
  - `StoreSecret` (from `cloudflare_secrets_store_secret`) — a secret within a store.
  - Available in both the cluster-scoped (`secretsstore.upjet-cloudflare.upbound.io`) and namespaced (`secretsstore.upjet-cloudflare.m.upbound.io`) API groups.
  - `StoreSecret.spec.forProvider.value` is provided as a Kubernetes `Secret` reference (`valueSecretRef`); the plaintext value is never stored in the resource spec.
  - `storeId` and `accountId` support cross-resource references/selectors (`storeIdRef`/`storeIdSelector`, `accountIdRef`/`accountIdSelector`).

### Changed

- Upgraded the upstream Cloudflare Terraform provider from `v5.18.0` to `v5.21.1`; provider schema, CRDs, and controllers were regenerated. Many resources gained new optional fields (additive).
- Hardened the runtime image build (`apk add`) with retries to tolerate transient Alpine mirror timeouts in the parallel publish matrix.
- **BREAKING** — `EmailSecurityTrustedDomains` (`email` group): `pattern` is now **required** (was optional); the `body` block was removed.
- **BREAKING** — `NetworkMonitoringRule` (`magic` group): `automaticAdvertisement`, `prefixes`, and `type` are now **required** (were optional).
- **BREAKING** — `Job` and `OwnershipChallenge` (`logpush` group): `destinationConf` is now **sensitive** and must be supplied via `destinationConfSecretRef` (a Kubernetes `Secret` reference) instead of a plaintext field.
- **BREAKING** — `TrustAccessAIControlsMcpServer` (`zero` group): `authCredentials` is now **sensitive** and must be supplied via `authCredentialsSecretRef`.
- `id` changed type from number to string on `EmailSecurityBlockSender`, `EmailSecurityImpersonationRegistry`, and `EmailSecurityTrustedDomains` (computed/status field; re-observed on reconcile).
- Relaxed from required to optional (non-breaking): `source`/`type` on `AI SearchInstance`, `zoneId` on `ApiShieldDiscoveryOperation` and `WorkersCustomDomain`, `ssl` on `CustomHostname`, `notificationUrl` on `StreamWebhook`.

### Removed

- **BREAKING** — `NetworkMonitoringRule` (`magic` group): the `bandwidth` field was removed; use `bandwidthThreshold` instead.
- `Watermark` (`stream` group): the previously required `file` field was removed (watermarks are now created from `url`). Existing, already-created watermarks are unaffected; manifests still setting `file` will have it pruned.
- `EmailSecurityTrustedDomains` (`email` group): the `body` block was removed (superseded by the flat `pattern`/`isRegex`/`isRecent`/`isSimilarity` fields).

### Migration notes

For providers that manage **live** resources of the affected kinds, reconcile the managed-resource manifests **before** upgrading, otherwise updates will fail schema/CEL validation. Full before/after examples are in
[`docs/migration/cloudflare-5.21-breaking-changes.md`](docs/migration/cloudflare-5.21-breaking-changes.md).

- `EmailSecurityTrustedDomains`: add `spec.forProvider.pattern`.
- `NetworkMonitoringRule`: add `automaticAdvertisement`, `prefixes`, `type`; replace `bandwidth` with `bandwidthThreshold`.
- `logpush` `Job` / `OwnershipChallenge`: create a Kubernetes `Secret` with the destination configuration and reference it via `destinationConfSecretRef` (remove the plaintext `destinationConf`).
- `zero` `TrustAccessAIControlsMcpServer`: move `authCredentials` into a Kubernetes `Secret` referenced by `authCredentialsSecretRef`.
- `stream` `Watermark`: no action for existing watermarks; drop `file` from manifests.

[Unreleased]: https://github.com/wildbitca/provider-upjet-cloudflare/compare/v0.2.10...HEAD
[0.2.10]: https://github.com/wildbitca/provider-upjet-cloudflare/compare/v0.2.6...v0.2.10
