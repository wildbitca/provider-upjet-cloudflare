package zonesetting

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/upjet/v2/pkg/config"
)

func TestConvertObjectSettings(t *testing.T) {
	type args struct {
		params map[string]any
		mode   config.Mode
	}
	type want struct {
		params map[string]any
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"SecurityHeaderPartialUserObjectIsCanonicalized": {
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "security_header",
					"zone_id":    "zone-1",
					"value": map[string]any{
						"strict_transport_security": map[string]any{
							"enabled": true,
							"max_age": float64(86400),
						},
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "security_header",
					"zone_id":    "zone-1",
					"value": map[string]any{
						"strict_transport_security": map[string]any{
							"enabled":            true,
							"include_subdomains": false,
							"max_age":            float64(86400),
							"nosniff":            false,
							"preload":            false,
						},
					},
				},
			},
		},
		"SecurityHeaderAlreadyFullObjectIsUnchanged": {
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "security_header",
					"value": map[string]any{
						"strict_transport_security": map[string]any{
							"enabled":            true,
							"include_subdomains": true,
							"max_age":            float64(31536000),
							"nosniff":            true,
							"preload":            true,
						},
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "security_header",
					"value": map[string]any{
						"strict_transport_security": map[string]any{
							"enabled":            true,
							"include_subdomains": true,
							"max_age":            float64(31536000),
							"nosniff":            true,
							"preload":            true,
						},
					},
				},
			},
		},
		"AegisPartialUserObjectIsCanonicalized": {
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "aegis",
					"value": map[string]any{
						"enabled": true,
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "aegis",
					"value": map[string]any{
						"enabled": true,
						"pool_id": "",
					},
				},
			},
		},
		"AutomaticPlatformOptimizationPartialUserObjectIsCanonicalized": {
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "automatic_platform_optimization",
					"value": map[string]any{
						"enabled":              true,
						"cache_by_device_type": true,
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "automatic_platform_optimization",
					"value": map[string]any{
						"enabled":              true,
						"cache_by_device_type": true,
						"cf":                   false,
						"hostnames":            []any{},
						"wordpress":            false,
						"wp_plugin":            false,
					},
				},
			},
		},
		"NelPartialUserObjectIsCanonicalized": {
			// nel's value only has one known key, so in practice a "partial"
			// object and the full object are the same shape - this is the
			// degenerate/trivial case of the fix, included for completeness.
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "nel",
					"value": map[string]any{
						"enabled": true,
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "nel",
					"value": map[string]any{
						"enabled": true,
					},
				},
			},
		},
		"UnknownSubKeyIsPreservedNotDropped": {
			// If Cloudflare adds a new sub-key before this template is
			// updated, we must not silently drop the user's value for it.
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "aegis",
					"value": map[string]any{
						"enabled":        true,
						"future_new_key": "some-value",
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "aegis",
					"value": map[string]any{
						"enabled":        true,
						"pool_id":        "",
						"future_new_key": "some-value",
					},
				},
			},
		},
		"FromTerraformModeIsNoOpEvenForObjectSettings": {
			args: args{
				mode: config.FromTerraform,
				params: map[string]any{
					"setting_id": "security_header",
					"value": map[string]any{
						"strict_transport_security": map[string]any{
							"enabled": true,
						},
					},
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "security_header",
					"value": map[string]any{
						"strict_transport_security": map[string]any{
							"enabled": true,
						},
					},
				},
			},
		},
		"MissingSettingIDIsNoOp": {
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"value": map[string]any{"enabled": true},
				},
			},
			want: want{
				params: map[string]any{
					"value": map[string]any{"enabled": true},
				},
			},
		},
		"MissingValueIsNoOp": {
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": "security_header",
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": "security_header",
				},
			},
		},
	}

	// Scalar/array setting_ids must be a complete no-op (regression guard
	// for the 58 non-object setting_ids): a representative sample covering
	// every value shape the generic dynamic-value conversion already
	// handles correctly (on/off string, numeric, and array-of-string),
	// plus an unrecognized setting_id, must all pass through unchanged.
	for _, tc := range []struct {
		name      string
		settingID string
		scalarVal any
	}{
		{name: "StringOnOffScalarSetting", settingID: "ssl", scalarVal: "strict"},
		{name: "NumericScalarSetting", settingID: "browser_cache_ttl", scalarVal: float64(14400)},
		{name: "ArrayOfStringSetting", settingID: "ciphers", scalarVal: []any{"ECDHE-ECDSA-AES128-GCM-SHA256"}},
		{name: "UnrecognizedSettingIDWithObjectValue", settingID: "some_future_setting", scalarVal: map[string]any{"anything": true}},
	} {
		cases["ScalarPassthrough_"+tc.name] = struct {
			args args
			want want
		}{
			args: args{
				mode: config.ToTerraform,
				params: map[string]any{
					"setting_id": tc.settingID,
					"value":      tc.scalarVal,
				},
			},
			want: want{
				params: map[string]any{
					"setting_id": tc.settingID,
					"value":      tc.scalarVal,
				},
			},
		}
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := New().Convert(tc.args.params, nil, tc.args.mode)
			if err != nil {
				t.Fatalf("Convert(...): unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want.params, got); diff != "" {
				t.Errorf("Convert(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestDeepMergeDefaultsDoesNotMutateInputs(t *testing.T) {
	template := map[string]any{"a": false, "b": map[string]any{"c": false}}
	user := map[string]any{"a": true}

	merged := deepMergeDefaults(user, template)
	if len(user) != 1 {
		t.Errorf("deepMergeDefaults(...) mutated its user input: %+v", user)
	}
	if merged["a"] != true {
		t.Errorf("deepMergeDefaults(...)[\"a\"] = %v, want true", merged["a"])
	}
	bv, ok := merged["b"].(map[string]any)
	if !ok || bv["c"] != false {
		t.Errorf("deepMergeDefaults(...)[\"b\"] = %v, want map[c:false]", merged["b"])
	}

	// The returned "b" must not alias template's "b": mutating it must not
	// corrupt the shared canonicalObjectTemplates-shaped input for the next
	// caller.
	bv["c"] = true
	if tv := template["b"].(map[string]any)["c"]; tv != false {
		t.Errorf("mutating deepMergeDefaults(...)[\"b\"] leaked into template: template[\"b\"][\"c\"] = %v, want false", tv)
	}
}
