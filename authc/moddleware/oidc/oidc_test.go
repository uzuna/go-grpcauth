package oidc

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	CtxKeyEmail = &contextKey{"email"}
	CtxKeySub   = &contextKey{"sub"}
)

type contextKey struct {
	Name string
}

func TestParseKey(t *testing.T) {
	f, err := os.Open("./testdata/keys.json")
	checkError(t, err)
	set, err := ParseKeys(f)
	checkError(t, err)
	key := set.LookupKeyID("1LTMzakihiRla_8z2BEJVXeWMqo")
	assert.Equal(t, 1, len(key))
	_, err = key[0].Materialize()
	checkError(t, err)

	keyn := set.LookupKeyID("eerr")
	assert.Equal(t, 0, len(keyn))
}
func TestVerifyToken(t *testing.T) {
	f, err := os.Open("./testdata/keys.json")
	checkError(t, err)

	var opts []ParseOption
	opts = append(opts,
		ClaimsMapOption("email", CtxKeyEmail),
		ClaimsMapOption("sub", CtxKeySub),
	)
	af, err := AuthFuncByJWK(f, opts...)
	checkError(t, err)

	// access token sample from https://docs.microsoft.com/en-us/azure/active-directory/develop/access-tokens
	b, err := ioutil.ReadFile("./testdata/accesstoken")
	checkError(t, err)

	ctx := context.Background()
	ctx, err = af(ctx, string(b))
	assert.Equal(t, "Unknown kid: i6lGk3FZzxRcUb2C3nEQ7syHJlY", err.Error())
	assert.Equal(t, "AbeLi@microsoft.com", ctx.Value(CtxKeyEmail))
	assert.Equal(t, "l3_roISQU222bULS9yi2k0XpqpOiMz5H3ZACo1GeXA", ctx.Value(CtxKeySub))
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
