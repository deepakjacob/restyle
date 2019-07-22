package storage

import (
	"context"
	"io"

	"github.com/deepakjacob/restyle/domain"
)

// Upload gets the user with the provided email.
func (cs *CloudStorage) Upload(ctx context.Context, user *domain.User,
	bucket string, attrs *domain.ImgAttrs, r io.Reader, prefix string) error {
	wc := cs.Bucket(bucket).Object(prefix).NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
