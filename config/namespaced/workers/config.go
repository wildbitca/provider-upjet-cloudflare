package workers

import "github.com/crossplane/upjet/v2/pkg/config"

func Configure(p *config.Provider) {
	p.AddResourceConfigurator("cloudflare_workers_script", func(r *config.Resource) {
		r.References["bindings.service"] = config.Reference{TerraformName: "cloudflare_workers_script"}
	})
}
