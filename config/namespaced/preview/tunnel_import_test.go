package preview

import "testing"

func TestCloudflaredTunnelImportID(t *testing.T) {
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
			got, err := cloudflaredTunnelImportID(t.Context(), test.externalName, map[string]any{"account_id": test.accountID}, nil)
			if (err != nil) != test.wantError || got != test.want {
				t.Fatalf("cloudflaredTunnelImportID() = %q, %v, want %q, error = %v", got, err, test.want, test.wantError)
			}
		})
	}
}
