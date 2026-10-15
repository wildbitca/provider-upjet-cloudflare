package zone

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/wildbitca/provider-upjet-cloudflare/config/zonesetting"
)

// Configure adds the zone group resource configurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_zone", func(r *config.Resource) {
		r.ShortGroup = "zone"
	})
	// Zone-scoped resources with zone_id reference to cloudflare_zone.
	for _, resName := range []string{
		"cloudflare_filter", "cloudflare_firewall_rule", "cloudflare_load_balancer",
		"cloudflare_load_balancer_monitor", "cloudflare_load_balancer_pool",
		"cloudflare_zone_setting", "cloudflare_zone_dnssec", "cloudflare_zone_cache_reserve",
		"cloudflare_zone_cache_variants", "cloudflare_zone_dns_settings",
		"cloudflare_zone_hold", "cloudflare_zone_lockdown", "cloudflare_zone_subscription",
		"cloudflare_healthcheck", "cloudflare_rate_limit", "cloudflare_page_rule",
		"cloudflare_custom_pages", "cloudflare_custom_ssl", "cloudflare_dns_firewall",
	} {
		name := resName
		p.AddResourceConfigurator(name, func(r *config.Resource) {
			r.References["zone_id"] = config.Reference{
				TerraformName: "cloudflare_zone",
			}
		})
	}
	// LateInitializer for resources where Cloudflare API returns default/computed values to avoid drift loops.
	p.AddResourceConfigurator("cloudflare_zone_setting", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"modified_on"},
		}
		// Canonicalize object-typed setting values (aegis,
		// automatic_platform_optimization, nel, security_header) to their
		// full known key set before upjet's auto-registered dynamic-value
		// conversion infers a Terraform type from forProvider.value, so
		// that Update can converge for these setting_ids. Must be
		// prepended, not appended: upjet's NewProvider() already appended
		// the generic dynamic-value conversion to this resource's
		// TerraformConversions before this configurator runs, and
		// execution order in that slice is the same for both ToTerraform
		// and FromTerraform. See config/zonesetting for the full
		// explanation and sources.
		r.TerraformConversions = append([]config.TerraformConversion{zonesetting.New()}, r.TerraformConversions...)
	})
	p.AddResourceConfigurator("cloudflare_zone_dnssec", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"modified_on"},
		}
	})
}
