# Migration guide — Cloudflare TF provider 5.18.0 → 5.21.1 (provider v0.2.10)

Provider **v0.2.10** upgrades the upstream Cloudflare Terraform provider from
`v5.18.0` to `v5.21.1`. This is additive for most resources, but **6 managed
resources changed incompatibly**. If you manage live resources of these kinds,
apply the changes below **before** upgrading the affected sub-providers,
otherwise the next reconcile will fail schema/CEL validation and the managed
resource will go `SYNCED=False`.

> **Identity is preserved.** All affected resources use
> `IdentifierFromProvider`, so the external resource is tracked by the
> `crossplane.io/external-name` annotation, not by these spec fields. As long as
> you **edit the existing managed resource in place** (never delete + recreate),
> Crossplane keeps managing the same Cloudflare object — these edits do not
> recreate anything.

Kinds appear in both scopes:
- Cluster-scoped group: `<group>.upjet-cloudflare.upbound.io/v1alpha1`
- Namespaced group: `<group>.upjet-cloudflare.m.upbound.io/v1alpha1`

`SecretRef` shape differs by scope:
- **Cluster** managed resources: `{ name, namespace, key }`
- **Namespaced** managed resources: `{ name, key }` (Secret must live in the MR's namespace)

Affected sub-providers to hold back until migrated: **`email`, `magic`, `logpush`, `zero`**. `stream` needs no action.

---

## 1. `logpush` `Job` and `OwnershipChallenge` — `destinationConf` is now sensitive

`destination_conf` became a **sensitive** field. It moves from a plaintext
`spec.forProvider.destinationConf` to `spec.forProvider.destinationConfSecretRef`
pointing at a Kubernetes `Secret`. (The `Job` `ownershipChallenge` field is
likewise referenced via `ownershipChallengeSecretRef`.)

Create the Secret with the existing destination string (e.g. an R2/S3/GCS URL
with its access token):

```bash
kubectl -n crossplane-system create secret generic logpush-destination \
  --from-literal=destinationConf='r2://my-bucket/{DATE}?account-id=...&access-key-id=...&secret-access-key=...'
```

**Before (5.18):**
```yaml
apiVersion: logpush.upjet-cloudflare.upbound.io/v1alpha1
kind: Job
metadata:
  name: example-logpush-job
spec:
  forProvider:
    zoneId: <zone_id>
    dataset: gateway_dns
    destinationConf: "r2://my-bucket/{DATE}?account-id=...&access-key-id=...&secret-access-key=..."
```

**After (5.21) — cluster:**
```yaml
apiVersion: logpush.upjet-cloudflare.upbound.io/v1alpha1
kind: Job
metadata:
  name: example-logpush-job
spec:
  forProvider:
    zoneId: <zone_id>
    dataset: gateway_dns
    destinationConfSecretRef:
      name: logpush-destination
      namespace: crossplane-system
      key: destinationConf
```

**After (5.21) — namespaced** (Secret in the same namespace as the MR; drop `namespace`):
```yaml
apiVersion: logpush.upjet-cloudflare.m.upbound.io/v1alpha1
kind: Job
metadata:
  name: example-logpush-job
  namespace: my-namespace
spec:
  forProvider:
    destinationConfSecretRef:
      name: logpush-destination
      key: destinationConf
```

`OwnershipChallenge` is identical — replace `destinationConf` with
`destinationConfSecretRef`.

---

## 2. `zero` `TrustAccessAIControlsMcpServer` — `authCredentials` is now sensitive

`auth_credentials` became **sensitive**. Move it into a Secret referenced by
`spec.forProvider.authCredentialsSecretRef`.

```bash
kubectl -n crossplane-system create secret generic mcp-auth \
  --from-literal=authCredentials='<existing auth credentials value>'
```

**Before → After (cluster):**
```yaml
# before
spec:
  forProvider:
    authCredentials: "<value>"
# after
spec:
  forProvider:
    authCredentialsSecretRef:
      name: mcp-auth
      namespace: crossplane-system
      key: authCredentials
```
(Namespaced: drop `namespace`; keep the Secret in the MR's namespace.)

---

## 3. `email` `SecurityTrustedDomains` — `pattern` now required, `body` removed

The `body` block was removed; the domain is now expressed with the flat
`pattern` field (now **required**), plus optional `isRegex` / `isRecent` /
`isSimilarity`.

**Before (5.18):**
```yaml
spec:
  forProvider:
    accountId: <account_id>
    body:
      - pattern: example.com
        isRegex: false
        isRecent: true
```

**After (5.21):**
```yaml
spec:
  forProvider:
    accountId: <account_id>
    pattern: example.com      # now required
    isRegex: false
    isRecent: true
    isSimilarity: false
```

---

## 4. `magic` `NetworkMonitoringRule` — new required fields, `bandwidth` renamed

- `automaticAdvertisement`, `prefixes`, and `type` are now **required**.
- `bandwidth` was **removed** → use `bandwidthThreshold`.

**Before (5.18):**
```yaml
spec:
  forProvider:
    accountId: <account_id>
    name: my_rule_1
    bandwidth: 1000
    duration: 1m
    packetThreshold: 10000
```

**After (5.21):**
```yaml
spec:
  forProvider:
    accountId: <account_id>
    name: my_rule_1
    bandwidthThreshold: 1000     # renamed from `bandwidth`
    automaticAdvertisement: true # now required
    prefixes:                    # now required
      - 203.0.113.1/32
    type: zscore                 # now required (zscore | threshold)
    duration: 1m
    packetThreshold: 10000
    prefixMatch: exact
    zscoreSensitivity: high
    zscoreTarget: bits
```

---

## 5. `stream` `Watermark` — `file` removed (no action for existing)

The required `file` field was removed (watermarks are now created from `url`).
Already-created watermarks are unaffected. For new manifests, drop `file`; if you
recreate a watermark, supply `url`. Any lingering `file:` in a manifest is
silently pruned by the CRD's structural schema.

---

## Rollout order (to release/stabilize safely)

1. Publish/upgrade the **unaffected** sub-providers and the **new** `secretsstore` first — they are safe.
2. Create the Secrets from §1 and §2.
3. Patch the live managed resources of the affected kinds in place (§1–§4).
4. Only then upgrade the `email`, `magic`, `logpush`, and `zero` sub-providers.
5. Verify: `kubectl get managed | grep -Ev 'True .*True'` should be empty for the affected kinds.
