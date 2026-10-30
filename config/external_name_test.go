package config

import (
	"context"
	"testing"
)

func TestWorkersRouteExternalName(t *testing.T) {
	tests := map[string]struct {
		externalName string
		parameters   map[string]interface{}
		want         string
		wantErr      bool
	}{
		"builds Terraform import ID": {
			externalName: "394ecb2c6017461c868c101d64ded9af",
			parameters:   map[string]interface{}{"zone_id": "d9b8169cf1a3c5cc37fe5040cc7b0db7"},
			want:         "d9b8169cf1a3c5cc37fe5040cc7b0db7/394ecb2c6017461c868c101d64ded9af",
		},
		"requires zone ID": {
			externalName: "394ecb2c6017461c868c101d64ded9af",
			parameters:   map[string]interface{}{},
			wantErr:      true,
		},
		"requires route ID": {
			parameters: map[string]interface{}{"zone_id": "d9b8169cf1a3c5cc37fe5040cc7b0db7"},
			wantErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := workersRouteExternalName.GetIDFn(context.Background(), tc.externalName, tc.parameters, nil)
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetIDFn() error = %v, wantErr %t", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("GetIDFn() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMonetizationRulesExternalName(t *testing.T) {
	cfg, ok := ExternalNameConfigs["cloudflare_monetization_rules"]
	if !ok {
		t.Fatal("cloudflare_monetization_rules is not configured")
	}
	got, err := cfg.GetExternalNameFn(map[string]interface{}{"zone_id": "d9b8169cf1a3c5cc37fe5040cc7b0db7"})
	if err != nil {
		t.Fatalf("GetExternalNameFn() error = %v", err)
	}
	if want := "d9b8169cf1a3c5cc37fe5040cc7b0db7"; got != want {
		t.Fatalf("GetExternalNameFn() = %q, want %q", got, want)
	}
}

func TestR2DomainExternalNames(t *testing.T) {
	// Terraform state for these resources carries no `id`.
	tests := map[string]struct {
		resource string
		tfstate  map[string]interface{}
		want     string
		wantErr  bool
	}{
		"custom domain uses domain": {
			resource: "cloudflare_r2_custom_domain",
			tfstate: map[string]interface{}{
				"account_id":  "023e105f4ecef8ad9ca31a8372d0c353",
				"bucket_name": "example-bucket",
				"domain":      "prefix.example-domain.com",
			},
			want: "prefix.example-domain.com",
		},
		"custom domain requires domain": {
			resource: "cloudflare_r2_custom_domain",
			tfstate:  map[string]interface{}{"bucket_name": "example-bucket"},
			wantErr:  true,
		},
		"managed domain uses bucket name": {
			resource: "cloudflare_r2_managed_domain",
			tfstate: map[string]interface{}{
				"account_id":  "023e105f4ecef8ad9ca31a8372d0c353",
				"bucket_name": "example-bucket",
				"domain":      "pub-0123456789abcdef.r2.dev",
			},
			want: "example-bucket",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, ok := ExternalNameConfigs[tc.resource]
			if !ok {
				t.Fatalf("%s is not configured", tc.resource)
			}
			got, err := cfg.GetExternalNameFn(tc.tfstate)
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetExternalNameFn() error = %v, wantErr %t", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("GetExternalNameFn() = %q, want %q", got, tc.want)
			}

			params := map[string]interface{}{"domain": "prefix.example-domain.com", "bucket_name": "example-bucket"}
			cfg.SetIdentifierArgumentFn(params, "")
			if params["domain"] != "prefix.example-domain.com" || params["bucket_name"] != "example-bucket" {
				t.Fatalf("SetIdentifierArgumentFn() overwrote spec arguments: %v", params)
			}
		})
	}
}

func TestRemovedResourcesAreNotConfigured(t *testing.T) {
	// cloudflare_registrar_domain was removed in Terraform provider v5.28.0.
	if _, ok := ExternalNameConfigs["cloudflare_registrar_domain"]; ok {
		t.Fatal("cloudflare_registrar_domain must not be configured")
	}
}
