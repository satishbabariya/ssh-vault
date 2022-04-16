package client

import (
	"context"
	"fmt"

	"github.com/satishbabariya/vault/pkg/server/gen"

	"github.com/cli/oauth"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zalando/go-keyring"
	"google.golang.org/grpc"
)

type VaultClient struct {
	conn   *grpc.ClientConn
	client gen.VaultClient
}

func NewVaultClient(conn *grpc.ClientConn) *VaultClient {
	return &VaultClient{
		conn:   conn,
		client: gen.NewVaultClient(conn),
	}
}

func (v *VaultClient) Login(ctx context.Context) (*string, error) {
	t, err := v.AuthenticateWithGithub(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to store token: %v", err)
	}

	err = keyring.Set(
		"vault", "token", *t,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to store token: %v", err)
	}

	return t, nil
}

func (v *VaultClient) AuthenticateWithGithub(ctx context.Context) (*string, error) {
	config, err := v.client.GetConfig(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	if config.GithubClientId == "" {
		return nil, fmt.Errorf("github client id is empty")
	}

	flow := &oauth.Flow{
		Host:     oauth.GitHubHost(config.GithubHost),
		ClientID: config.GithubClientId,
		Scopes: []string{
			"user:email",
		},
	}

	accessToken, err := flow.DeviceFlow()
	if err != nil {
		return nil, err
	}

	return &accessToken.Token, nil
}
