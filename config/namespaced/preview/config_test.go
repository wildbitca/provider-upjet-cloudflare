package preview_test

import (
	"encoding/json"
	"testing"

	preview "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/preview"
)

func TestTunnelConnectionSecret(t *testing.T) {
	attrs := map[string]any{"account_id": "account", "id": "tunnel", "tunnel_secret": "synthetic-secret"}
	out, err := preview.TunnelCredentials(attrs)
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
	out, err = preview.TunnelCredentials(attrs)
	if err != nil || len(out) != 0 {
		t.Fatal("published partial tunnel credentials")
	}
}
