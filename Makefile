# --- Settings ---
SHELL := /bin/bash
LOCAL_BIN := $(CURDIR)/bin
PROJECT_NAME := inventory_service

# Optional
export GO111MODULE=on
# Prepend local bin to PATH everywhere
export PATH := $(LOCAL_BIN):$(PATH)

# --- Include subtasks ---
include maketools/*.mk

RUN_ARGS:=
ifneq (,$(wildcard k8s/values_local.yaml))
    RUN_ARGS=--local-config=k8s/values_local.yaml
endif

$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

# Install all required CLI tools (pinned)
.PHONY: tools
tools: $(BUF_BIN) \
       $(PROTOC_GEN_GO_BIN) \
       $(PROTOC_GEN_GO_GRPC_BIN) \
       $(PROTOC_GEN_GRPC_GATEWAY_BIN) \
       $(PROTOC_GEN_OPENAPI_BIN) \
       $(GOLANGCI_BIN) \
       $(GOOSE_BIN)
	@echo "✓ tools ready"

.PHONY: proto-update
proto-update: $(BUF_BIN)
	cd api && $(BUF_BIN) dep update

.PHONY: buf-lint
buf-lint: $(BUF_BIN) proto-update
	cd api && $(BUF_BIN) lint

.PHONY: go-lint
go-lint: $(GOLANGCI_BIN)
	$(GOLANGCI_BIN) run ./...

.PHONY: lint
lint: go-lint buf-lint

.PHONY: generate
generate: .protodeps buf-lint
	$(BUF_BIN) generate -v

.PHONY: run
run:
	go run cmd/main.go $(RUN_ARGS)

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)
	find pkg/pb -type f -name '*.go' -delete || true
	rm -f api/api.swagger.json || true
