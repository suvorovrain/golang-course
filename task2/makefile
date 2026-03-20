LOCAL_BIN := $(CURDIR)/bin
BUF_BUILD := $(LOCAL_BIN)/buf

.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.buf-generate:
	set "PATH=$(LOCAL_BIN);%PATH%" && $(BUF_BUILD) generate

generate: .buf-generate

build:
	docker-compose up -d --build 
run:
	docker-compose up -d
stop:
	docker-compose down 
	
.tidy:
	go mod tidy

.PHONY: .bin-deps