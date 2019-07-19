package oauth

import (
	"context"

	"github.com/deepakjacob/restyle/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config oauth2 config
func Config(ctx context.Context) (*oauth2.Config, error) {
	env, err := config.Getenv(ctx)
	if err != nil {
		return nil, err
	}
	return &oauth2.Config{
		Endpoint:     google.Endpoint,
		ClientID:     env.ClientID,
		ClientSecret: env.ClientSecret,
		RedirectURL:  env.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}, nil
}
