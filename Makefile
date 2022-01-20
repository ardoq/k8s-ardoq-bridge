export KUBECONFIG ?= ${HOME}/.kube/config

# Active module mode, as we use Go modules to manage dependencies
export GO111MODULE=on
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
GINKGO=$(GOBIN)/ginkgo

SOURCES := $(shell find . -name '*.go')

.PHONY: get-ginkgo
## Installs Ginkgo CLI
get-ginkgo:
	@go install github.com/onsi/ginkgo/ginkgo

.PHONY: k8s-integration-tests
k8s-integration-tests: $(SOURCES) get-ginkgo
	@$(GINKGO) ./tests/integrations/k8s

.PHONY: ardoq-integration-tests
## Runs unit tests with code coverage enabled
ardoq-integration-tests: $(SOURCES) get-ginkgo
	@$(GINKGO) ./tests/integrations/ardoq
