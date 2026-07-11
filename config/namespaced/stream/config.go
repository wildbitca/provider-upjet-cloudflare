package stream

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds the stream group resource configurators.
//
// cloudflare_moq_relay (Media over QUIC relay, new in cloudflare provider
// 5.22.0) is an account-scoped media resource. It is homed in the existing
// "stream" media sub-provider rather than a new standalone group so no new
// family package / cluster wiring is required. An explicit Kind preserves the
// "MOQ" identity, since the default kind would strip the leading "moq"
// segment and yield an ambiguous "Relay".
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_moq_relay", func(r *config.Resource) {
		r.ShortGroup = "stream"
		r.Kind = "MOQRelay"
		r.References["account_id"] = config.Reference{
			TerraformName: "cloudflare_account",
		}
	})
}
