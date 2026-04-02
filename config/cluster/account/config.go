package account

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds the account group resource configurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_account", func(r *config.Resource) {
		r.ShortGroup = "account"
	})
	// Account-scoped resources with account_id reference to cloudflare_account.
	for _, name := range []string{
		"cloudflare_account_member", "cloudflare_account_subscription",
		"cloudflare_r2_bucket", "cloudflare_r2_bucket_cors", "cloudflare_r2_bucket_lifecycle",
		"cloudflare_workers_kv_namespace", "cloudflare_worker", "cloudflare_pages_project",
		"cloudflare_zero_trust_access_application", "cloudflare_zero_trust_access_group",
		"cloudflare_zero_trust_access_policy", "cloudflare_zero_trust_access_identity_provider",
		"cloudflare_api_token", "cloudflare_stream",
	} {
		resName := name
		p.AddResourceConfigurator(resName, func(r *config.Resource) {
			r.References["account_id"] = config.Reference{
				TerraformName: "cloudflare_account",
			}
		})
	}
}
