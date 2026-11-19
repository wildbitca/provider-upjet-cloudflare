# zone Setting experiments (T1)

Scope: terraform-provider-cloudflare v5.28.1 (scratch clone, nothing committed there), upjet v2.5.1. No network, no Cloudflare writes, no cluster access.
Reproduce E1: copy `zone-setting-e1_test.go.txt` to `internal/e1/e1_test.go` in a v5.28.1 clone and run
`GOMAXPROCS=2 go test -p 1 ./internal/e1 -run TestE1 -v` (importing the full `zone_setting` package exhausts a 4 GiB memory limit; the test therefore uses a verbatim copy of the model fields from `internal/services/zone_setting/model.go:19-26` and calls the same `apijson.MarshalForPatch` that `model.go:31-33` uses).

## Verdict: root cause NOT proven. No fix written.

## E1 - PATCH body (`resource.go:117` `MarshalJSONForUpdate` = `apijson.MarshalForPatch(plan, state)`)

| case (plan vs state) | patch body | full (Create) body |
|---|---|---|
| scalar `on` vs `off`, enabled unknown | `{"value":"on"}` | `{"value":"on"}` |
| scalar equal, enabled unknown | `{}` | `{"value":"on"}` |
| scalar equal, enabled known | `{}` | `{"enabled":true,"value":"on"}` |
| scalar differs, enabled known equal | `{"value":"on"}` | `{"enabled":true,"value":"on"}` |
| object differs in one field | `{"value":{"enabled":true}}` (partial) | `{"value":{"enabled":true,"max_age":100}}` |
| object equal | `{}` | full object |
| number 1.3 vs 1.2 | `{"value":1.3}` | `{"value":1.3}` |
| plan Int64 vs state Float64 (equal) | PANIC `reflect.Value.Int on float64` (`encoder.go:265`) | `{"value":1}` |
| plan String vs state Bool | `{"value":"on"}` | `{"value":"on"}` |

Findings:
- Claim "a differing scalar produces an empty/invalid PATCH body": **REFUTED**. A differing scalar yields `{"value":"on"}`, which is what Cloudflare accepts (`encoder.go:255-270` scalar omission only applies to equal values).
- Claim "equal scalar/object yields `{}` if Update runs": **CONFIRMED** (`encoder.go:255-270`, `522-540`). `{}` matches the observed `400 code 1007` on the settings that were equal to live.
- Claim "object diff yields a partial object": **CONFIRMED** (`{"value":{"enabled":true}}`), consistent with the `500` on `security_header`. Whether Cloudflare 500s on the partial object is not tested (no live write).
- New hazard: mismatched plan/state numeric types panic in the encoder (`encoder.go:265`), so a type difference is not silent.

Interpretation: the 400/500 bodies are explained IF Update runs for settings whose value equals live. This does not explain why differing settings are never attempted (symptom a in the issue); that depends on how the plan/state are built (E2/E3), which was not exercised.

## E2 - framework PlanResourceChange with fake transport: NOT RUN
Needs importing the full provider package; the compile exceeded the 4 GiB host memory limit on this host even with `GOMAXPROCS=2`/`-p 1`. Hypotheses H1 (state already holds desired value), H2 (type mismatch causes spurious drift) and the plan-modifier swallow: **UNDECIDED**.

## E3 - upjet config/state type equality: NOT RUN
Read-only finding only: `tf_conversion.go:92-125` (upjet v2.5.1) wraps the value as `{"value": v, "type": inferred}` in ToTerraform and does nothing in FromTerraform (`tf_conversion.go:93-97`). `inferTFTypeFromValue` maps every number to `tftypes.Number` and all lists by first element. H2 and H3 remain **UNDECIDED**.

## E4 - stale Synced=False: NOT RUN
Not read in depth (`external_tfpluginfw.go`, `nofork_store.go`). UNDECIDED.

## What is missing (smallest next experiments)
1. E2 on a machine/worker with more memory (or `go test -c` with `GOFLAGS=-p=1` in a cgroup-limited worker with >=8 GiB): run `PlanResourceChange` with prior state from a Read of `http3=off` and config `http3=on` produced by upjet's conversion; record whether the plan has a diff and which attributes differ (`value` vs `enabled`).
2. E3 in this repo: feed MR forProvider `on` / atProvider `off` through `Canonicalizer` + `NewTFDynamicValueConversion().Convert(ToTerraform)` and compare the tftypes with the state type produced by the TF Read.
3. Canary capture on `mkcloud.xyz` with debug logging (T6) to see the real plan diff; it is the cheapest way to separate H1 from H2.
