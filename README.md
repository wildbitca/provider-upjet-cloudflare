# Provider Cloudflare

`provider-cloudflare` is a [Crossplane](https://crossplane.io/) provider built using [Upjet](https://github.com/crossplane/upjet) that exposes XRM-conformant managed resources for the [Cloudflare API](https://developers.cloudflare.com/api/).

## Getting Started

This provider is available in two distribution modes:

| Mode | Package | CRDs | Use when |
|------|---------|------|----------|
| **Family** (recommended) | `provider-family-cloudflare` + individual sub-providers | Only what you install | You use a subset of Cloudflare services |
| **Monolith** | `provider-cloudflare` | All 425 | You need everything or want the simplest setup |

> For provider families background, see [Scalable Provider Families](https://blog.crossplane.io/crd-scaling-provider-families/) and the [Upbound Provider Families docs](https://docs.upbound.io/manuals/packages/providers/provider-families/).

### Install (Family — recommended)

Install only the Cloudflare services you need. The first sub-provider you install automatically pulls in `provider-family-cloudflare`, which manages the shared `ProviderConfig`.

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-dns
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-dns:v0.1.0
```

Need more services? Add more sub-providers:

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
```

See [docs/family/](docs/family/) for the full list of available sub-providers and configuration details.

### Install (Monolith)

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare:v0.1.0
```

### ProviderConfig

Create a Secret with Cloudflare credentials (API token recommended):

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
```

Create a ProviderConfig:

```yaml
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

### Rate limiting

Cloudflare’s API allows about **1,200 requests per five minutes**. The provider defaults to `--max-reconcile-rate=3` to reduce the risk of HTTP 429s. For large deployments:

- Set `--max-reconcile-rate` to **2–5** when deploying (e.g. in your Provider or Helm values).
- Keep `--sync` at **1h** (default) or higher to limit full drift-detection frequency.

Example override when running the provider:

```bash
provider --max-reconcile-rate=3 --sync=1h ...
```

## Developing

Fetch submodules first, then run the code-generation pipeline. `make generate` runs the Upjet generator (same as `go run cmd/generator/main.go "$PWD"`).

```bash
make submodules
go install golang.org/x/tools/cmd/goimports@latest
export PATH="$(go env GOPATH)/bin:$PATH"
make generate
```

Run code-generation pipeline (alternative):

```bash
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster (monolith):

```bash
make run
```

Build monolith binary:

```bash
make build
```

Build family sub-providers (all groups):

```bash
make build.family
```

Build specific sub-providers:

```bash
make build.family FAMILY_SUBPACKAGES="config dns zone"
```

**Fixing CI (check-diff / lint):** If the CI **check-diff** or **lint** job fails, regenerate locally and commit the result:

```bash
make submodules
make generate
make check-diff
git add -A && git commit -m "chore: regenerate for CI" && git push
```

## Supported resources

You can see the API reference at [doc.crds.dev](https://doc.crds.dev/github.com/wildbitca/provider-upjet-cloudflare) or on the [Upbound Marketplace](https://marketplace.upbound.io/providers/wildbitca/provider-cloudflare).

This provider exposes 198 managed resources from the [Cloudflare Terraform Provider](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs), including:

- **Zone**: `cloudflare_zone`, `cloudflare_zone_dnssec`, `cloudflare_zone_lockdown`, etc.
- **DNS**: `cloudflare_dns_record`, `cloudflare_dns_firewall`, etc.
- **Account**: `cloudflare_account`, `cloudflare_account_member`, etc.
- **Workers**: `cloudflare_worker`, `cloudflare_workers_kv`, `cloudflare_workers_route`, etc.
- **Load Balancer**: `cloudflare_load_balancer`, `cloudflare_load_balancer_pool`, `cloudflare_load_balancer_monitor`
- **Zero Trust**: `cloudflare_zero_trust_access_*`, `cloudflare_zero_trust_device_*`, etc.
- And more.

## Adding resources

- Find the resource in the [Cloudflare Terraform Provider docs](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs).
- Prefer **1 resource per PR**.
- Write a test case for new or modified resources.

## Report a Bug

Open an [issue](https://github.com/wildbitca/provider-upjet-cloudflare/issues).
