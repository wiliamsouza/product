# APIs

This folder contains REST and gRPC protocols definitions.

We use these definitions to generate API gateway, gRPC services and client libraries.

## Dependencies

You need to install the protocol buffer compiler `protoc` follow  this
installation [instructions](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation).

Install dependencies running:

```
make deps
```

## Validation

Lint proto definition running:

```
make lint
```

## Generation

Generate stub code running:

```
make generate
```
