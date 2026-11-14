package config

import (
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

func TestMCPRequiredIdentifiersReachTerraformConfig(t *testing.T) {
	for scope, provider := range map[string]*ujconfig.Provider{
		"cluster": GetProvider(), "namespaced": GetProviderNamespaced(),
	} {
		for _, name := range []string{
			"cloudflare_zero_trust_access_ai_controls_mcp_server",
			"cloudflare_zero_trust_access_ai_controls_mcp_portal",
		} {
			t.Run(scope+"/"+name, func(t *testing.T) {
				cfg := provider.Resources[name].ExternalName
				if cfg.DisableNameInitializer {
					t.Fatal("MCP identifiers must default to the Kubernetes resource name")
				}
				params := map[string]interface{}{}
				cfg.SetIdentifierArgumentFn(params, "mcp-test-id")
				if params["id"] != "mcp-test-id" {
					t.Fatalf("required Terraform id is missing from parameters: %#v", params)
				}
			})
		}
	}
}
