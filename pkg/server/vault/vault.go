package vault

import (
	"context"
	"ssh-vault/pkg/server/config"
	"ssh-vault/pkg/server/gen"
	"ssh-vault/pkg/server/store"

	"github.com/golang/protobuf/ptypes/empty"
)

type VaultServer struct {
	gen.UnimplementedVaultServer
	config *config.Config
	store  *store.Store
}

func NewVaultServer(config *config.Config, store *store.Store) *VaultServer {
	return &VaultServer{
		config: config,
		store:  store,
	}
}

func (v *VaultServer) GetConfig(context.Context, *empty.Empty) (*gen.AuthConfigResponse, error) {
	return &gen.AuthConfigResponse{
		GithubHost:     v.config.GitHubHost,
		GithubClientId: v.config.GithubClientID,
	}, nil
}
