package mocks

import (
	"context"
	"io"

	"github.com/deepakjacob/restyle/domain"
)

// FireStoreService upload service
type FireStoreService struct {
	UploadCall struct {
		Receives struct {
			Ctx      context.Context
			User     *domain.User
			ImgAttrs *domain.ImgAttrs
		}
		Returns struct {
			Err error
		}
	}
}

// Upload persists attributes of a file upload
func (fs *FireStoreService) Upload(ctx context.Context, user *domain.User,
	imgAttrs *domain.ImgAttrs, fileName string) error {
	fs.UploadCall.Receives.Ctx = ctx
	fs.UploadCall.Receives.User = user
	fs.UploadCall.Receives.ImgAttrs = imgAttrs
	return fs.UploadCall.Returns.Err
}

// CloudStorageService upload service
type CloudStorageService struct {
	UploadCall struct {
		Receives struct {
			Ctx      context.Context
			User     *domain.User
			Bucket   string
			ImgAttrs *domain.ImgAttrs
			Reader   io.Reader
			Prefix   string
		}
		Returns struct {
			Err error
		}
	}
}

// Upload persists the uploaded file
func (cs *CloudStorageService) Upload(ctx context.Context, user *domain.User,
	bucket string, imgAttrs *domain.ImgAttrs, r io.Reader, prefixed string) error {
	cs.UploadCall.Receives.Ctx = ctx
	cs.UploadCall.Receives.User = user
	cs.UploadCall.Receives.Bucket = bucket
	cs.UploadCall.Receives.ImgAttrs = imgAttrs
	cs.UploadCall.Receives.Reader = r
	cs.UploadCall.Receives.Prefix = prefixed
	return cs.UploadCall.Returns.Err
}
