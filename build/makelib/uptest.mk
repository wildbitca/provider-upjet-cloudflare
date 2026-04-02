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
# End to End Testing

ifndef UPTEST_LOCAL_DEPLOY_TARGET
  $(error UPTEST_LOCAL_DEPLOY_TARGET is not set. It is required for the provider or configuration \
		to be deployed before running the tests. \
		Example `local.xpkg.deploy.configuration.$(PROJECT_NAME)` for configurations or \
		`local.xpkg.deploy.provider.$(PROJECT_NAME)` for providers)
endif

UPTEST_ARGS ?=

UPTEST_SKIP_UPDATE ?= false
ifeq ($(UPTEST_SKIP_UPDATE),true)
    UPTEST_ARGS += --skip-update
endif

UPTEST_SKIP_IMPORT ?= false
ifeq ($(UPTEST_SKIP_IMPORT),true)
    UPTEST_ARGS += --skip-import
endif

UPTEST_SKIP_DELETE ?= false
ifeq ($(UPTEST_SKIP_DELETE),true)
    UPTEST_ARGS += --skip-delete
endif

UPTEST_DEFAULT_TIMEOUT ?=
ifdef UPTEST_DEFAULT_TIMEOUT
	UPTEST_ARGS += --default-timeout=$(UPTEST_DEFAULT_TIMEOUT)
endif

UPTEST_RENDER_ONLY ?= false
ifeq ($(UPTEST_RENDER_ONLY),true)
    UPTEST_ARGS += --render-only
endif

UPTEST_COMMAND = SKIP_DEPLOY_ARGO=$(SKIP_DEPLOY_ARGO) \
	KUBECTL=$(KUBECTL) \
	CHAINSAW=$(CHAINSAW) \
	CROSSPLANE_CLI=$(CROSSPLANE_CLI) \
	CROSSPLANE_NAMESPACE=$(CROSSPLANE_NAMESPACE) \
	YQ=$(YQ) \
	$(UPTEST) e2e $(UPTEST_INPUT_MANIFESTS) \
	--data-source="${UPTEST_DATASOURCE_PATH}" \
	--setup-script=$(UPTEST_SETUP_SCRIPT) \
	$(UPTEST_ARGS)

# This target requires the following environment variables to be set:
# - To ensure the proper functioning of the end-to-end test resource pre-deletion hook, it is crucial to arrange your resources appropriately.
#   You can check the basic implementation here: https://github.com/crossplane/uptest/blob/main/internal/templates/03-delete.yaml.tmpl
# - UPTEST_DATASOURCE_PATH (optional), see https://github.com/crossplane/uptest?tab=readme-ov-file#injecting-dynamic-values-and-datasource
UPTEST_SETUP_SCRIPT ?= test/setup.sh
uptest: $(UPTEST) $(KUBECTL) $(CHAINSAW) $(CROSSPLANE_CLI) $(YQ)
	@$(INFO) running automated tests
	$(UPTEST_COMMAND) || $(FAIL)
	@$(OK) running automated tests

# Run uptest together with all dependencies. Use `make e2e UPTEST_SKIP_DELET=true` to skip deletion of resources.
e2e: build controlplane.down controlplane.up $(UPTEST_LOCAL_DEPLOY_TARGET) uptest #

# Renders crossplane compositions in the current project
#
# Composition and Function must be defined, Environment and Observeed Resources are optionally available.
# The command discovers sees by parsing the `render.crossplane.io/*`-annotations of file found in the
# examples-directorry, usually: `/examples`. This folder can be overwritten: UPTEST_EXAMPLES_FOLDER='other-examples/'
#
# Possible values are:
#  - composition-path (Composition)
#  - function-path (Function)
#  - environment-path (Environment Configs)
#  - observe-path (Observed Resources)
#
#  Example: `render.crossplane.io/composition-path: "apis/kcl/composition.yaml"`
#
# Alternatively the user is able to specify `UPTEST_RENDER_FILES` with a comma or space-separated list
# of files to render.
#
#  Example: `make render UPTEST_RENDER_FILES="path/to/example.yaml,another-example.yaml`
UPTEST_RENDER_FILES ?=
UPTEST_EXAMPLES_FOLDER ?= ./examples
render: $(CROSSPLANE_CLI) ${YQ}
	@indir="$(UPTEST_EXAMPLES_FOLDER)"; \
	rm -rf "$(CACHE_DIR)/render"; \
	mkdir -p "$(CACHE_DIR)/render" || true; \
	if [ -z $${UPTEST_RENDER_FILES} ]; then \
		FILES=$$(find $$indir -type f -name '*.yaml' ); \
	else \
		FILES=$$(echo $${UPTEST_RENDER_FILES} | sed 's/,/ /g'); \
	fi; \
	$(INFO) Files to render: $${FILES}; \
	for file in $${FILES}; do \
		doc_count=$$(${YQ} eval 'documentIndex' $$file | wc -l); \
		for i in $$(seq 0 $$(($$doc_count - 1))); do \
			COMPOSITION=$$(${YQ} eval "select(documentIndex == $$i) | .metadata.annotations.\"render.crossplane.io/composition-path\"" $$file); \
			FUNCTION=$$(${YQ} eval "select(documentIndex == $$i) | .metadata.annotations.\"render.crossplane.io/function-path\"" $$file); \
			ENVIRONMENT=$$(${YQ} eval "select(documentIndex == $$i) | .metadata.annotations.\"render.crossplane.io/environment-path\"" $$file); \
			OBSERVE=$$(${YQ} eval "select(documentIndex == $$i) | .metadata.annotations.\"render.crossplane.io/observe-path\"" $$file); \
			IS_CROSSPLANE_OBJ=$$(${YQ} eval "select(documentIndex == $$i).apiVersion" $$file | grep -E "^pkg\.crossplane\.io|^apiextensions\.crossplane\.io"); \
			if [ ! -z "$$IS_CROSSPLANE_OBJ" ]; then \
				continue; \
			fi; \
			if [[ "$$ENVIRONMENT" == "null" ]]; then \
				ENVIRONMENT=""; \
			fi; \
			if [[ "$$OBSERVE" == "null" ]]; then \
				OBSERVE=""; \
			fi; \
			if [[ "$$COMPOSITION" == "null" || "$$FUNCTION" == "null" ]]; then \
				$(WARN) file $$file has document with no annotations for rendering, skipping; \
				continue; \
			fi; \
			OUT_FILE=$$(echo $$file | md5sum | head -c5); \
			ENVIRONMENT=$${ENVIRONMENT=="null" ? "" : $$ENVIRONMENT}; \
			OBSERVE=$${OBSERVE=="null" ? "" : $$OBSERVE}; \
			$(INFO) rendering $$file; \
			$(CROSSPLANE_CLI) render $$file $$COMPOSITION $$FUNCTION $${ENVIRONMENT:+-e $$ENVIRONMENT} $${OBSERVE:+-o $$OBSERVE} -x >> "$(CACHE_DIR)/render/$${OUT_FILE}.yaml"; \
			if [ $$? != 0 ]; then \
				$(ERR) fail rendering $$file; \
				exit 1; \
			fi; \
			$(OK) rendered $$file; \
		done; \
	done

# Prints the raw rendered yaml to stdout
#
# User can restrict the files to be shown with a comma or space-separated list
# of files to render.
#
#  Example: `make render.show UPTEST_RENDER_FILES="path/to/example.yaml,another-example.yaml`
render.show:
	@$(MAKE) render UPTEST_RENDER_FILES=$${UPTEST_RENDER_FILES} >/dev/null
	@if [ -z $${UPTEST_RENDER_FILES} ]; then \
		find "$(CACHE_DIR)/render" -type f -name "*.yaml" -exec cat {} \; ; \
	else \
		EXAMPLE_FILES=$$(echo $${UPTEST_RENDER_FILES} | sed 's/,/ /g'); \
		for file in $${EXAMPLE_FILES}; do \
			OUT_FILE=$$(echo $$file | md5sum | head -c5); \
			OUT_FILE="$(CACHE_DIR)/render/$${OUT_FILE}.yaml"; \
			if [ -f $${OUT_FILE} ]; then \
				cat "$${OUT_FILE}"; \
			fi; \
		done; \
	fi; \

# Validates the rendered output
#
# User can supply custom extensions-folder or file with UPTEST_VALIDATE_EXTENSIONS=path/to/extension
# Additionally there is the ability to restrict which files should be rendered. Samne rules apply as
# for `render` and `render.show` command
#
# Note: Extension in this context means:
#   - an XRD
#   - any provider of function-definition (providers.yaml, functions.yaml)
#   - a `crossplane.yaml`
#
#  Examle: `make render.validate UPTEST_VALIDATE_EXTENSIONS='some/folder`
UPTEST_VALIDATE_EXTENSIONS ?= crossplane.yaml
render.validate:
	@if [ -z $${UPTEST_RENDER_FILES} ]; then \
		EXAMPLE_FILES=$$(find $(UPTEST_EXAMPLES_FOLDER) -type f -name '*.yaml' ); \
	else \
		EXAMPLE_FILES=$$(echo $${UPTEST_RENDER_FILES} | sed 's/,/ /g'); \
	fi; \
	for file in $${EXAMPLE_FILES}; do \
		$(INFO) validating $$file; \
		RENDERED=$$($(MAKE) render.show UPTEST_RENDER_FILES=$${file}); \
		if [ -z "$${RENDERED}" ]; then \
			$(WARN) render produced empty output for: $$file; \
			continue; \
		fi; \
		echo "$${RENDERED}" | $(CROSSPLANE_CLI) beta validate $(UPTEST_VALIDATE_EXTENSIONS) - ; \
		if [ $$? -ne 0 ]; then \
			$(ERR) fail validating $$file; \
			exit 1; \
		fi; \
		$(OK) validated $$file; \
	done;


YAMLLINT_FOLDER ?= ./apis
yamllint: ## Static yamllint check
	@$(INFO) running yamllint
	@yamllint $(YAMLLINT_FOLDER) || $(FAIL)
	@$(OK) running yamllint

.PHONY: uptest e2e render yamllint


