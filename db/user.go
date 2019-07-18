package db

import (
	"context"
	"log"

	"github.com/deepakjacob/restyle/domain"
)

// Find user
func (fs *FireStore) Find(ctx context.Context, email string) (*domain.User, error) {
	docsnap, err := fs.Collection("users").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user *domain.User
	if err := docsnap.DataTo(&user); err != nil {
		return nil, err
	}
	user.Email = email
	log.Fatal(user)
	return user, nil
}
