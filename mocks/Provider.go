package mocks

import "golang.org/x/oauth2"

// Provider mock
type Provider struct {
	RedirectURLCall struct {
		Receives struct {
			State string
		}
		// not used since the config decides the url
		Returns struct {
			RedirectURL string
		}
	}
	Config *oauth2.Config
}

// RedirectURL mock
func (p *Provider) RedirectURL(state string) string {
	p.RedirectURLCall.Receives.State = state
	return p.Config.AuthCodeURL(state)
	//return p.RedirectURLCall.Returns.RedirectURL, p.RedirectURLCall.Returns.Error
}
