# Provider Cloudflare

`provider-cloudflare` is a [Crossplane](https://crossplane.io/) provider built using [Upjet](https://github.com/crossplane/upjet) that exposes XRM-conformant managed resources for the [Cloudflare API](https://developers.cloudflare.com/api/).

## Available as Provider Family

This provider is distributed as a **provider family**. Install only the Cloudflare services you need instead of all 425 CRDs.

| Package | Description | CRDs |
|---------|-------------|------|
| `provider-family-cloudflare` | ProviderConfig (auto-installed) | 5 |
| `provider-cloudflare-dns` | DNS records, DNSSEC, zone transfers | 14 |
| `provider-cloudflare-zone` | Zone settings, cache, subscriptions | 18 |
| `provider-cloudflare-cloudflare` | Rulesets, filters, lists, snippets | 24 |
| `provider-cloudflare-workers` | Workers scripts, KV, routes | 18 |
| `provider-cloudflare-r2` | R2 buckets, CORS, lifecycle | 16 |
| `provider-cloudflare-zero` | Zero Trust: Access, DLP, Gateway, Tunnel | 94 |

A monolith package (`provider-cloudflare`) with all resources is also available.

## Install

Install a sub-provider. The family provider is pulled automatically as a dependency:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-cloudflare-dns
spec:
  package: xpkg.upbound.io/wildbitca/provider-cloudflare-dns:v0.1.3
```

## ProviderConfig

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

Supported credential keys: `api_token`, `api_key`, `email`, `api_user_service_key`, `account_id`.

## Report a Bug

Open an [issue](https://github.com/wildbitca/provider-upjet-cloudflare/issues).
