GOLANGCI_BIN=$(LOCAL_BIN)/golangci-lint
GOLANGCI_TAG=latest

$(GOLANGCI_BIN):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN) $(GOLANGCI_TAG)

GCI_BIN=$(LOCAL_BIN)/gci
$(GCI_BIN):
	GOBIN=$(LOCAL_BIN) go install github.com/daixiang0/gci@latest
