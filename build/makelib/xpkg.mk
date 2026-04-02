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

ifeq ($(origin XPKG_DIR),undefined)
XPKG_DIR := $(ROOT_DIR)/package
endif

ifeq ($(origin XPKG_EXAMPLES_DIR),undefined)
XPKG_EXAMPLES_DIR := $(ROOT_DIR)/examples
endif

ifeq ($(origin XPKG_IGNORE),undefined)
XPKG_IGNORE := ''
endif

ifeq ($(origin XPKG_OUTPUT_DIR),undefined)
XPKG_OUTPUT_DIR := $(OUTPUT_DIR)/xpkg
endif

# shasum is not available on all systems. In that case, fall back to sha256sum.
ifneq ($(shell type shasum 2>/dev/null),)
SHA256SUM := shasum -a 256
else
SHA256SUM := sha256sum
endif

# a registry that is scoped to the current build tree on this host. this enables
# us to have isolation between concurrent builds on the same system, as in the case
# of multiple working directories or on a CI system with multiple executors. All images
# tagged with this build registry can safely be untagged/removed at the end of the build.
ifeq ($(origin BUILD_REGISTRY), undefined)
BUILD_REGISTRY := build-$(shell echo $(HOSTNAME)-$(ROOT_DIR) | $(SHA256SUM) | cut -c1-8)
endif

XPKG_REG_ORGS ?= xpkg.crossplane.io/crossplane
XPKG_REG_ORGS_NO_PROMOTE ?= xpkg.crossplane.io/crossplane
XPKG_LINUX_PLATFORMS := $(filter linux_%,$(PLATFORMS))
XPKG_ARCHS := $(subst linux_,,$(filter linux_%,$(PLATFORMS)))
XPKG_PLATFORMS := $(subst _,/,$(subst $(SPACE),$(COMMA),$(filter linux_%,$(PLATFORMS))))
XPKG_PLATFORMS_LIST := $(subst _,/,$(filter linux_%,$(PLATFORMS)))
XPKG_PLATFORM := $(subst _,/,$(PLATFORM))

XPKG_CLEANUP_EXAMPLES_ENABLED ?= false
XPKG_CLEANUP_EXAMPLES_VERSION ?= v0.12.1
XPKG_PROCESSED_EXAMPLES_DIR=$(XPKG_EXAMPLES_DIR)
ifeq ($(XPKG_CLEANUP_EXAMPLES_ENABLED),true)
	XPKG_PROCESSED_EXAMPLES_DIR=$(WORK_DIR)/xpkg-cleaned-examples
endif

# =====================================================================================
# XPKG Targets

# 1: xpkg
define xpkg.build.targets
xpkg.build.$(1):
ifeq ($(XPKG_CLEANUP_EXAMPLES_ENABLED),true)
	@rm -rf $(WORK_DIR)/xpkg-cleaned-examples
	@GOOS=$(HOSTOS) GOARCH=$(TARGETARCH) go run github.com/upbound/uptest/cmd/cleanupexamples@$(XPKG_CLEANUP_EXAMPLES_VERSION) $(XPKG_EXAMPLES_DIR) $(XPKG_PROCESSED_EXAMPLES_DIR) || $(FAIL)
endif
	@$(INFO) Building package $(1)-$(VERSION).xpkg for $(PLATFORM)
	@mkdir -p $(OUTPUT_DIR)/xpkg/$(PLATFORM)
	@controller_arg=$$$$(grep -E '^kind:\s+Provider\s*$$$$' $(XPKG_DIR)/crossplane.yaml > /dev/null && echo "--embed-runtime-image $(BUILD_REGISTRY)/$(1)-$(ARCH)"); \
	$(CROSSPLANE_CLI) xpkg build \
		$$$${controller_arg} \
		--package-root $(XPKG_DIR) \
		--examples-root $(XPKG_PROCESSED_EXAMPLES_DIR) \
		--ignore $(XPKG_IGNORE) \
		--package-file $(XPKG_OUTPUT_DIR)/$(PLATFORM)/$(1)-$(VERSION).xpkg || $(FAIL)
	@$(OK) Built package $(1)-$(VERSION).xpkg for $(PLATFORM)
ifeq ($(XPKG_CLEANUP_EXAMPLES_ENABLED),true)
	@rm -rf $(WORK_DIR)/xpkg-cleaned-examples
endif
xpkg.build: xpkg.build.$(1)
endef
$(foreach x,$(XPKGS),$(eval $(call xpkg.build.targets,$(x))))

# 1: registry/org 2: repo
define xpkg.release.targets
xpkg.release.publish.$(1).$(2):
	@$(INFO) Pushing package $(1)/$(2):$(VERSION)
	@$(CROSSPLANE_CLI) xpkg push \
		$(foreach p,$(XPKG_LINUX_PLATFORMS),--package-files $(XPKG_OUTPUT_DIR)/$(p)/$(2)-$(VERSION).xpkg ) \
		$(1)/$(2):$(VERSION) || $(FAIL)
	@$(OK) Pushed package $(1)/$(2):$(VERSION)
xpkg.release.publish: xpkg.release.publish.$(1).$(2)

xpkg.release.promote.$(1).$(2):
	@$(INFO) Promoting package from $(1)/$(2):$(VERSION) to $(1)/$(2):$(CHANNEL)
	@docker buildx imagetools create -t $(1)/$(2):$(CHANNEL) $(1)/$(2):$(VERSION)
	@[ "$(CHANNEL)" = "master" ] || docker buildx imagetools create -t $(1)/$(2):$(VERSION)-$(CHANNEL) $(1)/$(2):$(VERSION)
	@$(OK) Promoted package from $(1)/$(2):$(VERSION) to $(1)/$(2):$(CHANNEL)
xpkg.release.promote: xpkg.release.promote.$(1).$(2)
endef
$(foreach r,$(XPKG_REG_ORGS), $(foreach x,$(XPKGS),$(eval $(call xpkg.release.targets,$(r),$(x)))))

# ====================================================================================
# Common Targets

do.build.xpkgs: $(foreach i,$(XPKGS),xpkg.build.$(i))
do.skip.xpkgs:
	@$(OK) Skipping xpkg build for unsupported platform $(IMAGE_PLATFORM)

ifneq ($(filter $(XPKG_PLATFORM),$(XPKG_PLATFORMS_LIST)),)
build.artifacts.platform: do.build.xpkgs
else
build.artifacts.platform: do.skip.xpkgs
endif

# only publish package for main / master and release branches by default
RELEASE_BRANCH_FILTER ?= main master release-%
ifneq ($(filter $(RELEASE_BRANCH_FILTER),$(BRANCH_NAME)),)
publish.artifacts: $(foreach r,$(XPKG_REG_ORGS), $(foreach x,$(XPKGS),xpkg.release.publish.$(r).$(x)))
endif

# NOTE(hasheddan): promotion fails using buildx imagetools create with some
# registries, so a NO_PROMOTE list is supported here. Additionally, channels may
# not be used on some registries that infer vanity tags.
# https://github.com/containerd/containerd/issues/5978
# https://github.com/estesp/manifest-tool/issues/122
# https://github.com/moby/buildkit/issues/2438
promote.artifacts: $(foreach r,$(filter-out $(XPKG_REG_ORGS_NO_PROMOTE),$(XPKG_REG_ORGS)), $(foreach x,$(XPKGS),xpkg.release.promote.$(r).$(x)))
