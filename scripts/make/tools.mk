# Tools installation and management

# Tool versions
BUF_VERSION := 1.28.1
GOLANGCI_LINT_VERSION := 1.55.2
PROTOC_GEN_GO_VERSION := 1.31.0
PROTOC_GEN_GO_GRPC_VERSION := 1.3.0
PROTOC_GEN_GRPC_GATEWAY_VERSION := 2.19.0
PROTOC_GEN_OPENAPIV2_VERSION := 2.19.0
PROTOC_GEN_DOC_VERSION := 1.5.1
PROTOC_GEN_VALIDATE_VERSION := 1.0.4

# Local binary directory
LOCAL_BIN := $(CURDIR)/bin
export PATH := $(LOCAL_BIN):$(PATH)

# Binary paths
BUF := $(LOCAL_BIN)/buf
GOLANGCI_LINT := $(LOCAL_BIN)/golangci-lint
PROTOC_GEN_GO := $(LOCAL_BIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(LOCAL_BIN)/protoc-gen-go-grpc
PROTOC_GEN_GRPC_GATEWAY := $(LOCAL_BIN)/protoc-gen-grpc-gateway
PROTOC_GEN_OPENAPIV2 := $(LOCAL_BIN)/protoc-gen-openapiv2
PROTOC_GEN_DOC := $(LOCAL_BIN)/protoc-gen-doc
PROTOC_GEN_VALIDATE := $(LOCAL_BIN)/protoc-gen-validate

$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

.PHONY: install-tools
install-tools: $(LOCAL_BIN) install-buf install-protoc-plugins install-golangci-lint ## Install all required tools

.PHONY: install-buf
install-buf: $(LOCAL_BIN) ## Install buf CLI
	@echo "Installing buf $(BUF_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v$(BUF_VERSION)
	@echo "buf installed successfully at $(BUF)"
	@$(BUF) --version

.PHONY: install-protoc-plugins
install-protoc-plugins: $(LOCAL_BIN) ## Install protoc Go plugins
	@echo "Installing protoc-gen-go $(PROTOC_GEN_GO_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v$(PROTOC_GEN_GO_VERSION)
	@echo "Installing protoc-gen-go-grpc $(PROTOC_GEN_GO_GRPC_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v$(PROTOC_GEN_GO_GRPC_VERSION)
	@echo "Installing protoc-gen-grpc-gateway $(PROTOC_GEN_GRPC_GATEWAY_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v$(PROTOC_GEN_GRPC_GATEWAY_VERSION)
	@echo "Installing protoc-gen-openapiv2 $(PROTOC_GEN_OPENAPIV2_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v$(PROTOC_GEN_OPENAPIV2_VERSION)
	@echo "Installing protoc-gen-doc $(PROTOC_GEN_DOC_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v$(PROTOC_GEN_DOC_VERSION)
	@echo "Installing protoc-gen-validate $(PROTOC_GEN_VALIDATE_VERSION)..."
	@GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v$(PROTOC_GEN_VALIDATE_VERSION)
	@echo "All protoc plugins installed successfully"

.PHONY: install-golangci-lint
install-golangci-lint: $(LOCAL_BIN) ## Install golangci-lint
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) v$(GOLANGCI_LINT_VERSION)
	@echo "golangci-lint installed successfully at $(GOLANGCI_LINT)"
	@$(GOLANGCI_LINT) --version

.PHONY: check-tools
check-tools: ## Check if all required tools are installed
	@echo "Checking installed tools..."
	@command -v buf >/dev/null 2>&1 || { echo "buf is not installed. Run 'make install-buf'"; exit 1; }
	@command -v protoc-gen-go >/dev/null 2>&1 || { echo "protoc-gen-go is not installed. Run 'make install-protoc-plugins'"; exit 1; }
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "protoc-gen-go-grpc is not installed. Run 'make install-protoc-plugins'"; exit 1; }
	@command -v protoc-gen-grpc-gateway >/dev/null 2>&1 || { echo "protoc-gen-grpc-gateway is not installed. Run 'make install-protoc-plugins'"; exit 1; }
	@command -v protoc-gen-openapiv2 >/dev/null 2>&1 || { echo "protoc-gen-openapiv2 is not installed. Run 'make install-protoc-plugins'"; exit 1; }
	@command -v protoc-gen-doc >/dev/null 2>&1 || { echo "protoc-gen-doc is not installed. Run 'make install-protoc-plugins'"; exit 1; }
	@command -v protoc-gen-validate >/dev/null 2>&1 || { echo "protoc-gen-validate is not installed. Run 'make install-protoc-plugins'"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed. Run 'make install-golangci-lint'"; exit 1; }
	@echo "All tools are installed ✓"

