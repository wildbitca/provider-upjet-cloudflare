package address

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds the address group resource configurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_address_map", func(r *config.Resource) {
		r.ShortGroup = "address"
		r.Kind = "AddressMap" // Avoid package name "map" (Go reserved keyword)
	})
}
