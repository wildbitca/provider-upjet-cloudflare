package workerartifact

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"oras.land/oras-go/v2/registry/remote/auth"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"
)

// Resource is implemented by generated Script API adapters.
type Resource interface {
	xpresource.Managed
	GetWorkerArtifact() *api.Source
	GetWorkerArtifactFiles() *api.Files
	SetWorkerArtifactFiles(*api.Files)
	SetResolvedWorkerArtifact(string)
	ValidateWorkerArtifactInputs() error
}

// Resolver owns disposable files, isolated by managed resource UID and digest.
// Per-resource locks coalesce concurrent reconciliations without serializing
// independent previews. Nothing in this cache is required to delete a Worker.
type Resolver struct {
	Root      string
	Fetcher   Fetcher
	locks     sync.Map
	sweepMu   sync.Mutex
	lastSweep time.Time
	// Retention defaults to 24 hours since a cache entry's last use.
	Retention time.Duration
}

var defaultResolver = &Resolver{
	Root: filepath.Join(os.TempDir(), "provider-cloudflare-worker-artifacts"), Fetcher: RegistryFetcher{},
}

// NewInitializer attaches bundle preparation to the existing managed reconciler.
func NewInitializer(kube client.Client) managed.Initializer {
	return defaultResolver.Initializer(kube)
}

// Initializer is also used with isolated resolvers in integration tests.
func (r *Resolver) Initializer(kube client.Client) managed.Initializer {
	return managed.InitializerFn(func(ctx context.Context, mg xpresource.Managed) error {
		resource, ok := mg.(Resource)
		if !ok || resource.GetWorkerArtifact() == nil || mg.GetDeletionTimestamp() != nil {
			return nil
		}
		if err := resource.ValidateWorkerArtifactInputs(); err != nil {
			return err
		}
		source := resource.GetWorkerArtifact()
		ref, err := parseReference(source.Ref)
		if err != nil {
			return err
		}
		credential, err := credentials(ctx, kube, mg.GetNamespace(), ref.Registry, source.PullSecretRef)
		if err != nil {
			return err
		}
		files, err := r.materialize(ctx, string(mg.GetUID()), source.Ref, credential)
		if err != nil {
			return err
		}
		resource.SetWorkerArtifactFiles(files)
		resource.SetResolvedWorkerArtifact(ref.Reference)
		return nil
	})
}

func credentials(ctx context.Context, kube client.Client, namespace, host string, ref *api.SecretReference) (auth.Credential, error) {
	if ref == nil {
		return auth.EmptyCredential, nil
	}
	if ref.Name == "" || (namespace != "" && ref.Namespace != "" && ref.Namespace != namespace) {
		return auth.EmptyCredential, fmt.Errorf("artifact pull Secret must have a name and stay in the managed resource namespace")
	}
	if ref.Namespace != "" {
		namespace = ref.Namespace
	}
	if namespace == "" {
		return auth.EmptyCredential, fmt.Errorf("cluster-scoped Worker artifact pull Secret requires a namespace")
	}
	secret := &corev1.Secret{}
	if err := kube.Get(ctx, types.NamespacedName{Namespace: namespace, Name: ref.Name}, secret); err != nil {
		return auth.EmptyCredential, fmt.Errorf("read artifact pull Secret: %w", err)
	}
	if secret.Type != corev1.SecretTypeDockerConfigJson {
		return auth.EmptyCredential, fmt.Errorf("artifact pull Secret must have type kubernetes.io/dockerconfigjson")
	}
	return dockerCredential(secret.Data[corev1.DockerConfigJsonKey], host)
}

func (r *Resolver) materialize(ctx context.Context, uid, ref string, credential auth.Credential) (*api.Files, error) {
	if uid == "" {
		return nil, fmt.Errorf("worker artifact requires a managed resource UID")
	}
	key := sha256.Sum256([]byte(uid))
	keyString := hex.EncodeToString(key[:])
	value, _ := r.locks.LoadOrStore(keyString, &sync.Mutex{})
	mu := value.(*sync.Mutex)
	mu.Lock()
	defer r.collect()
	defer mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	refHash := sha256.Sum256([]byte(ref))
	parent := filepath.Join(r.Root, keyString)
	root := filepath.Join(parent, hex.EncodeToString(refHash[:]))
	if files, err := cachedBundle(root, ref); err == nil {
		now := time.Now()
		if err := os.Chtimes(root, now, now); err != nil {
			return nil, err
		}
		return files, nil
	}
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return nil, err
	}
	stage, err := os.MkdirTemp(parent, ".download-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(stage) //nolint:errcheck // Best-effort cleanup of disposable staging files.
	files, err := r.prepareBundle(ctx, stage, ref, credential)
	if err != nil {
		return nil, err
	}
	extracted := files.Root
	files.Root = root
	if err := publishBundle(extracted, ref, files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *Resolver) prepareBundle(ctx context.Context, stage, ref string, credential auth.Credential) (*api.Files, error) {
	archive, err := os.CreateTemp(stage, "bundle-*.tgz")
	if err != nil {
		return nil, err
	}
	defer archive.Close() //nolint:errcheck // Read/write operation errors are checked below.
	if err := r.Fetcher.Fetch(ctx, ref, credential, archive); err != nil {
		return nil, err
	}
	if _, err := archive.Seek(0, 0); err != nil {
		return nil, err
	}
	extracted := filepath.Join(stage, "files")
	if err := os.Mkdir(extracted, 0o700); err != nil {
		return nil, err
	}
	if err := extract(archive, extracted); err != nil {
		return nil, err
	}
	return inspectBundle(extracted)
}

func publishBundle(extracted, ref string, files *api.Files) error {
	marker, err := json.Marshal(ready{Ref: ref, Files: files})
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(extracted, markerName), marker, 0o600); err != nil {
		return err
	}
	if err := os.RemoveAll(files.Root); err != nil {
		return err
	}
	return os.Rename(extracted, files.Root)
}

// collect retires unused cache versions and abandoned downloads. Taking the
// same UID lock as materialize prevents removal during cache recovery or fetch.
// Active entries are touched on every reconciliation; retention greatly exceeds
// the managed reconciler and HTTP request timeouts.
func (r *Resolver) collect() {
	if !r.sweepMu.TryLock() {
		return
	}
	defer r.sweepMu.Unlock()
	now := time.Now()
	if now.Sub(r.lastSweep) < time.Hour {
		return
	}
	r.lastSweep = now
	retention := r.Retention
	if retention <= 0 {
		retention = 24 * time.Hour
	}
	entries, err := os.ReadDir(r.Root)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() || len(entry.Name()) != 64 {
			continue
		}
		r.collectParent(entry.Name(), now, retention)
	}
}

func (r *Resolver) collectParent(key string, now time.Time, retention time.Duration) {
	value, _ := r.locks.LoadOrStore(key, &sync.Mutex{})
	mu := value.(*sync.Mutex)
	if !mu.TryLock() {
		return
	}
	defer mu.Unlock()
	parent := filepath.Join(r.Root, key)
	versions, err := os.ReadDir(parent)
	if err != nil {
		return
	}
	for _, version := range versions {
		info, statErr := version.Info()
		if statErr == nil && version.IsDir() && now.Sub(info.ModTime()) > retention {
			_ = os.RemoveAll(filepath.Join(parent, version.Name())) // Disposable cache; retry on a later sweep.
		}
	}
	_ = os.Remove(parent) // Remove empty UID directories only.

}

type ready struct {
	Ref   string     `json:"ref"`
	Files *api.Files `json:"files"`
}

func cachedBundle(root, ref string) (*api.Files, error) {
	f, err := os.OpenInRoot(root, markerName)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck // Read-only cache marker.
	data, err := io.ReadAll(io.LimitReader(f, 16<<20))
	if err != nil {
		return nil, err
	}
	var marker ready
	if err := json.Unmarshal(data, &marker); err != nil {
		return nil, err
	}
	if marker.Ref != ref || marker.Files == nil || marker.Files.Root != root {
		return nil, fmt.Errorf("stale Worker artifact cache")
	}
	actual, err := inspectBundle(root)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(actual, marker.Files) {
		return nil, fmt.Errorf("worker artifact cache contents changed")
	}
	return actual, nil
}
