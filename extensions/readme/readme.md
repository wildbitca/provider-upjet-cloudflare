# Provider Upjet Cloudflare

`provider-upjet-cloudflare` is a [Crossplane](https://crossplane.io/) provider built using [Upjet](https://github.com/crossplane/upjet) that exposes XRM-conformant managed resources for the [Cloudflare API](https://developers.cloudflare.com/api/).

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) 1.24+
- [OpenTofu](https://opentofu.org/) 1.9+ (required for schema generation; auto-installed by `make generate`)
- Kubernetes cluster (e.g. [Kind](https://kind.sigs.k8s.io/))

### Install

```bash
up ctp provider install wildbitca/provider-upjet-cloudflare:v0.1.0
```

Or declaratively:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-cloudflare
spec:
  package: xpkg.upbound.io/wildbitca/provider-upjet-cloudflare:v0.1.0
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

Cloudflare's API allows about **1,200 requests per five minutes**. The provider defaults to `--max-reconcile-rate=3` to reduce the risk of HTTP 429s. For large deployments:

- Set `--max-reconcile-rate` to **2–5** when deploying (e.g. in your Provider or Helm values).
- Keep `--sync` at **1h** (default) or higher to limit full drift-detection frequency.

## Supported resources

You can see the API reference at [doc.crds.dev](https://doc.crds.dev/github.com/wildbitca/provider-upjet-cloudflare) or on the [Upbound Marketplace](https://marketplace.upbound.io/providers/wildbitca/provider-upjet-cloudflare).

This provider exposes managed resources from the [Cloudflare Terraform Provider](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs), including Zones, DNS, Workers, Load Balancers, Zero Trust, and more.

## Report a Bug

Open an [issue](https://github.com/wildbitca/provider-upjet-cloudflare/issues).
