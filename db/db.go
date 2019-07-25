package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/deepakjacob/restyle/config"
	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
)

// FireStore connection
type FireStore struct {
	*firestore.Client
}

// New return firestore client connection
func New(ctx context.Context) (*FireStore, error) {
	env, err := config.Getenv(ctx)
	if err != nil {
		logger.Log.Error("error in pulling env vars details for firestore", zap.Error(err))
		return nil, err
	}
	client, err := firestore.NewClient(ctx, env.ProjectID)
	if err != nil {
		return nil, err
	}
	return &FireStore{client}, nil
}
