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
	kubectl apply --wait=true -Rf tests/integrations/k8s/manifests
	kubectl wait --for=condition=ready --timeout=180s pod -l app=nginx
	@helm upgrade --install k8s-ardoq-bridge ./chart --wait --set "ardoq.baseUri='${ARDOQ_BASEURI}',ardoq.org='${ARDOQ_ORG}',ardoq.workspaceId='${ARDOQ_WORKSPACE_ID}',ardoq.apiKey=${ARDOQ_APIKEY},ardoq.cluster='${ARDOQ_CLUSTER}'"
	kubectl wait --for=condition=ready --timeout=180s pod -l app.kubernetes.io/name=k8s-ardoq-bridge
	@$(GINKGO) ./tests/integrations/k8s
	kubectl delete --wait=true -Rf tests/integrations/k8s/manifests
	kubectl wait --for=delete --timeout=180s pod -l app=nginx
	echo "Waiting for cleanup"
	sleep 5
	helm delete k8s-ardoq-bridge --wait

.PHONY: ardoq-integration-tests
## Runs unit tests with code coverage enabled
ardoq-integration-tests: $(SOURCES) get-ginkgo
	@$(GINKGO) ./tests/integrations/ardoq
