package oauth

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config oauth2 config
func Config() (*oauth2.Config, error) {
	// redirectUrl := "http://localhost:8000/auth/callback"
	// TODO: accept the following params from outside
	endPoint := google.Endpoint
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectTo := os.Getenv("REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectTo == "" {
		return nil, errors.New("oauth2 config incorrect")
	}
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	return &oauth2.Config{
		Endpoint:     endPoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectTo,
		Scopes:       scopes,
	}, nil
}
