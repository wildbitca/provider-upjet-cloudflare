package config

import (
	"testing"

	"github.com/crossplane/upjet/v2/pkg/resource"
	cluster "github.com/wildbitca/provider-upjet-cloudflare/apis/cluster/zero/v1alpha1"
	namespaced "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/zero/v1alpha1"
)

func TestServiceTokenObservationDoesNotRequestRotation(t *testing.T) {
	for name, tr := range map[string]resource.Terraformed{
		"cluster":    &cluster.TrustAccessServiceToken{},
		"namespaced": &namespaced.TrustAccessServiceToken{},
	} {
		t.Run(name, func(t *testing.T) {
			// A newly created token reports version 1 without a previous-secret expiry.
			// Sending that computed version back as configuration fails native validation.
			if _, err := tr.LateInitialize([]byte(`{"client_secret_version":1}`)); err != nil {
				t.Fatal(err)
			}
			params, err := tr.GetParameters()
			if err != nil {
				t.Fatal(err)
			}
			if _, exists := params["client_secret_version"]; exists {
				t.Fatal("observed secret version became an unrequested rotation input")
			}
			// Explicitly configured rotations must remain available.
			expiry := "2026-11-18T12:00:00Z"
			if err := tr.SetParameters(map[string]any{
				"client_secret_version": float64(2), "previous_client_secret_expires_at": expiry,
			}); err != nil {
				t.Fatal(err)
			}
			if _, err := tr.LateInitialize([]byte(`{"client_secret_version":1}`)); err != nil {
				t.Fatal(err)
			}
			params, err = tr.GetParameters()
			if err != nil {
				t.Fatal(err)
			}
			if params["client_secret_version"] != float64(2) || params["previous_client_secret_expires_at"] != expiry {
				t.Fatal("explicit rotation configuration changed")
			}
		})
	}
}
