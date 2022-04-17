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
if [ -d pkg/proto ]; then
    rm -rf pkg/proto
fi

mkdir -p pkg/proto

if ! which protoc-gen-twirp > /dev/null; then
  go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

protoc --twirp_out=. --go_out=. *.proto