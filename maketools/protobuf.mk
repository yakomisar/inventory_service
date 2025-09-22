# --- Protobuf / Buf toolchain (modern, pinned) ---

# Пины версий (можно обновлять при необходимости)
PROTOBUF_GO_VER := v1.36.6
GRPC_GO_VER := v1.5.1
GRPC_GATEWAY_VER := v2.26.3
BUF_BIN_TAG := v1.53.0

# protoc обычно не обязателен при buf generate
PROTOC_BIN := protoc

# --- protoc plugins ---

# protoc-gen-go
PROTOC_GEN_GO_BIN := $(LOCAL_BIN)/protoc-gen-go
$(PROTOC_GEN_GO_BIN):
	GOBIN="$(LOCAL_BIN)" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOBUF_GO_VER)

# protoc-gen-go-grpc
PROTOC_GEN_GO_GRPC_BIN := $(LOCAL_BIN)/protoc-gen-go-grpc
$(PROTOC_GEN_GO_GRPC_BIN):
	GOBIN="$(LOCAL_BIN)" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(GRPC_GO_VER)

# grpc-gateway v2
PROTOC_GEN_GRPC_GATEWAY_BIN := $(LOCAL_BIN)/protoc-gen-grpc-gateway
$(PROTOC_GEN_GRPC_GATEWAY_BIN):
	GOBIN="$(LOCAL_BIN)" go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@$(GRPC_GATEWAY_VER)

# openapi v2 (из grpc-gateway v2)
PROTOC_GEN_OPENAPI_BIN := $(LOCAL_BIN)/protoc-gen-openapiv2
$(PROTOC_GEN_OPENAPI_BIN):
	GOBIN="$(LOCAL_BIN)" go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@$(GRPC_GATEWAY_VER)

# --- Buf CLI ---

BUF_BIN := $(LOCAL_BIN)/buf
$(BUF_BIN):
	GOBIN="$(LOCAL_BIN)" go install github.com/bufbuild/buf/cmd/buf@$(BUF_BIN_TAG)

# Список обязательных плагинов для генерации
.PHONY: .protodeps
.protodeps: $(PROTOC_GEN_GO_BIN) \
            $(PROTOC_GEN_GO_GRPC_BIN) \
            $(PROTOC_GEN_GRPC_GATEWAY_BIN) \
            $(PROTOC_GEN_OPENAPI_BIN) \
            $(BUF_BIN)
	@echo "✓ proto deps installed"
