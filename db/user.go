package db

import (
	"context"

	"github.com/deepakjacob/restyle/domain"
)

// User gets the user with the provided email.
func (fs *FireStore) User(ctx context.Context, email string) (*domain.User, error) {
	docsnap, err := fs.Collection("users").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user *domain.User
	if err := docsnap.DataTo(&user); err != nil {
		return nil, err
	}
	user.Email = email
	return user, nil
}
