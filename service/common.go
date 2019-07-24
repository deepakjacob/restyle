package service

import (
	"context"
	"io"

	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/storage"
)

// FireStoreService firestore operations
type FireStoreService interface {
	Upload(context.Context, *domain.User, *domain.ImgAttrs, string) error
	List(context.Context, *domain.User, *domain.ImgSearch, string) (*domain.ImgSearchResult, error)
}

// CloudStorageService cloud storage operations
type CloudStorageService interface {
	Upload(context.Context, *domain.User, string, *domain.ImgAttrs, io.Reader, string) error
}

// FireStoreServiceImpl impl for firestore
type FireStoreServiceImpl struct {
	FireStore *db.FireStore
}

// Upload persists attributes of a file upload
func (fs *FireStoreServiceImpl) Upload(ctx context.Context, user *domain.User,
	imgAttrs *domain.ImgAttrs, fileName string) error {
	return fs.FireStore.Upload(ctx, user, imgAttrs, fileName)
}

// List persists attributes of a file upload
func (fs *FireStoreServiceImpl) List(ctx context.Context, user *domain.User,
	imgAttrs *domain.ImgSearch, pattern string) (*domain.ImgSearchResult, error) {
	return fs.FireStore.List(ctx, user, imgAttrs, pattern)
}

// CloudStorageServiceImpl impl for firestore
type CloudStorageServiceImpl struct {
	CloudStorage *storage.CloudStorage
}

// Upload persists the uploaded file
func (cs *CloudStorageServiceImpl) Upload(ctx context.Context, user *domain.User,
	bucket string, imgAttrs *domain.ImgAttrs, r io.Reader, prefixed string) error {
	return cs.Upload(ctx, user, bucket, imgAttrs, r, prefixed)
}
