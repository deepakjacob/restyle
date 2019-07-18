package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/oauth"
	"github.com/deepakjacob/restyle/service"
	"github.com/deepakjacob/restyle/templates"
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

type httpClient struct {
	Client *http.Client
}

func (h *httpClient) Get(url string) (*oauth.GoogleUser, error) {
	response, err := h.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	var user oauth.GoogleUser
	if err := json.Unmarshal([]byte(contents), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user json to struct %v", err.Error())
	}
	return &user, nil
}

func setupRouteHandlers() *mux.Router {
	cfg, err := Config()
	if err != nil {
		logger.Log.Fatal("application missing mandatory params", zap.Error(err))
	}

	fsClient, _ := db.New(context.Background(), os.Getenv("GOOGLE_PROJECT_ID"))

	client := &httpClient{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	userServiceImpl := &service.UserServiceImpl{fsClient}

	auth := &handlers.OAuth2{
		Provider: &oauth.Provider{
			HTTPClient: client,
			Config:     &oauth.OAuth2Configurer{Config: cfg},
		},
		UserService: userServiceImpl,
		RandStr:     util.RandStr,
	}

	r := mux.NewRouter()
	r.HandleFunc("/auth", auth.Handle)
	r.HandleFunc("/auth/callback", auth.HandleCallback)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("index").Parse(templates.Index)
		t.Execute(w, nil)
	})
	r.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("error").Parse(templates.Error)
		t.Execute(w, nil)
	})

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
