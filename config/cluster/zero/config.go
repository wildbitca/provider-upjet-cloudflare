package zero

import ujconfig "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds Zero Trust resource configurators.
func Configure(p *ujconfig.Provider) {
	// CF GET for self-hosted apps returns both self_hosted_domains and
	// destinations (backward compat). Late-init would write both to spec,
	// causing TF ConflictsWith on the next reconcile. Ignore the CF-internal
	// destinations field so only self_hosted_domains is managed via spec.
	p.AddResourceConfigurator("cloudflare_zero_trust_access_application", func(r *ujconfig.Resource) {
		r.LateInitializer = ujconfig.LateInitializer{
			IgnoredFields: []string{"destinations"},
		}
	})

	// CF GET for any IdP type includes scim_config in the response, but
	// the TF schema disallows scim_config for type=onetimepin. Late-init
	// would write scim_config to spec, causing ConflictsWith on next reconcile.
	p.AddResourceConfigurator("cloudflare_zero_trust_access_identity_provider", func(r *ujconfig.Resource) {
		r.LateInitializer = ujconfig.LateInitializer{
			IgnoredFields: []string{"scim_config"},
		}
	})
}
