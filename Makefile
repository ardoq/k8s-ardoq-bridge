.ONESHELL:

VERSION=`cat VERSION`
export KUBECONFIG ?= ${HOME}/.kube/config

# Active module mode, as we use Go modules to manage dependencies
export GO111MODULE=on
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
GINKGO=$(GOBIN)/ginkgo

SOURCES := $(shell find . -name '*.go')

#include .env

install-git-hooks:
	pre-commit install --hook-type pre-commit --hook-type pre-push

.PHONY: get-ginkgo
## Installs Ginkgo CLI
get-ginkgo:
	@go get github.com/onsi/ginkgo/v2/ginkgo/internal@v2
	@go install github.com/onsi/ginkgo/v2/ginkgo

bootstrap: get-ginkgo
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.5.0
	golangci-lint version

.PHONY: k8s-integration-tests
k8s-integration-tests: $(SOURCES) get-ginkgo
	@$(GINKGO) -v  --show-node-events ./tests/integrations/k8s

.PHONY: ardoq-integration-tests
## Runs ardoq integration tests
ardoq-integration-tests: $(SOURCES) get-ginkgo
	@$(GINKGO) -v  --show-node-events ./tests/integrations/ardoq

.PHONY: unit-tests
## Runs ardoq integration tests
unit-tests: $(SOURCES) get-ginkgo
	@$(GINKGO) -v  --show-node-events -p  ./tests/unit

cleanup:
	go mod tidy

.PHONY: all-tests
## Runs all tests
all-tests: $(SOURCES) unit-tests ardoq-integration-tests k8s-integration-tests cleanup

kind-up:
	kind create cluster --name=kind --config ./kind/config.yaml

kind-down:
	kind delete cluster --name=kind

docker-build:
	docker build -t ardoq/k8s-ardoq-bridge:devel .

kind-load: docker-build
	kind load docker-image ardoq/k8s-ardoq-bridge:devel --name=kind

run:
	go run main.go --kubeconfig=$(HOME)/.kube/config
