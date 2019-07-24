package db

import (
	"context"
	"fmt"

	"github.com/deepakjacob/restyle/domain"
	"google.golang.org/api/iterator"
)

// List look for images matching search criteria
func (fs *FireStore) List(ctx context.Context, user *domain.User,
	attrs *domain.ImgSearch, pattern string) (*domain.ImgSearchResult, error) {
	objects := make([]string, 0)

	iter := fs.Collection("images").Where("obj_type", "==", attrs.ObjType).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		m := doc.Data()
		objects = append(objects, fmt.Sprintf("%v", m["file"]))
	}
	res := &domain.ImgSearchResult{
		List: objects,
	}
	return res, nil
}
