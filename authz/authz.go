package authz

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServiceAuthorize ServiceがAuthorizeMethodを実装していることを期待する
type ServiceAuthorize interface {
	Authorize(context.Context, string) error
}

// UnaryServerInterceptor is call ServiceAuthorize method
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		if srv, ok := info.Server.(ServiceAuthorize); ok {
			err = srv.Authorize(ctx, info.FullMethod)
		} else {
			return nil, fmt.Errorf("each service should implement an authorization")
		}
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return handler(ctx, req)
	}
}
