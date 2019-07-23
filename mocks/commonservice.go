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
			Error error
		}
	}
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
			Error error
		}
	}
}
