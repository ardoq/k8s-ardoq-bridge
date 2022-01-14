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

.PHONY: gingko-tests
integration-tests: get-ginkgo
	@$(GINKGO) \
	-coverprofile=coverage.txt \
	./tests/integrations/ardoq

.PHONY: unit-tests
## Runs unit tests with code coverage enabled
unit-tests: $(SOURCES) get-ginkgo
	go test -v -short -race -timeout 15m ./...
