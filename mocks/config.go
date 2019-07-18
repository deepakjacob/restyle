package mocks

import (
	"context"

	"golang.org/x/oauth2"
)

// Config mock oauth2 config
type Config struct {
	AuthCodeURLCall struct {
		Receives struct {
			State string
		}
		Returns struct {
			URL string
		}
	}
	ExchangeCall struct {
		Receives struct {
			Ctx  context.Context
			Code string
		}
		Returns struct {
			Token *oauth2.Token
			Error error
		}
	}
}

// AuthCodeURL mock
func (c *Config) AuthCodeURL(state string) string {
	c.AuthCodeURLCall.Receives.State = state
	return c.AuthCodeURLCall.Returns.URL
}

// Exchange mock
func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	c.ExchangeCall.Receives.Ctx = ctx
	c.ExchangeCall.Receives.Code = code
	return c.ExchangeCall.Returns.Token, c.ExchangeCall.Returns.Error

}
