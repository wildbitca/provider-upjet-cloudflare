package config

import (
	"context"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const originPullsResource = "cloudflare_authenticated_origin_pulls"

// originPullsInitializer seeds status.atProvider.hostname, IN MEMORY ONLY, for
// a managed resource that has no observation yet.
//
// Why: the upstream Read() of cloudflare_authenticated_origin_pulls reads the
// top-level computed `hostname` from the prior state (not config[0].hostname).
// upjet builds that prior state from the parameters and the observation
// (status.atProvider); for a brand-new resource neither carries a top-level
// hostname, so the first refresh fails with "missing required hostname
// parameter". `hostname` is computed/read-only, so it must NOT be injected into
// the parameters (they become main.tf.json and Terraform rejects it with
// "Invalid Configuration for Read-Only Attribute"). The observation is only
// ever used as prior state, never as configuration, so it is the safe carrier.
//
// The object is not persisted by this initializer: the managed reconciler
// keeps using the same in-memory object for Connect/Observe and then overwrites
// the observation with the refreshed Terraform state.
func originPullsInitializer(_ client.Client) managed.Initializer {
	return managed.InitializerFn(func(_ context.Context, mg xpresource.Managed) error {
		tr, ok := mg.(resource.Terraformed)
		if !ok {
			return nil
		}
		obs, err := tr.GetObservation()
		if err != nil {
			return err
		}
		if h, _ := obs["hostname"].(string); h != "" {
			return nil
		}
		params, err := tr.GetParameters()
		if err != nil {
			return err
		}
		h := originPullsConfigHostname(params)
		if h == "" {
			return nil
		}
		if obs == nil {
			obs = map[string]any{}
		}
		obs["hostname"] = h
		return tr.SetObservation(obs)
	})
}

// originPullsConfigHostname returns config[0].hostname from the Terraform
// parameters, or "" if it is not set.
func originPullsConfigHostname(params map[string]any) string {
	var first map[string]any
	switch c := params["config"].(type) {
	case []any:
		if len(c) > 0 {
			first, _ = c[0].(map[string]any)
		}
	case []map[string]any:
		if len(c) > 0 {
			first = c[0]
		}
	case map[string]any:
		first = c
	}
	h, _ := first["hostname"].(string)
	return h
}
