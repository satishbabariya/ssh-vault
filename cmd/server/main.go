package main

import (
	"context"
	"fmt"
	"net"
	"ssh-vault/pkg/server/config"
	"ssh-vault/pkg/server/gen"
	"ssh-vault/pkg/server/interceptor"
	"ssh-vault/pkg/server/store"
	"ssh-vault/pkg/server/vault"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// Initialize config
	config, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// Create a new store
	store, err := store.NewStore(config.DatabaseURL)
	if err != nil {
		logrus.Fatal(err)
	}
	defer store.Close()

	// Initialize postgres store with create table statements
	err = store.Init(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	// interceptors will be used to validate ghitub token
	interceptor := interceptor.NewInterceptor(config)

	// Create new gRPC server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryInterceptor),
		grpc.StreamInterceptor(interceptor.StreamInterceptor),
	)

	// Register the server with the gRPC server
	vault := vault.NewVaultServer(config, store)

	// Register the services with the gRPC server.
	gen.RegisterVaultServer(server, vault)

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
