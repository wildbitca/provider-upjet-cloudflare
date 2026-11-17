// Package preview supplies references and connection outputs for private Worker previews.
package preview

import (
	"encoding/json"

	"github.com/crossplane/upjet/v2/pkg/config"
)

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_connectivity_directory_service", func(r *config.Resource) {
		r.References["host.network.tunnel_id"] = config.Reference{TerraformName: "cloudflare_zero_trust_tunnel_cloudflared"}
		r.References["host.resolver_network.tunnel_id"] = config.Reference{TerraformName: "cloudflare_zero_trust_tunnel_cloudflared"}
	})
	p.AddResourceConfigurator("cloudflare_workers_script", func(r *config.Resource) {
		r.References["bindings.service_id"] = config.Reference{TerraformName: "cloudflare_connectivity_directory_service"}
	})
	p.AddResourceConfigurator("cloudflare_zero_trust_access_application", func(r *config.Resource) {
		r.References["policies.include.service_token.token_id"] = config.Reference{TerraformName: "cloudflare_zero_trust_access_service_token"}
	})
	p.AddResourceConfigurator("cloudflare_zero_trust_tunnel_cloudflared", func(r *config.Resource) {
		r.Sensitive.AdditionalConnectionDetailsFn = TunnelCredentials
	})
	p.AddResourceConfigurator("cloudflare_zero_trust_access_service_token", func(r *config.Resource) {
		r.Sensitive.AdditionalConnectionDetailsFn = func(attrs map[string]any) (map[string][]byte, error) {
			out := map[string][]byte{}
			for _, key := range []string{"client_id", "client_secret"} {
				if value, ok := attrs[key].(string); ok && value != "" {
					out[key] = []byte(value)
				}
			}
			return out, nil
		}
	})
}

// TunnelCredentials exports the native cloudflared credential file through its connection Secret.
func TunnelCredentials(attrs map[string]any) (map[string][]byte, error) {
	values := map[string]string{}
	for field, key := range map[string]string{"account_id": "AccountTag", "id": "TunnelID", "tunnel_secret": "TunnelSecret"} {
		value, _ := attrs[field].(string)
		if value == "" {
			return nil, nil
		}
		values[key] = value
	}
	data, err := json.Marshal(values)
	return map[string][]byte{"credentials.json": data, "tunnel_id": []byte(values["TunnelID"])}, err
}
