# Prebuilt Worker artifacts

`Script` (`workers.upjet-cloudflare[.m].upbound.io/v1alpha1`) accepts an optional
`spec.artifact`. It selects an immutable OCI bundle containing prebuilt Worker
modules and optional Workers Static Assets. CI builds the application; the
provider fetches it and supplies local files to the existing upload resources.

Cloudflare added native Terraform uploads of Workers and static assets in
[Terraform provider 5.11.0](https://developers.cloudflare.com/changelog/post/2025-10-09-assets-terraform/).
This extension supplies release artifacts to that existing implementation. It
does not run Wrangler or introduce a separate OpenTofu Workspace/state backend.
The provider's existing Terraform execution and managed-resource lifecycle apply.

For the complete CI-to-deployment flow, see
[GitHub Actions and GitOps](github-actions-gitops.md): application workflow YAML,
bundle publication, digest outputs, preview promotion, registry/Cloudflare
credential ownership and deployment verification. The sections below define the
bundle format and managed-resource inputs used by that workflow.

## Build and publish

The OCI image manifest must have artifact type
`application/vnd.gitkb.cloudflare.worker.v1` and exactly one layer of type
`application/vnd.gitkb.cloudflare.worker.bundle.v1.tar+gzip`. No container image
entrypoint, application build, or executable installation is involved.

The layer contains `worker-artifact.json` at its root:

```json
{
  "schemaVersion": 1,
  "mainModule": "worker.js",
  "modules": [
    {"name": "worker.js", "contentType": "application/javascript+module"},
    {"name": "chunks/helper.js", "contentType": "application/javascript+module"}
  ],
  "assets": "public"
}
```

Module names are bundle-relative file paths and retain those names in the Worker
upload. The main module must be in the module list. A static-only application can
omit `mainModule` and `modules`, leaving `schemaVersion` and `assets`. A Worker
without static assets can omit `assets`. Modules must be outside the public asset
directory. Compatibility dates/flags, bindings, secrets, account/Worker identity,
asset routing configuration, routes and access policies belong in the resource
or surrounding application chart, rather than the bundle manifest.

For example, after producing the build output and manifest in `dist/`:

```sh
tar -C dist -czf worker-bundle.tgz worker-artifact.json worker.js chunks public
oras push ghcr.io/example/worker-bundles:build-123 \
  --artifact-type application/vnd.gitkb.cloudflare.worker.v1 \
  worker-bundle.tgz:application/vnd.gitkb.cloudflare.worker.bundle.v1.tar+gzip
oras resolve ghcr.io/example/worker-bundles:build-123
```

Promote the resulting manifest digest, not the mutable build tag. Keep published
bundles available for recovery and rollback for as long as the GitOps revision is
supported. The provider verifies both manifest and compressed-layer digests.

## Select the bundle

```yaml
apiVersion: workers.upjet-cloudflare.m.upbound.io/v1alpha1
kind: Script
metadata:
  name: preview-worker
  namespace: preview-123
spec:
  artifact:
    ref: ghcr.io/example/worker-bundles@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
    pullSecretRef:
      name: worker-registry
  forProvider:
    accountId: "your-cloudflare-account-id"
    scriptName: preview-123-worker
    compatibilityDate: "2026-11-16"
    compatibilityFlags: [nodejs_compat]
    bindings:
      - name: ASSETS
        type: assets
      - name: SITE
        type: service
        service: preview-123-site
  providerConfigRef:
    kind: ClusterProviderConfig
    name: cloudflare
```

Use an existing Cloudflare provider configuration with the permissions needed for
your Worker resources. The example's `pullSecretRef` is optional for a public
registry. A private registry uses a `kubernetes.io/dockerconfigjson` Secret,
independently of the Cloudflare API credentials. The provider supports explicit
username/password, encoded `auth`, or identity-token credentials in `auths` for
the exact registry host. Docker credential helper executables are not run.
Use External Secrets or another credential manager to refresh expiring registry
credentials, including ECR authorization tokens.

Namespaced resources can only select pull Secrets from their own namespace.
Cluster-scoped resources must supply `pullSecretRef.namespace`. Registry traffic
uses HTTPS with normal certificate verification. The controller needs a writable
temporary directory and network access to the selected registry.

Artifact mode owns code/module inputs and the asset upload directory. It rejects
explicit `content`, `contentFile`, `contentSha256`, `contentType`, `mainModule`,
`files`, `modules`, `bodyPart`, `keepAssets`, `assets.directory` and
`assets.jwtSecretRef` in either `forProvider` or `initProvider`.
Other inputs, including `assets.config`, remain available. Without `spec.artifact`,
the existing input path is used.

`Script` uploads the main module through `content_file` and extra modules through
the existing `files` map. New bundles can update the Script pin in place.

The separate `Version` resource does not accept `spec.artifact`. An OrbStack test
against Cloudflare exposed a native refresh failure for a new Version without an
ID; native module changes also require replacement, conflicting with Upjet
`prevent_destroy`. Version artifact support needs a separate lifecycle design.
`status.resolvedArtifactDigest` reports a successfully fetched bundle, not a
successful Cloudflare deployment. Use the resource's normal conditions and
application smoke checks to verify deployment behavior.

## Recovery and cache lifecycle

The cache is disposable and isolated by managed-resource UID and artifact
reference. Files are verified before becoming visible to reconciliation. Cache
hits validate file hashes; missing or corrupted files are downloaded again.
The transient `ArtifactFiles` field is excluded from Kubernetes serialization.
Derived paths are supplied to Terraform and filtered out of late initialization,
so they do not become desired inputs in the Kubernetes spec. A pod restart can
therefore reconstruct the deployment inputs from GitOps intent.

Inactive cache versions and abandoned downloads become eligible for removal
after 24 hours. Later materializations sweep at most once per hour, taking the
same per-resource locks as downloads. Removing a cache never deletes Cloudflare
resources. Deleting a managed resource uses the existing provider lifecycle and
does not require fetching the artifact or reading its registry Secret. The
provider supplies a minimal inline module only to satisfy native configuration
validation during refresh/destroy; the managed deletion path never uploads it. Normal
Cloudflare credentials and the resource's management/deletion policy still apply.

Controller replacement while a Script is already deleting is not restart-safe:
Upjet skips Terraform-state reconstruction during deletion, and the managed
resource may disappear while its remote Worker remains. The live synthetic test
reproduced this limit. Keep the controller running through normal deletion,
verify the Worker is absent in Cloudflare afterward, and remove a confirmed
orphan manually if necessary. Cache-independent deletion does not establish
restart-safe deletion.

Archives accept regular files and directories, with at most 10,000 entries,
128 MiB compressed and 256 MiB expanded. Symlinks, hard links, device files,
duplicate paths and paths escaping the bundle are rejected. The reserved
`.artifact-ready.json` cache marker cannot be supplied by an archive.

This feature covers Workers and Workers Static Assets. It does not implement the
separate Cloudflare Pages content-deployment API or select environment DNS,
Cloudflare Access policies, VPC bindings or public ingress automatically.

## Private Worker resources

Namespaced references connect Worker VPC bindings to Directory services,
Directory services to tunnels, and Access policies to service tokens. Tunnel
and service-token connection Secrets support the tunnel connector and its clients.

Access service tokens return a computed `client_secret_version` even when no
rotation was requested. Late initialization excludes that field in both API
scopes: sending the returned version back as configuration requires a
`previous_client_secret_expires_at` value and otherwise blocks the next refresh.
Explicit rotations can still configure both inputs. If an earlier controller
already copied the version into `spec.forProvider.clientSecretVersion`, remove
only that automatically populated field after installing the corrected package.
Keep the existing external identifier and connection Secret so the same token
continues to reconcile.

Artifact-backed Scripts resolve references using incremental merge patches with
a resource-version precondition. This preserves binding Secret selectors and
earlier resolved fields when a later reference changes, without overwriting
concurrent GitOps edits. Reference updates also retain the controller-local
bundle attachment prepared by the initializer, because it is intentionally
absent from API responses. The normal reference resolver remains selected for
Scripts without an artifact; other controllers are unchanged. The implementation
uses the managed reconciler's reference-resolver extension point and changes no
Crossplane or Upjet dependencies. If an earlier controller removed binding
selectors, resync the existing GitOps manifest after installing the corrected
Workers package to restore the declared sources.
