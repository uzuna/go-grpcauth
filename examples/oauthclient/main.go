package oauthclient

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

// Configure for OpenIDConnect id_token flow
type Config struct {
	ClientID        string
	ClientSecret    string
	RedirectAddress string
	AuthURL         string
	TokenURL        string
	Token           string
}

func LoadConfig() Config {
	c := Config{
		ClientID:        os.Getenv("CLIENT_ID"),
		ClientSecret:    os.Getenv("CLIENT_SECRET"),
		RedirectAddress: os.Getenv("REDIRECT_ADDR"),
		AuthURL:         os.Getenv("OAUTH_AUTH_URL"),
		TokenURL:        os.Getenv("OAUTH_TOKEN_URL"),
		Token:           os.Getenv("TOKEN"),
	}
	return c
}

// Get Access Token
func GetAccessToken(config Config) (string, error) {
	// precheck
	if len(config.AuthURL) < 1 {
		return "", errors.Errorf("Must set AuthURL [env:OAUTH_AUTH_URL]")
	}
	if len(config.TokenURL) < 1 {
		return "", errors.Errorf("Must set TokenURL [env:OAUTH_TOKEN_URL]")
	}

	l, err := net.Listen("tcp", config.RedirectAddress)
	if err != nil {
		return "", err
	}
	defer l.Close()

	oauthConfig := &oauth2.Config{
		Scopes: []string{
			"openid",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("http://%s", config.RedirectAddress),
	}

	stateBytes := make([]byte, 16)
	_, err = rand.Read(stateBytes)
	if err != nil {
		return "", err
	}

	state := fmt.Sprintf("%x", stateBytes)
	err = open.Start(oauthConfig.AuthCodeURL(state,
		oauth2.SetAuthURLParam("response_type", "id_token"),
		oauth2.SetAuthURLParam("nonce", "0011223"),
	))
	if err != nil {
		return "", err
	}

	quit := make(chan string)
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			w.Write([]byte(`<script>location.href = "/close?" + location.hash.substring(1);</script>`))
		} else {
			w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
			w.(http.Flusher).Flush()
			log.Println(req.URL.Query())
			quit <- req.URL.Query().Get("access_token")
		}
	}))
	return <-quit, nil
}
