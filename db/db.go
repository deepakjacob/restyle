package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/deepakjacob/restyle/config"
)

// FireStore connection
type FireStore struct {
	*firestore.Client
}

// New return firestore client connection
func New(ctx context.Context) (*FireStore, error) {
	env, err := config.Getenv(ctx)
	if err != nil {
		return nil, err
	}
	client, err := firestore.NewClient(ctx, env.ProjectID)
	if err != nil {
		return nil, err
	}
	return &FireStore{client}, nil
}
