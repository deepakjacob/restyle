package service

import (
	"context"
	"errors"
	"io"

	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/oauth"
	"github.com/deepakjacob/restyle/storage"
	"go.uber.org/zap"
)

// UploadService contracts
type UploadService interface {
	Upload(context.Context, *domain.ImgAttrs, io.Reader) error
}

// UploadServiceImpl impl for interface
type UploadServiceImpl struct {
	FireStore    *db.FireStore
	CloudStorage *storage.CloudStorage
	RandStr      func() string
}

// Upload for uploading image and associated data
func (u *UploadServiceImpl) Upload(
	ctx context.Context, attrs *domain.ImgAttrs, r io.Reader) error {
	user, err := oauth.UserFromCtx(ctx)
	if err != nil {
		logger.Log.Error("service:upload", zap.Error(err))
		return errors.New("user not found in context")
	}
	err = u.FireStore.Upload(ctx, user, attrs)
	if err != nil {
		// TODO: add user details in logging preferably a proxy id
		logger.Log.Error("service:upload:firestore", zap.Error(err))
		return errors.New("error in saving image attrbutes for db")
	}
	bucket := getDefaultBucket()
	prefixed := derivePrefix(attrs, u.RandStr)
	err = u.CloudStorage.Upload(ctx, user, bucket, attrs, r, prefixed)
	if err != nil {
		// TODO: add user details in logging preferably a proxy id
		logger.Log.Error("service:upload:storage", zap.Error(err))
		return errors.New("error in saving image in cloud")
	}
	return nil
}

// TODO: optimize this
func derivePrefix(attrs *domain.ImgAttrs, f func() string) string {
	return (attrs.ObjType + "/" + attrs.Material + "/" + f())
}

// TODO: get the project id from the env config
func getDefaultBucket() string {
	return "project_up"
}
