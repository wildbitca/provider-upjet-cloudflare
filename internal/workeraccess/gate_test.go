package workeraccess

import (
	"encoding/json"
	"testing"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	workers "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
	zero "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/zero/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestEndpointRevalidatesCachedAccessReference(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := zero.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	var endpoint workers.ScriptSubdomain
	if err := json.Unmarshal([]byte(`{"metadata":{"namespace":"preview"},"spec":{"forProvider":{"accountId":"account-a","enabled":true,"scriptName":"worker","accessApplicationRef":{"name":"access"}}}}`), &endpoint); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name, account, domain, policy string
		missing, wantError            bool
	}{
		{name: "matching", account: "account-a", domain: "worker.example.workers.dev", policy: "non_identity"},
		{name: "wrong account", account: "account-b", domain: "worker.example.workers.dev", policy: "non_identity", wantError: true},
		{name: "deleted application", missing: true, wantError: true},
		{name: "bypass policy", account: "account-a", domain: "worker.example.workers.dev", policy: "bypass", wantError: true},
		{name: "different Worker", account: "account-a", domain: "other.example.workers.dev", policy: "non_identity", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(scheme)
			if !tc.missing {
				app := &zero.TrustAccessApplication{ObjectMeta: metav1.ObjectMeta{Name: "access", Namespace: "preview"}}
				app.SetConditions(xpv2.Available(), xpv2.ReconcileSuccess())
				rule := map[string]any{"service_token": map[string]any{"token_id": "relay"}}
				if tc.policy == "bypass" {
					rule = map[string]any{"everyone": map[string]any{}}
				}
				if err := app.SetObservation(map[string]any{
					"account_id": tc.account, "domain": tc.domain,
					"policies": []any{map[string]any{"decision": tc.policy, "include": []any{rule}}},
				}); err != nil {
					t.Fatal(err)
				}
				builder.WithObjects(app)
			}
			err := NewInitializer(builder.Build()).Initialize(t.Context(), endpoint.DeepCopy())
			if (err != nil) != tc.wantError {
				t.Fatalf("Initialize() = %v, want error = %v", err, tc.wantError)
			}
		})
	}
}
