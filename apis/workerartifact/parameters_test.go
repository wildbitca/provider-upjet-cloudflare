package workerartifact

import "testing"

func TestDeletionNeedsNoArtifactFiles(t *testing.T) {
	p := map[string]any{
		"script_name": "preview", "account_id": "account",
		"content_file": "/lost/worker.js", "files": map[string]any{"helper.js": "/lost/helper.js"},
		"assets": map[string]any{"directory": "/lost/public"},
	}
	got, err := Parameters(&Source{Ref: "registry.example/missing"}, nil, true, "cloudflare_workers_script", p)
	if err != nil {
		t.Fatal(err)
	}
	if got["content"] == nil || got["script_name"] != "preview" || got["account_id"] != "account" {
		t.Fatal("deletion configuration must satisfy native content validation and preserve identity")
	}
	for _, key := range []string{"content_file", "content_sha256", "files", "assets", "body_part", "keep_assets"} {
		if _, ok := got[key]; ok {
			t.Fatalf("deletion still depends on artifact field %s", key)
		}
	}
}

func TestNonArtifactDeletionPreservesInputs(t *testing.T) {
	got, err := Parameters(nil, nil, true, "cloudflare_workers_script", map[string]any{"content_file": "/existing/input.js"})
	if err != nil || got["content_file"] != "/existing/input.js" || got["content"] != nil {
		t.Fatalf("non-artifact deletion inputs changed: %v", err)
	}
}
