package workerartifact

import (
	"context"
	"fmt"
	"reflect"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewReferenceResolver preserves configured binding sources while resolving
// artifact-backed Script references. Other inputs retain the default resolver.
func NewReferenceResolver(kube client.Client) managed.ReferenceResolver {
	return &scriptReferenceResolver{kube: kube, fallback: managed.NewAPISimpleReferenceResolver(kube)}
}

type scriptReferenceResolver struct {
	kube     client.Client
	fallback managed.ReferenceResolver
}

func (r *scriptReferenceResolver) ResolveReferences(ctx context.Context, mg xpresource.Managed) error {
	artifact, ok := mg.(Resource)
	if !ok || artifact.GetWorkerArtifact() == nil {
		return r.fallback.ResolveReferences(ctx, mg)
	}
	refs, ok := mg.(interface {
		ResolveReferences(context.Context, client.Reader) error
	})
	if !ok {
		return nil
	}
	before := mg.DeepCopyObject().(client.Object)
	if err := refs.ResolveReferences(ctx, r.kube); err != nil {
		return fmt.Errorf("resolve Worker references: %w", err)
	}
	if reflect.DeepEqual(before, mg) {
		return nil
	}
	// This is an incremental update, not a complete field-ownership declaration.
	// A merge patch retains earlier reference updates and Secret selectors; the
	// version precondition prevents overwriting concurrent GitOps changes.
	patch := client.MergeFromWithOptions(before, client.MergeFromWithOptimisticLock{})
	// The API response contains no controller-local attachment. Keep the bundle
	// prepared by the initializer for the connector that runs after this update.
	files := artifact.GetWorkerArtifactFiles()
	defer artifact.SetWorkerArtifactFiles(files)
	return r.kube.Patch(ctx, mg, patch, client.FieldOwner("cloudflare-worker-artifact-references"))
}
