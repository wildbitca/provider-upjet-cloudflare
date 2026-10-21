package config

import (
	"context"
	"testing"
)

func TestWorkersRouteExternalName(t *testing.T) {
	tests := map[string]struct {
		externalName string
		parameters   map[string]interface{}
		want          string
		wantErr       bool
	}{
		"builds Terraform import ID": {
			externalName: "394ecb2c6017461c868c101d64ded9af",
			parameters:   map[string]interface{}{"zone_id": "d9b8169cf1a3c5cc37fe5040cc7b0db7"},
			want:          "d9b8169cf1a3c5cc37fe5040cc7b0db7/394ecb2c6017461c868c101d64ded9af",
		},
		"requires zone ID": {
			externalName: "394ecb2c6017461c868c101d64ded9af",
			parameters:   map[string]interface{}{},
			wantErr:       true,
		},
		"requires route ID": {
			parameters: map[string]interface{}{"zone_id": "d9b8169cf1a3c5cc37fe5040cc7b0db7"},
			wantErr:   true,
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
