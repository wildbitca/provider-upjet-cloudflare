# Publish a Worker bundle with GitHub Actions

This example publishes a Worker from an application repository to GitHub Container
Registry (GHCR). A GitOps repository records the manifest digest and Argo CD applies
the Script resource. Crossplane fetches the bundle and uploads it to Cloudflare.

## Create and publish the bundle

Add `worker.js` to the application repository:

```js
export default {
  fetch() {
    return new Response("Hello from a digest-pinned Worker");
  },
};
```

Save this workflow as `.github/workflows/publish-worker.yml` in the same repository.
Run it from the Actions tab for the commit you want to publish. For a compiled
application, run its build before packaging and list the resulting modules and
optional asset directory in the [bundle manifest](worker-artifacts.md#build-and-publish).

```yaml
name: Publish Worker bundle

on:
  workflow_dispatch:

permissions: {}

jobs:
  publish:
    runs-on: ubuntu-24.04
    permissions:
      contents: read
      packages: write
    outputs:
      ref: ${{ steps.publish.outputs.ref }}
    steps:
      - uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1
        with:
          persist-credentials: false
      - uses: oras-project/setup-oras@22ce207df3b08e061f537244349aac6ae1d214f6
        with:
          version: 1.3.0
      - name: Package the Worker
        run: |
          mkdir bundle
          cp worker.js bundle/
          cat > bundle/worker-artifact.json <<'JSON'
          {
            "schemaVersion": 1,
            "mainModule": "worker.js",
            "modules": [
              {"name": "worker.js", "contentType": "application/javascript+module"}
            ]
          }
          JSON
          tar -C bundle -czf worker-bundle.tgz worker-artifact.json worker.js
      - name: Publish to GHCR
        id: publish
        env:
          GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          printf '%s' "$GHCR_TOKEN" | oras login ghcr.io -u "$GITHUB_ACTOR" --password-stdin
          repository="ghcr.io/${GITHUB_REPOSITORY,,}-worker"
          target="$repository:$GITHUB_SHA-$GITHUB_RUN_ID-$GITHUB_RUN_ATTEMPT"
          oras push "$target" \
            --artifact-type application/vnd.gitkb.cloudflare.worker.v1 \
            worker-bundle.tgz:application/vnd.gitkb.cloudflare.worker.bundle.v1.tar+gzip
          digest="$(oras resolve "$target")"
          [[ "$digest" =~ ^sha256:[0-9a-f]{64}$ ]]
          ref="$repository@$digest"
          printf 'ref=%s\n' "$ref" >> "$GITHUB_OUTPUT"
          printf '    artifact:\n      ref: "%s"\n' "$ref" >> "$GITHUB_STEP_SUMMARY"
```

The repository's workflow token needs write access to the GHCR package. For an
existing package, grant the application repository Actions access in the package
settings. The workflow publishes a unique tag for each run attempt and exposes
the immutable reference in its summary and through
`${{ needs.publish.outputs.ref }}` for a downstream job.

## Promote the digest

Install the Workers provider and configure Cloudflare credentials in a
`ClusterProviderConfig` named `cloudflare`. For a private GHCR package, provision
a `kubernetes.io/dockerconfigjson` pull Secret named `worker-registry` in the
Script's namespace. Its credential needs read access to the package. Use the
exact registry hostname, `ghcr.io`, in the Docker config. For a public package,
omit `pullSecretRef`.

Commit this resource to the path watched by your Argo CD Application. Create the
`preview` namespace first and replace the example reference with the workflow
output and the account ID with your own.

```yaml
apiVersion: workers.upjet-cloudflare.m.upbound.io/v1alpha1
kind: Script
metadata:
  name: example-worker
  namespace: preview
spec:
  providerConfigRef:
    kind: ClusterProviderConfig
    name: cloudflare
  artifact:
    ref: ghcr.io/example/application-worker@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
    pullSecretRef:
      name: worker-registry
  forProvider:
    accountId: "your-cloudflare-account-id"
    scriptName: example-worker
    compatibilityDate: "2025-02-14"
```

If you use Helm, pass the workflow output as `artifact.ref` in the environment
values and render it under `spec.artifact` in the Script. A promotion job can
open a GitOps PR with that values change. Give it credentials scoped to the
environment repository. Keep registry publishing credentials in CI and registry
pull and Cloudflare API credentials in the cluster.

## Verify and roll back

After Argo syncs, compare `status.resolvedArtifactDigest` with the selected digest
and check the Script's `Synced` and `Ready` conditions. Configure a route or
serving endpoint separately according to your environment's access policy, then
send a request and expect `Hello from a digest-pinned Worker`. To roll back,
restore an earlier digest in GitOps and retain the corresponding registry bundle
so the controller can fetch it again.

The workflow publishes application files. Install and update the Crossplane
provider package separately. See the [Worker artifact guide](worker-artifacts.md)
for input validation, cache recovery and deletion behavior.
