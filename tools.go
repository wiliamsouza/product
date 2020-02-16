// +build tools

package product

import (
	_ "github.com/gogo/protobuf/protoc-gen-gogo"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/uber/prototool/cmd/prototool"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "golang.org/x/tools/cmd/benchcmp"
)
