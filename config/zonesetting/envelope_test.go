package zonesetting

import (
	"reflect"
	"testing"

	"github.com/crossplane/upjet/v2/pkg/config"
)

func TestConvertUnwrapsDynamicEnvelope(t *testing.T) {
	envelope := func(inner any) map[string]any {
		return map[string]any{
			"type":  []any{"object", map[string]any{"type": "string", "value": "string"}},
			"value": inner,
		}
	}
	cases := map[string]struct {
		mode config.Mode
		in   any
		want any
	}{
		"FromTerraformScalar": {config.FromTerraform, envelope("on"), "on"},
		"ToTerraformStale":    {config.ToTerraform, envelope("on"), "on"},
		"NullValue":           {config.FromTerraform, envelope(nil), nil},
		"PlainString":         {config.FromTerraform, "on", "on"},
		"PlainObject": {config.FromTerraform,
			map[string]any{"strict_transport_security": map[string]any{"enabled": true}},
			map[string]any{"strict_transport_security": map[string]any{"enabled": true}}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := New().Convert(map[string]any{"setting_id": "0rtt", "value": tc.in}, nil, tc.mode)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got["value"], tc.want) {
				t.Fatalf("got %#v, want %#v", got["value"], tc.want)
			}
		})
	}
}
