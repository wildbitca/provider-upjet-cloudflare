package v1alpha1

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestScriptServiceReference(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	target := func(namespace, workerName string) *Script {
		script := &Script{ObjectMeta: metav1.ObjectMeta{Name: "gateway", Namespace: namespace}}
		meta.SetExternalName(script, workerName)
		return script
	}
	client := fake.NewClientBuilder().WithScheme(scheme).WithObjects(
		target("preview", "preview-gateway-worker"),
		target("shared", "shared-gateway-worker"),
	).Build()

	caller := &Script{
		ObjectMeta: metav1.ObjectMeta{Name: "caller", Namespace: "preview"},
		Spec: ScriptSpec{ForProvider: ScriptParameters{Bindings: []BindingsParameters{
			{Name: reference.ToPtrValue("LOCAL"), Type: reference.ToPtrValue("service"), ServiceRef: &xpv2.NamespacedReference{Name: "gateway"}},
			{Name: reference.ToPtrValue("SHARED"), Type: reference.ToPtrValue("service"), ServiceRef: &xpv2.NamespacedReference{Name: "gateway", Namespace: "shared"}},
			{Name: reference.ToPtrValue("EXTERNAL"), Type: reference.ToPtrValue("service"), Service: reference.ToPtrValue("external-worker")},
		}}},
	}

	if err := caller.ResolveReferences(t.Context(), client); err != nil {
		t.Fatal(err)
	}

	params, err := caller.GetParameters()
	if err != nil {
		t.Fatal(err)
	}
	bindings := params["bindings"].([]any)
	for i, want := range []string{"preview-gateway-worker", "shared-gateway-worker", "external-worker"} {
		binding := bindings[i].(map[string]any)
		if got := binding["service"]; got != want {
			t.Errorf("binding %d service = %v, want %q", i, got, want)
		}
		if _, ok := binding["serviceRef"]; ok {
			t.Errorf("binding %d includes the Kubernetes serviceRef in Terraform parameters", i)
		}
	}
}
