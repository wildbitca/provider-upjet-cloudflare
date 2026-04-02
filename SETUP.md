# Setup Instructions

## Prerequisites

- **Go** 1.24 or later
- **Terraform** 1.5.x (required for schema generation; must be &lt; 1.6 due to BSL)
- **goimports**: `go install golang.org/x/tools/cmd/goimports@latest`

## Generate CRDs and Controllers

After cloning, run:

```bash
make submodules
make generate
```

This will:

1. Fetch the crossplane/build submodule
2. Generate Cloudflare provider schema (via `tofu init`)
3. Clone Cloudflare provider docs for scraping
4. Generate all 198 Crossplane managed resources, CRDs, controllers, and examples

## Run Locally

```bash
make run
```

## Build Package

```bash
make all
```
