package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/deepakjacob/restyle/config"
	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/handlers"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/oauth"
	"github.com/deepakjacob/restyle/service"
	"github.com/deepakjacob/restyle/templates"
	"github.com/deepakjacob/restyle/util"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func setupRouteHandlers() *mux.Router {
	auth, err := setupAuth()
	if err != nil {
		logger.Log.Fatal("initialization err", zap.Error(err))
		return nil
	}
	r := mux.NewRouter()

	r.HandleFunc("/auth", auth.Handle)
	r.HandleFunc("/auth/callback", auth.HandleCallback)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/error", errorHandler)

	return r
}
func main() {
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

func setupAuth() (*handlers.OAuth2, error) {
	ctx := config.BootstrapCtx(context.Background())
	cfg, err := oauth.Config(ctx)
	if err != nil {
		return nil, err
	}
	fsClient, _ := db.New(ctx)
	userServiceImpl := &service.UserServiceImpl{fsClient}
	auth := &handlers.OAuth2{
		Provider: &oauth.Provider{
			HTTPClient: oauth.Client,
			Config:     &oauth.OAuth2Configurer{Config: cfg},
		},
		UserService: userServiceImpl,
		RandStr:     util.RandStr,
	}
	return auth, nil
}

var indexHandler = func(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("index").Parse(templates.Index)
	t.Execute(w, nil)
}

var errorHandler = func(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("error").Parse(templates.Error)
	t.Execute(w, nil)
}
