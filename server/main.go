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
	"github.com/deepakjacob/restyle/signer"
	"github.com/deepakjacob/restyle/storage"
	"github.com/deepakjacob/restyle/templates"
	"github.com/deepakjacob/restyle/util"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func setupRouteHandlers() *mux.Router {
	logger.Log.Info("bootstrapping context")
	ctx := config.BootstrapCtx(context.Background())
	authconfig, err := oauth.Config(ctx)
	if err != nil {
		logger.Log.Fatal("botstrapping context", zap.Error(err))
		return nil
	}
	logger.Log.Info("init connections to firestore")
	fsClient, err := db.New(ctx)
	if err != nil {
		logger.Log.Fatal("firestore", zap.Error(err))
		return nil
	}
	logger.Log.Info("init connections to cloud storage")
	csClient, err := storage.New(ctx)
	if err != nil {
		logger.Log.Fatal("cloud storage", zap.Error(err))
		return nil
	}
	firestoreService := &service.FireStoreServiceImpl{fsClient}
	logger.Log.Debug("main::init", zap.Any("firestore.service", firestoreService))

	cloudStorageService := &service.CloudStorageServiceImpl{csClient}
	userService := &service.UserServiceImpl{fsClient}
	uploadService := &service.UploadServiceImpl{
		FireStoreService:    firestoreService,
		CloudStorageService: cloudStorageService,
		RandStr:             util.RandStr,
	}

	listService := &service.ListServiceImpl{
		FireStoreService: firestoreService,
	}

	logger.Log.Debug("main::init", zap.Any("list.service", listService))

	auth, err := auth(authconfig, userService)
	if err != nil {
		logger.Log.Fatal("error in oauth setup", zap.Error(err))
		return nil
	}
	upload := upload(uploadService)
	list := list(listService)

	r := mux.NewRouter()

	r.HandleFunc("/auth", auth.Handle)
	r.HandleFunc("/auth/callback", auth.HandleCallback)
	r.HandleFunc("/error", errorHandler)

	s := r.PathPrefix("/api").Subrouter()
	s.Use(auth.Middleware)
	s.HandleFunc("/", indexHandler)
	s.HandleFunc("/upload", upload.Handle)
	s.HandleFunc("/list", list.Handle)

	return r
}
func main() {
	initLogger()

	r := setupRouteHandlers()
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Log.Info("server about to listen @", zap.String("address", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}

func upload(uploadService service.UploadService) *handlers.Upload {
	upload := &handlers.Upload{UploadService: uploadService}
	return upload
}

func list(listService service.ListService) *handlers.List {
	list := &handlers.List{ListService: listService}
	return list
}

func auth(config *oauth2.Config, userService service.UserService) (*handlers.OAuth2, error) {
	signer := &signer.Signer{}
	auth := &handlers.OAuth2{
		Provider: &oauth.ProviderImpl{
			HTTPClient: oauth.Client,
			Config:     &oauth.OAuth2Configurer{Config: config},
		},
		UserService: userService,
		RandStr:     util.RandStr,
		Signer:      signer,
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

func initLogger() {
	if err := logger.Init(-1, "2006-01-02T15:04:05Z07:00"); err != nil {
		log.Fatal(err)
	}
}
