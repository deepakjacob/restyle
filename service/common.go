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
	List(context.Context, *domain.User, *domain.ImgSearch,
		string) (*domain.ImgSearchResult, error)
	// TODO: VerifyCode to return VerificationStatus struct insead of bool
	VerifyCode(context.Context, string, string) (bool, error)
	RegisterMobileUser(context.Context, string, string,
		string) (*domain.RegistrationStatus, error)
}

// CloudStorageService cloud storage operations
type CloudStorageService interface {
	Upload(context.Context, *domain.User, string, *domain.ImgAttrs, io.Reader, string) error
}

// FireStoreServiceImpl impl for firestore
type FireStoreServiceImpl struct {
	FireStore *db.FireStore
}

// VerifyCode verfied user provided code
func (fs *FireStoreServiceImpl) VerifyCode(
	ctx context.Context, mobileNumber string,
	verificationCode string) (bool, error) {
	return fs.FireStore.VerifyCode(ctx, mobileNumber, verificationCode)
}

// RegisterMobileUser verfied user provided code
func (fs *FireStoreServiceImpl) RegisterMobileUser(
	ctx context.Context, attrs *domain.RegistrationAttrs) (*domain.RegistrationStatus, error) {
	return fs.FireStore.RegisterMobileUser(ctx, attrs)
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
	return cs.CloudStorage.Upload(ctx, user, bucket, imgAttrs, r, prefixed)
}
