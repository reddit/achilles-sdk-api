SHELL:=/bin/bash

PWD := $(PWD)
CONTROLLER_GEN := $(PWD)/bin/controller-gen
CONTROLLER_GEN_CMD := $(CONTROLLER_GEN)
GOSIMPORTS := $(PWD)/bin/gosimports
GOSIMPORTS_CMD := $(GOSIMPORTS)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install -modfile=tools/go.mod $(2) ;\
}
endef

.PHONY: manifests
manifests: $(CONTROLLER_GEN) $(GOSIMPORTS)
	$(CONTROLLER_GEN_CMD) object paths="./api/..."
	# avoid diff from controller-gen generated code
	$(GOSIMPORTS_CMD) -local github.com/reddit/achilles-sdk-api -l -w .

.PHONY: generate
generate: manifests
	go generate ./...

.PHONY: lint
lint: $(STATICCHECK) $(GOSIMPORTS)
	cd tools && go mod tidy
	go mod tidy
	go fmt ./...
	go list ./... | xargs go vet
	go list ./... | xargs $(STATICCHECK_CMD)
	$(GOSIMPORTS_CMD) -local github.com/reddit/achilles-sdk-api -l -w .

$(CONTROLLER_GEN):
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen)

$(GOSIMPORTS):
	$(call go-get-tool,$(GOSIMPORTS),github.com/rinchsan/gosimports/cmd/gosimports)

$(STATICCHECK):
	$(call go-get-tool,$(STATICCHECK),honnef.co/go/tools/cmd/staticcheck)
