# Copyright 2025 The Crossplane Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ====================================================================================
# Options

# Optional. The Go Binary to use
GO ?= go

# The go project including repo name, for example, github.com/rook/rook
GO_PROJECT ?= $(PROJECT_REPO)

# Optional. These are subdirs that we look for all go files to test, vet, and fmt
GO_SUBDIRS ?= cmd pkg

# Optional. Additional subdirs used for integration or e2e testings
GO_INTEGRATION_TESTS_SUBDIRS ?=

# Optional build flags passed to go tools
GO_BUILDFLAGS ?=
GO_LDFLAGS ?=
GO_TAGS ?=
GO_TEST_FLAGS ?=
GO_TEST_SUITE ?=
GO_COVER_MODE ?= count
GO_CGO_ENABLED ?= 0

# ====================================================================================
# Setup go environment

# turn on more verbose build when V=1
ifeq ($(V),1)
GO_LDFLAGS += -v -n
GO_BUILDFLAGS += -x
endif

# whether to generate debug information in binaries. this includes DWARF and symbol tables.
ifeq ($(DEBUG),0)
GO_LDFLAGS += -s -w
endif

# set GOOS and GOARCH
GOOS := $(OS)
GOARCH := $(ARCH)
export GOOS GOARCH

# set GOHOSTOS and GOHOSTARCH
GOHOSTOS := $(HOSTOS)
GOHOSTARCH := $(TARGETARCH)

GO_PACKAGES := $(foreach t,$(GO_SUBDIRS),$(GO_PROJECT)/$(t)/...)
GO_INTEGRATION_TEST_PACKAGES := $(foreach t,$(GO_INTEGRATION_TESTS_SUBDIRS),$(GO_PROJECT)/$(t)/integration)

ifneq ($(GO_TEST_PARALLEL),)
GO_TEST_FLAGS += -p $(GO_TEST_PARALLEL)
endif

ifneq ($(GO_TEST_SUITE),)
GO_TEST_FLAGS += -run '$(GO_TEST_SUITE)'
endif

GOPATH := $(shell $(GO) env GOPATH)

GOHOST := GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) $(GO)
GO_VERSION := $(shell $(GO) version | sed -ne 's/[^0-9]*\(\([0-9]\.\)\{0,4\}[0-9][^.]\).*/\1/p')

# We use a consistent version of golangci-lint to ensure everyone gets the same
# linters.
GOLANGCILINT_VERSION ?= 1.64.8
GOLANGCILINT := $(TOOLS_HOST_DIR)/golangci-lint-v$(GOLANGCILINT_VERSION)

GO_BIN_DIR := $(abspath $(OUTPUT_DIR)/bin)
GO_OUT_DIR := $(GO_BIN_DIR)/$(PLATFORM)
GO_TEST_DIR := $(abspath $(OUTPUT_DIR)/tests)
GO_TEST_OUTPUT := $(GO_TEST_DIR)/$(PLATFORM)
GO_LINT_DIR := $(abspath $(OUTPUT_DIR)/lint)
GO_LINT_OUTPUT := $(GO_LINT_DIR)/$(PLATFORM)

ifeq ($(GOOS),windows)
GO_OUT_EXT := .exe
endif

GO_COMMON_FLAGS = $(GO_BUILDFLAGS) -tags '$(GO_TAGS)' -trimpath
GO_STATIC_FLAGS = $(GO_COMMON_FLAGS) -installsuffix static -ldflags '$(GO_LDFLAGS)'
GO_GENERATE_FLAGS = $(GO_BUILDFLAGS) -tags 'generate $(GO_TAGS)'

# ====================================================================================
# Go Targets

go.build:
	@$(INFO) go build $(PLATFORM)
	$(foreach p,$(GO_STATIC_PACKAGES),@CGO_ENABLED=0 $(GO) build -v -o $(GO_OUT_DIR)/$(lastword $(subst /, ,$(p)))$(GO_OUT_EXT) $(GO_STATIC_FLAGS) $(p) || $(FAIL) ${\n})
	$(foreach p,$(GO_TEST_PACKAGES) $(GO_LONGHAUL_TEST_PACKAGES),@CGO_ENABLED=0 $(GO) test -c -o $(GO_TEST_OUTPUT)/$(lastword $(subst /, ,$(p)))$(GO_OUT_EXT) $(GO_STATIC_FLAGS) $(p) || $(FAIL) ${\n})
	@$(OK) go build $(PLATFORM)

go.install:
	@$(INFO) go install $(PLATFORM)
	$(foreach p,$(GO_STATIC_PACKAGES),@CGO_ENABLED=0 $(GO) install -v $(GO_STATIC_FLAGS) $(p) || $(FAIL) ${\n})
	@$(OK) go install $(PLATFORM)

go.test.unit:
	@$(INFO) go test unit-tests
	@mkdir -p $(GO_TEST_OUTPUT)
	@CGO_ENABLED=$(GO_CGO_ENABLED) $(GOHOST) test -cover $(GO_STATIC_FLAGS) $(GO_PACKAGES) || $(FAIL)
	@CGO_ENABLED=$(GO_CGO_ENABLED) $(GOHOST) test -v -covermode=$(GO_COVER_MODE) -coverprofile=$(GO_TEST_OUTPUT)/coverage.txt $(GO_TEST_FLAGS) $(GO_STATIC_FLAGS) $(GO_PACKAGES) 2>&1 | tee $(GO_TEST_OUTPUT)/unit-tests.log || $(FAIL)
	@$(OK) go test unit-tests

# TODO(negz): Do providers use this target? crossplane/crossplane doesn't.
# https://github.com/crossplane/crossplane/blob/master/Makefile#L135
go.test.integration:
	@$(INFO) go test integration-tests
	@mkdir -p $(GO_TEST_OUTPUT) || $(FAIL)
	@CGO_ENABLED=0 $(GOHOST) test $(GO_STATIC_FLAGS) $(GO_INTEGRATION_TEST_PACKAGES) || $(FAIL)
	@CGO_ENABLED=0 $(GOHOST) test $(GO_TEST_FLAGS) $(GO_STATIC_FLAGS) $(GO_INTEGRATION_TEST_PACKAGES) $(TEST_FILTER_PARAM) 2>&1 | tee $(GO_TEST_OUTPUT)/integration-tests.log || $(FAIL)
	@$(OK) go test integration-tests

go.lint: $(GOLANGCILINT)
	@$(INFO) golangci-lint
	@mkdir -p $(GO_LINT_OUTPUT)
	@$(GOLANGCILINT) run $(GO_LINT_ARGS) || $(FAIL)
	@$(OK) golangci-lint

go.modules.check:
	@$(INFO) verify go modules dependencies are tidy
	@$(GO) mod tidy
	@changed=$$(git diff --exit-code --name-only go.mod go.sum 2>&1) && [ -z "$${changed}" ] || (echo "go.mod is not tidy. Please run 'go mod tidy' and stage the changes" 1>&2; $(FAIL))
	@$(OK) go modules are tidy
	@$(INFO) verify go modules dependencies have expected content
	@$(GO) mod verify || $(FAIL)
	@$(OK) go modules dependencies verified

go.modules.download:
	@$(INFO) mod download
	@$(GO) mod download || $(FAIL)
	@$(OK) mod download

go.clean:
	@$(GO) clean -cache -testcache -modcache
	@rm -fr $(GO_BIN_DIR) $(GO_TEST_DIR)

go.generate:
	@$(INFO) go generate $(PLATFORM)
	@CGO_ENABLED=0 $(GOHOST) generate $(GO_GENERATE_FLAGS) $(GO_PACKAGES) $(GO_INTEGRATION_TEST_PACKAGES) || $(FAIL)
	@$(OK) go generate $(PLATFORM)
	@$(INFO) go mod tidy
	@$(GOHOST) mod tidy || $(FAIL)
	@$(OK) go mod tidy

.PHONY: go.build go.install go.test.unit go.test.integration go.lint go.vendor go.vendor.check go.clean go.generate

# ====================================================================================
# Common Targets

build.code.platform: go.build
clean: go.clean
lint.run: go.lint
test.run: go.test.unit
generate.run: go.generate

# ====================================================================================
# Special Targets

modules.download: go.modules.download
modules.check: go.modules.check

define GO_HELPTEXT
Go Targets:
    modules.download Download Go modules.
    modules.check    Fail the build if Go modules have changed.
endef
export GO_HELPTEXT

go.help:
	@echo "$$GO_HELPTEXT"

help-special: go.help

.PHONY: modules.download modules.check go.help

# ====================================================================================
# Tools install targets

$(GOLANGCILINT):
	@$(INFO) installing golangci-lint-v$(GOLANGCILINT_VERSION) $(SAFEHOSTPLATFORM)
	@mkdir -p $(TOOLS_HOST_DIR)/tmp-golangci-lint || $(FAIL)
	@curl -fsSL https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCILINT_VERSION)/golangci-lint-$(GOLANGCILINT_VERSION)-$(SAFEHOSTPLATFORM).tar.gz | tar -xz --strip-components=1 -C $(TOOLS_HOST_DIR)/tmp-golangci-lint || $(FAIL)
	@mv $(TOOLS_HOST_DIR)/tmp-golangci-lint/golangci-lint $(GOLANGCILINT) || $(FAIL)
	@rm -fr $(TOOLS_HOST_DIR)/tmp-golangci-lint
	@$(OK) installing golangci-lint-v$(GOLANGCILINT_VERSION) $(SAFEHOSTPLATFORM)

# These targets are deprecated but are still used by some existing repos.
# The modules targets should be used instead.
vendor: modules.download
vendor.check: modules.check
