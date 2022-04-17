package vault

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/v43/github"
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

func (v *VaultServer) ListRemoteHosts(ctx context.Context, in *empty.Empty) (*proto.RemoteHostListResponse, error) {

	gh_user := ctx.Value("github_user").(*github.User)

	identity, err := v.store.GetIdentity(ctx, gh_user.GetLogin())
	if err != nil {
		return nil, err
	}

	permissions, err := v.store.GetPermissions(ctx, identity)
	if err != nil {
		return nil, err
	}

	remoteIds := make([]int64, len(permissions))
	for i, permission := range permissions {
		remoteIds[i] = permission.RemoteID
	}

	remotes, err := v.store.ListRemotes(ctx, remoteIds)
	if err != nil {
		return nil, err
	}

	hosts := make([]*proto.RemoteHost, len(remotes))
	for i, remote := range remotes {
		hosts[i] = &proto.RemoteHost{
			Host: remote.Host,
			Name: remote.Name,
			Port: int64(remote.Port),
		}
	}

	return &proto.RemoteHostListResponse{
		RemoteHosts: hosts,
	}, nil
}
