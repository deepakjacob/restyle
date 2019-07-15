package oauth

import (
	"golang.org/x/oauth2"
)

// Provider implements Provider interface
type Provider struct {
	Config *oauth2.Config
}

//RedirectURL returns redirect url to oauth2 provider signin page
func (p *Provider) RedirectURL(state string) string {
	return p.Config.AuthCodeURL(state)
}
