package config

import (
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

func TestCloudflaredTunnelImportID(t *testing.T) {
	resource := &ujconfig.Resource{Name: "cloudflare_zero_trust_tunnel_cloudflared"}
	ExternalNameConfigurations()(resource)
	for _, test := range []struct {
		name, externalName, accountID, want string
		wantError                           bool
	}{
		{name: "bare tunnel ID", externalName: "tunnel-id", accountID: "account-id", want: "account-id/tunnel-id"},
		{name: "qualified ID", externalName: "account-id/tunnel-id", want: "account-id/tunnel-id"},
		{name: "before creation"},
		{name: "missing account", externalName: "tunnel-id", wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := resource.ExternalName.GetImportIDFn(t.Context(), test.externalName, map[string]any{"account_id": test.accountID}, nil)
			if (err != nil) != test.wantError || got != test.want {
				t.Fatalf("import ID = %q, %v, want %q, error = %v", got, err, test.want, test.wantError)
			}
		})
	}
	got, err := resource.ExternalName.GetIDFn(t.Context(), "tunnel-id", map[string]any{"account_id": "account-id"}, nil)
	if err != nil || got != "tunnel-id" {
		t.Fatalf("state ID = %q, %v, want the tunnel ID", got, err)
	}
}
