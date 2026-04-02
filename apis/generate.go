//go:build generate
// +build generate

// NOTE: See the below link for details on what is happening here.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

// Remove existing CRDs
//go:generate rm -rf ../package/crds

// Remove generated files (|| true so already-missing files on bind mounts do not fail)
//go:generate bash -c "find . -iname 'zz_*' ! -iname 'zz_generated.managed*.go' -delete || true"
//go:generate bash -c "find . -type d -empty -delete"
//go:generate bash -c "find ../internal/controller -iname 'zz_*' -delete || true"
//go:generate bash -c "find ../internal/controller -type d -empty -delete"
//go:generate rm -rf ../examples-generated

// Generate documentation from Terraform docs. Use -mod=mod so the scraper
// can be resolved when the workspace uses -mod=vendor (scraper is not in vendor).
//go:generate go run -mod=mod github.com/crossplane/upjet/v2/cmd/scraper -n ${TERRAFORM_PROVIDER_SOURCE} -r ../.work/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_DOCS_PATH} -o ../config/provider-metadata.yaml

// Generate config/generated.lst (JSON array of resource names) for schema-version-diff.
//go:generate python3 ../scripts/generate_lst.py ../config/provider-metadata.yaml ../config/generated.lst

// Run Upjet generator
//go:generate go run ../cmd/generator/main.go ..

// Generate deepcopy methodsets and CRD manifests
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile=../hack/boilerplate.go.txt paths=./... crd:allowDangerousTypes=true,crdVersions=v1 output:artifacts:config=../package/crds

// Generate crossplane-runtime methodsets (resource.Claim, etc)
//go:generate go run -tags generate github.com/crossplane/crossplane-tools/cmd/angryjet generate-methodsets --header-file=../hack/boilerplate.go.txt ./...

package apis

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen" //nolint:typecheck

	_ "github.com/crossplane/crossplane-tools/cmd/angryjet" //nolint:typecheck
)
