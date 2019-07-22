package oauth

import (
	"context"
	"fmt"

	"github.com/deepakjacob/restyle/domain"
	"golang.org/x/oauth2"
)

// Provider interface
type Provider interface {
	RedirectURL(string) string
	GoogleUser(context.Context, string) (*GoogleUser, error)
}

// ProviderImpl implements Provider interface
type ProviderImpl struct {
	Config     OAuth2Config
	HTTPClient HTTPClient
}

//RedirectURL returns redirect url to oauth2 provider signin page
func (p *ProviderImpl) RedirectURL(state string) string {
	return p.Config.AuthCodeURL(state)
}

//GoogleUser returns google user
func (p *ProviderImpl) GoogleUser(ctx context.Context, code string) (*GoogleUser, error) {
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

// OAuth2Config interface for google struct
type OAuth2Config interface {
	AuthCodeURL(string) string
	Exchange(context.Context, string) (*oauth2.Token, error)
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

type userKey string

var userCtxKey userKey

// UserToCtx sets user in call context
func UserToCtx(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

//UserFromCtx returns logged in user
func UserFromCtx(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(userCtxKey).(*domain.User)
	if !ok {
		return nil, fmt.Errorf("auth: user lookup failed for supplied context")
	}
	return user, nil
}
