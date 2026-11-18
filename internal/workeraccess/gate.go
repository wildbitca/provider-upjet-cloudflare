package workeraccess

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	workers "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
	zero "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/zero/v1alpha1"
	configaccess "github.com/wildbitca/provider-upjet-cloudflare/config/namespaced/workeraccess"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type accessSource struct {
	account  *string
	name     *string
	ref      *xpv2.NamespacedReference
	selectBy *xpv2.NamespacedSelector
}

// NewInitializer validates Access before the connector can contact Cloudflare.
// It deliberately resolves with no cached Worker name.
func NewInitializer(kube client.Reader) managed.Initializer {
	return managed.InitializerFn(func(ctx context.Context, mg xpresource.Managed) error {
		if mg.GetDeletionTimestamp() != nil {
			return nil
		}
		var desired, initial accessSource
		switch worker := mg.(type) {
		case *workers.Script:
			p, i := worker.Spec.ForProvider, worker.Spec.InitProvider
			desired = accessSource{p.AccountID, p.ScriptName, p.AccessApplicationRef, p.AccessApplicationSelector}
			initial = accessSource{i.AccountID, i.ScriptName, i.AccessApplicationRef, i.AccessApplicationSelector}
		case *workers.ScriptSubdomain:
			p, i := worker.Spec.ForProvider, worker.Spec.InitProvider
			desired = accessSource{p.AccountID, p.ScriptName, p.AccessApplicationRef, p.AccessApplicationSelector}
			initial = accessSource{i.AccountID, i.ScriptName, i.AccessApplicationRef, i.AccessApplicationSelector}
		default:
			return nil
		}

		account := reference.FromPtrValue(desired.account)
		if account == "" {
			account = reference.FromPtrValue(initial.account)
		}
		if err := validate(ctx, kube, mg, account, desired); err != nil {
			return err
		}
		if err := validate(ctx, kube, mg, account, initial); err != nil {
			return err
		}
		return nil
	})
}

func validate(ctx context.Context, kube client.Reader, mg xpresource.Managed, account string, source accessSource) error {
	if source.ref == nil && source.selectBy == nil {
		return nil
	}
	if account == "" {
		return fmt.Errorf("an Access application reference requires a Cloudflare account ID")
	}
	result, err := reference.NewAPINamespacedResolver(kube, mg).Resolve(ctx, reference.NamespacedResolutionRequest{
		Namespace: mg.GetNamespace(),
		Reference: source.ref,
		Selector:  source.selectBy,
		To: reference.To{
			List:    &zero.TrustAccessApplicationList{},
			Managed: &zero.TrustAccessApplication{},
		},
		Extract: func(resource xpresource.Managed) string {
			app, ok := resource.(*zero.TrustAccessApplication)
			if !ok || reference.FromPtrValue(app.Status.AtProvider.AccountID) != account {
				return ""
			}
			return configaccess.AccessWorkerName()(app)
		},
	})
	if err != nil {
		return fmt.Errorf("resolve same-account service-token Access application: %w", err)
	}
	if result.ResolvedValue == "" {
		return fmt.Errorf("same-account service-token Access application is unavailable")
	}
	if name := reference.FromPtrValue(source.name); name != "" && name != result.ResolvedValue {
		return fmt.Errorf("configured Worker name does not match the referenced Access application")
	}
	return nil
}
