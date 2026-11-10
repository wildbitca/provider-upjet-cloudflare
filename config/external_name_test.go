package config

import (
	"context"
	"testing"

	"github.com/crossplane/upjet/v2/pkg/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestR2DomainsAreIDLessAndYieldExternalName(t *testing.T) {
	p := GetProvider()
	for res, attr := range map[string]string{
		"cloudflare_r2_custom_domain":  "domain",
		"cloudflare_r2_managed_domain": "bucket_name",
	} {
		cfg := p.Resources[res]
		if cfg == nil {
			t.Fatalf("%s not in provider", res)
		}
		if _, ok := cfg.TerraformResource.Schema["id"]; ok {
			t.Fatalf("%s now has an id attribute; drop the custom external name", res)
		}
		state := map[string]any{"account_id": "a", "bucket_name": "bkt", "domain": "cdn.example.com", attr: "want-" + attr}
		cd := &metav1.ObjectMeta{}
		if _, err := resource.SetCriticalAnnotations(cd, cfg, state, ""); err != nil {
			t.Fatalf("%s: SetCriticalAnnotations() = %v", res, err)
		}
		if got := cd.GetAnnotations()["crossplane.io/external-name"]; got != "want-"+attr {
			t.Fatalf("%s: external-name = %q", res, got)
		}
	}
}

func TestRemovedResourcesAreNotConfigured(t *testing.T) {
	// cloudflare_registrar_domain was removed in Terraform provider v5.28.0.
	if _, ok := ExternalNameConfigs["cloudflare_registrar_domain"]; ok {
		t.Fatal("cloudflare_registrar_domain must not be configured")
	}
}

func TestOriginPullsSeedsHostname(t *testing.T) {
	tests := map[string]struct {
		base map[string]interface{}
		want interface{}
	}{
		"seeds hostname from config list": {
			base: map[string]interface{}{"config": []interface{}{map[string]interface{}{"hostname": "dapi.tumed.dev", "cert_id": "c"}}},
			want: "dapi.tumed.dev",
		},
		"no config leaves hostname unset": {
			base: map[string]interface{}{"zone_id": "z"},
			want: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			originPullsExternalName.SetIdentifierArgumentFn(tc.base, "aop-x")
			if got := tc.base["hostname"]; got != tc.want {
				t.Fatalf("hostname = %v, want %v", got, tc.want)
			}
		})
	}
}
