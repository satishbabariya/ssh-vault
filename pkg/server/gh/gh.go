package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func GetGithubUserFromToken(ctx context.Context, token string) (*github.User, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	if !resp.TokenExpiration.IsZero() {
		return nil, fmt.Errorf("token expired at %v", resp.TokenExpiration)
	}

	return user, nil
}
