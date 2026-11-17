// Package workerartifact materializes immutable OCI bundles for Worker resources.
package workerartifact

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

const (
	// ArtifactType identifies the immutable OCI bundle format.
	ArtifactType = "application/vnd.gitkb.cloudflare.worker.v1"
	// LayerType identifies its single gzip-compressed tar layer.
	LayerType       = "application/vnd.gitkb.cloudflare.worker.bundle.v1.tar+gzip"
	maxManifestSize = 1 << 20
	maxArchiveSize  = 128 << 20
)

// Fetcher writes a verified, compressed bundle to dst.
type Fetcher interface {
	Fetch(context.Context, string, auth.Credential, io.Writer) error
}

// RegistryFetcher uses the OCI distribution protocol, including bearer challenges.
type RegistryFetcher struct {
	Client *http.Client
}

func parseReference(ref string) (registry.Reference, error) {
	r, err := registry.ParseReference(ref)
	if err != nil {
		return r, fmt.Errorf("invalid OCI artifact reference: %w", err)
	}
	d, err := digest.Parse(r.Reference)
	if err != nil || d.Algorithm() != digest.SHA256 || d.Validate() != nil {
		return r, fmt.Errorf("worker artifact ref must be pinned as registry/repository@sha256:<64 lowercase hex characters>")
	}
	return r, nil
}

// Fetch verifies both the manifest pin and the bundle blob before it can be used.
func (f RegistryFetcher) Fetch(ctx context.Context, ref string, credential auth.Credential, dst io.Writer) error {
	repo, parsed, err := f.repository(ref, credential)
	if err != nil {
		return err
	}
	desc, err := resolveManifest(ctx, repo, parsed.Reference)
	if err != nil {
		return err
	}
	data, err := verifiedManifest(ctx, repo, desc)
	if err != nil {
		return err
	}
	layer, err := bundleLayer(data)
	if err != nil {
		return err
	}
	return downloadLayer(ctx, repo, layer, dst)
}

func (f RegistryFetcher) repository(ref string, credential auth.Credential) (*remote.Repository, registry.Reference, error) {
	parsed, err := parseReference(ref)
	if err != nil {
		return nil, parsed, err
	}
	repo, err := remote.NewRepository(parsed.Registry + "/" + parsed.Repository)
	if err != nil {
		return nil, parsed, fmt.Errorf("open OCI repository: %w", err)
	}
	client := f.Client
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Minute}
	}
	repo.Client = &auth.Client{Client: client, Cache: auth.NewCache(), Credential: auth.StaticCredential(parsed.Registry, credential)}
	return repo, parsed, nil
}

func resolveManifest(ctx context.Context, repo *remote.Repository, expected string) (ocispec.Descriptor, error) {
	desc, err := repo.Resolve(ctx, expected)
	if err != nil {
		return desc, fmt.Errorf("resolve worker artifact: %w", err)
	}
	if desc.Digest.String() != expected || desc.Size < 0 || desc.Size > maxManifestSize || desc.MediaType != ocispec.MediaTypeImageManifest {
		return desc, fmt.Errorf("worker artifact must be a bounded OCI image manifest at the requested digest")
	}
	return desc, nil
}

func verifiedManifest(ctx context.Context, repo *remote.Repository, desc ocispec.Descriptor) ([]byte, error) {
	rc, err := repo.Fetch(ctx, desc)
	if err != nil {
		return nil, fmt.Errorf("fetch OCI manifest: %w", err)
	}
	data, readErr := io.ReadAll(io.LimitReader(rc, maxManifestSize+1))
	closeErr := rc.Close()
	if readErr != nil || closeErr != nil {
		return nil, fmt.Errorf("read OCI manifest: %w", firstError(readErr, closeErr))
	}
	if int64(len(data)) != desc.Size || digest.FromBytes(data) != desc.Digest {
		return nil, fmt.Errorf("OCI manifest digest or size mismatch")
	}
	return data, nil
}

func bundleLayer(data []byte) (ocispec.Descriptor, error) {
	var manifest ocispec.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("decode OCI manifest: %w", err)
	}
	if manifest.SchemaVersion != 2 || manifest.ArtifactType != ArtifactType || len(manifest.Layers) != 1 {
		return ocispec.Descriptor{}, fmt.Errorf("expected %s with exactly one bundle layer", ArtifactType)
	}
	layer := manifest.Layers[0]
	return layer, validateLayer(layer)
}

func validateLayer(layer ocispec.Descriptor) error {
	if layer.MediaType != LayerType || layer.Size <= 0 || layer.Size > maxArchiveSize || len(layer.URLs) != 0 || layer.Digest.Algorithm() != digest.SHA256 || layer.Digest.Validate() != nil {
		return fmt.Errorf("invalid worker bundle layer type, size, digest or external URL")
	}
	return nil
}

func downloadLayer(ctx context.Context, repo *remote.Repository, layer ocispec.Descriptor, dst io.Writer) error {
	blob, err := repo.Fetch(ctx, layer)
	if err != nil {
		return fmt.Errorf("fetch worker bundle: %w", err)
	}
	defer blob.Close() //nolint:errcheck // A read-only response body; copy errors are checked below.
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(dst, h), io.LimitReader(blob, layer.Size+1))
	if err != nil {
		return fmt.Errorf("read worker bundle: %w", err)
	}
	if n != layer.Size || hex.EncodeToString(h.Sum(nil)) != layer.Digest.Encoded() {
		return fmt.Errorf("worker bundle digest or size mismatch")
	}
	return nil
}

func firstError(a, b error) error {
	if a != nil {
		return a
	}
	return b
}

func dockerCredential(data []byte, host string) (auth.Credential, error) {
	var config struct {
		Auths map[string]struct {
			Username      string `json:"username"`
			Password      string `json:"password"`
			Auth          string `json:"auth"`
			IdentityToken string `json:"identitytoken"`
		} `json:"auths"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return auth.EmptyCredential, fmt.Errorf("invalid dockerconfigjson registry credentials")
	}
	for key, entry := range config.Auths {
		registryHost := strings.TrimSuffix(strings.TrimPrefix(key, "https://"), "/")
		if host == "docker.io" && registryHost == "index.docker.io/v1" {
			registryHost = "docker.io"
		}
		if registryHost != host {
			continue
		}
		credential := auth.Credential{Username: entry.Username, Password: entry.Password, RefreshToken: entry.IdentityToken}
		if entry.Auth != "" {
			decoded, err := base64.StdEncoding.DecodeString(entry.Auth)
			if err != nil {
				return credential, fmt.Errorf("invalid encoded registry auth")
			}
			var ok bool
			credential.Username, credential.Password, ok = strings.Cut(string(decoded), ":")
			if !ok {
				return credential, fmt.Errorf("registry auth must encode username:password")
			}
		}
		if credential == auth.EmptyCredential {
			return credential, fmt.Errorf("empty credentials for artifact registry")
		}
		return credential, nil
	}
	return auth.EmptyCredential, fmt.Errorf("pull Secret has no credentials for artifact registry %s", host)
}
