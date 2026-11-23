package workerconnectivity_test

import (
	"encoding/json"
	"testing"

	"github.com/wildbitca/provider-upjet-cloudflare/config"
	"github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/workerconnectivity"
)

func TestTunnelConnectionSecret(t *testing.T) {
	attrs := map[string]any{"account_id": "account", "id": "tunnel", "tunnel_secret": "synthetic-secret"}
	out, err := workerconnectivity.TunnelCredentials(attrs)
	if err != nil {
		t.Fatal(err)
	}
	var credentials map[string]string
	if err := json.Unmarshal(out["credentials.json"], &credentials); err != nil {
		t.Fatal(err)
	}
	if credentials["AccountTag"] != "account" || credentials["TunnelID"] != "tunnel" || credentials["TunnelSecret"] != "synthetic-secret" {
		t.Fatal("incomplete cloudflared credentials")
	}
	delete(attrs, "tunnel_secret")
	out, err = workerconnectivity.TunnelCredentials(attrs)
	if err != nil || len(out) != 0 {
		t.Fatal("published partial tunnel credentials")
	}
}

func TestPrivateWorkerReferenceAndConnectionConfiguration(t *testing.T) {
	resources := config.GetProviderNamespaced().Resources
	for _, test := range []struct{ resource, field, target string }{
		{"cloudflare_connectivity_directory_service", "host.network.tunnel_id", "cloudflare_zero_trust_tunnel_cloudflared"},
		{"cloudflare_connectivity_directory_service", "host.resolver_network.tunnel_id", "cloudflare_zero_trust_tunnel_cloudflared"},
		{"cloudflare_workers_script", "bindings.service_id", "cloudflare_connectivity_directory_service"},
		{"cloudflare_zero_trust_access_application", "policies.include.service_token.token_id", "cloudflare_zero_trust_access_service_token"},
	} {
		if got := resources[test.resource].References[test.field].TerraformName; got != test.target {
			t.Errorf("%s.%s refers to %q, want %q", test.resource, test.field, got, test.target)
		}
	}
	tunnel := resources["cloudflare_zero_trust_tunnel_cloudflared"].Sensitive.AdditionalConnectionDetailsFn
	if tunnel == nil {
		t.Fatal("tunnel connection credentials are not configured")
	}
	token := resources["cloudflare_zero_trust_access_service_token"].Sensitive.AdditionalConnectionDetailsFn
	if token == nil {
		t.Fatal("service token connection credentials are not configured")
	}
	details, err := token(map[string]any{"client_id": "test-client", "client_secret": "test-value", "id": "test-token-id"})
	if err != nil || len(details) != 2 || string(details["client_id"]) != "test-client" || string(details["client_secret"]) != "test-value" {
		t.Fatal("service token did not export the expected connection keys")
	}
}
