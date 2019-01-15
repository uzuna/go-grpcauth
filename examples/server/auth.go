package main

import (
	"context"
	"log"
	"net/http"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/pkg/errors"
	"github.com/uzuna/go-grpcauth/authc/moddleware/oidc"
)

var (
	CtxKeyEmail = &contextKey{"email"}
	CtxKeySub   = &contextKey{"sub"}
)

// DebugAuthentication 目視Debug用
func DebugAuthentication() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		log.Println("token")
		// tokenの目視確認
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		_ = token
		// log.Println(token)

		return ctx, nil
	}
}

// AccessTokenAuthentication jwtベースの認証
func AccessTokenAuthentication(jwkurl string) (grpc_auth.AuthFunc, error) {
	// jwtの検証には指定したpathにあるjwk setを使う
	res, err := http.Get(jwkurl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	// create auth function
	// jwtから後段で利用するParameterをContextにMapさせる
	opts := []oidc.ParseOption{
		oidc.ClaimsMapOption("email", CtxKeyEmail),
		oidc.ClaimsMapOption("sub", CtxKeySub),
	}
	af, err := oidc.AuthFuncByJWK(res.Body, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Return Auth Function
	return func(ctx context.Context) (context.Context, error) {
		// Authenticate Bearer MetadataからTokenを取得
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		// Token検証
		ctx, err = af(ctx, token)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return ctx, nil
	}, nil
}
