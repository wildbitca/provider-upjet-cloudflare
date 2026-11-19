package workerartifact

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	api "github.com/wildbitca/provider-upjet-cloudflare/apis/workerartifact"
)

const (
	manifestName    = "worker-artifact.json"
	markerName      = ".artifact-ready.json"
	maxExpandedSize = 256 << 20
	maxFiles        = 10000
)

func validPath(name string) bool {
	return fs.ValidPath(name) && name != "." && !strings.ContainsAny(name, "\\\x00") && name != markerName
}

func extract(archive io.Reader, target string) error {
	gz, err := gzip.NewReader(archive)
	if err != nil {
		return fmt.Errorf("open worker bundle gzip: %w", err)
	}
	defer gz.Close() //nolint:errcheck // Decompression errors are checked while reading.
	root, err := os.OpenRoot(target)
	if err != nil {
		return err
	}
	defer root.Close() //nolint:errcheck // Individual file operations are checked.
	// Bound payload, metadata and padding, including ignored tar entries.
	limited := &io.LimitedReader{R: gz, N: maxExpandedSize + maxFiles*2048 + 1}
	return extractTar(tar.NewReader(limited), limited, root)
}

func extractTar(tr *tar.Reader, limited *io.LimitedReader, root *os.Root) error {
	seen := map[string]bool{}
	var expanded int64
	for count := 0; ; count++ {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return finishGzip(limited)
		}
		if err != nil {
			return fmt.Errorf("read worker bundle tar: %w", err)
		}
		if err := validateHeader(count, hdr); err != nil {
			return err
		}
		if err := extractEntry(root, tr, hdr, seen, &expanded); err != nil {
			return err
		}
	}
}

func finishGzip(limited *io.LimitedReader) error {
	// Reading to EOF validates gzip's checksum; tar may leave trailing padding.
	if _, err := io.Copy(io.Discard, limited); err != nil {
		return fmt.Errorf("finish worker bundle gzip: %w", err)
	}
	if limited.N == 0 {
		return fmt.Errorf("worker bundle exceeds decompression limit")
	}
	return nil
}

func validateHeader(count int, hdr *tar.Header) error {
	if count >= maxFiles || hdr.Size < 0 || (hdr.Typeflag == tar.TypeDir && hdr.Size != 0) {
		return fmt.Errorf("invalid or excessive worker bundle entries")
	}
	return nil
}

func extractEntry(root *os.Root, tr *tar.Reader, hdr *tar.Header, seen map[string]bool, expanded *int64) error {
	if (hdr.Name == "." || hdr.Name == "./") && hdr.Typeflag == tar.TypeDir {
		return nil
	}
	name := strings.TrimSuffix(strings.TrimPrefix(hdr.Name, "./"), "/")
	if !validPath(name) || seen[name] {
		return fmt.Errorf("invalid or duplicate worker bundle entry %q", hdr.Name)
	}
	seen[name] = true
	if hdr.Typeflag == tar.TypeDir {
		return root.MkdirAll(name, 0o700)
	}
	if hdr.Typeflag != tar.TypeReg || hdr.Size > maxExpandedSize-*expanded {
		return fmt.Errorf("worker bundle only accepts bounded regular files and directories")
	}
	*expanded += hdr.Size
	return writeEntry(root, tr, name, hdr.Size)
}

func writeEntry(root *os.Root, reader io.Reader, name string, size int64) error {
	if err := root.MkdirAll(path.Dir(name), 0o700); err != nil {
		return err
	}
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	_, copyErr := io.CopyN(f, reader, size)
	return firstError(copyErr, f.Close())
}

func readManifest(root string) (api.Manifest, error) {
	var manifest api.Manifest
	f, err := os.OpenInRoot(root, manifestName)
	if err != nil {
		return manifest, fmt.Errorf("open %s: %w", manifestName, err)
	}
	defer f.Close() //nolint:errcheck // Read-only local manifest.
	info, err := f.Stat()
	if err != nil {
		return manifest, err
	}
	if info.Size() > maxManifestSize {
		return manifest, fmt.Errorf("worker bundle manifest exceeds its size limit")
	}
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return manifest, fmt.Errorf("decode worker bundle manifest: %w", err)
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return manifest, fmt.Errorf("worker bundle manifest must contain exactly one JSON object")
	}
	return manifest, validateManifestShape(manifest)
}

func validateManifestShape(manifest api.Manifest) error {
	if manifest.SchemaVersion != 1 {
		return fmt.Errorf("worker bundle schemaVersion must be 1")
	}
	if len(manifest.Modules) == 0 && (manifest.Assets == "" || manifest.MainModule != "") {
		return fmt.Errorf("worker bundle must provide modules or assets")
	}
	return nil
}

func inspectBundle(root string) (*api.Files, error) {
	manifest, err := readManifest(root)
	if err != nil {
		return nil, err
	}
	files := &api.Files{Root: root, Manifest: manifest, Hashes: map[string]string{}}
	hasher := &bundleHasher{root: root, hashes: files.Hashes}
	if err := filepath.WalkDir(root, hasher.walk); err != nil {
		return nil, err
	}
	if err := validateModules(manifest, files.Hashes); err != nil {
		return nil, err
	}
	files.ContentSHA256 = files.Hashes[manifest.MainModule]
	if err := validateAssets(root, manifest); err != nil {
		return nil, err
	}
	return files, nil
}

type bundleHasher struct {
	root   string
	hashes map[string]string
	size   int64
}

func (b *bundleHasher) walk(p string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if entry.IsDir() || (entry.Name() == markerName && filepath.Dir(p) == b.root) {
		return nil
	}
	if !entry.Type().IsRegular() {
		return fmt.Errorf("non-regular cached bundle file")
	}
	return b.record(p, entry)
}

func (b *bundleHasher) record(p string, entry fs.DirEntry) error {
	info, err := entry.Info()
	if err != nil {
		return err
	}
	if len(b.hashes) >= maxFiles || info.Size() > maxExpandedSize-b.size {
		return fmt.Errorf("worker artifact cache exceeds bundle limits")
	}
	b.size += info.Size()
	rel, err := filepath.Rel(b.root, p)
	if err != nil {
		return err
	}
	hash, err := hashFile(b.root, rel)
	if err == nil {
		b.hashes[filepath.ToSlash(rel)] = hash
	}
	return err
}

func validateModules(manifest api.Manifest, hashes map[string]string) error {
	seen := map[string]bool{}
	for _, module := range manifest.Modules {
		if err := validateModule(module, hashes, seen); err != nil {
			return err
		}
		seen[module.Name] = true
	}
	if len(manifest.Modules) > 0 && (!seen[manifest.MainModule] || manifest.MainModule == "") {
		return fmt.Errorf("mainModule must name a bundled module")
	}
	return nil
}

func validateModule(module api.Module, hashes map[string]string, seen map[string]bool) error {
	_, _, mimeErr := mime.ParseMediaType(module.ContentType)
	if !validPath(module.Name) || module.Name == manifestName || mimeErr != nil || seen[module.Name] || hashes[module.Name] == "" {
		return fmt.Errorf("invalid, duplicate or missing worker module %q", module.Name)
	}
	return nil
}

func validateAssets(root string, manifest api.Manifest) error {
	if manifest.Assets == "" {
		return nil
	}
	if !validPath(manifest.Assets) {
		return fmt.Errorf("invalid assets directory")
	}
	stat, err := os.Stat(filepath.Join(root, manifest.Assets))
	if err != nil || !stat.IsDir() {
		return fmt.Errorf("bundle assets path must name a directory")
	}
	for _, module := range manifest.Modules {
		if strings.HasPrefix(module.Name, manifest.Assets+"/") {
			return fmt.Errorf("worker modules must be outside the public assets directory")
		}
	}
	return nil
}

func hashFile(root, rel string) (string, error) {
	f, err := os.OpenInRoot(root, rel)
	if err != nil {
		return "", err
	}
	defer f.Close() //nolint:errcheck // Read-only hash input.
	h := sha256.New()
	if _, err := io.Copy(h, io.LimitReader(f, maxExpandedSize+1)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
