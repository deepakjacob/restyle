package storage

import (
	"context"

	"cloud.google.com/go/storage"
)

// CloudStorage connection
type CloudStorage struct {
	*storage.Client
}

// New return cloud storage connection
func New(ctx context.Context) (*CloudStorage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &CloudStorage{client}, nil
}
