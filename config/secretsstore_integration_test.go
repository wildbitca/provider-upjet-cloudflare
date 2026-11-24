package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	ssv1alpha1 "github.com/wildbitca/provider-upjet-cloudflare/apis/cluster/secretsstore/v1alpha1"
)

type fixedSecret struct{ val []byte }

func (f fixedSecret) GetSecretData(context.Context, *xpv2.SecretReference) (map[string][]byte, error) {
	return nil, nil
}
func (f fixedSecret) GetSecretValue(context.Context, xpv2.SecretKeySelector) ([]byte, error) {
	return f.val, nil
}

// TestSecretsStoreRefresh verifies against the real Terraform provider binary
// that Store and StoreSecret Refresh/Import work whether external-name is the
// bare leaf ID or the legacy composite path (<account_id>/<store_id> or
// <account_id>/<store_id>/<secret_id>), and that a new resource without an
// external name refreshes as missing.
func TestSecretsStoreRefresh(t *testing.T) {
	bin := os.Getenv("ORIGINPULLS_PROVIDER_BIN")
	if bin == "" {
		bin = os.Getenv("CLOUDFLARE_PROVIDER_BIN")
	}
	if bin == "" {
		t.Skip("ORIGINPULLS_PROVIDER_BIN / CLOUDFLARE_PROVIDER_BIN not set")
	}
	setupTofuEnv(t, bin)

	const (
		accID    = "00000000000000000000000000000001"
		storeID  = "00000000000000000000000000000002"
		secretID = "00000000000000000000000000000003"
	)
	wantStorePath := fmt.Sprintf("/accounts/%s/secrets_store/stores/%s", accID, storeID)
	wantSecretPath := fmt.Sprintf("/accounts/%s/secrets_store/stores/%s/secrets/%s", accID, storeID, secretID)

	var mu sync.Mutex
	var reqs []string
	resetReqs := func() {
		mu.Lock()
		reqs = nil
		mu.Unlock()
	}
	snapshotReqs := func() []string {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(reqs)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		reqs = append(reqs, r.Method+" "+r.URL.Path)
		mu.Unlock()
		if r.Method == http.MethodGet && r.URL.Path == wantStorePath {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success":true,"errors":[],"messages":[],"result":{"id":%q,"account_id":%q,"name":"default_store","created":"2026-07-13T17:49:59Z","modified":"2026-10-27T01:34:50Z"}}`, storeID, accID)
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == wantSecretPath {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success":true,"errors":[],"messages":[],"result":{"id":%q,"account_id":%q,"store_id":%q,"name":"SZ_DEV_TEST","scopes":["workers"],"status":"active","created":"2026-07-13T17:49:59Z","modified":"2026-10-27T01:34:50Z"}}`, secretID, accID, storeID)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"success":false,"errors":[{"code":404,"message":"not found %s %s"}]}`, r.Method, r.URL.Path)
	}))
	t.Cleanup(srv.Close)

	p := GetProvider()
	storeCfg := p.Resources["cloudflare_secrets_store"]
	if storeCfg == nil {
		t.Fatal("cloudflare_secrets_store not configured")
	}
	secretCfg := p.Resources["cloudflare_secrets_store_secret"]
	if secretCfg == nil {
		t.Fatal("cloudflare_secrets_store_secret not configured")
	}

	newStoreCR := func(uid, extName string, policies ...xpv2.ManagementAction) resource.Terraformed {
		cr := &ssv1alpha1.Store{}
		cr.SetName("st-" + uid)
		cr.SetUID(types.UID(uid))
		cr.SetManagementPolicies(policies)
		acc, nm := accID, "default_store"
		cr.Spec.ForProvider.AccountID = &acc
		cr.Spec.ForProvider.Name = &nm
		if extName != "" {
			meta.SetExternalName(cr, extName)
		}
		cr.CreationTimestamp = metav1.Now()
		return cr
	}

	newSecretCR := func(uid, extName string, policies ...xpv2.ManagementAction) resource.Terraformed {
		cr := &ssv1alpha1.StoreSecret{}
		cr.SetName("ss-" + uid)
		cr.SetUID(types.UID(uid))
		cr.SetManagementPolicies(policies)
		acc, st, nm, sc := accID, storeID, "SZ_DEV_TEST", "workers"
		cr.Spec.ForProvider.AccountID = &acc
		cr.Spec.ForProvider.StoreID = &st
		cr.Spec.ForProvider.Name = &nm
		cr.Spec.ForProvider.Scopes = []*string{&sc}
		cr.Spec.ForProvider.ValueSecretRef = xpv2.SecretKeySelector{Key: "k", SecretReference: xpv2.SecretReference{Name: "s", Namespace: "ns"}}
		if extName != "" {
			meta.SetExternalName(cr, extName)
		}
		cr.CreationTimestamp = metav1.Now()
		return cr
	}

	ws := terraform.NewWorkspaceStore(logging.NewNopLogger())
	ts := terraform.Setup{
		Version:       "1.5.7",
		Requirement:   terraform.ProviderRequirement{Source: "cloudflare/cloudflare", Version: "5.28.1"},
		Configuration: map[string]any{"api_token": strings.Repeat("a", 40), "base_url": srv.URL},
	}

	assertStateID := func(t *testing.T, attrs []byte, wantID string) map[string]any {
		t.Helper()
		var state map[string]any
		if err := json.Unmarshal(attrs, &state); err != nil {
			t.Fatalf("unmarshal state: %v", err)
		}
		if state["id"] != wantID {
			t.Fatalf("state id = %v, want %q", state["id"], wantID)
		}
		return state
	}

	runCases := func(
		t *testing.T,
		prefix string,
		cfg *ujconfig.Resource,
		wantID string,
		newCR func(uid, extName string, policies ...xpv2.ManagementAction) resource.Terraformed,
		cases map[string]string,
	) {
		t.Helper()
		t.Run("new CR without external-name refreshes as missing", func(t *testing.T) {
			cr := newCR(prefix+"-new", "", xpv2.ManagementActionObserve, xpv2.ManagementActionCreate, xpv2.ManagementActionUpdate, xpv2.ManagementActionDelete)
			w, err := ws.Workspace(context.Background(), fixedSecret{val: []byte("dummy")}, cr, ts, cfg)
			if err != nil {
				t.Fatalf("Workspace() = %v", err)
			}
			res, err := w.Refresh(context.Background())
			if err != nil {
				t.Fatalf("Refresh() = %v", err)
			}
			if res.Exists {
				t.Fatal("expected Exists=false for new CR without external-name")
			}
		})

		for name, extName := range cases {
			t.Run(name, func(t *testing.T) {
				resetReqs()
				cr := newCR(prefix+"-ref-"+name, extName, xpv2.ManagementActionObserve, xpv2.ManagementActionCreate, xpv2.ManagementActionUpdate, xpv2.ManagementActionDelete)
				w, err := ws.Workspace(context.Background(), fixedSecret{val: []byte("dummy")}, cr, ts, cfg)
				if err != nil {
					t.Fatalf("Workspace() = %v", err)
				}
				res, err := w.Refresh(context.Background())
				if err != nil || !res.Exists {
					t.Fatalf("Refresh() = (exists=%v, err=%v), reqs=%v", res.Exists, err, snapshotReqs())
				}
				state := assertStateID(t, res.State.GetAttributes(), wantID)
				if _, err := resource.SetCriticalAnnotations(cr, cfg, state, string(res.State.GetPrivateRaw())); err != nil {
					t.Fatalf("SetCriticalAnnotations() = %v", err)
				}
				if got := meta.GetExternalName(cr); got != wantID {
					t.Fatalf("external-name = %q, want %q", got, wantID)
				}

				resetReqs()
				crImp := newCR(prefix+"-imp-"+name, extName, xpv2.ManagementActionObserve)
				wImp, err := ws.Workspace(context.Background(), fixedSecret{val: []byte("dummy")}, crImp, ts, cfg)
				if err != nil {
					t.Fatalf("Import Workspace() = %v", err)
				}
				impRes, err := wImp.Import(context.Background(), crImp)
				if err != nil || !impRes.Exists {
					t.Fatalf("Import() = (exists=%v, err=%v), reqs=%v", impRes.Exists, err, snapshotReqs())
				}
				assertStateID(t, impRes.State.GetAttributes(), wantID)
			})
		}
	}

	t.Run("Store", func(t *testing.T) {
		runCases(t, "st", storeCfg, storeID, newStoreCR, map[string]string{
			"bare_id":  storeID,
			"two_part": accID + "/" + storeID,
		})
	})

	t.Run("StoreSecret", func(t *testing.T) {
		runCases(t, "ss", secretCfg, secretID, newSecretCR, map[string]string{
			"bare_id":    secretID,
			"three_part": accID + "/" + storeID + "/" + secretID,
		})
	})
}
