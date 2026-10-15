// Package zonesetting holds a resource-scoped config.TerraformConversion for
// cloudflare_zone_setting that canonicalizes the dynamically-typed `value`
// attribute for the subset of setting_ids whose value is an object, rather
// than a scalar or array.
//
// Background: cloudflare_zone_setting.value is a Terraform
// DynamicPseudoType (it can hold any shape depending on setting_id). Upjet
// generates it schemaless on both forProvider and atProvider
// (x-kubernetes-preserve-unknown-fields) and automatically registers a
// generic dynamic-value conversion (upjet's
// pkg/config.NewTFDynamicValueConversion / dynamicValueConversion) for any
// resource with a Dynamic attribute. That generic conversion infers a
// Terraform type for the value purely by reflecting over whatever JSON the
// user supplied in forProvider.value.
//
// The real terraform-provider-cloudflare uses its own
// customfield.NormalizedDynamicValue type for `value`
// (github.com/cloudflare/terraform-provider-cloudflare/internal/customfield),
// whose DynamicSemanticEquals falls back to Go's maps.EqualFunc when
// comparing two object values - which requires both sides to have the exact
// same set of keys to ever report equal. If a user's manifest only sets a
// subset of an object setting's keys (very normal - most people don't fill
// in every optional sub-key), the inferred Terraform type for
// forProvider.value has fewer attributes than the full object Cloudflare's
// API returns in state. The two values then can never be semantically
// equal, so the provider's own PlanResourceChange always reports a
// difference and Update never converges, regardless of whether the field
// the user actually cares about changed.
//
// This package deep-merges the user-supplied partial object onto a
// canonical, full-key-set template (sourced from Cloudflare's published API
// reference) before the generic dynamic-value conversion runs, so the
// inferred Terraform type is stable across reconciles and can converge.
//
// Known limitation: config.TerraformConversion.Convert only receives
// spec.forProvider parameters, not status.atProvider, so the canonical
// template values are static zero-value defaults rather than the resource's
// last-observed server state. For automatic_platform_optimization in
// particular, some sub-keys (wordpress, wp_plugin, hostnames) are described
// by Cloudflare as auto-detected signals; because Terraform's model always
// sends the *entire* value on every apply (value is
// `tfsdk:"value" json:"value,required"` in the real provider's
// ZoneSettingModel - never a partial PATCH), this already matches what a
// user writing the full object directly in HCL would have had to do, but it
// means a reconcile can still overwrite an auto-detected field back to its
// zero-value default if the user's manifest never sets it explicitly. This
// is a bounded, detectable follow-up (same observable symptom as the bug
// this package fixes: "Update never converges for setting X"), not
// something this package can avoid given the extension point it has to work
// with.
package zonesetting

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Terraform attribute names (tf tags) on cloudflare_zone_setting.
const (
	settingIDKey = "setting_id"
	valueKey     = "value"
)

// canonicalObjectTemplates holds the full, known sub-attribute shape (with
// zero-value defaults) for every cloudflare_zone_setting setting_id whose
// value is an object instead of a scalar/array.
//
// Source: https://developers.cloudflare.com/api/resources/zones/subresources/settings/methods/edit/
// (fetched 2026-10-15), cross-checked against
// github.com/cloudflare/terraform-provider-cloudflare's published docs for
// cloudflare_zone_setting (security_header example) and
// cloudflare_nel_setting (dedicated resource mirroring the `nel`
// setting_id's value shape).
//
// NOTE: this list was re-derived from the current provider during this
// implementation. The provider's setting_id table currently has 62 entries
// and exactly 4 Object-typed ones (aegis, automatic_platform_optimization,
// nel, security_header). There is no `minify` setting_id in the current
// provider version. If Cloudflare adds a new sub-key to one of these
// objects, or adds a new Object-typed setting_id, this map needs a
// corresponding update - until then, the affected sub-key/setting falls
// back to today's pre-existing (broken-for-objects, fine-for-scalars)
// behavior rather than causing an error.
var canonicalObjectTemplates = map[string]map[string]any{
	"aegis": {
		"enabled": false,
		"pool_id": "",
	},
	"automatic_platform_optimization": {
		"cache_by_device_type": false,
		"cf":                   false,
		"enabled":              false,
		"hostnames":            []any{},
		"wordpress":            false,
		"wp_plugin":            false,
	},
	"nel": {
		"enabled": false,
	},
	"security_header": {
		"strict_transport_security": map[string]any{
			"enabled":            false,
			"include_subdomains": false,
			"max_age":            float64(0),
			"nosniff":            false,
			"preload":            false,
		},
	},
}

// Canonicalizer is a config.TerraformConversion that deep-merges a partial
// user-supplied cloudflare_zone_setting value for known object-typed
// setting_ids onto that setting's full canonical key set. It must run
// before upjet's auto-registered dynamic-value conversion infers a
// Terraform type from forProvider.value (callers should prepend it to
// Resource.TerraformConversions, not append).
//
// It is a no-op:
//   - in FromTerraform mode (status.atProvider is left exactly as the real
//     provider returned it);
//   - for any setting_id not present in canonicalObjectTemplates, which
//     covers every scalar/array setting_id (the overwhelming majority);
//   - when forProvider.value isn't present or isn't shaped like an object.
type Canonicalizer struct{}

// New returns a new Canonicalizer.
func New() Canonicalizer {
	return Canonicalizer{}
}

// Convert implements config.TerraformConversion.
func (Canonicalizer) Convert(params map[string]any, _ *config.Resource, mode config.Mode) (map[string]any, error) {
	if mode != config.ToTerraform {
		return params, nil
	}
	settingID, ok := params[settingIDKey].(string)
	if !ok || settingID == "" {
		return params, nil
	}
	template, ok := canonicalObjectTemplates[settingID]
	if !ok {
		return params, nil
	}
	userValue, ok := params[valueKey].(map[string]any)
	if !ok {
		// value is absent/nil, or not shaped like an object (e.g. we don't
		// recognize this setting_id as object-typed) - nothing to do.
		return params, nil
	}
	params[valueKey] = deepMergeDefaults(userValue, template)
	return params, nil
}

// deepMergeDefaults returns a new map containing every key from template,
// recursively overridden by any matching key present in user. Keys present
// in user but absent from template are preserved as-is, so a Cloudflare
// sub-key this template doesn't yet know about doesn't get dropped.
//
// template's own composite values (maps/slices) are deep-copied rather than
// aliased into the result: canonicalObjectTemplates is a package-level,
// unsynchronized singleton read by every concurrent reconcile of every
// cloudflare_zone_setting resource in the process, so the result must never
// share mutable storage with it.
func deepMergeDefaults(user, template map[string]any) map[string]any {
	merged := make(map[string]any, len(template)+len(user))
	for k, tv := range template {
		merged[k] = deepCopyValue(tv)
	}
	for k, uv := range user {
		tv, hasDefault := template[k]
		if !hasDefault {
			merged[k] = uv
			continue
		}
		tvMap, tvIsMap := tv.(map[string]any)
		uvMap, uvIsMap := uv.(map[string]any)
		if tvIsMap && uvIsMap {
			merged[k] = deepMergeDefaults(uvMap, tvMap)
			continue
		}
		merged[k] = uv
	}
	return merged
}

// deepCopyValue returns a deep copy of v when it is one of the composite
// shapes canonicalObjectTemplates is built from (map[string]any or []any);
// any other value (string, bool, float64, nil, ...) is already immutable
// from the caller's point of view and is returned as-is.
func deepCopyValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		cp := make(map[string]any, len(t))
		for k, vv := range t {
			cp[k] = deepCopyValue(vv)
		}
		return cp
	case []any:
		cp := make([]any, len(t))
		for i, vv := range t {
			cp[i] = deepCopyValue(vv)
		}
		return cp
	default:
		return v
	}
}
