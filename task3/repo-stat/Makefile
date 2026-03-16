SHELL := /bin/sh

GOBIN := $(shell go env GOPATH)/bin

PROTOLINT := $(GOBIN)/protolint
GOIMPORTS := $(GOBIN)/goimports
GOLANGCI_LINT := $(GOBIN)/golangci-lint
PROTOC_GEN_GO := $(GOBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOBIN)/protoc-gen-go-grpc

PROTOC_GEN_GO_VERSION := v1.32.0
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1
GOLANGCI_LINT_VERSION := v2.11.2

.PHONY: lint fix protobuf protolint gofmtcheck golint tools test

lint: protolint gofmtcheck golint

fix:
	$(PROTOLINT) lint -fix .
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	$(GOLANGCI_LINT) run --fix --timeout=2m -E gocritic -v ./...

protobuf:
	PATH="$(GOBIN):$$PATH" protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/subscriber/subscriber.proto

protolint:
	$(PROTOLINT) .

gofmtcheck:
	find . -type f -name '*.go' ! -name '*.pb.go' -print0 | xargs -0 $(GOIMPORTS) -w
	git diff --exit-code

golint:
	$(GOLANGCI_LINT) run --timeout=2m -E gocritic -v ./...

tools:
	go install github.com/yoheimuta/protolint/cmd/protolint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

test:
	go test -race -coverprofile cover.out \
		$(shell go list ./... | grep -E -v 'repo-stat/(proto)')
	go tool cover -html=cover.out -o cover.html