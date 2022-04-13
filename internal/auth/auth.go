package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"ssh-vault/internal/config"
	"ssh-vault/internal/proto"
	"ssh-vault/internal/store"

	"github.com/sirupsen/logrus"
)

type AuthServiceServer struct {
	proto.UnimplementedAuthServiceServer
	config *config.Config
	store  *store.Store
}

func NewAuthServiceServer(config *config.Config, store *store.Store) *AuthServiceServer {
	return &AuthServiceServer{
		config: config,
		store:  store,
	}
}

func (s *AuthServiceServer) GetConfig(context.Context, *proto.Empty) (*proto.AuthConfigResponse, error) {
	fmt.Println(s.config)
	return &proto.AuthConfigResponse{
		GithubHost:     s.config.GitHubHost,
		GithubClientId: s.config.GithubClientID,
	}, nil
}

func (s *AuthServiceServer) Authenticate(ctx context.Context, in *proto.AuthenticateRequest) (*proto.AuthenticateResponse, error) {
	url := "https://api.github.com/user"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	req.Header.Add("Authorization", "token "+in.Token)

	res, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	fmt.Println(string(body))

	logrus.Println(res)

	return nil, nil
}
