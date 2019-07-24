package service

import (
	"context"
	"errors"
	"io"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
)

// UploadService contracts
type UploadService interface {
	Upload(context.Context, *domain.ImgAttrs, io.Reader) error
}

// UploadServiceImpl impl of upload. In additition to upload functionality, the
// implementation also provides a mechanism for getting the user who perform the
// upload and to generate a random string. The reason for injecting these methods
// on to struct rather than the interface is that during testing we can effectively
// provide mocks instead of bypassing any existing functionality.
type UploadServiceImpl struct {
	User                func(context.Context) (*domain.User, error)
	RandStr             func() string
	FireStoreService    FireStoreService
	CloudStorageService CloudStorageService
}

// Upload for uploading image and associated data
func (u *UploadServiceImpl) Upload(
	ctx context.Context, attrs *domain.ImgAttrs, r io.Reader) error {
	user, err := u.User(ctx)
	if err != nil {
		logger.Log.Error("service:upload", zap.Error(err))
		return errors.New("user not found in context")
	}
	fileName := u.RandStr()
	prefixed := storagePattern(attrs, fileName)
	err = u.FireStoreService.Upload(ctx, user, attrs, prefixed)
	if err != nil {
		// TODO: add user details in logging preferably a proxy id
		logger.Log.Error("service:upload:firestore", zap.Error(err))
		return errors.New("error in saving image attrbutes for db")
	}
	bucket := getDefaultBucket()
	err = u.CloudStorageService.Upload(ctx, user, bucket, attrs, r, prefixed)
	if err != nil {
		// TODO: add user details in logging preferably a proxy id
		logger.Log.Error("service:upload:storage", zap.Error(err))
		return errors.New("error in saving image in cloud")
	}
	return nil
}

// TODO: optimize this
func storagePattern(attrs *domain.ImgAttrs, name string) string {
	return attrs.ObjType + "/" + attrs.Material + "/" + name
}

// TODO: get the project id from the env config
func getDefaultBucket() string {
	return "project_up"
}
