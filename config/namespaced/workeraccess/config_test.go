package workeraccess_test

import (
	"testing"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	zero "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/zero/v1alpha1"
	"github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/workeraccess"
)

func TestWorkerEndpointRequiresReadyAccessPolicy(t *testing.T) {
	app := &zero.TrustAccessApplication{}
	if err := app.SetObservation(map[string]any{"domain": "example-worker.example-account.workers.dev", "policies": []any{map[string]any{"decision": "non_identity", "include": []any{map[string]any{"service_token": map[string]any{"token_id": "relay"}}}}}}); err != nil {
		t.Fatal(err)
	}
	extract := workeraccess.AccessWorkerName()
	if got := extract(app); got != "" {
		t.Fatal("unready Access application enabled a Worker endpoint")
	}
	app.SetConditions(xpv2.Available(), xpv2.ReconcileSuccess())
	if got := extract(app); got != "example-worker" {
		t.Fatalf("wrong Worker name: %q", got)
	}
	for _, domain := range []string{"*.example-account.workers.dev", "gateway.example.com", "gateway.example-account.workers.dev/path"} {
		if err := app.SetObservation(map[string]any{"domain": domain, "policies": []any{map[string]any{"decision": "non_identity", "include": []any{map[string]any{"service_token": map[string]any{"token_id": "relay"}}}}}}); err != nil {
			t.Fatal(err)
		}
		if got := extract(app); got != "" {
			t.Fatalf("accepted an unprotected or ambiguous hostname %q", domain)
		}
	}
}
