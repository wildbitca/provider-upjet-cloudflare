package config

// Integration test for cloudflare_zone_setting Refresh/Plan/Apply through the
// real upjet workspace code path (terraform.WorkspaceStore -> FileProducer ->
// tofu) and the real terraform-provider-cloudflare binary against a local fake
// Cloudflare HTTP API.
//
// Run:
//   ORIGINPULLS_PROVIDER_BIN=/path/terraform-provider-cloudflare_v5.28.1 \
//     go test ./config -run TestZoneSettingUpdate -v -count=1
// Requires `tofu` on PATH. Skipped when both CLOUDFLARE_PROVIDER_BIN and
// ORIGINPULLS_PROVIDER_BIN are unset.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	zonev1alpha1 "github.com/wildbitca/provider-upjet-cloudflare/apis/cluster/zone/v1alpha1"
)

const zsZone = "fa595c1d547d59f66377093bedfde4e6"

type fakeSettingServer struct {
	srv     *httptest.Server
	mu      sync.Mutex
	values  map[string]any
	patches map[string][]map[string]any
}

func newFakeSettingServer(t *testing.T, initial map[string]any) *fakeSettingServer {
	t.Helper()
	f := &fakeSettingServer{
		values:  make(map[string]any, len(initial)),
		patches: make(map[string][]map[string]any),
	}
	for k, v := range initial {
		f.values[k] = v
	}
	prefix := "/zones/" + zsZone + "/settings/"
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, prefix) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		settingID := strings.TrimPrefix(r.URL.Path, prefix)
		w.Header().Set("Content-Type", "application/json")

		f.mu.Lock()
		defer f.mu.Unlock()

		switch r.Method {
		case http.MethodGet:
			val, ok := f.values[settingID]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success":  true,
				"errors":   []any{},
				"messages": []any{},
				"result": map[string]any{
					"id":          settingID,
					"value":       val,
					"editable":    true,
					"modified_on": "2026-01-01T00:00:00Z",
				},
			})
		case http.MethodPatch:
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var body map[string]any
			if err := json.Unmarshal(raw, &body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// Reject the upjet dynamicValueConversion envelope just like Cloudflare's API.
			if m, isMap := body["value"].(map[string]any); isMap && len(m) == 2 && m["type"] != nil && m["value"] != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"errors":  []any{map[string]any{"code": 1007, "message": "Invalid value for zone setting " + settingID}},
				})
				return
			}
			f.patches[settingID] = append(f.patches[settingID], body)
			if v, ok := body["value"]; ok {
				f.values[settingID] = mergePatchValue(f.values[settingID], v)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success":  true,
				"errors":   []any{},
				"messages": []any{},
				"result": map[string]any{
					"id":          settingID,
					"value":       f.values[settingID],
					"editable":    true,
					"modified_on": "2026-01-02T00:00:00Z",
				},
			})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func mergePatchValue(cur, patch any) any {
	curMap, curOK := cur.(map[string]any)
	patchMap, patchOK := patch.(map[string]any)
	if !curOK || !patchOK {
		return patch
	}
	out := make(map[string]any, len(curMap)+len(patchMap))
	for k, v := range curMap {
		out[k] = v
	}
	for k, v := range patchMap {
		out[k] = mergePatchValue(out[k], v)
	}
	return out
}

func newZoneSettingCR(uid, settingID, specJSON string, obsJSON *string) *zonev1alpha1.Setting {
	z := zsZone
	sid := settingID
	editable := true
	cr := &zonev1alpha1.Setting{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "zs-" + uid,
			UID:               types.UID(uid),
			CreationTimestamp: metav1.Now(),
		},
		Spec: zonev1alpha1.SettingSpec{
			ClusterManagedResourceSpec: xpv2.ClusterManagedResourceSpec{
				ManagementPolicies: xpv2.ManagementPolicies{
					xpv2.ManagementActionObserve,
					xpv2.ManagementActionCreate,
					xpv2.ManagementActionUpdate,
					xpv2.ManagementActionDelete,
					xpv2.ManagementActionLateInitialize,
				},
			},
			ForProvider: zonev1alpha1.SettingParameters{
				ZoneID:    &z,
				SettingID: &sid,
				Value:     &apiextensionsv1.JSON{Raw: []byte(specJSON)},
			},
		},
	}
	meta.SetExternalName(cr, settingID)
	if obsJSON != nil {
		// Intentionally omit enabled/modifiedOn/timeRemaining to match live gke-iac CRs.
		cr.Status.AtProvider = zonev1alpha1.SettingObservation{
			Editable:  &editable,
			ID:        &sid,
			SettingID: &sid,
			ZoneID:    &z,
			Value:     &apiextensionsv1.JSON{Raw: []byte(*obsJSON)},
		}
	}
	return cr
}

func TestZoneSettingUpdate(t *testing.T) {
	bin := os.Getenv("CLOUDFLARE_PROVIDER_BIN")
	if bin == "" {
		bin = os.Getenv("ORIGINPULLS_PROVIDER_BIN")
	}
	if bin == "" {
		t.Skip("CLOUDFLARE_PROVIDER_BIN / ORIGINPULLS_PROVIDER_BIN not set")
	}
	setupTofuEnv(t, bin)

	secOff := map[string]any{
		"strict_transport_security": map[string]any{
			"enabled":            false,
			"include_subdomains": false,
			"max_age":            float64(0),
			"nosniff":            false,
			"preload":            false,
		},
	}
	srv := newFakeSettingServer(t, map[string]any{
		"always_use_https":  "off",
		"challenge_ttl":     float64(1800),
		"browser_cache_ttl": float64(14400),
		"security_header":   secOff,
		"0rtt":              "off",
	})

	p := GetProvider()
	cfg := p.Resources["cloudflare_zone_setting"]
	ws := terraform.NewWorkspaceStore(logging.NewNopLogger())
	ts := terraform.Setup{
		Version:       "1.5.7",
		Requirement:   terraform.ProviderRequirement{Source: "cloudflare/cloudflare", Version: "5.28.1"},
		Configuration: map[string]any{"api_token": strings.Repeat("a", 40), "base_url": srv.srv.URL},
	}

	strPtr := func(s string) *string { return &s }
	cases := []struct {
		name      string
		settingID string
		specJSON  string
		obsJSON   *string
		wantVal   any
	}{
		{
			name:      "scalar_string_with_live_atprovider_shape",
			settingID: "always_use_https",
			specJSON:  `"on"`,
			obsJSON:   strPtr(`"off"`),
			wantVal:   "on",
		},
		{
			name:      "numeric_nonzero",
			settingID: "challenge_ttl",
			specJSON:  `3600`,
			obsJSON:   strPtr(`1800`),
			wantVal:   float64(3600),
		},
		{
			name:      "numeric_zero",
			settingID: "browser_cache_ttl",
			specJSON:  `0`,
			obsJSON:   strPtr(`14400`),
			wantVal:   float64(0),
		},
		{
			name:      "partial_object_security_header",
			settingID: "security_header",
			specJSON:  `{"strict_transport_security":{"enabled":true,"max_age":15552000}}`,
			obsJSON:   strPtr(`{"strict_transport_security":{"enabled":false,"include_subdomains":false,"max_age":0,"nosniff":false,"preload":false}}`),
			wantVal: map[string]any{
				"strict_transport_security": map[string]any{
					"enabled":            true,
					"include_subdomains": false,
					"max_age":            float64(15552000),
					"nosniff":            false,
					"preload":            false,
				},
			},
		},
		{
			name:      "empty_atprovider_fallback",
			settingID: "0rtt",
			specJSON:  `"on"`,
			obsJSON:   nil,
			wantVal:   "on",
		},
		{
			name:      "stale_envelope_in_atprovider",
			settingID: "0rtt",
			specJSON:  `"off"`,
			obsJSON:   strPtr(`{"type":"string","value":"on"}`),
			wantVal:   "off",
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cr := newZoneSettingCR(fmt.Sprintf("zs-uid-%d", i), tc.settingID, tc.specJSON, tc.obsJSON)
			ctx := context.Background()

			w, err := ws.Workspace(ctx, noopSecrets{}, cr, ts, cfg)
			if err != nil {
				t.Fatalf("Workspace() = %v", err)
			}
			ref, err := w.Refresh(ctx)
			if err != nil {
				t.Fatalf("Refresh() = %v", err)
			}
			if !ref.Exists {
				t.Fatalf("Refresh() reported Exists=false")
			}

			planBefore, err := w.Plan(ctx)
			if err != nil {
				t.Fatalf("Plan() before apply = %v", err)
			}
			if planBefore.UpToDate {
				t.Fatalf("Plan() before apply reported UpToDate=true, expected update")
			}

			applyRes, err := w.Apply(ctx)
			if err != nil {
				t.Fatalf("Apply() = %v", err)
			}

			// Verify FromTerraform unwrapping works on applied state.
			var rawObs map[string]any
			if err := json.Unmarshal(applyRes.State.GetAttributes(), &rawObs); err != nil {
				t.Fatalf("Unmarshal(apply state) = %v", err)
			}
			obs, err := cfg.ApplyTFConversions(rawObs, config.FromTerraform)
			if err != nil {
				t.Fatalf("ApplyTFConversions(FromTerraform) = %v", err)
			}
			if !reflect.DeepEqual(obs["value"], tc.wantVal) {
				t.Fatalf("post-apply unwrapped value = %#v, want %#v", obs["value"], tc.wantVal)
			}

			// Verify fake Cloudflare received the exact plain PATCH body.
			srv.mu.Lock()
			patches := srv.patches[tc.settingID]
			edgeVal := srv.values[tc.settingID]
			srv.mu.Unlock()
			if len(patches) == 0 {
				t.Fatalf("expected PATCH request for %s, got none", tc.settingID)
			}
			if !reflect.DeepEqual(edgeVal, tc.wantVal) {
				t.Fatalf("fake Cloudflare edge value = %#v, want %#v (patches: %#v)", edgeVal, tc.wantVal, patches)
			}

			// Verify subsequent Plan converges (no drift loop).
			planAfter, err := w.Plan(ctx)
			if err != nil {
				t.Fatalf("Plan() after apply = %v", err)
			}
			if !planAfter.UpToDate {
				t.Fatalf("Plan() after apply reported UpToDate=false (update did not converge)")
			}
		})
	}
}
