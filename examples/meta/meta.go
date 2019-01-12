package meta

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationKey = "authorization"
)

// ontextから値を取り出す
func Authorization(ctx context.Context) (string, error) {
	return fromMeta(ctx, AuthorizationKey)
}

func fromMeta(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("not found metadata")
	}
	vs := md[key]
	if len(vs) == 0 {
		return "", fmt.Errorf("not found %s in metadata", key)
	}
	return vs[0], nil
}
