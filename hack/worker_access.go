package hack

import (
	"strings"

	"github.com/crossplane/upjet/v2/pkg/pipeline/templates"
)

// WorkerAccessControllerTemplate adds an uncached Access check to namespaced
// Worker controllers before they connect to Cloudflare.
func WorkerAccessControllerTemplate() string {
	const worker = `or (eq .ResourceType "cloudflare_workers_script") (eq .ResourceType "cloudflare_workers_script_subdomain")`
	const imports = "\t{{ .Imports }}"
	const initializers = "\tvar initializers managed.InitializerChain"

	template := templates.ControllerTemplate
	if strings.Count(template, imports) != 1 || strings.Count(template, initializers) != 1 {
		panic("unexpected Upjet controller template")
	}
	template = strings.Replace(template, imports, imports+`{{ if `+worker+` }}
 workeraccess "github.com/wildbitca/provider-upjet-cloudflare/internal/workeraccess"{{ end }}`, 1)
	return strings.Replace(template, initializers, initializers+`{{ if `+worker+` }}
	initializers = append(initializers, workeraccess.NewInitializer(mgr.GetAPIReader())){{ end }}`, 1)
}
