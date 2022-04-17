package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cli/oauth"
	"github.com/cli/oauth/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/satishbabariya/vault/pkg/client/interceptor"
	"github.com/satishbabariya/vault/pkg/proto"
	"github.com/twitchtv/twirp"
	"github.com/zalando/go-keyring"
)

type VaultClient struct {
	client proto.Vault
}

func NewClient(ctx context.Context, baseURL string) (*VaultClient, error) {
	token, err := keyring.Get("vault", "token")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	interceptor := interceptor.NewClientInterceptor(token)
	client := proto.NewVaultProtobufClient(baseURL, &http.Client{}, twirp.WithClientInterceptors(interceptor.AuthInterceptor()))

	return &VaultClient{
		client: client,
	}, nil
}

func NewClientUnsafe(ctx context.Context, baseURL string) (*VaultClient, error) {
	client := proto.NewVaultProtobufClient(baseURL, &http.Client{})

	return &VaultClient{
		client: client,
	}, nil
}

func (v *VaultClient) Login(ctx context.Context) error {
	t, err := v.authenticateWithGithub(ctx)
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

func (v *VaultClient) authenticateWithGithub(ctx context.Context) (*api.AccessToken, error) {
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

func (v *VaultClient) ListRemoteHosts(ctx context.Context) ([]*proto.RemoteHost, error) {
	remotes, err := v.client.ListRemoteHosts(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	return remotes.RemoteHosts, nil
}
