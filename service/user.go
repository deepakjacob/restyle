package service

import (
	"context"

	"github.com/deepakjacob/restyle/db"
	"github.com/deepakjacob/restyle/domain"
)

// UserServiceImpl impl for interface
type UserServiceImpl struct {
	FireStore *db.FireStore
}

// Find access user
func (u *UserServiceImpl) Find(ctx context.Context, email string) (*domain.User, error) {
	return u.FireStore.Find(ctx, email)
}
