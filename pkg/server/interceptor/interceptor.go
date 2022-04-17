/*
 * (C) Copyright 2022 Satish Babariya (https://satishbabariya.com/) and others.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Contributors:
 *     satish babariya (satish.babariya@gmail.com)
 *
 */

package interceptor

import (
	"context"

	"github.com/satishbabariya/vault/pkg/server/config"
	"github.com/satishbabariya/vault/pkg/server/gh"
	"github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"
	"google.golang.org/grpc/metadata"
)

type Interceptor struct {
	config        *config.Config
	PublicMethods map[string]bool
}

func NewInterceptor(config *config.Config) *Interceptor {
	return &Interceptor{
		config: config,
		PublicMethods: map[string]bool{
			// "/vault.Vault/GetConfig": true,
			"Vault/GetConfig": true,
		},
	}
}

func (interceptor *Interceptor) IsPublic(method string) bool {
	logrus.Info("method: ", method)
	_, ok := interceptor.PublicMethods[method]
	return ok
}

func (interceptor *Interceptor) NewVaultInterceptor() twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			service, ok := twirp.ServiceName(ctx)
			if !ok {
				return nil, twirp.NewError(twirp.Unauthenticated, "service name not found")
			}

			method, ok := twirp.MethodName(ctx)
			if !ok {
				return nil, twirp.NewError(twirp.BadRoute, "method not found")
			}

			fullMethod := service + "/" + method

			err := interceptor.authorize(ctx, fullMethod, req)
			if err != nil {
				return nil, err
			}
			return next(ctx, req)
		}
	}
}

func (interceptor *Interceptor) authorize(ctx context.Context, method string, payload interface{}) error {

	if interceptor.IsPublic(method) {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return twirp.NewError(twirp.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return twirp.NewError(twirp.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]

	gh_user, err := gh.GetGithubUserFromToken(ctx, accessToken)
	if err != nil {
		return twirp.NewError(twirp.Unauthenticated, "invalid access token")
	}

	ctx = context.WithValue(ctx, "gh_user", gh_user)

	return nil
}
