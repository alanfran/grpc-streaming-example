GOARCH:=amd64
GOOS:=darwin

PRODUCT:=grpc-streaming-example
VERSION:=0.1.0
BUILD:=$(shell date +%s)
RELEASE:=0

REPO:=github.com/alanfran

export GO111MODULE=on

LDFLAGS:=
ifeq ($(GOOS),windows)
BUILDFLAGS:=
else
BUILDFLAGS:=-race
endif

################################################################################
# BUILD
################################################################################
.PHONY: build-deps
build-deps:
	go mod tidy

.PHONY: build
build: build-deps ## build the application
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENVS) go build -v -ldflags \
		"-X main.Version=$(VERSION) -X main.Build=$(BUILD) $(LDFLAGS)" \
		$(BUILDFLAGS) -o $(PRODUCT) "$(REPO)/$(PRODUCT)/cmd/$(PRODUCT)"

.PHONY: vendor
vendor: ## create vendor directory of dependencies (used for docker images and for debugging)
	go mod vendor

.PHONY: clean
clean: ## clean up build products and rpms
	go clean ./...
	rm -f $(PRODUCT)
	rm -rf vendor

################################################################################
# LINT
################################################################################

.PHONY: lint-deps
lint-deps: ## get linter for testing
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.18

.PHONY: lint
lint: lint-deps ## get linter for testing
	golangci-lint run

################################################################################
# TEST
################################################################################
.PHONY: test
test: ## run unit tests and code coverage
	go vet ./...
	go test ./... -v
	go test ./... -cover

################################################################################
# API
################################################################################

.PHONY: proto-deps
proto-deps:
	go get github.com/golang/protobuf/protoc-gen-go
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

.PHONY: proto
proto: proto-deps ## build language specific bindings
	protoc \
		-I ./api \
		-I ./third_party/proto \
		--go_out=plugins=grpc,paths=source_relative:./pkg \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:./pkg \
		--swagger_out=logtostderr=true:docs/generated \
		api/example/v1/*.proto

################################################################################
# HELP
################################################################################

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help