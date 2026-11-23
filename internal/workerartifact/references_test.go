package workerartifact

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	scriptapi "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
	api "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type resolvingScript struct {
	*scriptapi.Script
	ResolvedName string `json:"-"`
	ResolveError error  `json:"-"`
}

func (s *resolvingScript) DeepCopyObject() runtime.Object {
	clone := *s
	clone.Script = s.DeepCopy()
	return &clone
}

func (s *resolvingScript) ResolveReferences(context.Context, client.Reader) error {
	if s.ResolveError != nil {
		return s.ResolveError
	}
	name := s.ResolvedName
	s.Spec.ForProvider.ScriptName = &name
	return nil
}

type referencePatch struct {
	patchType types.PatchType
	data      []byte
}

type referencePatchClient struct {
	client.Client
	patches  []referencePatch
	patchErr error
}

func (c *referencePatchClient) Patch(_ context.Context, object client.Object, patch client.Patch, _ ...client.PatchOption) error {
	data, err := patch.Data(object)
	if err != nil {
		return err
	}
	c.patches = append(c.patches, referencePatch{patchType: patch.Type(), data: data})
	if c.patchErr == nil {
		// An API patch response carries only the persisted fields.
		object.(Resource).SetWorkerArtifactFiles(nil)
	}
	return c.patchErr
}

func referenceFixture(t *testing.T) *resolvingScript {
	t.Helper()
	cr := &scriptapi.Script{}
	if err := json.Unmarshal([]byte(`{"spec":{"forProvider":{"bindings":[{"name":"ORIGIN","type":"plain_text","textSecretRef":{"name":"worker-env","key":"ORIGIN"}}]}}}`), cr); err != nil {
		t.Fatal(err)
	}
	cr.SetGroupVersionKind(scriptapi.Script_GroupVersionKind)
	cr.SetName("worker")
	cr.SetNamespace("preview")
	cr.SetResourceVersion("7")
	cr.Spec.Artifact = &api.Source{Ref: "registry.example/worker@sha256:" + strings.Repeat("a", 64)}
	return &resolvingScript{Script: cr, ResolvedName: "worker-v1"}
}

func TestArtifactReferenceUpdatesKeepBindingSources(t *testing.T) {
	cr := referenceFixture(t)
	files := &api.Files{Root: t.TempDir()}
	cr.Spec.ArtifactFiles = files
	kube := &referencePatchClient{}
	resolver := NewReferenceResolver(kube)
	for _, name := range []string{"worker-v1", "worker-v2"} {
		cr.ResolvedName = name
		if err := resolver.ResolveReferences(t.Context(), cr); err != nil {
			t.Fatal(err)
		}
		if len(kube.patches) == 0 || cr.Spec.ArtifactFiles != files {
			t.Fatal("reference update must patch the resource and retain its local bundle")
		}
		patch := kube.patches[len(kube.patches)-1]
		var document struct {
			Metadata struct {
				ResourceVersion string `json:"resourceVersion"`
			} `json:"metadata"`
			Spec struct {
				ForProvider map[string]any `json:"forProvider"`
			} `json:"spec"`
		}
		if err := json.Unmarshal(patch.data, &document); err != nil {
			t.Fatal(err)
		}
		if patch.patchType != types.MergePatchType || document.Metadata.ResourceVersion != "7" {
			t.Fatal("reference update lacks merge semantics or optimistic locking")
		}
		if len(document.Spec.ForProvider) != 1 || document.Spec.ForProvider["scriptName"] != name {
			t.Fatal("reference update patched fields beyond the resolved Worker name")
		}
	}
	if err := resolver.ResolveReferences(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if len(kube.patches) != 2 {
		t.Fatal("expected one API write per change to the resolved name")
	}
}

func TestNonArtifactScriptKeepsDefaultReferenceResolver(t *testing.T) {
	cr := referenceFixture(t)
	cr.Spec.Artifact = nil
	kube := &referencePatchClient{}
	if err := NewReferenceResolver(kube).ResolveReferences(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if len(kube.patches) != 1 || kube.patches[0].patchType != types.ApplyPatchType {
		t.Fatal("non-artifact Script did not retain the default resolver")
	}
}

func TestReferencePatchConflictIsReturned(t *testing.T) {
	cr := referenceFixture(t)
	kube := &referencePatchClient{
		patchErr: apierrors.NewConflict(schema.GroupResource{Group: scriptapi.Script_GroupVersionKind.Group, Resource: "scripts"}, cr.Name, errors.New("concurrent update")),
	}
	if err := NewReferenceResolver(kube).ResolveReferences(t.Context(), cr); !apierrors.IsConflict(err) {
		t.Fatal("concurrent update conflict was swallowed")
	}
}

func TestUnresolvedReferenceDoesNotPatch(t *testing.T) {
	cr := referenceFixture(t)
	cr.ResolveError = errors.New("reference could not be resolved")
	kube := &referencePatchClient{}
	if err := NewReferenceResolver(kube).ResolveReferences(t.Context(), cr); !errors.Is(err, cr.ResolveError) || len(kube.patches) != 0 {
		t.Fatal("unresolved reference must return its error without patching")
	}
}
