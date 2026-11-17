package workerartifact

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// ValidateRawInputs also sees sensitive Secret references omitted by TFParser.
func ValidateRawInputs(value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var p map[string]any
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	for _, key := range []string{"content", "contentFile", "contentSha256", "contentType", "mainModule", "modules", "files", "bodyPart", "keepAssets"} {
		if v, ok := p[key]; ok && v != nil {
			return fmt.Errorf("artifact mode cannot be combined with %s", key)
		}
	}
	if assets, ok := p["assets"].(map[string]any); ok {
		for _, key := range []string{"directory", "jwtSecretRef"} {
			if v, ok := assets[key]; ok && v != nil {
				return fmt.Errorf("artifact mode cannot be combined with assets.%s", key)
			}
		}
	}
	return nil
}

// ValidateInputs prevents two owners for the bundle's code and asset files.
// Bindings, compatibility settings and assets.config remain normal inputs.
func ValidateInputs(p map[string]any) error {
	for _, key := range []string{"content", "content_file", "content_sha256", "content_type", "main_module", "modules", "files", "body_part", "keep_assets"} {
		if v, ok := p[key]; ok && v != nil {
			return fmt.Errorf("artifact mode cannot be combined with %s", key)
		}
	}
	if assets, ok := p["assets"].(map[string]any); ok {
		for _, key := range []string{"directory", "jwt"} {
			if v, ok := assets[key]; ok && v != nil {
				return fmt.Errorf("artifact mode cannot be combined with assets.%s", key)
			}
		}
	}
	return nil
}

// Parameters injects file inputs into Terraform parameters, never into the
// Kubernetes spec. Deletion deliberately needs neither the bundle nor registry.
func Parameters(source *Source, files *Files, deleting bool, kind string, p map[string]any) (map[string]any, error) {
	if source == nil {
		return p, nil
	}
	if kind != "cloudflare_workers_script" {
		return nil, fmt.Errorf("unsupported artifact resource %q", kind)
	}
	if deleting {
		// The native provider validates content/content_file/assets even during
		// refresh and destroy. Supply a self-contained configuration without
		// reading disposable files. The managed deletion path never uploads it.
		for _, key := range []string{"content_file", "content_sha256", "files", "assets", "body_part", "keep_assets"} {
			delete(p, key)
		}
		p["content"] = "export default {};"
		p["content_type"] = "application/javascript+module"
		p["main_module"] = "worker.js"
		return p, nil
	}
	if err := ValidateInputs(p); err != nil {
		return nil, err
	}
	if files == nil {
		return nil, fmt.Errorf("worker artifact has not been materialized")
	}
	m := files.Manifest
	scriptParameters(files, p)
	if m.Assets != "" {
		assets, _ := p["assets"].(map[string]any)
		if assets == nil {
			assets = map[string]any{}
		}
		assets["directory"] = filepath.Join(files.Root, m.Assets)
		p["assets"] = assets
	}
	return p, nil
}

func scriptParameters(files *Files, p map[string]any) {
	additional := map[string]any{}
	for _, module := range files.Manifest.Modules {
		if module.Name == files.Manifest.MainModule {
			p["content_file"] = filepath.Join(files.Root, module.Name)
			p["content_type"] = module.ContentType
			p["main_module"] = files.Manifest.MainModule
			p["content_sha256"] = files.ContentSHA256
		} else {
			additional[module.Name] = map[string]any{"content_file": filepath.Join(files.Root, module.Name), "content_type": module.ContentType}
		}
	}
	if len(additional) > 0 {
		p["files"] = additional
	}
}
