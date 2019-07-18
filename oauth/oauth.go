package oauth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

//RedirectURL returns redirect url to oauth2 provider signin page
func (p *Provider) RedirectURL(state string) string {
	return p.Config.AuthCodeURL(state)
}

//GoogleUser returns google user
func (p *Provider) GoogleUser(ctx context.Context, code string) (*GoogleUser, error) {
	token, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	url := fmt.Sprintf(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s",
		token.AccessToken,
	)
	return p.HTTPClient.Get(url)
}

// HTTPClient for http operations
type HTTPClient interface {
	Get(string) (*GoogleUser, error)
}

// OAuth2Config interface for google struct
type OAuth2Config interface {
	AuthCodeURL(string) string
	Exchange(context.Context, string) (*oauth2.Token, error)
}

// Provider implements Provider interface
type Provider struct {
	Config     OAuth2Config
	HTTPClient HTTPClient
}

// OAuth2Configurer abstracts away the oauth2 functionality
type OAuth2Configurer struct {
	Config *oauth2.Config
}

// AuthCodeURL provides signin screen url
func (c *OAuth2Configurer) AuthCodeURL(state string) string {
	return c.Config.AuthCodeURL(state)
}

// Exchange exchanges token for code
func (c *OAuth2Configurer) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.Config.Exchange(ctx, code)
}
