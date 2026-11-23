package workerartifact

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry/remote/auth"
)

func bundle(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o600, Typeflag: tar.TypeReg, Size: int64(len(content))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func standardBundle(t *testing.T) []byte {
	t.Helper()
	return bundle(t, map[string]string{
		manifestName:        `{"schemaVersion":1,"mainModule":"worker.js","modules":[{"name":"worker.js","contentType":"application/javascript+module"}],"assets":"public"}`,
		"worker.js":         `export default {fetch() { return new Response("test"); }}`,
		"public/index.html": "<h1>artifact</h1>",
	})
}

type testRegistry struct {
	server          *httptest.Server
	ref             string
	blob            []byte
	blobFetches     atomic.Int32
	unauthenticated atomic.Int32
	corrupt         atomic.Bool
}

func newRegistry(t *testing.T, compressed []byte, private bool) *testRegistry {
	t.Helper()
	layer := ocispec.Descriptor{MediaType: LayerType, Digest: digest.FromBytes(compressed), Size: int64(len(compressed))}
	manifest := ocispec.Manifest{Versioned: specs.Versioned{SchemaVersion: 2}, MediaType: ocispec.MediaTypeImageManifest, ArtifactType: ArtifactType, Layers: []ocispec.Descriptor{layer}}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	manifestDigest := digest.FromBytes(data).String()
	r := &testRegistry{blob: compressed}
	r.server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if private {
			user, password, ok := req.BasicAuth()
			if !ok || user != "worker" || password != "registry-password" {
				r.unauthenticated.Add(1)
				w.Header().Set("WWW-Authenticate", `Basic realm="worker-artifacts"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		switch req.URL.Path {
		case "/v2/bundles/manifests/" + manifestDigest:
			w.Header().Set("Content-Type", ocispec.MediaTypeImageManifest)
			w.Header().Set("Docker-Content-Digest", manifestDigest)
			w.Header().Set("Content-Length", fmt.Sprint(len(data)))
			if req.Method != http.MethodHead {
				_, _ = w.Write(data)
			}
		case "/v2/bundles/blobs/" + layer.Digest.String():
			r.blobFetches.Add(1)
			w.Header().Set("Content-Length", fmt.Sprint(len(compressed)))
			if r.corrupt.Load() {
				_, _ = w.Write(bytes.Repeat([]byte("x"), len(compressed)))
			} else {
				_, _ = w.Write(compressed)
			}
		default:
			http.NotFound(w, req)
		}
	}))
	r.ref = strings.TrimPrefix(r.server.URL, "https://") + "/bundles@" + manifestDigest
	t.Cleanup(r.server.Close)
	return r
}

func TestRegistryFetch(t *testing.T) {
	for _, private := range []bool{false, true} {
		t.Run(fmt.Sprint("private=", private), func(t *testing.T) {
			r := newRegistry(t, standardBundle(t), private)
			credential := auth.EmptyCredential
			if private {
				credential = auth.Credential{Username: "worker", Password: "registry-password"}
			}
			var out bytes.Buffer
			if err := (RegistryFetcher{Client: r.server.Client()}).Fetch(t.Context(), r.ref, credential, &out); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(out.Bytes(), r.blob) || r.blobFetches.Load() != 1 {
				t.Fatal("bundle was not fetched intact")
			}
		})
	}
}

func TestRegistryRejectsCorruptBlobAndMissingAuth(t *testing.T) {
	r := newRegistry(t, standardBundle(t), true)
	f := RegistryFetcher{Client: r.server.Client()}
	if err := f.Fetch(t.Context(), r.ref, auth.EmptyCredential, io.Discard); err == nil {
		t.Fatal("private registry accepted anonymous client")
	}
	r.corrupt.Store(true)
	if err := f.Fetch(t.Context(), r.ref, auth.Credential{Username: "worker", Password: "registry-password"}, io.Discard); err == nil {
		t.Fatal("corrupt blob accepted")
	}
}

func TestPinnedReferences(t *testing.T) {
	for _, ref := range []string{"ghcr.io/org/app:latest", "ghcr.io/org/app", "https://ghcr.io/org/app@sha256:" + strings.Repeat("a", 64), "ghcr.io/org/app@sha256:bad", "ghcr.io/org/app@sha512:" + strings.Repeat("a", 128)} {
		if _, err := parseReference(ref); err == nil {
			t.Errorf("accepted %q", ref)
		}
	}
	if _, err := parseReference("ghcr.io/org/app@sha256:" + strings.Repeat("a", 64)); err != nil {
		t.Fatal(err)
	}
}

func TestDockerCredentials(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("worker:password:with:colons"))
	data := []byte(`{"auths":{"https://ghcr.io/":{"auth":"` + encoded + `"}}}`)
	c, err := dockerCredential(data, "ghcr.io")
	if err != nil || c.Username != "worker" || c.Password != "password:with:colons" {
		t.Fatalf("credential decode failed: %v", err)
	}
	for _, data := range []string{`{`, `{"auths":{}}`, `{"auths":{"ghcr.io":{"auth":"bad"}}}`, `{"auths":{"ghcr.io":{}}}`} {
		if _, err := dockerCredential([]byte(data), "ghcr.io"); err == nil {
			t.Error("invalid credentials accepted")
		}
	}
	if _, err := dockerCredential([]byte(`{"auths":{"other.example":{"username":"a","password":"b"}}}`), "ghcr.io"); err == nil {
		t.Fatal("cross-host credentials accepted")
	}
}

func TestCanceledFetch(t *testing.T) {
	r := newRegistry(t, standardBundle(t), false)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := (RegistryFetcher{Client: r.server.Client()}).Fetch(ctx, r.ref, auth.EmptyCredential, io.Discard); err == nil {
		t.Fatal("canceled fetch succeeded")
	}
}

func TestDockerHubCredentialAlias(t *testing.T) {
	data := []byte(`{"auths":{"https://index.docker.io/v1/":{"username":"worker","password":"password"}}}`)
	c, err := dockerCredential(data, "docker.io")
	if err != nil || c.Username != "worker" || c.Password != "password" {
		t.Fatalf("Docker Hub credential alias not resolved: %v", err)
	}
	for _, host := range []string{"ghcr.io", "index.docker.io.evil.example"} {
		if _, err := dockerCredential(data, host); err == nil {
			t.Fatalf("Docker Hub credentials leaked to %s", host)
		}
	}
}
