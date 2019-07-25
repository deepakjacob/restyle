package db

import (
	"context"
	"fmt"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

// List look for images matching search criteria
func (fs *FireStore) List(ctx context.Context, user *domain.User,
	attrs *domain.ImgSearch, pattern string) (*domain.ImgSearchResult, error) {
	objects := make([]string, 0)

	logger.Log.Debug("searching for images")
	images := fs.Collection("images")
	logger.Log.Debug("adding criteria", zap.Any("object type", attrs.ObjType))
	query := images.Where("obj_type", "==", attrs.ObjType)
	query = query.Where("material", "==", "silk")
	iter := query.Documents(ctx)
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
	logger.Log.Debug("returning images", zap.Int("count", len(res.List)))
	return res, nil
}
