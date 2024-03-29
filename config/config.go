package config

import (
	"context"
	"errors"
	"os"

	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
)

type envKey string

// EnvKey key under which env vars are saved in context
var EnvKey envKey

// Env struct for env vars
type Env struct {
	ProjectID        string
	AppCredentials   string
	ClientID         string
	ClientSecret     string
	JWTKey           string
	JWTEncKey        string
	RedirectURL      string
	TwilioNumber     string
	TwilioAccountSid string
	TwilioAuthToken  string
}

// BootstrapCtx context initialization
func BootstrapCtx(parent context.Context) context.Context {
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("REDIRECT_URL")
	twilioNumber := os.Getenv("TWILIO_NUMBER")
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")

	logger.Log.Info("env vars",
		zap.String("projectid", twilioAccountSid))

	env := &Env{
		ProjectID:    projectID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// for twilio
		TwilioNumber:     twilioNumber,
		TwilioAccountSid: twilioAccountSid,
		TwilioAuthToken:  twilioAuthToken,
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
