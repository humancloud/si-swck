#
# Copyright (c) 2019 Stackinsights to present
# All rights reserved
#

mk_path  := $(abspath $(lastword $(MAKEFILE_LIST)))
mk_dir   := $(dir $(mk_path))
root_dir := $(mk_dir)../..
tool_bin := $(root_dir)/bin

include $(root_dir)/hack/build/base.mk

CONTROLLER_GEN_VERSION := v0.7.0
KUSTOMIZE_VERSION := v4.5.6
GOLANGCI_LINT_VERSION := v1.53.3

##@ Code quality and integrity

# The goimports tool does not arrange imports in 3 blocks if there are already more than three blocks.
# To avoid that, before running it, we collapse all imports in one block, then run the formatter.
format: goimports ## Format all Go code
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	done
	$(GOIMPORTS) -w -local github.com/humancloud/si-swck .

## Check that the status is consistent with CI.
check: ## Check that the status
	mkdir -p /tmp/artifacts
	git diff >/tmp/artifacts/check.diff 2>&1
	@go mod tidy &> /dev/null
	@if [ ! -z "`git status -s`" ]; then \
		echo "Following files are not consistent with CI:"; \
		git status -s; \
		cat /tmp/artifacts/check.diff; \
		exit 1; \
	fi

.PHONY: lint
lint: golangci-lint ## Lint codes
	$(GOLANGCILINT) -v run --config $(root_dir)/golangci.yml 

CONTROLLER_GEN = $(tool_bin)/controller-gen
.PHONY: controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION))

KUSTOMIZE = $(tool_bin)/kustomize
.PHONY: kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v4@$(KUSTOMIZE_VERSION))

ENVTEST = $(tool_bin)/setup-envtest
.PHONY: envtest
envtest: ## Download envtest-setup locally if necessary.
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)
	

GOIMPORTS = $(tool_bin)/goimports
.PHONY: goimports
goimports: ## Download goimports locally if necessary.
	$(call go-install-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports@latest)

GOLANGCILINT= $(tool_bin)/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-install-tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))
	

.PHONY: dependency-check
dependency-check: licenseeye ## Check dependencies
	$(LICENSEEYE) -c $(module_dir)/.dep.licenserc.yaml dep check

.PHONY: dependency-resolve
dependency-resolve: licenseeye ## Check dependencies
	$(LICENSEEYE) -c $(module_dir)/.dep.licenserc.yaml dep resolve -o $(module_dir)/dist/licenses -s $(module_dir)/dist/LICENSE.tpl
