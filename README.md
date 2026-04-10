# Provider Cloudflare

`provider-cloudflare` is a [Crossplane](https://crossplane.io/) provider built using [Upjet](https://github.com/crossplane/upjet) that exposes XRM-conformant managed resources for the [Cloudflare API](https://developers.cloudflare.com/api/).

## Getting Started

This provider uses the **family** distribution model: install only the sub-providers you need. The first sub-provider you install automatically pulls in `provider-family-cloudflare`, which manages the shared `ProviderConfig`.

> For provider families background, see [Scalable Provider Families](https://blog.crossplane.io/crd-scaling-provider-families/) and the [Upbound Provider Families docs](https://docs.upbound.io/manuals/packages/providers/provider-families/).

### Install

Install only the Cloudflare services you need:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-dns
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-dns:v0.2.0
```

Need more services? Add more sub-providers:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-zone
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-zone:v0.2.0
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-workers
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-workers:v0.2.0
```

### Available sub-providers

| Sub-provider | Resources |
|-------------|-----------|
| `provider-cloudflare-dns` | DNS records, firewall, zone transfers |
| `provider-cloudflare-zone` | Zone settings, DNSSEC, cache variants, subscriptions |
| `provider-cloudflare-cloudflare` | Rulesets, healthchecks, filters, lists |
| `provider-cloudflare-workers` | Workers scripts, KV namespaces, routes, cron triggers |
| `provider-cloudflare-r2` | R2 buckets, CORS, lifecycle, custom domains |
| `provider-cloudflare-bot` | Bot management |
| `provider-cloudflare-managed` | Managed transforms |
| `provider-cloudflare-leaked` | Leaked credential checks |
| `provider-cloudflare-tiered` | Tiered caching |
| `provider-cloudflare-universal` | Universal SSL settings |
| `provider-cloudflare-url` | URL normalization |
| `provider-cloudflare-web` | Web analytics |
| `provider-cloudflare-notification` | Notification policies |
| And 47 more... | See full list in `cmd/provider/` |

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

Cloudflare's API allows about **1,200 requests per five minutes**. For large deployments, use a `DeploymentRuntimeConfig` to throttle:

```yaml
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: cloudflare-throttled
spec:
  deploymentTemplate:
    spec:
      selector: {}
      template:
        spec:
          containers:
            - name: package-runtime
              args:
                - --poll=30m
                - --max-reconcile-rate=1
                - --sync=4h
```

## Developing

```bash
make submodules
go install golang.org/x/tools/cmd/goimports@latest
export PATH="$(go env GOPATH)/bin:$PATH"
make generate
```

Run against a Kubernetes cluster:

```bash
make run
```

Build family sub-providers:

```bash
make build.family
make build.family FAMILY_SUBPACKAGES="config dns zone"
```

## Supported resources

This provider exposes 198 managed resources from the [Cloudflare Terraform Provider v5.18.0](https://registry.terraform.io/providers/cloudflare/cloudflare/5.18.0/docs).

API reference: [doc.crds.dev](https://doc.crds.dev/github.com/wildbitca/provider-upjet-cloudflare)

## Report a Bug

Open an [issue](https://github.com/wildbitca/provider-upjet-cloudflare/issues).
