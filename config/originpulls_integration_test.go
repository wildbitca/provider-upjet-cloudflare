package config

// Integration test for the OriginPulls refresh. It drives the REAL upjet
// workspace code path (terraform.WorkspaceStore -> FileProducer -> tofu) against
// the REAL terraform-provider-cloudflare binary, with the Cloudflare HTTP API
// replaced by a local fake (provider `base_url`). Only GET requests are served,
// so nothing can write.
//
// Run:
//   ORIGINPULLS_PROVIDER_BIN=/path/terraform-provider-cloudflare_v5.28.1 \
//     go test ./config -run TestOriginPullsRefresh -v -count=1
// Requires `tofu` on PATH. Skipped when ORIGINPULLS_PROVIDER_BIN is unset.

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
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

	v1alpha1 "github.com/wildbitca/provider-upjet-cloudflare/apis/cluster/authenticated/v1alpha1"
)

const opZone = "023e105f4ecef8ad9ca31a8372d0c353"
const opHost = "app.example.com"

type noopSecrets struct{}

func (noopSecrets) GetSecretData(context.Context, *xpv2.SecretReference) (map[string][]byte, error) {
	return nil, nil
}

func (noopSecrets) GetSecretValue(context.Context, xpv2.SecretKeySelector) ([]byte, error) {
	return nil, nil
}

func newOriginPulls(uid string, existing bool, extName string) *v1alpha1.OriginPulls {
	cr := &v1alpha1.OriginPulls{}
	cr.SetName("op-" + uid)
	cr.SetUID(types.UID(uid))
	cr.SetManagementPolicies(xpv2.ManagementPolicies{xpv2.ManagementActionObserve, xpv2.ManagementActionCreate, xpv2.ManagementActionUpdate, xpv2.ManagementActionDelete})
	h, z, c, e := opHost, opZone, "2458ce5a-0c35-4c7f-82c7-8e9487d3ff60", true
	cr.Spec.ForProvider.Config = []v1alpha1.ConfigParameters{{Hostname: &h, CertID: &c, Enabled: &e}}
	cr.Spec.ForProvider.ZoneID = &z
	if extName != "" {
		meta.SetExternalName(cr, extName)
	}
	if existing {
		// What a CR that already reconciled looks like (v1.0.0 era).
		meta.SetExternalName(cr, opHost)
		cr.Status.AtProvider.Hostname = &h
		cr.Status.AtProvider.ID = &h
		cr.Status.AtProvider.Enabled = &e
		cr.Status.AtProvider.CertID = &c
	}
	cr.CreationTimestamp = metav1.Now()
	return cr
}

func fakeCloudflare(t *testing.T) (*httptest.Server, *[]string) {
	var mu sync.Mutex
	var reqs []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		reqs = append(reqs, r.Method+" "+r.URL.Path)
		mu.Unlock()
		prefix := "/zones/" + opZone + "/origin_tls_client_auth/hostnames/"
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, prefix) && strings.TrimPrefix(r.URL.Path, prefix) != "" {
			h := strings.TrimPrefix(r.URL.Path, prefix)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success":true,"errors":[],"messages":[],"result":{"hostname":%q,"cert_id":"2458ce5a-0c35-4c7f-82c7-8e9487d3ff60","enabled":true,"status":"active","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z","expires_on":"2030-01-01T00:00:00Z","issuer":"x","signature":"x","serial_number":"1","certificate":"x"}}`, h)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"success":false,"errors":[{"code":1,"message":"unexpected %s %s"}]}`, r.Method, r.URL.Path)
	}))
	t.Cleanup(srv.Close)
	return srv, &reqs
}

func setupTofuEnv(t *testing.T, providerBin string) {
	dir := t.TempDir()
	// upjet runs a binary named "terraform"; point it at tofu.
	tofu, err := exec.LookPath("tofu")
	if err != nil {
		t.Skip("tofu not on PATH")
	}
	bin := filepath.Join(dir, "bin")
	_ = os.MkdirAll(bin, 0o755)
	_ = os.Symlink(tofu, filepath.Join(bin, "terraform"))
	plug := filepath.Join(dir, "plugin")
	_ = os.MkdirAll(plug, 0o755)
	_ = os.Symlink(providerBin, filepath.Join(plug, "terraform-provider-cloudflare"))
	rc := filepath.Join(dir, "tofurc")
	body := fmt.Sprintf("provider_installation {\n  dev_overrides {\n    \"registry.opentofu.org/cloudflare/cloudflare\" = %q\n    \"registry.terraform.io/cloudflare/cloudflare\" = %q\n    \"cloudflare/cloudflare\" = %q\n  }\n  direct {}\n}\n", plug, plug, plug)
	if err := os.WriteFile(rc, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("TF_CLI_CONFIG_FILE", rc)
	t.Setenv("TOFU_CLI_CONFIG_FILE", rc)
	t.Setenv("TF_IN_AUTOMATION", "1")
	t.Setenv("TMPDIR", dir)
}

// runRefresh mimics one reconcile of the generated controller: run the
// resource initializers, build the workspace (writes main.tf.json + prior
// tfstate) and refresh. `seed` additionally reproduces the v1.0.3 behaviour.
func runRefresh(t *testing.T, base string, cr *v1alpha1.OriginPulls, useInitializer, v103Seed bool) (string, error) {
	p := GetProvider()
	cfg := p.Resources[originPullsResource]
	if cfg == nil {
		t.Fatal("resource not configured")
	}
	if v103Seed {
		cfg.ExternalName.SetIdentifierArgumentFn = func(b map[string]any, _ string) {
			if h := originPullsConfigHostname(b); h != "" {
				b["hostname"] = h
			}
		}
	}
	if useInitializer {
		if len(cfg.InitializerFns) == 0 {
			t.Fatal("OriginPulls has no initializer wired")
		}
		for _, fn := range cfg.InitializerFns {
			if err := fn(nil).Initialize(context.Background(), cr); err != nil {
				t.Fatal(err)
			}
		}
	}
	ws := terraform.NewWorkspaceStore(logging.NewNopLogger())
	ts := terraform.Setup{
		Version:       "1.5.7",
		Requirement:   terraform.ProviderRequirement{Source: "cloudflare/cloudflare", Version: "5.28.1"},
		Configuration: map[string]any{"api_token": strings.Repeat("a", 40), "base_url": base},
	}
	w, err := ws.Workspace(context.Background(), noopSecrets{}, cr, ts, cfg)
	if err != nil {
		return "", err
	}
	res, err := w.Refresh(context.Background())
	if err != nil {
		return "", err
	}
	if !res.Exists {
		return "", fmt.Errorf("refresh reported resource does not exist")
	}
	return string(res.State.GetAttributes()), nil
}

func TestOriginPullsRefresh(t *testing.T) {
	bin := os.Getenv("ORIGINPULLS_PROVIDER_BIN")
	if bin == "" {
		t.Skip("ORIGINPULLS_PROVIDER_BIN not set")
	}
	setupTofuEnv(t, bin)
	srv, reqs := fakeCloudflare(t)

	cases := []struct {
		name           string
		existing       bool
		extName        string
		initializer    bool
		v103           bool
		wantErrContain string // "" => must succeed
	}{
		// Negative controls: prove the harness enforces the real provider rules.
		{name: "control/new CR, no fix => missing hostname", existing: false, initializer: false, wantErrContain: "missing required hostname parameter"},
		{name: "control/existing CR, v1.0.3 seed => read-only attribute", existing: true, initializer: false, v103: true, wantErrContain: "Read-Only Attribute"},
		// Subject under test.
		{name: "fix/new CR (no status, no external-name)", existing: false, initializer: true},
		{name: "fix/new CR (external-name = hostname, no status)", existing: false, extName: opHost, initializer: true},
		{name: "fix/existing CR (populated atProvider)", existing: true, initializer: true},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cr := newOriginPulls(fmt.Sprintf("uid-%d", i), tc.existing, tc.extName)
			state, err := runRefresh(t, srv.URL, cr, tc.initializer, tc.v103)
			if tc.wantErrContain != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrContain) {
					t.Fatalf("want error containing %q, got err=%v state=%s", tc.wantErrContain, err, state)
				}
				t.Logf("expected failure reproduced: %.300s", err.Error())
				return
			}
			if err != nil {
				t.Fatalf("refresh failed: %v", err)
			}
			if !strings.Contains(state, `"hostname":"`+opHost+`"`) {
				t.Fatalf("refreshed state lacks hostname: %s", state)
			}
			// a refresh must never send anything but GET
			for _, r := range *reqs {
				if !strings.HasPrefix(r, "GET ") {
					t.Fatalf("non-GET request issued: %s", r)
				}
			}
		})
	}
	var _ resource.Terraformed = (*v1alpha1.OriginPulls)(nil)
}
