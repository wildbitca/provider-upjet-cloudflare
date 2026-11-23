package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/crossplane/upjet/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	clusterscript "github.com/wildbitca/provider-upjet-cloudflare/apis/cluster/workers/v1alpha1"
	script "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
	api "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"
	artifact "github.com/wildbitca/provider-upjet-cloudflare/internal/workerartifact"
)

func TestWorkerArtifactsReachRealUpjetFileProducer(t *testing.T) {
	ref := "registry.example/worker@sha256:" + strings.Repeat("a", 64)
	for _, tc := range []struct {
		name       string
		tr         resource.Terraformed
		namespaced bool
	}{
		{"namespaced Script", &script.Script{Spec: script.ScriptSpec{Artifact: &api.Source{Ref: ref}}}, true},
		{"cluster Script", &clusterscript.Script{Spec: clusterscript.ScriptSpec{Artifact: &api.Source{Ref: ref}}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.tr.SetName("preview-worker")
			tc.tr.SetUID("preview-worker-uid")
			files := &api.Files{Root: t.TempDir(), ContentSHA256: strings.Repeat("b", 64), Manifest: api.Manifest{SchemaVersion: 1, MainModule: "worker.js", Modules: []api.Module{{Name: "helper.wasm", ContentType: "application/wasm"}, {Name: "worker.js", ContentType: "application/javascript+module"}}, Assets: "public"}}
			adapter, ok := tc.tr.(artifact.Resource)
			if !ok {
				t.Fatal("generated resource has no artifact adapter")
			}
			adapter.SetWorkerArtifactFiles(files)
			provider := GetProvider()
			if tc.namespaced {
				provider = GetProviderNamespaced()
			}
			cfg := provider.Resources[tc.tr.GetTerraformResourceType()]
			if len(cfg.InitializerFns) != 1 {
				t.Fatal("artifact initializer not wired into resource configuration")
			}
			producer, err := terraform.NewFileProducer(t.Context(), noopSecrets{}, t.TempDir(), tc.tr, terraform.Setup{}, cfg)
			if err != nil {
				t.Fatal(err)
			}
			out := producer.BuildMainTF()
			p := out["resource"].(map[string]any)[tc.tr.GetTerraformResourceType()].(map[string]any)[tc.tr.GetName()].(map[string]any)
			if p["artifact"] != nil {
				t.Fatal("Kubernetes artifact input leaked into Terraform schema")
			}
			if p["assets"].(map[string]any)["directory"] != filepath.Join(files.Root, "public") {
				t.Fatal("static assets lost before Terraform upload")
			}
			if tc.tr.GetTerraformResourceType() == "cloudflare_workers_script" {
				if p["content_file"] != filepath.Join(files.Root, "worker.js") || len(p["files"].(map[string]any)) != 1 {
					t.Fatal("Script main module or extra module lost")
				}
			}
			serialized, err := json.Marshal(tc.tr)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(serialized), files.Root) {
				t.Fatal("Terraform parameters mutated Kubernetes desired spec")
			}
			if _, err := tc.tr.LateInitialize([]byte(`{"files":{"helper.wasm":{"content_file":"/temporary/file"}},"content_type":"application/javascript+module"}`)); err != nil {
				t.Fatal(err)
			}
			after, err := json.Marshal(tc.tr)
			if err != nil || strings.Contains(string(after), "/temporary/file") {
				t.Fatal("late initialization persisted multipart module paths")
			}
			tc.tr.SetDeletionTimestamp(&metav1.Time{})
			adapter.SetWorkerArtifactFiles(nil)
			if _, err := terraform.NewFileProducer(t.Context(), noopSecrets{}, t.TempDir(), tc.tr, terraform.Setup{}, cfg); err != nil {
				t.Fatalf("deletion requires artifact files: %v", err)
			}
		})
	}
}

func TestWorkerArtifactCRDSchemas(t *testing.T) {
	for _, file := range []string{"workers.upjet-cloudflare.m.upbound.io_scripts.yaml", "workers.upjet-cloudflare.upbound.io_scripts.yaml"} {
		data, err := os.ReadFile(filepath.Join("..", "package", "crds", file))
		if err != nil {
			t.Fatal(err)
		}
		var crd map[string]any
		if err := yaml.Unmarshal(data, &crd); err != nil {
			t.Fatal(err)
		}
		versions := crd["spec"].(map[string]any)["versions"].([]any)
		for _, v := range versions {
			properties := v.(map[string]any)["schema"].(map[string]any)["openAPIV3Schema"].(map[string]any)["properties"].(map[string]any)
			spec := properties["spec"].(map[string]any)["properties"].(map[string]any)
			if spec["artifactFiles"] != nil || spec["-"] != nil {
				t.Fatal("transient files exposed in CRD")
			}
			ref := spec["artifact"].(map[string]any)["properties"].(map[string]any)["ref"].(map[string]any)
			re, err := regexp.Compile(ref["pattern"].(string))
			if err != nil {
				t.Fatal(err)
			}
			if !re.MatchString("ghcr.io/org/worker@sha256:"+strings.Repeat("a", 64)) || re.MatchString("ghcr.io/org/worker:latest") {
				t.Fatal("CRD does not enforce immutable artifact pin")
			}
		}
	}
}
