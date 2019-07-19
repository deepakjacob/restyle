package config

import (
	"context"
	"errors"
	"os"
)

type envKey string

// EnvKey key under which env vars are saved in context
var EnvKey envKey

// Env struct for env vars
type Env struct {
	ProjectID      string
	AppCredentials string
	ClientID       string
	ClientSecret   string
	JWTKey         string
	JWTEncKey      string
	RedirectURL    string
}

// BootstrapCtx context initialization
func BootstrapCtx(parent context.Context) context.Context {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("REDIRECT_URL")
	env := &Env{
		ProjectID:    projectID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
	return context.WithValue(parent, EnvKey, env)
}

// Getenv get env vars
func Getenv(ctx context.Context) (*Env, error) {
	env, ok := ctx.Value(EnvKey).(*Env)
	if !ok {
		return nil, errors.New("oauth:config::unable to get env vars from context")
	}
	return env, nil
}
