.PHONY: receiver validator data mailing client
receiver: r_build
	@./receiver
r_build:
	@go build -o receiver ./cmd/receiver/receiver.go

validator: v_build
	@./validator
v_build:
	@go build -o validator ./cmd/validator/validator.go

data: d_build
	@./data
d_build:
	@go build -o data ./cmd/data/data.go

mailing: m_build
	@./mailing
m_build:
	@go build -o mailing ./cmd/mailing/mailing.go

client:
	@go run ./cmd/client/client.go


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
	@rm -rf /tmp/cover.out
integration:
	@go test -short ./tests/integration/ --tags=integration
