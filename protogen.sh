#!/bin/sh

# check if protoc is installed
if ! which protoc > /dev/null; then
  echo "Protobuf compiler not found in PATH."
fi

# Update PATH so that the protoc compiler can find the plugins
export PATH="$PATH:$(go env GOPATH)/bin"

# check if protoc-gen-go is installed
if ! which protoc-gen-go > /dev/null; then
  echo "Protobuf Go compiler not found in PATH."
  echo "Installing it..."
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
fi

# remove internal/proto directory if it exists
if [ -d internal/proto ]; then
    rm -rf internal/proto
fi

mkdir -p internal/proto

protoc --proto_path=. --go_out=internal/proto --go_opt=paths=source_relative \
    --go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative \
    vault.proto