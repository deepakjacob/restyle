package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/oauth"
	"github.com/deepakjacob/restyle/util"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Config get the config for authentication
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

func setupRouteHandlers() *mux.Router {

	r := mux.NewRouter()

	// ----------------------- Login / Logout    -----------------------
	cfg, err := Config()
	if err != nil {
		logger.Log.Fatal("application missing mandatory params", zap.Error(err))
	}
	auth := &handlers.OAuth2{
		Provider: &oauth.Provider{
			Config: cfg,
		},
		RandStr: util.RandStr,
	}

	r.HandleFunc("/auth", auth.Handle)
	r.HandleFunc("/auth/callback", auth.HandleCallback)

	return r
}
func main() {

	// Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	if err := logger.Init(-1, "2006-01-02T15:04:05Z07:00"); err != nil {
		log.Fatal(err)
	}

	r := setupRouteHandlers()
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
