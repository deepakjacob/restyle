package mocks

import "github.com/deepakjacob/restyle/oauth"

// GoogleClient mock
type GoogleClient struct {
	GetCall struct {
		Receives struct {
			URL string
		}
		Returns struct {
			GoogleUser *oauth.GoogleUser
			Err        error
		}
	}
}

// Get call mock
func (g *GoogleClient) Get(url string) (*oauth.GoogleUser, error) {
	g.GetCall.Receives.URL = url
	return g.GetCall.Returns.GoogleUser, g.GetCall.Returns.Err
}
