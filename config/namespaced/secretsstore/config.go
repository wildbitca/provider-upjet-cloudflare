package secretsstore

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure adds the secretsstore group resource configurators.
//
// Groups the Cloudflare Secrets Store resources:
//   - cloudflare_secrets_store        -> Store (account-scoped secrets store)
//   - cloudflare_secrets_store_secret -> Secret (a secret within a store)
//
// The secret's store_id references the Store resource so that the store can be
// resolved via a Crossplane reference/selector.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_secrets_store", func(r *config.Resource) {
		r.ShortGroup = "secretsstore"
		r.References["account_id"] = config.Reference{
			TerraformName: "cloudflare_account",
		}
	})
	p.AddResourceConfigurator("cloudflare_secrets_store_secret", func(r *config.Resource) {
		r.ShortGroup = "secretsstore"
		r.References["account_id"] = config.Reference{
			TerraformName: "cloudflare_account",
		}
		r.References["store_id"] = config.Reference{
			TerraformName: "cloudflare_secrets_store",
		}
	})
}
