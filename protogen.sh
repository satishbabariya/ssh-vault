#!/bin/sh

# check if protoc is installed
if ! which protoc > /dev/null; then
  echo "Protobuf compiler not found in PATH."

  # install protoc
  echo "Installing Protobuf compiler..."

  # if on Mac OS X
  if which brew > /dev/null; then
    brew install protobuf
  fi

  # if on Linux and apt-get is available
  if which apt-get > /dev/null; then
      sudo apt-get install protobuf-compiler -y
  elif which yum > /dev/null; then
    sudo yum install protobuf-compiler -y
  else 
    echo "No package manager found. Please install protobuf-compiler manually."
  fi
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

# remove pkg/server/gen directory if it exists
if [ -d pkg/server/gen ]; then
    rm -rf pkg/server/gen
fi

mkdir -p pkg/server/gen

protoc --proto_path=. --go_out=pkg/server/gen --go_opt=paths=source_relative \
    --go-grpc_out=pkg/server/gen --go-grpc_opt=paths=source_relative \
    *.proto