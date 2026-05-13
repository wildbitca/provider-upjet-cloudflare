package zero

import ujconfig "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds Zero Trust resource configurators.
func Configure(p *ujconfig.Provider) {
	// CF GET for self-hosted apps always returns both destinations and
	// self_hosted_domains (they mirror each other). Late-init must ignore both
	// so upjet doesn't write them into spec alongside whichever field the user
	// chose, which would cause TF ConflictsWith on every refresh.
	// The user picks exactly one of the two fields in their spec; the other
	// must be absent so TF only sees a single field in main.tf.json.
	p.AddResourceConfigurator("cloudflare_zero_trust_access_application", func(r *ujconfig.Resource) {
		r.LateInitializer = ujconfig.LateInitializer{
			IgnoredFields: []string{"destinations", "self_hosted_domains"},
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
