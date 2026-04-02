# Provider Family Configuration

This provider is distributed as a **provider family**, allowing you to install only the Cloudflare services you need instead of all 425 CRDs.

For background on provider families, see:
- [Solving the Crossplane Provider CRD Scaling Problem](https://blog.crossplane.io/crd-scaling-provider-families/)
- [Upbound Provider Families documentation](https://docs.upbound.io/manuals/packages/providers/provider-families/)

## Monolith vs Family

| | Monolith | Family |
|---|---|---|
| **Package** | `provider-cloudflare` | `provider-family-cloudflare` + sub-providers |
| **CRDs installed** | 425 (all resources) | Only what you install |
| **Pods** | 1 | 1 per sub-provider + 1 family |
| **Memory** | High (all controllers loaded) | Low (only active controllers) |
| **Best for** | Simple setups, testing | Production, resource-constrained clusters |

> **Important**: The monolith and family providers **cannot coexist** in the same cluster. If migrating from monolith to family, see [Migration](#migrating-from-monolith-to-family) below.

## Prerequisites

- Crossplane v1.14+ or Universal Crossplane (UXP) v1.14+
- `kubectl` access to the cluster

## Installing Family Providers

### Step 1: Install a sub-provider

Install the sub-providers for the Cloudflare services you use. The **first sub-provider** you install automatically pulls in `provider-family-cloudflare`, which manages the shared `ProviderConfig`.

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-dns
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-dns:v0.1.0
```

After installation, `kubectl get providers` shows both the sub-provider and the auto-installed family provider:

```
NAME                          INSTALLED   HEALTHY   PACKAGE                                                        AGE
provider-cloudflare-dns       True        True      xpkg.upbound.io/wildbitca/provider-cloudflare-dns:v0.1.0      60s
provider-family-cloudflare    True        True      xpkg.upbound.io/wildbitca/provider-family-cloudflare:v0.1.0   60s
```

### Step 2: Add more sub-providers as needed

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-zone
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-zone:v0.1.0
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-workers
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-workers:v0.1.0
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-r2
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-r2:v0.1.0
```

### Step 3: Configure credentials

The `ProviderConfig` is managed by `provider-family-cloudflare` and shared across all sub-providers. Create it once:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloudflare-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "api_token": "YOUR_CLOUDFLARE_API_TOKEN"
    }
---
apiVersion: upjet-cloudflare.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: cloudflare-creds
      namespace: crossplane-system
      key: credentials
```

Supported credential keys: `api_token`, `api_key`, `email`, `api_user_service_key`, `account_id`, `api_hostname`, `api_base_path`, `api_client_logging`.

## Available Sub-Providers

Each sub-provider maps to a Cloudflare API group. CRD counts include both cluster-scoped and namespaced variants.

| Group | Package | CRDs | Description |
|-------|---------|------|-------------|
| `config` | `provider-family-cloudflare` | 5 | ProviderConfig (auto-installed) |
| `dns` | `provider-cloudflare-dns` | 14 | DNS records, DNSSEC, zone transfers |
| `zone` | `provider-cloudflare-zone` | 18 | Zone settings, cache, subscriptions |
| `cloudflare` | `provider-cloudflare-cloudflare` | 24 | Rulesets, filters, lists, snippets, images |
| `workers` | `provider-cloudflare-workers` | 18 | Workers scripts, KV, routes, cron triggers |
| `r2` | `provider-cloudflare-r2` | 16 | R2 buckets, CORS, lifecycle, custom domains |
| `zero` | `provider-cloudflare-zero` | 94 | Zero Trust: Access, DLP, Gateway, Tunnel, Device |
| `email` | `provider-cloudflare-email` | 16 | Email routing, email security |
| `magic` | `provider-cloudflare-magic` | 20 | Magic Transit/WAN, GRE/IPsec tunnels |
| `stream` | `provider-cloudflare-stream` | 14 | Stream video, live inputs, watermarks |
| `account` | `provider-cloudflare-account` | 12 | Account settings, members, tokens |
| `load` | `provider-cloudflare-load` | 6 | Load balancers, pools, monitors |
| `api` | `provider-cloudflare-api` | 14 | API Shield, schema validation |
| `authenticated` | `provider-cloudflare-authenticated` | 8 | Authenticated origin pulls |
| `custom` | `provider-cloudflare-custom` | 8 | Custom hostnames, custom SSL |
| `cloudforce` | `provider-cloudflare-cloudforce` | 8 | CloudForce One requests |
| `waiting` | `provider-cloudflare-waiting` | 8 | Waiting rooms |
| `pages` | `provider-cloudflare-pages` | 4 | Pages projects, domains |
| `schema` | `provider-cloudflare-schema` | 6 | Schema validation settings |
| `logpush` | `provider-cloudflare-logpush` | 4 | Logpush jobs |
| `notification` | `provider-cloudflare-notification` | 4 | Notification policies |
| `regional` | `provider-cloudflare-regional` | 4 | Regional hostnames, tiered cache |
| `content` | `provider-cloudflare-content` | 4 | Content scanning |
| `ai` | `provider-cloudflare-ai` | 4 | AI search instances |
| `argo` | `provider-cloudflare-argo` | 4 | Argo smart routing, tiered caching |
| `calls` | `provider-cloudflare-calls` | 4 | Calls SFU/TURN apps |
| `leaked` | `provider-cloudflare-leaked` | 4 | Leaked credential checks |
| `page` | `provider-cloudflare-page` | 4 | Page rules, page shield |
| `token` | `provider-cloudflare-token` | 4 | Token validation |
| `web` | `provider-cloudflare-web` | 4 | Web analytics |
| `d1` | `provider-cloudflare-d1` | 2 | D1 databases |
| `queue` | `provider-cloudflare-queue` | 2 | Queue consumers |
| `worker` | `provider-cloudflare-worker` | 2 | Worker versions |

<details>
<summary>Remaining sub-providers (2 CRDs each)</summary>

| Group | Package | Description |
|-------|---------|-------------|
| `access` | `provider-cloudflare-access` | Access rules |
| `address` | `provider-cloudflare-address` | Address maps |
| `bot` | `provider-cloudflare-bot` | Bot management |
| `byo` | `provider-cloudflare-byo` | BYO IP prefixes |
| `certificate` | `provider-cloudflare-certificate` | Certificate packs |
| `cloud` | `provider-cloudflare-cloud` | Cloud connector rules |
| `connectivity` | `provider-cloudflare-connectivity` | Connectivity directory |
| `firewall` | `provider-cloudflare-firewall` | Firewall rules |
| `hostname` | `provider-cloudflare-hostname` | Hostname TLS settings |
| `hyperdrive` | `provider-cloudflare-hyperdrive` | Hyperdrive configs |
| `image` | `provider-cloudflare-image` | Image variants |
| `keyless` | `provider-cloudflare-keyless` | Keyless certificates |
| `list` | `provider-cloudflare-list` | List items |
| `logpull` | `provider-cloudflare-logpull` | Logpull retention |
| `managed` | `provider-cloudflare-managed` | Managed transforms |
| `mtls` | `provider-cloudflare-mtls` | mTLS certificates |
| `observatory` | `provider-cloudflare-observatory` | Observatory scheduled tests |
| `organization` | `provider-cloudflare-organization` | Organization profiles |
| `origin` | `provider-cloudflare-origin` | Origin CA certificates |
| `rate` | `provider-cloudflare-rate` | Rate limits |
| `registrar` | `provider-cloudflare-registrar` | Registrar domains |
| `snippet` | `provider-cloudflare-snippet` | Snippet rules |
| `spectrum` | `provider-cloudflare-spectrum` | Spectrum applications |
| `sso` | `provider-cloudflare-sso` | SSO connectors |
| `tiered` | `provider-cloudflare-tiered` | Tiered cache |
| `total` | `provider-cloudflare-total` | Total TLS |
| `turnstile` | `provider-cloudflare-turnstile` | Turnstile widgets |
| `universal` | `provider-cloudflare-universal` | Universal SSL settings |
| `url` | `provider-cloudflare-url` | URL normalization |
| `user` | `provider-cloudflare-user` | User agent blocking |
| `web3` | `provider-cloudflare-web3` | Web3 hostnames |

</details>

## Common Combinations

Here are typical sub-provider sets for common use cases:

**Website management** (DNS + WAF + Zone settings):
```
provider-cloudflare-dns + provider-cloudflare-zone + provider-cloudflare-cloudflare
```
→ 56 CRDs (vs 425 monolith)

**Workers platform** (Workers + KV + R2 + D1 + Queues):
```
provider-cloudflare-workers + provider-cloudflare-r2 + provider-cloudflare-d1 + provider-cloudflare-queue
```
→ 38 CRDs

**Zero Trust**:
```
provider-cloudflare-zero
```
→ 94 CRDs

**Full stack** (DNS + Zone + Security + Workers + R2):
```
provider-cloudflare-dns + provider-cloudflare-zone + provider-cloudflare-cloudflare + provider-cloudflare-workers + provider-cloudflare-r2
```
→ 90 CRDs

## Migrating from Monolith to Family

> The monolith and family providers **cannot coexist**. The monolith provider owns the CRDs; installing a family sub-provider for the same resources will conflict.

### Migration steps

1. **Export existing managed resources** for backup:
   ```bash
   kubectl get managed -o yaml > backup-managed-resources.yaml
   ```

2. **Delete the monolith provider** (this removes its CRDs):
   ```bash
   kubectl delete provider provider-cloudflare
   ```

3. **Install family sub-providers** for the services you use (see examples above).

4. **Recreate ProviderConfig** with the same name and credentials. The `ProviderConfig` schema is identical between monolith and family.

5. **Verify** — Claims and XRs will automatically reconcile and recreate managed resources once the corresponding CRDs are available:
   ```bash
   kubectl get providers
   kubectl get managed
   ```

### What stays the same

- `ProviderConfig` name, schema, and credentials
- CRD API groups and versions (`dns.upjet-cloudflare.upbound.io/v1alpha1`, etc.)
- Managed resource specs and external names
- Compositions and XRDs that reference these resources

### What changes

- Multiple provider Pods instead of one
- Each sub-provider has its own leader election
- `provider-family-cloudflare` Pod handles authentication for all sub-providers

## Versioning

All sub-providers share the same version tag. A single `git tag v0.1.0` publishes the entire family. There is no independent versioning per sub-provider.
