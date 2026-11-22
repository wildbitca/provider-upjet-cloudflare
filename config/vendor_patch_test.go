package config_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/spf13/afero"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dnsv1alpha1 "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/dns/v1alpha1"
	zonev1alpha1 "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/zone/v1alpha1"
	"github.com/wildbitca/provider-upjet-cloudflare/config"
)

// Fails unless patches/upjet-skip-state-without-id.patch is applied, i.e.
// unless the tests run after `make vendor`.
func TestNewResourceGetsNoSyntheticState(t *testing.T) {
	ctx := context.Background()
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	name := "example.com"
	zone := "zone-id"
	tr := &dnsv1alpha1.Record{
		ObjectMeta: metav1.ObjectMeta{Name: "apex-spf", Namespace: "cloudflare"},
		Spec:       dnsv1alpha1.RecordSpec{ForProvider: dnsv1alpha1.RecordParameters{Name: &name, ZoneID: &zone}},
	}
	cfg := config.GetProviderNamespaced().Resources["cloudflare_dns_record"]
	fp, err := terraform.NewFileProducer(ctx, nil, "ws", tr, terraform.Setup{}, cfg, terraform.WithFileSystem(fs))
	if err != nil {
		t.Fatalf("NewFileProducer() = %v", err)
	}
	if err := fp.EnsureTFState(ctx, ""); err != nil {
		t.Fatalf("EnsureTFState() = %v", err)
	}
	if s, _ := fs.ReadFile(filepath.Join("ws", "terraform.tfstate")); len(s) != 0 {
		t.Fatalf("synthetic tfstate written for a Record with no external name; Refresh will call Read with an empty id "+
			"(is the vendored upjet patched? run `make vendor`):\n%s", s)
	}
}

// Fails unless patches/upjet-dynamic-config-no-envelope.patch is applied, i.e.
// unless the tests run after `make vendor`.
func TestDynamicSettingMainTFHasNoEnvelope(t *testing.T) {
	ctx := context.Background()
	cfg := config.GetProviderNamespaced().Resources["cloudflare_zone_setting"]

	cases := []struct {
		name      string
		settingID string
		wantSpec  string
		obsValue  string
		wantMain  any
	}{
		{
			name:      "scalar string",
			settingID: "always_use_https",
			wantSpec:  `"on"`,
			obsValue:  `"off"`,
			wantMain:  "on",
		},
		{
			name:      "canonicalized partial object",
			settingID: "security_header",
			wantSpec:  `{"strict_transport_security":{"enabled":true,"max_age":15552000}}`,
			obsValue:  `{"strict_transport_security":{"enabled":false,"include_subdomains":false,"max_age":0,"nosniff":false,"preload":false}}`,
			wantMain: map[string]any{
				"strict_transport_security": map[string]any{
					"enabled":            true,
					"include_subdomains": false,
					"max_age":            float64(15552000),
					"nosniff":            false,
					"preload":            false,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.Afero{Fs: afero.NewMemMapFs()}
			zone := "zone-id"
			sid := tc.settingID
			tr := &zonev1alpha1.Setting{
				ObjectMeta: metav1.ObjectMeta{Name: "s", Namespace: "cloudflare"},
				Spec: zonev1alpha1.SettingSpec{
					ForProvider: zonev1alpha1.SettingParameters{
						ZoneID:    &zone,
						SettingID: &sid,
						Value:     &apiextensionsv1.JSON{Raw: []byte(tc.wantSpec)},
					},
				},
				Status: zonev1alpha1.SettingStatus{
					AtProvider: zonev1alpha1.SettingObservation{
						ID:        &sid,
						ZoneID:    &zone,
						SettingID: &sid,
						Value:     &apiextensionsv1.JSON{Raw: []byte(tc.obsValue)},
					},
				},
			}
			meta.SetExternalName(tr, sid)

			fp, err := terraform.NewFileProducer(ctx, nil, "ws", tr, terraform.Setup{}, cfg, terraform.WithFileSystem(fs))
			if err != nil {
				t.Fatalf("NewFileProducer() = %v", err)
			}
			if err := fp.EnsureTFState(ctx, sid); err != nil {
				t.Fatalf("EnsureTFState() = %v", err)
			}
			if _, err := fp.WriteMainTF(); err != nil {
				t.Fatalf("WriteMainTF() = %v", err)
			}

			rawMain, err := fs.ReadFile(filepath.Join("ws", "main.tf.json"))
			if err != nil {
				t.Fatalf("ReadFile(main.tf.json) = %v", err)
			}
			var mainTF struct {
				Resource map[string]map[string]map[string]any `json:"resource"`
			}
			if err := json.Unmarshal(rawMain, &mainTF); err != nil {
				t.Fatalf("Unmarshal(main.tf.json) = %v", err)
			}
			gotVal := mainTF.Resource["cloudflare_zone_setting"]["s"]["value"]
			if !reflect.DeepEqual(gotVal, tc.wantMain) {
				t.Fatalf("main.tf.json value = %#v, want %#v (is the vendored upjet patched? run `make vendor`)", gotVal, tc.wantMain)
			}

			rawState, err := fs.ReadFile(filepath.Join("ws", "terraform.tfstate"))
			if err != nil {
				t.Fatalf("ReadFile(terraform.tfstate) = %v", err)
			}
			var tfState struct {
				Resources []struct {
					Instances []struct {
						Attributes map[string]any `json:"attributes"`
					} `json:"instances"`
				} `json:"resources"`
			}
			if err := json.Unmarshal(rawState, &tfState); err != nil {
				t.Fatalf("Unmarshal(terraform.tfstate) = %v", err)
			}
			stateVal, ok := tfState.Resources[0].Instances[0].Attributes["value"].(map[string]any)
			if !ok || stateVal["type"] == nil || stateVal["value"] == nil {
				t.Fatalf("terraform.tfstate value should still carry {type,value} envelope, got %#v", tfState.Resources[0].Instances[0].Attributes["value"])
			}
		})
	}
}
