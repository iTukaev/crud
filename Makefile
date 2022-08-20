.PHONY: run
run: build
	./app
build:
	go build -o app ./cmd/app/app.go

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .deps buf
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

buf:
	buf generate api


MIGRATION_DIR:=./migrations
.PHONY: create migrate
create:
	goose -dir=$(MIGRATION_DIR) create $(NAME) sql

migrate:
	./migrate.sh


PKG_LIST := $(shell go list ./... \
	| grep -v -E '/internal/api/(validator)' \
	| grep -v -E '/(postgresql)' \
	| grep -v -E '/(cmd)' \
	| grep -v -E '/(pkg)')
.PHONY: test cover integration
test:
	@go test -short ${PKG_LIST}

cover:
	@go test -short ${PKG_LIST} -coverprofile=/tmp/cover.out -covermode=count -coverpkg=./...
	@go tool cover -html=/tmp/cover.out

integration:
	@go test -short ${PKG_LIST} ./tests/integration/ --tags=integration
