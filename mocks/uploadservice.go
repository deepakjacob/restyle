package mocks

import (
	"context"
	"io"

	"github.com/deepakjacob/restyle/domain"
)

// UploadService upload service
type UploadService struct {
	UploadCall struct {
		Receives struct {
			Ctx      context.Context
			ImgAttrs *domain.ImgAttrs
			Reader   io.Reader
		}
		Returns struct {
			Error error
		}
	}
}

// Upload mock
func (u *UploadService) Upload(ctx context.Context, imgAttrs *domain.ImgAttrs, r io.Reader) error {
	u.UploadCall.Receives.Ctx = ctx
	u.UploadCall.Receives.ImgAttrs = imgAttrs
	u.UploadCall.Receives.Reader = r
	return u.UploadCall.Returns.Error
}
