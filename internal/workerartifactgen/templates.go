// Package workerartifactgen extends Upjet templates only for Worker uploads.
// Keeping the extension here makes regeneration reproducible without patching Upjet.
package workerartifactgen

import (
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/pipeline/templates"
)

const apiCondition = `and (eq .CRD.Kind "Script") (or (eq .CRD.Group "workers.upjet-cloudflare.upbound.io") (eq .CRD.Group "workers.upjet-cloudflare.m.upbound.io"))`
const tfCondition = `eq .Terraform.ResourceType "cloudflare_workers_script"`
const artifactImport = `workerartifact "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"`

// Configure extends the templates at generation time, not the runtime provider.
func Configure() {
	controllerCondition := `eq .ResourceType "cloudflare_workers_script"`
	templates.ControllerTemplate = replace(templates.ControllerTemplate, "\t{{ .Imports }}", "\t{{ .Imports }}{{ if "+controllerCondition+" }}\n\tworkerartifactrefs \"github.com/wildbitca/provider-upjet-cloudflare/internal/workerartifact\"{{ end }}")
	templates.ControllerTemplate = replace(templates.ControllerTemplate, "\t\tmanaged.WithInitializers(initializers),", "\t\tmanaged.WithInitializers(initializers),{{ if "+controllerCondition+" }}\n\t\tmanaged.WithReferenceResolver(workerartifactrefs.NewReferenceResolver(mgr.GetClient())),{{ end }}")
	templates.CRDTypesTemplate = replace(templates.CRDTypesTemplate, "\t{{ .Imports }}", "\t{{ .Imports }}{{ if "+apiCondition+" }}\n\t"+artifactImport+"{{ end }}")
	templates.CRDTypesTemplate = replace(templates.CRDTypesTemplate, "\tForProvider       {{ .CRD.ForProviderType }} `json:\"forProvider\"`", `{{ if `+apiCondition+` }}
	// Artifact optionally selects a digest-pinned prebuilt Worker bundle.
	// +optional
	Artifact *workerartifact.Source `+"`json:\"artifact,omitempty\"`"+`
	// ArtifactFiles is disposable local state and is never written to Kubernetes.
	ArtifactFiles *workerartifact.Files `+"`json:\"-\"`"+`
{{ end }}	ForProvider       {{ .CRD.ForProviderType }} `+"`json:\"forProvider\"`")
	templates.CRDTypesTemplate = replace(templates.CRDTypesTemplate, "\tAtProvider          {{ .CRD.AtProviderType }} `json:\"atProvider,omitempty\"`", "\tAtProvider          {{ .CRD.AtProviderType }} `json:\"atProvider,omitempty\"`\n{{ if "+apiCondition+" }}\t// ResolvedArtifactDigest identifies the fetched bundle, not deployment health.\n\t// +optional\n\tResolvedArtifactDigest string `json:\"resolvedArtifactDigest,omitempty\"`\n{{ end }}")
	templates.TerraformedTemplate = replace(templates.TerraformedTemplate, "\t{{ .Imports }}", "\t{{ .Imports }}{{ if "+tfCondition+" }}\n\t"+artifactImport+"{{ end }}")
	anchor := `    return base, json.TFParser.Unmarshal(p, &base)
}

// SetParameters`
	replacement := `{{ if ` + tfCondition + ` }}
    if err := json.TFParser.Unmarshal(p, &base); err != nil { return nil, err }
    if tr.Spec.Artifact != nil && tr.GetDeletionTimestamp() == nil {
        if err := tr.ValidateWorkerArtifactInputs(); err != nil { return nil, err }
    }
    return workerartifact.Parameters(tr.Spec.Artifact, tr.Spec.ArtifactFiles, tr.GetDeletionTimestamp() != nil, "{{ .Terraform.ResourceType }}", base)
{{ else }}    return base, json.TFParser.Unmarshal(p, &base)
{{ end }}}

// SetParameters`
	templates.TerraformedTemplate = replace(templates.TerraformedTemplate, anchor, replacement)
	templates.TerraformedTemplate = replace(templates.TerraformedTemplate, "    li := resource.NewGenericLateInitializer(opts...)", `{{ if `+tfCondition+` }}
    if tr.Spec.Artifact != nil {
        for _, name := range []string{"Content", "ContentFile", "ContentSha256", "ContentType", "MainModule", "Modules", "Files", "BodyPart", "Assets", "KeepAssets"} {
            opts = append(opts, resource.WithNameFilter(name))
        }
    }
{{ end }}    li := resource.NewGenericLateInitializer(opts...)`)
	templates.TerraformedTemplate += `
{{ if ` + tfCondition + ` }}
// GetWorkerArtifact returns the optional immutable bundle selection.
func (tr *{{ .CRD.Kind }}) GetWorkerArtifact() *workerartifact.Source { return tr.Spec.Artifact }
// GetWorkerArtifactFiles returns the controller-local bundle attachment.
func (tr *{{ .CRD.Kind }}) GetWorkerArtifactFiles() *workerartifact.Files { return tr.Spec.ArtifactFiles }
// SetWorkerArtifactFiles attaches controller-local files without persisting paths.
func (tr *{{ .CRD.Kind }}) SetWorkerArtifactFiles(files *workerartifact.Files) { tr.Spec.ArtifactFiles = files }
// SetResolvedWorkerArtifact records which immutable bundle was fetched.
func (tr *{{ .CRD.Kind }}) SetResolvedWorkerArtifact(digest string) { tr.Status.ResolvedArtifactDigest = digest }
// ValidateWorkerArtifactInputs rejects conflicting code and asset input owners.
func (tr *{{ .CRD.Kind }}) ValidateWorkerArtifactInputs() error {
    if err := workerartifact.ValidateRawInputs(tr.Spec.ForProvider); err != nil { return err }
    return workerartifact.ValidateRawInputs(tr.Spec.InitProvider)
}
{{ end }}`
}

func replace(template, anchor, replacement string) string {
	if strings.Count(template, anchor) != 1 {
		panic(fmt.Sprintf("Worker artifact template extension expected one anchor %q", anchor))
	}
	return strings.Replace(template, anchor, replacement, 1)
}
