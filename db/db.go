package db

import (
	"context"

	"cloud.google.com/go/firestore"
)

// FireStore specifics
type FireStore struct {
	*firestore.Client
}

// New return a fire client connection
func New(ctx context.Context, projectID string) (*FireStore, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FireStore{client}, nil
}
