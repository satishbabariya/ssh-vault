package main

import (
	"context"
	"fmt"
	"net"
	"ssh-vault/internal/config"
	"ssh-vault/internal/store"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	store, err := store.NewStore(config.DatabaseURL)
	if err != nil {
		logrus.Fatal(err)
	}
	defer store.Close()

	err = store.Init(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	// Create new gRPC server
	server := grpc.NewServer()

	// Register the services with the gRPC server.

	// listen on the port
	port := fmt.Sprintf("0.0.0.0:%s", config.Port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logrus.Fatal(err)
	}

	// start the server
	logrus.Info("Starting gRPC server on port: ", port)
	logrus.Fatal(server.Serve(listener))
}
