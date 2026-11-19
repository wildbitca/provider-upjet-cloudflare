package workerartifact

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	connectivity "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/connectivity/v1alpha1"
	scriptapi "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
	api "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type referencePatchClient struct {
	client.Client
	types             []types.PatchType
	documents         []map[string]any
	patchErr          error
	discardLocalFiles bool
}

func (c *referencePatchClient) Patch(ctx context.Context, object client.Object, patch client.Patch, options ...client.PatchOption) error {
	data, err := patch.Data(object)
	if err != nil {
		return err
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		return err
	}
	c.types = append(c.types, patch.Type())
	c.documents = append(c.documents, document)
	if c.patchErr != nil {
		return c.patchErr
	}
	if patch.Type() == types.ApplyPatchType {
		return nil // The fake client does not implement server-side apply ownership.
	}
	if err := c.Client.Patch(ctx, object, patch, options...); err != nil {
		return err
	}
	if c.discardLocalFiles {
		object.(Resource).SetWorkerArtifactFiles(nil)
	}
	return nil
}

func referenceFixture(t *testing.T) (*scriptapi.Script, *referencePatchClient) {
	t.Helper()
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{scriptapi.AddToScheme, connectivity.AddToScheme} {
		if err := add(scheme); err != nil {
			t.Fatal(err)
		}
	}
	cr := &scriptapi.Script{}
	if err := json.Unmarshal([]byte(`{"spec":{"forProvider":{"bindings":[{"name":"ORIGIN","type":"plain_text","textSecretRef":{"name":"worker-env","key":"ORIGIN"}},{"name":"PRIVATE_API","type":"vpc_service","serviceIdRef":{"name":"origin","policy":{"resolve":"Always"}}}]}}}`), cr); err != nil {
		t.Fatal(err)
	}
	cr.SetGroupVersionKind(scriptapi.Script_GroupVersionKind)
	cr.SetName("worker")
	cr.SetNamespace("preview")
	cr.SetResourceVersion("7")
	cr.Spec.Artifact = &api.Source{Ref: "registry.example/worker@sha256:" + strings.Repeat("a", 64)}
	dir := &connectivity.DirectoryService{ObjectMeta: metav1.ObjectMeta{Name: "origin", Namespace: "preview", Annotations: map[string]string{"crossplane.io/external-name": "service-v1"}}}
	kube := fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr, dir).Build()
	return cr, &referencePatchClient{Client: kube}
}

func TestArtifactReferenceUpdatesKeepBindingSources(t *testing.T) {
	cr, kube := referenceFixture(t)
	files := &api.Files{Root: t.TempDir()}
	cr.Spec.ArtifactFiles = files
	kube.discardLocalFiles = true
	resolver := NewReferenceResolver(kube)
	if err := resolver.ResolveReferences(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if kube.types[0] != types.MergePatchType || kube.documents[0]["metadata"].(map[string]any)["resourceVersion"] != "7" {
		t.Fatal("reference update lacks merge semantics or optimistic locking")
	}
	if cr.Spec.ArtifactFiles != files {
		t.Fatal("reference update lost the bundle needed by the connector")
	}
	// A change to the referenced service must retain the binding inputs.
	dir := &connectivity.DirectoryService{}
	if err := kube.Get(t.Context(), client.ObjectKey{Namespace: "preview", Name: "origin"}, dir); err != nil {
		t.Fatal(err)
	}
	dir.Annotations["crossplane.io/external-name"] = "service-v2"
	if err := kube.Update(t.Context(), dir); err != nil {
		t.Fatal(err)
	}
	if err := resolver.ResolveReferences(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if kube.types[1] != types.MergePatchType {
		t.Fatal("later reference update changed patch semantics")
	}
	if cr.Spec.ArtifactFiles != files {
		t.Fatal("later reference update lost the controller-local bundle")
	}
	stored := &scriptapi.Script{}
	if err := kube.Get(t.Context(), client.ObjectKeyFromObject(cr), stored); err != nil {
		t.Fatal(err)
	}
	if len(stored.Spec.ForProvider.Bindings) != 2 || stored.Spec.ForProvider.Bindings[0].TextSecretRef == nil || stored.Spec.ForProvider.Bindings[1].ServiceIDRef == nil {
		t.Fatal("later reference update lost a binding source")
	}
	before := len(kube.types)
	if err := resolver.ResolveReferences(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if len(kube.types) != before {
		t.Fatal("unchanged references caused an API write")
	}
}

func TestNonArtifactScriptKeepsDefaultReferenceResolver(t *testing.T) {
	cr, kube := referenceFixture(t)
	cr.Spec.Artifact = nil
	if err := NewReferenceResolver(kube).ResolveReferences(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if len(kube.types) != 1 || kube.types[0] != types.ApplyPatchType {
		t.Fatal("non-artifact Script did not retain the default resolver")
	}
}

func TestReferencePatchConflictIsReturned(t *testing.T) {
	cr, kube := referenceFixture(t)
	kube.patchErr = apierrors.NewConflict(schema.GroupResource{Group: scriptapi.Script_GroupVersionKind.Group, Resource: "scripts"}, cr.Name, errors.New("concurrent update"))
	if err := NewReferenceResolver(kube).ResolveReferences(t.Context(), cr); !apierrors.IsConflict(err) {
		t.Fatal("concurrent update conflict was swallowed")
	}
}

func TestUnresolvedReferenceDoesNotPatch(t *testing.T) {
	cr, kube := referenceFixture(t)
	cr.Spec.ForProvider.Bindings[1].ServiceIDRef.Name = "missing"
	if err := NewReferenceResolver(kube).ResolveReferences(t.Context(), cr); err == nil || len(kube.types) != 0 {
		t.Fatal("unresolved reference caused an API update")
	}
}
