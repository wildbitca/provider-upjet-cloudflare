package workerartifact

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"oras.land/oras-go/v2/registry/remote/auth"
)

func TestUnsafeArchives(t *testing.T) {
	for _, tc := range []struct {
		name string
		kind byte
		size int64
	}{
		{"../escape", tar.TypeReg, 0},
		{"/absolute", tar.TypeReg, 0},
		{"nested/../../escape", tar.TypeReg, 0},
		{"back\\slash", tar.TypeReg, 0},
		{"link", tar.TypeSymlink, 0},
		{"hardlink", tar.TypeLink, 0},
		{"device", tar.TypeChar, 0},
		{"huge", tar.TypeReg, maxExpandedSize + 1},
		{markerName, tar.TypeReg, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			tw := tar.NewWriter(gz)
			if err := tw.WriteHeader(&tar.Header{Name: tc.name, Typeflag: tc.kind, Size: tc.size, Linkname: "../escape"}); err != nil {
				t.Fatal(err)
			}
			_ = tw.Close()
			_ = gz.Close()
			if err := extract(bytes.NewReader(buf.Bytes()), t.TempDir()); err == nil {
				t.Fatal("unsafe archive accepted")
			}
		})
	}
}

func TestBundleManifestValidation(t *testing.T) {
	for _, manifest := range []string{
		`{"schemaVersion":2,"assets":"public"}`,
		`{"schemaVersion":1}`,
		`{"schemaVersion":1,"mainModule":"missing.js","modules":[{"name":"worker.js","contentType":"application/javascript+module"}]}`,
		`{"schemaVersion":1,"mainModule":"worker.js","modules":[{"name":"../worker.js","contentType":"application/javascript+module"}]}`,
		`{"schemaVersion":1,"assets":"../outside"}`,
		`{"schemaVersion":1,"assets":"worker.js"}`,
		`{"schemaVersion":1,"mainModule":"public/worker.js","modules":[{"name":"public/worker.js","contentType":"application/javascript+module"}],"assets":"public"}`,
		`{"schemaVersion":1,"assets":"public","unknown":true}`,
	} {
		t.Run(manifest, func(t *testing.T) {
			root := t.TempDir()
			data := bundle(t, map[string]string{manifestName: manifest, "worker.js": "worker", "public/worker.js": "worker"})
			if err := extract(bytes.NewReader(data), root); err != nil {
				t.Fatal(err)
			}
			if _, err := inspectBundle(root); err == nil {
				t.Fatal("invalid manifest accepted")
			}
		})
	}
}

func TestExtractionRejectsDuplicatePaths(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for range 2 {
		if err := tw.WriteHeader(&tar.Header{Name: "duplicate", Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
	}
	_ = tw.Close()
	_ = gz.Close()
	if err := extract(bytes.NewReader(buf.Bytes()), t.TempDir()); err == nil {
		t.Fatal("duplicate paths accepted")
	}
}

func TestFailedMaterializationDoesNotPublishPartialCache(t *testing.T) {
	fetcher := &memoryFetcher{data: bundle(t, map[string]string{"worker.js": "missing manifest"})}
	r := &Resolver{Root: t.TempDir(), Fetcher: fetcher}
	if _, err := r.materialize(t.Context(), "resource", "registry.example/bundle@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", auth.EmptyCredential); err == nil {
		t.Fatal("invalid bundle materialized")
	}
	if err := filepath.WalkDir(r.Root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == markerName || (d.IsDir() && len(d.Name()) > 10 && d.Name()[:10] == ".download-") {
			t.Errorf("partial cache survived at %s", p)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGzipChecksumIsVerified(t *testing.T) {
	data := standardBundle(t)
	data[len(data)-1] ^= 0xff
	if err := extract(bytes.NewReader(data), t.TempDir()); err == nil {
		t.Fatal("corrupt gzip trailer was accepted")
	}
}

func TestExtractionAcceptsRootDirectoryHeaders(t *testing.T) {
	for _, name := range []string{".", "./"} {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			tw := tar.NewWriter(gz)
			if err := tw.WriteHeader(&tar.Header{Name: name, Typeflag: tar.TypeDir}); err != nil {
				t.Fatal(err)
			}
			if err := tw.WriteHeader(&tar.Header{Name: "./worker.js", Typeflag: tar.TypeReg, Size: 6}); err != nil {
				t.Fatal(err)
			}
			if _, err := tw.Write([]byte("worker")); err != nil {
				t.Fatal(err)
			}
			if err := tw.Close(); err != nil {
				t.Fatal(err)
			}
			if err := gz.Close(); err != nil {
				t.Fatal(err)
			}
			root := t.TempDir()
			if err := extract(bytes.NewReader(buf.Bytes()), root); err != nil {
				t.Fatal(err)
			}
			if data, err := os.ReadFile(filepath.Join(root, "worker.js")); err != nil || string(data) != "worker" {
				t.Fatalf("regular file following root directory not extracted: %v", err)
			}
		})
	}
}
