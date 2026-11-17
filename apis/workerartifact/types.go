package workerartifact

// Source selects a prebuilt Worker bundle. Environment-specific settings remain
// in forProvider; no application build runs in the provider.
type Source struct {
	// Ref is an OCI manifest reference pinned by SHA-256 digest, without a URL scheme.
	// +kubebuilder:validation:Pattern=`^[^[:space:]@]+@sha256:[a-f0-9]{64}$`
	Ref string `json:"ref"`
	// PullSecretRef selects a kubernetes.io/dockerconfigjson Secret. Namespaced
	// resources can only use their own namespace. Cluster resources must specify it.
	// +optional
	PullSecretRef *SecretReference `json:"pullSecretRef,omitempty"`
}

// SecretReference locates registry credentials, separately from Cloudflare credentials.
type SecretReference struct {
	// Name identifies the Secret.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Namespace defaults to the managed resource's namespace.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// Manifest describes files inside a Worker bundle, not its deployment settings.
type Manifest struct {
	SchemaVersion int      `json:"schemaVersion"`
	MainModule    string   `json:"mainModule,omitempty"`
	Modules       []Module `json:"modules,omitempty"`
	Assets        string   `json:"assets,omitempty"`
}

// Module identifies a file's Worker module name and MIME type.
type Module struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
}

// Files is disposable controller-local state. It is never serialized in a
// managed resource. The immutable Source is sufficient to recover it.
type Files struct {
	Root          string            `json:"root"`
	Manifest      Manifest          `json:"manifest"`
	ContentSHA256 string            `json:"contentSHA256,omitempty"`
	Hashes        map[string]string `json:"hashes"`
}
