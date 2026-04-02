# Developing provider-cloudflare

## Prerequisites

- [Go](https://go.dev/doc/install) 1.24+
- [OpenTofu](https://opentofu.org/) 1.9+ (required for schema generation; auto-installed by `make generate`)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports): `go install golang.org/x/tools/cmd/goimports@latest` and ensure it is on your PATH

## Code generation

Fetch submodules first, then run the code-generation pipeline. `make generate` runs the Upjet generator (same as `go run cmd/generator/main.go "$PWD"`).

```bash
make submodules
export PATH="$(go env GOPATH)/bin:$PATH"
make generate
```

Alternative:

```bash
go run cmd/generator/main.go "$PWD"
```

## Running locally

Run against a Kubernetes cluster:

```bash
make run
```

Build, push, and install:

```bash
make all
```

Build binary only:

```bash
make build
```

## Fixing CI (check-diff / lint)

If the CI **check-diff** or **lint** job fails, regenerate locally and commit the result:

```bash
make submodules
make generate
make check-diff
git add -A && git commit -m "chore: regenerate for CI" && git push
```

For CI parity on Apple Silicon or when your Go version differs from CI, use the dev container:

```bash
make check-diff.ci
```

## Adding resources

- Find the resource in the [Cloudflare Terraform Provider docs](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs).
- Prefer **1 resource per PR**.
- Write a test case for new or modified resources.
