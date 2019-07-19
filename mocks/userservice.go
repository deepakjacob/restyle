package mocks

import (
	"context"

	"github.com/deepakjacob/restyle/domain"
)

// UserService user srevice
type UserService struct {
	FindCall struct {
		Receives struct {
			Ctx   context.Context
			Email string
		}
		Returns struct {
			User  *domain.User
			Error error
		}
	}
}

// Find mock
func (u *UserService) Find(ctx context.Context, email string) (*domain.User, error) {
	u.FindCall.Receives.Ctx = ctx
	u.FindCall.Receives.Email = email
	return u.FindCall.Returns.User, u.FindCall.Returns.Error
}
