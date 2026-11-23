package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
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
// that StoreSecret Refresh/Import work whether external-name is the bare
// secret ID or the legacy 3-part <account_id>/<store_id>/<secret_id> path,
// and that a new resource without an external name refreshes as missing.
func TestSecretsStoreRefresh(t *testing.T) {
	bin := os.Getenv("ORIGINPULLS_PROVIDER_BIN")
	if bin == "" {
		t.Skip("ORIGINPULLS_PROVIDER_BIN not set")
	}
	setupTofuEnv(t, bin)

	const (
		accID    = "00000000000000000000000000000001"
		storeID  = "00000000000000000000000000000002"
		secretID = "00000000000000000000000000000003"
	)
	wantPath := fmt.Sprintf("/accounts/%s/secrets_store/stores/%s/secrets/%s", accID, storeID, secretID)

	var mu sync.Mutex
	var reqs []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		reqs = append(reqs, r.Method+" "+r.URL.Path)
		mu.Unlock()
		if r.Method == http.MethodGet && r.URL.Path == wantPath {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success":true,"errors":[],"messages":[],"result":{"id":%q,"account_id":%q,"store_id":%q,"name":"SZ_DEV_TEST","scopes":["workers"],"status":"active","created":"2026-07-13T17:49:59Z","modified":"2026-10-27T01:34:50Z"}}`, secretID, accID, storeID)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"success":false,"errors":[{"code":404,"message":"not found %s %s"}]}`, r.Method, r.URL.Path)
	}))
	t.Cleanup(srv.Close)

	p := GetProvider()
	cfg := p.Resources["cloudflare_secrets_store_secret"]
	if cfg == nil {
		t.Fatal("cloudflare_secrets_store_secret not configured")
	}

	newCR := func(uid, extName string, policies ...xpv2.ManagementAction) *ssv1alpha1.StoreSecret {
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

	t.Run("new CR without external-name refreshes as missing", func(t *testing.T) {
		cr := newCR("ss-new", "", xpv2.ManagementActionObserve, xpv2.ManagementActionCreate, xpv2.ManagementActionUpdate, xpv2.ManagementActionDelete)
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

	for name, extName := range map[string]string{
		"bare_id":    secretID,
		"three_part": accID + "/" + storeID + "/" + secretID,
	} {
		t.Run(name, func(t *testing.T) {
			mu.Lock()
			reqs = nil
			mu.Unlock()

			cr := newCR("ss-ref-"+name, extName, xpv2.ManagementActionObserve, xpv2.ManagementActionCreate, xpv2.ManagementActionUpdate, xpv2.ManagementActionDelete)
			w, err := ws.Workspace(context.Background(), fixedSecret{val: []byte("dummy")}, cr, ts, cfg)
			if err != nil {
				t.Fatalf("Workspace() = %v", err)
			}
			res, err := w.Refresh(context.Background())
			if err != nil || !res.Exists {
				t.Fatalf("Refresh() = (exists=%v, err=%v), reqs=%v", res.Exists, err, reqs)
			}
			var state map[string]any
			if err := json.Unmarshal(res.State.GetAttributes(), &state); err != nil {
				t.Fatalf("unmarshal state: %v", err)
			}
			if state["id"] != secretID {
				t.Fatalf("state id = %v, want %q", state["id"], secretID)
			}
			if _, err := resource.SetCriticalAnnotations(cr, cfg, state, string(res.State.GetPrivateRaw())); err != nil {
				t.Fatalf("SetCriticalAnnotations() = %v", err)
			}
			if got := meta.GetExternalName(cr); got != secretID {
				t.Fatalf("external-name = %q, want %q", got, secretID)
			}

			mu.Lock()
			reqs = nil
			mu.Unlock()
			crImp := newCR("ss-imp-"+name, extName, xpv2.ManagementActionObserve)
			wImp, err := ws.Workspace(context.Background(), fixedSecret{val: []byte("dummy")}, crImp, ts, cfg)
			if err != nil {
				t.Fatalf("Import Workspace() = %v", err)
			}
			impRes, err := wImp.Import(context.Background(), crImp)
			if err != nil || !impRes.Exists {
				t.Fatalf("Import() = (exists=%v, err=%v), reqs=%v", impRes.Exists, err, reqs)
			}
		})
	}
}
