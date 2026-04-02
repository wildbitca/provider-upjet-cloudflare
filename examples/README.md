# Examples

Apply resources in order:

1. **Secret** (credentials): Create a Secret with key `credentials` containing JSON with `api_token` (and optionally `account_id`). Use `cluster/providerconfig/secret.yaml.tmpl` or `namespaced/providerconfig/secret.yaml.tmpl` as a template; replace `YOUR_CLOUDFLARE_API_TOKEN` with a real token.

2. **ProviderConfig**: Apply `cluster/providerconfig/providerconfig.yaml` (or namespaced equivalent). It references the secret (e.g. `example-creds`) and is named `default`.

3. **Zone**: Apply `cluster/zone/zone.yaml`. Set `spec.forProvider.name` and `spec.forProvider.accountId` (or `accountIdRef`) for your domain and account. The example uses `providerConfigRef.name: default`.

4. **DNS Record**: Apply `cluster/dns/dns-record.yaml` after the Zone exists. It references the zone via `zoneIdRef.name: example-zone` and uses the same ProviderConfig.

For local testing:

- Apply CRDs: `kubectl apply -f package/crds` (or the path where the Makefile outputs CRDs).
- Start the provider: `make run` (set `TERRAFORM_VERSION`, `TERRAFORM_PROVIDER_SOURCE`, `TERRAFORM_PROVIDER_VERSION` if needed).
- Apply the example manifests above, then run `kubectl describe <resource>` to confirm reconciliation and the `crossplane.io/external-name` annotation.
