package workerartifactgen

import (
	"bytes"
	"go/format"
	"strings"
	"testing"
	"text/template"

	"github.com/crossplane/upjet/v2/pkg/pipeline/templates"
)

func render(t *testing.T, source string, values map[string]any) []byte {
	t.Helper()
	tmpl, err := template.New("test").Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, values); err != nil {
		t.Fatal(err)
	}
	formatted, err := format.Source(out.Bytes())
	if err != nil {
		t.Fatalf("invalid generated Go: %v\n%s", err, out.String())
	}
	return formatted
}

func TestReferenceResolverIsSelectedOnlyForScriptControllers(t *testing.T) {
	apiOriginal, tfOriginal, controllerOriginal := templates.CRDTypesTemplate, templates.TerraformedTemplate, templates.ControllerTemplate
	t.Cleanup(func() {
		templates.CRDTypesTemplate = apiOriginal
		templates.TerraformedTemplate = tfOriginal
		templates.ControllerTemplate = controllerOriginal
	})
	Configure()
	values := map[string]any{
		"Header": "// fixture", "GenStatement": "// generated", "Package": "test", "Imports": "",
		"CRD": map[string]string{"Kind": "Example"}, "TypePackageAlias": "v1alpha1.", "FeaturesPackageAlias": "features.",
	}
	for _, resource := range []string{"cloudflare_dns_record", "cloudflare_worker_version", "cloudflare_r2_bucket"} {
		values["ResourceType"] = resource
		if !bytes.Equal(render(t, controllerOriginal, values), render(t, templates.ControllerTemplate, values)) {
			t.Fatalf("changed unrelated controller %s", resource)
		}
	}
	values["ResourceType"] = "cloudflare_workers_script"
	if !strings.Contains(string(render(t, templates.ControllerTemplate, values)), "managed.WithReferenceResolver(workerartifactrefs.NewReferenceResolver(mgr.GetClient()))") {
		t.Fatal("Script controller does not select the artifact-aware resolver")
	}
}

func TestUnrelatedResourcesGenerateIdentically(t *testing.T) {
	apiOriginal, tfOriginal, controllerOriginal := templates.CRDTypesTemplate, templates.TerraformedTemplate, templates.ControllerTemplate
	t.Cleanup(func() {
		templates.CRDTypesTemplate = apiOriginal
		templates.TerraformedTemplate = tfOriginal
		templates.ControllerTemplate = controllerOriginal
	})
	Configure()
	values := map[string]any{
		"Header": "// Test fixture", "GenStatement": "// Code generated.", "Imports": "", "Types": "",
		"CRD":                        map[string]string{"Kind": "Record", "APIVersion": "v1alpha1", "Group": "dns.upjet-cloudflare.m.upbound.io", "Scope": "Namespaced", "ForProviderType": "RecordParameters", "InitProviderType": "RecordInitParameters", "AtProviderType": "RecordObservation", "ParametersTypeName": "RecordParameters"},
		"XPV2CommonAPIsPackageAlias": "v2.", "XPCommonAPIsPackageAlias": "v2.", "Provider": map[string]string{"ShortName": "cloudflare"},
		"APIVersion": "v1alpha1", "Terraform": map[string]any{"ResourceType": "cloudflare_dns_record", "ResourceSchema": map[string]int{"SchemaVersion": 0}},
		"Sensitive": map[string]any{}, "LateInitializer": map[string]any{"IgnoredFields": []string{}, "ConditionalIgnoredFields": []string{}},
	}
	values["CRD"].(map[string]string)["ValidationRules"] = ""
	if !bytes.Equal(render(t, apiOriginal, values), render(t, templates.CRDTypesTemplate, values)) {
		t.Fatal("Worker extension changes unrelated resource API")
	}
	if !bytes.Equal(render(t, tfOriginal, values), render(t, templates.TerraformedTemplate, values)) {
		t.Fatal("Worker extension changes unrelated Terraform adapter")
	}
	values["CRD"].(map[string]string)["Kind"] = "Version"
	values["CRD"].(map[string]string)["Group"] = "worker.upjet-cloudflare.m.upbound.io"
	values["Terraform"].(map[string]any)["ResourceType"] = "cloudflare_worker_version"
	if !bytes.Equal(render(t, apiOriginal, values), render(t, templates.CRDTypesTemplate, values)) ||
		!bytes.Equal(render(t, tfOriginal, values), render(t, templates.TerraformedTemplate, values)) {
		t.Fatal("Script artifact extension changes the unsupported Version resource")
	}
}
