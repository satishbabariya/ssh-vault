package vault

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/satishbabariya/vault/pkg/proto"
	"github.com/satishbabariya/vault/pkg/server/config"
	"github.com/satishbabariya/vault/pkg/server/store"
)

type VaultServer struct {
	// gen.UnimplementedVaultServer
	config *config.Config
	store  *store.Store
}

func NewVaultServer(config *config.Config, store *store.Store) *VaultServer {
	return &VaultServer{
		config: config,
		store:  store,
	}
}

func (v *VaultServer) GetConfig(context.Context, *empty.Empty) (*proto.AuthConfigResponse, error) {
	return &proto.AuthConfigResponse{
		GithubHost:     v.config.GitHubHost,
		GithubClientId: v.config.GithubClientID,
	}, nil
}
