package service

import (
	"context"
	"errors"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/oauth"
	"go.uber.org/zap"
)

// ListService contracts
type ListService interface {
	// TODO: create domain.SearchAttrs instead of using domain.ImgAttrs
	List(context.Context, *domain.ImgSearch) (*domain.ImgSearchResult, error)
}

// ListServiceImpl impl for interface
type ListServiceImpl struct {
	FireStoreService FireStoreService
	// CloudStorageService CloudStorageService
}

// List of images based on search attrbutes
func (u *ListServiceImpl) List(
	ctx context.Context, attrs *domain.ImgSearch) (*domain.ImgSearchResult, error) {
	user, err := oauth.UserFromCtx(ctx)
	if err != nil {
		logger.Log.Error("service:list", zap.Error(err))
		return nil, errors.New("user not found in context")
	}
	pattern := searchPattern(attrs)
	results, err := u.FireStoreService.List(ctx, user, attrs, pattern)
	if err != nil {
		// TODO: add user details in logging preferably a proxy id
		logger.Log.Error("service:list:firestore", zap.Error(err))
		return nil, errors.New("error in saving image attrbutes for db")
	}
	return results, nil
}

// TODO: optimize this
func searchPattern(attrs *domain.ImgSearch) string {
	return attrs.ObjType + "/" + attrs.Material
}
