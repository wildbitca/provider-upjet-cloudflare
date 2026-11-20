# Release notes and specifics: provider-upjet-cloudflare

The release and rollout procedure is shared across the Upjet providers. See the org-gitops runbook
`docs/runbooks/upjet-provider-release-and-rollout.md` (https://github.com/wildbitca/org-gitops).

## Repo specifics

- Terraform provider and OpenTofu versions are pinned in two places: `Makefile` and `.github/workflows/publish.yml`.
- Commit `examples-generated` after regeneration.
- v2.0.0 is a MAJOR (CRD-breaking removals); see `docs/migration/` and the org-gitops runbook for the list. Only
  `TrustAccessApplication` and `TrustOrganization` are used by wildbit and were unaffected.

## Version history, 2026-11-20

- v1.0.6: bugfix/dependency release.
- v2.0.0: feature release with the Terraform provider bump to 5.29.0 (and OpenTofu 1.13.9).
