ENV = /usr/bin/env

.SHELLFLAGS = -c # Run commands in a -c flag
.SILENT: ; # no need for @
.ONESHELL: ; # recipes execute in same shell
.NOTPARALLEL: ; # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

.PHONY: all # All targets are accessible for user
.DEFAULT: help # Running Make will run the help target

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean:
	find . -iname '*.swp' -delete
	rm -rf dist
	rm -f count.out

lint: ## Run linter
	golangci-lint  run

deps: ## Install dependencies
	go install github.com/gogo/protobuf/protoc-gen-gogo
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/goreleaser/goreleaser
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go install github.com/uber/prototool/cmd/prototool
	go install github.com/vektra/mockery/cmd/mockery
	go install golang.org/x/tools/cmd/benchcmp


mock: ## Generate mocks for repositories and services interfaces
	rm -rf mocks
	mockery -name DataStore
	mockery -name UseCase

test: ## Run unit tests
	go test -covermode=count -coverprofile=count.out -v ./...
	go tool cover -func=count.out

build: ## Compile project
	goreleaser --snapshot --skip-publish --rm-dist

migrate: ## Run database migrations 
	migrate -path scripts/migrations/ -database postgres://cataloging:cataloging@127.0.0.1:5432/cataloging?sslmode=disable up 1

populate: ## Populate database with mocked data
	migrate -path scripts/migrations/ -database postgres://cataloging:cataloging@127.0.0.1:5432/cataloging?sslmode=disable up 1
