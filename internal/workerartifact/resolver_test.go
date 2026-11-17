package workerartifact

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"oras.land/oras-go/v2/registry/remote/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	scriptapi "github.com/wildbitca/provider-upjet-cloudflare/apis/namespaced/workers/v1alpha1"
	api "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"
)

type memoryFetcher struct {
	data  []byte
	calls atomic.Int32
}

func (f *memoryFetcher) Fetch(ctx context.Context, _ string, _ auth.Credential, w io.Writer) error {
	f.calls.Add(1)
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err := w.Write(f.data)
	return err
}

func newScript(uid, ref string) *scriptapi.Script {
	return &scriptapi.Script{ObjectMeta: metav1.ObjectMeta{Name: "gateway", Namespace: "preview", UID: types.UID(uid)}, Spec: scriptapi.ScriptSpec{Artifact: &api.Source{Ref: ref}}}
}

func fakeKube(t *testing.T) *fake.ClientBuilder {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	return fake.NewClientBuilder().WithScheme(scheme)
}

func TestInitializerUsesPrivateRegistryAndKeepsPathsOutOfSpec(t *testing.T) {
	registry := newRegistry(t, standardBundle(t), true)
	host := strings.TrimPrefix(registry.server.URL, "https://")
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "registry", Namespace: "preview"}, Type: corev1.SecretTypeDockerConfigJson, Data: map[string][]byte{corev1.DockerConfigJsonKey: []byte(`{"auths":{"` + host + `":{"username":"worker","password":"registry-password"}}}`)}}
	kube := fakeKube(t).WithObjects(secret).Build()
	resolver := &Resolver{Root: t.TempDir(), Fetcher: RegistryFetcher{Client: registry.server.Client()}}
	cr := newScript("private-worker", registry.ref)
	cr.Spec.Artifact.PullSecretRef = &api.SecretReference{Name: "registry"}
	if err := resolver.Initializer(kube).Initialize(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	p, err := cr.GetParameters()
	if err != nil {
		t.Fatal(err)
	}
	if p["content_file"] == nil || p["content_sha256"] == "" || p["assets"] == nil {
		t.Fatalf("file parameters missing: %v", p)
	}
	serialized, err := json.Marshal(cr)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(serialized, []byte(resolver.Root)) || cr.Spec.ForProvider.ContentFile != nil || bytes.Contains(serialized, []byte("registry-password")) {
		t.Fatal("controller-local paths or credentials persisted in resource")
	}
	if cr.Status.ResolvedArtifactDigest == "" {
		t.Fatal("resolved digest missing")
	}
	if _, err := cr.DeepCopy().GetParameters(); err != nil {
		t.Fatalf("materialized files lost during deep copy: %v", err)
	}
	if _, err := cr.LateInitialize([]byte(`{"content_file":"/transient/path","main_module":"worker.js","assets":{"directory":"/transient/assets"}}`)); err != nil {
		t.Fatal(err)
	}
	if cr.Spec.ForProvider.ContentFile != nil || cr.Spec.ForProvider.MainModule != nil || cr.Spec.ForProvider.Assets != nil {
		t.Fatal("artifact paths were late-initialized into spec")
	}
}

func TestCacheConcurrencyIsolationCorruptionAndRestart(t *testing.T) {
	fetcher := &memoryFetcher{data: standardBundle(t)}
	r := &Resolver{Root: t.TempDir(), Fetcher: fetcher}
	ref := "registry.example/bundle@sha256:" + strings.Repeat("a", 64)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			if _, err := r.materialize(t.Context(), "same-resource", ref, auth.EmptyCredential); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if fetcher.calls.Load() != 1 {
		t.Fatalf("concurrent fetches were not coalesced: %d", fetcher.calls.Load())
	}
	a, err := r.materialize(t.Context(), "same-resource", ref, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.materialize(t.Context(), "different-preview", ref, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	if a.Root == b.Root {
		t.Fatal("previews share mutable directories")
	}
	restarted := &Resolver{Root: r.Root, Fetcher: fetcher}
	if _, err := restarted.materialize(t.Context(), "same-resource", ref, auth.EmptyCredential); err != nil {
		t.Fatal(err)
	}
	if fetcher.calls.Load() != 2 {
		t.Fatal("valid cache was not reused after resolver restart")
	}
	if err := os.WriteFile(filepath.Join(a.Root, "worker.js"), []byte("corrupt"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := restarted.materialize(t.Context(), "same-resource", ref, auth.EmptyCredential); err != nil {
		t.Fatal(err)
	}
	if fetcher.calls.Load() != 3 {
		t.Fatal("corrupt cache was not re-fetched")
	}
	if err := os.RemoveAll(a.Root); err != nil {
		t.Fatal(err)
	}
	if _, err := restarted.materialize(t.Context(), "same-resource", ref, auth.EmptyCredential); err != nil {
		t.Fatal(err)
	}
	if _, err := cachedBundle(b.Root, ref); err != nil {
		t.Fatalf("one preview damaged another's cache: %v", err)
	}
}

func TestDeletionSkipsArtifactRegistryAndSecret(t *testing.T) {
	cr := newScript("deleting", "invalid-reference")
	now := metav1.Now()
	cr.DeletionTimestamp = &now
	cr.Spec.Artifact.PullSecretRef = &api.SecretReference{Name: "missing"}
	r := &Resolver{Root: t.TempDir(), Fetcher: &memoryFetcher{}}
	if err := r.Initializer(fakeKube(t).Build()).Initialize(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	if _, err := cr.GetParameters(); err != nil {
		t.Fatalf("deletion depends on unavailable artifact: %v", err)
	}
}

func TestCredentialNamespacesAndMissingSecret(t *testing.T) {
	kube := fakeKube(t).Build()
	for _, tc := range []struct {
		namespace string
		ref       api.SecretReference
	}{
		{"preview", api.SecretReference{Name: "x", Namespace: "other-preview"}},
		{"", api.SecretReference{Name: "x"}},
		{"preview", api.SecretReference{Name: "missing"}},
	} {
		if _, err := credentials(t.Context(), kube, tc.namespace, "registry.example", &tc.ref); err == nil {
			t.Fatal("invalid Secret selection accepted")
		}
	}
}

func TestExistingInputsAndArtifactConflict(t *testing.T) {
	cr := newScript("existing", "")
	cr.Spec.Artifact = nil
	content := "export default {}"
	cr.Spec.ForProvider.Content = &content
	if err := NewInitializer(fakeKube(t).Build()).Initialize(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	p, err := cr.GetParameters()
	if err != nil || p["content"] != content {
		t.Fatalf("existing inline input changed: %v", err)
	}
	cr.Spec.Artifact = &api.Source{Ref: "registry.example/bundle@sha256:" + strings.Repeat("a", 64)}
	if err := NewInitializer(fakeKube(t).Build()).Initialize(t.Context(), cr); err == nil {
		t.Fatal("conflicting artifact and content accepted")
	}
	cr.Spec.ForProvider.Content = nil
	cr.Spec.InitProvider.Content = &content
	if err := cr.ValidateWorkerArtifactInputs(); err == nil {
		t.Fatal("conflicting initProvider content accepted")
	}
}

func TestScriptGetsAllModules(t *testing.T) {
	m := api.Manifest{SchemaVersion: 1, MainModule: "worker.js", Modules: []api.Module{{Name: "worker.js", ContentType: "application/javascript+module"}, {Name: "helper.wasm", ContentType: "application/wasm"}}}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	fetcher := &memoryFetcher{data: bundle(t, map[string]string{manifestName: string(data), "worker.js": "export default {}", "helper.wasm": "wasm"})}
	resolver := &Resolver{Root: t.TempDir(), Fetcher: fetcher}
	ref := "registry.example/bundle@sha256:" + strings.Repeat("a", 64)
	initializer := resolver.Initializer(fakeKube(t).Build())
	s := newScript("script", ref)
	if err := initializer.Initialize(t.Context(), s); err != nil {
		t.Fatal(err)
	}
	sp, err := s.GetParameters()
	if err != nil {
		t.Fatal(err)
	}
	if len(sp["files"].(map[string]any)) != 1 || sp["main_module"] != "worker.js" {
		t.Fatal("Script lost multipart modules")
	}
}

func TestStaticOnlyBundle(t *testing.T) {
	fetcher := &memoryFetcher{data: bundle(t, map[string]string{manifestName: `{"schemaVersion":1,"assets":"public"}`, "public/index.html": "static"})}
	r := &Resolver{Root: t.TempDir(), Fetcher: fetcher}
	cr := newScript("static", "registry.example/bundle@sha256:"+strings.Repeat("a", 64))
	if err := r.Initializer(fakeKube(t).Build()).Initialize(t.Context(), cr); err != nil {
		t.Fatal(err)
	}
	p, err := cr.GetParameters()
	if err != nil {
		t.Fatal(err)
	}
	if p["content_file"] != nil || p["assets"] == nil {
		t.Fatal(fmt.Sprint("invalid static-only parameters: ", p))
	}
}

func TestCacheSweepPreservesActiveAndLockedEntries(t *testing.T) {
	fetcher := &memoryFetcher{data: standardBundle(t)}
	r := &Resolver{Root: t.TempDir(), Fetcher: fetcher, Retention: time.Hour}
	ref := "registry.example/bundle@sha256:" + strings.Repeat("a", 64)
	stale, err := r.materialize(t.Context(), "stale", ref, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	active, err := r.materialize(t.Context(), "active", ref, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(stale.Root, old, old); err != nil {
		t.Fatal(err)
	}
	key := filepath.Base(filepath.Dir(stale.Root))
	value, ok := r.locks.Load(key)
	if !ok {
		t.Fatal("missing resource lock")
	}
	mu := value.(*sync.Mutex)
	mu.Lock()
	r.lastSweep = time.Time{}
	r.collect()
	if _, err := os.Stat(stale.Root); err != nil {
		t.Fatal("sweep removed an in-use resource cache")
	}
	mu.Unlock()
	r.lastSweep = time.Time{}
	r.collect()
	if _, err := os.Stat(stale.Root); !os.IsNotExist(err) {
		t.Fatal("unused cache did not expire")
	}
	if _, err := cachedBundle(active.Root, ref); err != nil {
		t.Fatalf("active cache expired: %v", err)
	}
}

func TestArtifactPinChangeAndRollback(t *testing.T) {
	fetcher := &memoryFetcher{data: standardBundle(t)}
	r := &Resolver{Root: t.TempDir(), Fetcher: fetcher}
	firstRef := "registry.example/bundle@sha256:" + strings.Repeat("a", 64)
	secondRef := "registry.example/bundle@sha256:" + strings.Repeat("b", 64)
	first, err := r.materialize(t.Context(), "worker", firstRef, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	fetcher.data = bundle(t, map[string]string{manifestName: `{"schemaVersion":1,"mainModule":"worker.js","modules":[{"name":"worker.js","contentType":"application/javascript+module"}]}`, "worker.js": "export default {fetch(){return new Response('new')}}"})
	second, err := r.materialize(t.Context(), "worker", secondRef, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	rolledBack, err := r.materialize(t.Context(), "worker", firstRef, auth.EmptyCredential)
	if err != nil {
		t.Fatal(err)
	}
	if first.Root == second.Root || first.ContentSHA256 == second.ContentSHA256 || rolledBack.ContentSHA256 != first.ContentSHA256 || fetcher.calls.Load() != 2 {
		t.Fatal("artifact pin update or cached rollback selected the wrong bytes")
	}
}
