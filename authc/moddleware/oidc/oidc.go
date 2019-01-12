package oidc

import (
	"bytes"
	"context"
	"fmt"
	"io"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
)

// AuthAccessToken
type AuthAccessToken func(ctx context.Context, token string) (context.Context, error)

// AuthFuncByJWK make AuthFunction by jwk info
func AuthFuncByJWK(r io.Reader, opts ...ParseOption) (AuthAccessToken, error) {
	set, err := ParseKeys(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	parser := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			kid, ok := token.Header["kid"]
			if !ok {
				return nil, fmt.Errorf("Hasnot kid property")
			}
			key := set.LookupKeyID(kid.(string))
			if len(key) < 1 {
				return nil, fmt.Errorf("Unknown kid: %s", kid)
			}
			return key[0].Materialize()
		}
		return nil, fmt.Errorf("Ignore algorithem [%s]", token.Header["alg"])
	}

	// set options
	opt := newOption()
	for _, v := range opts {
		v(opt)
	}

	return func(ctx context.Context, tokenString string) (context.Context, error) {

		// parse token
		token, err := jwt.Parse(tokenString, parser)
		if token == nil {
			// jwtのformatではないもの
			if err != nil {
				return ctx, errors.WithStack(err)
			}
		}

		// jwtのformatならerrでもtokenの中身を解析する
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if err != nil {
				return ctx, errors.WithStack(err)
			}
			return ctx, errors.Errorf("Fail parse jwt claims")
		}

		// Parse Option
		// Assign claims info to context
		if len(opt.claimsMap) > 0 {
			for k, v := range opt.claimsMap {
				ctx = context.WithValue(ctx, v, claims[k])
			}
		}

		if err != nil {
			return ctx, errors.WithStack(err)
		}
		if !token.Valid {
			return ctx, errors.Errorf("Token is invalid")
		}
		return ctx, nil
	}, nil
}

// ParseOption Interface  of parse option
type ParseOption func(o *option)

// ClaimsMapOption is setup context value
func ClaimsMapOption(name string, contextKey interface{}) ParseOption {
	return func(o *option) {
		o.claimsMap[name] = contextKey
	}
}

func newOption() *option {
	return &option{
		claimsMap: make(map[string]interface{}),
	}
}

type option struct {
	claimsMap map[string]interface{}
}

// parse json web key from io.Reader
func ParseKeys(r io.Reader) (*jwk.Set, error) {
	b := new(bytes.Buffer)
	_, err := io.Copy(b, r)
	if err != nil {
		return nil, err
	}
	return jwk.Parse(b.Bytes())
}
