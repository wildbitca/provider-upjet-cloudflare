package config_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dnsv1alpha1 "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/dns/v1alpha1"
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
