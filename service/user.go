package service

import (
	"context"
	"fmt"

	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
)

// UserService get user from firestore
type UserService interface {
	Find(context.Context, string) (*domain.User, error)
}

// UserServiceImpl impl for interface
type UserServiceImpl struct {
	FireStore *db.FireStore
}

// Find user with the provided email. The error thrown need not be an actual error.
// For e.g,this is valid when a new user tries to signup for the first time
func (u *UserServiceImpl) Find(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.FireStore.User(ctx, email)
	if err != nil {
		logger.Log.Error("service:user", zap.Error(err))

		return nil, fmt.Errorf("No use with email - %s", email)
	}
	return user, nil
}
