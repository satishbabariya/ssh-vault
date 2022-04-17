package main

import (
	"context"
	"net/http"

	"github.com/satishbabariya/vault/pkg/proto"
	"github.com/satishbabariya/vault/pkg/server/config"
	"github.com/satishbabariya/vault/pkg/server/interceptor"
	"github.com/satishbabariya/vault/pkg/server/store"
	"github.com/satishbabariya/vault/pkg/server/vault"
	"github.com/twitchtv/twirp"

	"github.com/sirupsen/logrus"
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

	// Initialize vault service
	vault := vault.NewVaultServer(config, store)

	interceptor := interceptor.NewInterceptor(config)

	// Create a new twirp server
	handler := proto.NewVaultServer(vault, twirp.WithServerInterceptors(interceptor.NewVaultInterceptor()))

	// Create a new http server
	logrus.Info("Starting vault server on port: ", config.Port)
	err = http.ListenAndServe("0.0.0.0:"+config.Port, handler)
	if err != nil {
		logrus.Fatal(err)
	}
}
