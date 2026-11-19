# Plan: zone `Setting` Update fix (provider-cloudflare-zone v1.0.5)

Source issue: wildbitca/org-gitops#372 (comment "Update 2026-11-19"). Versions: terraform-provider-cloudflare v5.28.1, upjet v2.5.1.
Goal: make `Setting` (cloudflare_zone_setting) reconcile correctly: differing values get PATCHed with a valid body, equal values stay Synced, and Synced=False does not stick until pod restart.

## Sizing
- mode: light | resume_test: fail (sessions) | risk_items: 3 (contracts others consume: CRD/atProvider shape), 5 (delivery: provider image rollout), 6 (no existing test reproduces the live behaviour)
- forecast: lines~150 units=1 surfaces=[config/zonesetting, patches/, examples] | route: implementer+tester+code-reviewer+verifier | verification: independent verifier | spawn_budget: 4
- divisible: no | reason: independent-verification (T1 experiments gate T3/T4; units share state)
- scorecard: resumability 2, blast_radius 2 (all zone Settings of all zones), risk 2, uncertainty 3 (root cause unproven), verifiability 2, spread 1 -> STANDARD

## Root-cause status: HYPOTHESIS (evidence collected, not proven)

### What the sources show (file refs are in terraform-provider-cloudflare v5.28.1 unless noted)
1. `internal/services/zone_setting/model.go`: `Value` is `customfield.NormalizedDynamicValue` (`json:"value,required"`); `Enabled` is `computed_optional`. `Update` (`resource.go`, `Update`) builds the body with `data.MarshalJSONForUpdate(*state)` = `apijson.MarshalForPatch(plan, state)`, i.e. a *diff* body, not the full `{"value": v}` that `Create` sends via `MarshalJSON`.
2. `internal/apijson/encoder.go` (~L522 DynamicValuable branch; L248-286 scalar branches): in patch mode a scalar equal to the state is omitted. So a PATCH body is `{}` or `{"enabled":...}` whenever the planned `value` is considered equal to the state's, or contains only computed_optional `enabled`. Cloudflare answers such bodies with `400 code 1007 Invalid value for zone setting <id>`. This matches symptom (a) ("plain `{"value": ...}` works"), and a nested object diff (security_header) producing a partial object matches the `500`.
3. `internal/customfield/dynamic.go`: `NormalizingDynamicPlanModifier` replaces the plan with the state when `semanticEquals` says equal; `semanticEquals` compares via `t.IntValue/FloatValue/ChildAttributes` and falls back to `maps.EqualFunc`. Both sides must carry the same type after upjet's inference.
4. upjet v2.5.1 `pkg/config/tf_conversion.go` `dynamicValueConversion.Convert`: ToTerraform wraps `value` as `{"value": v, "type": inferTFTypeFromValue(v)}` and FromTerraform is a no-op ("Terraform does not return dynamic pseudo-types in state"). `inferTFTypeFromValue` is type-by-reflection (e.g. every JSON number -> Number, lists typed by first element, empty list -> list(object)). Existing repo code (`config/zonesetting/canonicalize.go`) already proved the premise "FromTerraform is a no-op" is false: the envelope leaks into atProvider and is unwrapped by `unwrapDynamicEnvelope`.
5. Symptom (b) (differing value never PATCHed, equal value PATCHed) is the inverse of what step 3 predicts for a healthy state. Candidate explanations, none proven:
   - H1: the prior state upjet feeds to PlanResourceChange (reconstructed from `status.atProvider`/tfstate annotation) already contains the *desired* value (or a stale/enveloped one), so plan==state and `ProposedNewState` carries no diff -> Observe says up to date; equal values pass through a code path that does diff against live.
   - H2: type mismatch between inferred config type and Read-derived state type (e.g. `bool`/`string` vs `{"type": ..}`-wrapped, or `number` vs int64) makes `semanticEquals` false for equal values (spurious drift) while the plan modifier's replacement-by-state hides real drift.
   - H3: late-initialization copies the observed `value` into forProvider, so the desired value is overwritten by live.
6. Symptom (c) (Synced=False until restart) is consistent with a cached `LastAsyncOperation=False`/stale tf state annotation in the upjet no-fork external client (`pkg/controller/external_tfpluginfw.go`, `nofork_store.go`); not read in depth yet.

Not proven: which of H1-H3 (or a combination of the patch-diff body in item 2) is the live cause. T1 below decides this with unit-level experiments only (no Cloudflare writes).

## Tasks

### T1 - Experiments to prove the root cause (read-only on the repo; scratch only)
Run in a scratch clone of `terraform-provider-cloudflare` at v5.28.1 (`$AI_RESOURCES_SCRATCH/clone`; the TF module path is not v5-suffixed so it cannot be `go get`-ed from this repo) and in this repo. Outputs go to `docs/agent-history/implementer/zone-setting-experiments.md`.
- E1 (patch body): in the scratch clone add a throwaway test calling `ZoneSettingModel.MarshalJSONForUpdate` for: scalar differing (`off`->`on`), scalar equal, plan value dynamic-typed differently from state (String vs Dynamic-wrapped), `enabled` unknown/known, object (`security_header`). Record exact JSON. Confirms or refutes item 2 (empty/partial body).
- E2 (plan diff, in-process): drive the framework resource schema (`zone_setting.NewResource().Schema` + `PlanResourceChange` via `providerserver`/`tfprotov6`, fake HTTP transport, no network) with prior state = Read result and config = upjet output of `dynamicValueConversion` for `http3 on` vs live `off`. Record whether the plan has a diff and the planned `value`.
- E3 (upjet state path): in this repo write a test that feeds the sequence MR (forProvider.value=on, atProvider.value=off with and without envelope) through `canonicalize` + upjet `dynamicValueConversion` and logs the tftypes of config and reconstructed state; assert type equality of both sides.
- E4 (staleness, c): read `external_tfpluginfw.go`/`nofork_store.go` for what is cached across reconciles (opTracker, `LastAsyncOperation`) and reproduce with a fake external client: first Update fails, then desired spec corrected; assert Synced recovers without restart.
- Acceptance: each of H1-H3 and the patch-body claim is marked CONFIRMED/REFUTED with captured JSON/types. Root-cause status in this document is updated to PROVEN only after this.

### T2 - Regression tests (must fail before the fix, pass after) - `config/zonesetting/`
- `TestPatchBodyScalarDiffersFromLive` (new file `config/zonesetting/update_body_test.go`): the converted params for `http3` on vs live off yield a non-empty planned diff and a body containing `value`.
- `TestScalarEqualToLiveIsNoop`: equal desired/live yields no diff and no PATCH.
- `TestConvertFromTerraformNeverLeaksEnvelope` (extends `envelope_test.go`): atProvider.value is the plain value for scalar, object, list, null.
- `TestStateAndConfigTypesMatch` (E3 as a permanent test): inferred config type == type rebuilt from atProvider for each of: bool-like `on/off` strings, numbers (`min_tls_version` "1.3", `0rtt`), lists, `security_header` object.
- `TestFailedUpdateDoesNotStickSynced` (E4): after a failed Update, a corrected spec reconciles to Synced=True without client restart.
These tests must be recorded failing on the current tree (`go test ./config/zonesetting/...` output saved by the tester) before T3.

### T3 - Fix (only after T1 proves the cause; shape depends on result)
Preferred options, in order (decision recorded in T1's report):
1. Extend `Canonicalizer.Convert` (`config/zonesetting/canonicalize.go`) to normalize *both* directions so config type and state type are identical (unwrap envelope in atProvider; for ToTerraform ensure scalar numbers/bools/strings use the same representation the TF Read produces).
2. If the diff-body is the cause (item 2): force a full `{"value": v}` body by making the plan always differ from state for a changed value (do not let `NormalizingDynamicPlanModifier` swallow it), or patch upjet via `patches/*.patch` + `patches/apply.sh` (see PR #34 pattern) so the dynamic conversion preserves the type the TF Read produced. A vendored patch needs a test in `config/vendor_patch_test.go` style.
3. If the staleness (c) is in upjet's nofork store: add a vendored patch resetting cached async-operation state on spec generation change.
Also in this task: drop `minify` from docs/examples if present (the composition itself lives in org-gitops; see T6).
- Files: `config/zonesetting/canonicalize.go`, optionally `patches/*.patch`, `examples/` for the CR.
- Acceptance: all T2 tests pass; `go build ./...`, `go vet ./config/...`, `go test ./config/...` green; no changes to generated CRDs unless the schema changes (state it explicitly).

### T4 - Handoff to tester
Run `go test ./config/... ./internal/...` and `make lint` (or the repo CI equivalent) and attach the failing-before/passing-after evidence for every T2 test.

### T5 - Review and verify
code-reviewer: check that no behaviour changes for the 4 object setting_ids, that `canonicalObjectTemplates` semantics hold (`canonicalize_test.go` unchanged and green), and that any vendored patch applies idempotently via `patches/apply.sh`. verifier: map each acceptance criterion below to evidence.

### T6 - Rollout (owner: user/org-gitops, not the implementer)
1. Release provider-cloudflare-zone vX (next patch after v1.0.5) with CHANGELOG entry.
2. Canary zone `mkcloud.xyz` only: pin the canary Settings to the new package, enable Update for 3 Settings: `http3` (scalar, differing), `min_tls_version` (equal), `security_header` (object). Provider debug log on; capture request bodies.
3. Success gate: differing setting gets PATCHed once with a `{"value": ...}` body and converges; equal settings produce no PATCH; no 400/500; no Synced=False after forced error recovery without a pod restart.
4. Then widen to the remaining zones (wildbit-org etc.) in batches. Rollback: revert package pin and set managementPolicies back to Observe-only.
5. Separate small org-gitops PR: remove the retired `minify` Setting from the composition.

## Acceptance criteria (Given/When/Then)
- AC1: Given a Setting `http3` with desired `on` and live `off`, when reconciled, then exactly one PATCH is issued with body containing `value:"on"` and the resource becomes Synced=True. (TestPatchBodyScalarDiffersFromLive)
- AC2: Given desired == live, when reconciled, then no PATCH. (TestScalarEqualToLiveIsNoop)
- AC3: Given atProvider rebuilt after a pod restart, when refreshed, then no "Unable to Read Previously Saved State" and no envelope in atProvider. (TestConvertFromTerraformNeverLeaksEnvelope)
- AC4: Given a failed Update followed by a corrected spec, when reconciled, then Synced=True without restarting the provider. (TestFailedUpdateDoesNotStickSynced)
- AC5: Existing `canonicalize_test.go` and `envelope_test.go` cases still pass unchanged.
- AC6: Canary zone `mkcloud.xyz` meets the T6 success gate before wider rollout.

## Risks
- Root cause unproven; T3 is gated by T1 (do not implement from the hypothesis).
- Wide blast radius: every Setting on every zone; a wrong body can change live zone configuration (TLS, security_header HSTS). Mitigation: canary zone and Observe-only rollback.
- Vendored upjet patches add maintenance cost on upjet upgrades (must stay idempotent in `patches/apply.sh`).
- `canonicalObjectTemplates` zero-defaults can overwrite auto-detected fields (known limitation in `canonicalize.go`).
- The TF provider cannot be imported as a Go module from this repo, so TF-side experiments run in a scratch clone only and must not be committed.
