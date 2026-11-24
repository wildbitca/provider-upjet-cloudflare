package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"k8s.io/apimachinery/pkg/types"

	workers "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
)

const (
	testRouteID = "394ecb2c6017461c868c101d64ded9af"
	testZoneID  = "d9b8169cf1a3c5cc37fe5040cc7b0db7"
)

func TestWorkersRouteWorkspace(t *testing.T) {
	cfg := GetProvider().Resources["cloudflare_workers_route"]
	// Crossplane configures a local provider mirror. This test exercises the
	// actual workspace and state production without running any CLI commands.
	setup := terraform.Setup{
		Version:     "1.13.9",
		Requirement: terraform.ProviderRequirement{Source: "cloudflare/cloudflare", Version: "5.29.0"},
	}
	for _, script := range []string{"", "example-worker"} {
		name := "without_script"
		if script != "" {
			name = "with_script"
		}
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			route := &workers.Route{}
			route.SetName("example-route")
			route.SetUID(types.UID("example-route"))
			route.SetNamespace("default")
			zone, pattern := testZoneID, "example.com/*"
			route.Spec.ForProvider.ZoneID = &zone
			route.Spec.ForProvider.Pattern = &pattern
			if script != "" {
				route.Spec.ForProvider.Script = &script
			}

			dir := t.TempDir()
			t.Setenv("TMPDIR", dir)
			store := terraform.NewWorkspaceStore(logging.NewNopLogger(), terraform.WithDisableInit(true))
			workspace, err := store.Workspace(ctx, nil, route, setup, cfg)
			if err != nil {
				t.Fatalf("Workspace(new Route with no external name): %v", err)
			}
			result, err := workspace.Import(ctx, route)
			if err != nil || result.Exists {
				t.Fatalf("Import(new Route) = %+v, %v; want not found without invoking Terraform", result, err)
			}
			workspaceDir := filepath.Join(dir, string(route.GetUID()))
			if _, err := os.Stat(filepath.Join(workspaceDir, "terraform.tfstate")); !os.IsNotExist(err) {
				t.Fatalf("new Route must not have synthetic state; stat error = %v", err)
			}
			var main struct {
				Resource map[string]map[string]map[string]any `json:"resource"`
			}
			readRouteWorkspaceJSON(t, filepath.Join(workspaceDir, "main.tf.json"), &main)
			args := main.Resource[route.GetTerraformResourceType()][route.GetName()]
			if args["zone_id"] != zone || args["pattern"] != pattern {
				t.Fatalf("Terraform arguments = %v; want zone %q and pattern %q", args, zone, pattern)
			}
			gotScript, hasScript := args["script"]
			if (script == "" && hasScript) || (script != "" && gotScript != script) {
				t.Fatalf("Terraform script = %v (present %t), want %q or omitted when empty", gotScript, hasScript, script)
			}

			// Simulate a successful create followed by pod replacement:
			// the Kubernetes object persists, but the workspace does not.
			id := testRouteID
			meta.SetExternalName(route, id)
			route.Status.AtProvider.ID = &id
			restartDir := t.TempDir()
			t.Setenv("TMPDIR", restartDir)
			restarted := terraform.NewWorkspaceStore(logging.NewNopLogger(), terraform.WithDisableInit(true))
			if _, err := restarted.Workspace(ctx, nil, route, setup, cfg); err != nil {
				t.Fatalf("Workspace(existing Route after restart): %v", err)
			}
			var state struct {
				Resources []struct {
					Instances []struct {
						Attributes map[string]any `json:"attributes"`
					} `json:"instances"`
				} `json:"resources"`
			}
			readRouteWorkspaceJSON(t, filepath.Join(restartDir, string(route.GetUID()), "terraform.tfstate"), &state)
			if len(state.Resources) != 1 || len(state.Resources[0].Instances) != 1 {
				t.Fatalf("expected one Route instance in the state used by refresh: %+v", state)
			}
			attrs := state.Resources[0].Instances[0].Attributes
			if attrs["id"] != id || attrs["zone_id"] != zone {
				t.Fatalf("state used by refresh has id %v and zone %v; want bare ID %q and zone %q", attrs["id"], attrs["zone_id"], id, zone)
			}
		})
	}
}

func readRouteWorkspaceJSON(t *testing.T, path string, into any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, into); err != nil {
		t.Fatal(err)
	}
}
