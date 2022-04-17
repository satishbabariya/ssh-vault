package client

import (
	"context"
	"fmt"

	"github.com/satishbabariya/vault/pkg/client/interceptor"
	"github.com/satishbabariya/vault/pkg/server/gen"

	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zalando/go-keyring"
	"google.golang.org/grpc"
)

type VaultClient struct {
	conn   *grpc.ClientConn
	client gen.VaultClient
}

func NewClient(ctx context.Context) (*VaultClient, error) {
	token, err := keyring.Get("vault", "token")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	interceptor := interceptor.NewClientInterceptor(token)
	conn, err := grpc.DialContext(ctx, "localhost:1203", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor))
	if err != nil {
		return nil, err
	}

	return NewVaultClient(conn), nil
}

func NewClientUnsafe(ctx context.Context) (*VaultClient, error) {
	conn, err := grpc.DialContext(ctx, "localhost:1203", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return NewVaultClient(conn), nil
}

func NewVaultClient(conn *grpc.ClientConn) *VaultClient {
	return &VaultClient{
		conn:   conn,
		client: gen.NewVaultClient(conn),
	}
}

func (v *VaultClient) Close() error {
	return v.conn.Close()
}

func (v *VaultClient) Login(ctx context.Context) error {
	t, err := v.AuthenticateWithGithub(ctx)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	err = keyring.Set(
		"vault", "token", t.Token,
	)

	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}

func (v *VaultClient) AuthenticateWithGithub(ctx context.Context) (*api.AccessToken, error) {
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

	return accessToken, nil
}
