package interceptor

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/grpc/metadata"
)

// ClientInterceptor is a gRPC interceptor that adds the access token to the request
type ClientInterceptor struct {
	accessToken string
}

func NewClientInterceptor(accessToken string) *ClientInterceptor {
	return &ClientInterceptor{accessToken: accessToken}
}

func (interceptor *ClientInterceptor) AuthInterceptor() twirp.Interceptor {
	return func(next twirp.Method) twirp.Method {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			}

			md.Set("authorization", interceptor.accessToken)

			ctx = metadata.NewOutgoingContext(ctx, md)
			return next(ctx, req)
		}
	}
}
