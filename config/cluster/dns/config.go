package dns

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds the dns group resource configurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_dns_record", func(r *config.Resource) {
		r.ShortGroup = "dns"
		r.References["zone_id"] = config.Reference{
			TerraformName: "cloudflare_zone",
		}
	})
}
