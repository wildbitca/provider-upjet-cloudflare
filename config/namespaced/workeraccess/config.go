package workeraccess

import (
	"regexp"
	"strings"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/resource"
	corev1 "k8s.io/api/core/v1"
)

// Configure links namespaced Workers and endpoints to their Cloudflare Access application.
func Configure(p *config.Provider) {
	for _, resource := range []string{"cloudflare_workers_script", "cloudflare_workers_script_subdomain"} {
		p.AddResourceConfigurator(resource, func(r *config.Resource) {
			r.References["script_name"] = config.Reference{
				TerraformName: "cloudflare_zero_trust_access_application",
				RefFieldName:  "AccessApplicationRef", SelectorFieldName: "AccessApplicationSelector",
				Extractor: "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/workeraccess.AccessWorkerName()",
			}
		})
	}
}

// AccessWorkerName resolves only a ready policy protecting exactly one workers.dev hostname.
var workerName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)

func AccessWorkerName() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		if mg.GetDeletionTimestamp() != nil || mg.GetCondition(xpv2.TypeReady).Status != corev1.ConditionTrue || mg.GetCondition(xpv2.TypeSynced).Status != corev1.ConditionTrue {
			return ""
		}
		tr, ok := mg.(resource.Terraformed)
		if !ok {
			return ""
		}
		attrs, err := tr.GetObservation()
		if err != nil || !serviceTokenPolicies(attrs) {
			return ""
		}
		domain, _ := attrs["domain"].(string)
		parts := strings.Split(domain, ".")
		if len(parts) != 4 || strings.Join(parts[2:], ".") != "workers.dev" || !workerName.MatchString(parts[0]) {
			return ""
		}
		return parts[0]
	}
}

// serviceTokenPolicies requires observed inline policies that admit only
// explicit service tokens. Empty, bypass, and uninspectable policies fail closed.
func serviceTokenPolicies(attrs map[string]any) bool {
	policies, ok := attrs["policies"].([]any)
	if !ok || len(policies) == 0 {
		return false
	}
	for _, item := range policies {
		policy, ok := item.(map[string]any)
		if !ok || policy["decision"] != "non_identity" {
			return false
		}
		if !serviceTokenIncludes(policy["include"]) {
			return false
		}
	}
	return true
}

func serviceTokenIncludes(value any) bool {
	includes, ok := value.([]any)
	if !ok || len(includes) == 0 {
		return false
	}
	for _, item := range includes {
		rule, ok := item.(map[string]any)
		if !ok || len(rule) != 1 {
			return false
		}
		token, ok := rule["service_token"].(map[string]any)
		if !ok {
			return false
		}
		id, ok := token["token_id"].(string)
		if !ok || id == "" {
			return false
		}
	}
	return true
}
