.PHONY: run
run: build
	./app
build:
	go build -o app ./cmd/app/app.go

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: buf
buf:
	buf generate api

MIGRATION_DIR:=./migrations
.PHONY: create
create:
	goose -dir=$(MIGRATION_DIR) create $(NAME) sql

.PHONY: migrate
migrate:
	./migrate.sh
