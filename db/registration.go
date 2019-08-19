package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deepakjacob/restyle/domain"
)

// RegisterMobileUser registers a mobile user
func (fs *FireStore) RegisterMobileUser(ctx context.Context, attrs *domain.RegistrationAttrs) (*domain.RegistrationStatus, error) {
	doc := make(map[string]interface{})
	doc["verification_code"] = attrs.VerificationCode
	doc["pin"] = attrs.Pin
	doc["first_generated_on"] = time.Now()
	doc["first_verified_on"] = nil

	ref := fs.Collection("musers").Doc(attrs.MobileNumber)
	_, err := ref.Set(ctx, doc)
	if err != nil {
		return nil, err
	}
	return &domain.RegistrationStatus{
		Pin:        attrs.Pin,
		StatusCd:   "200",
		StatusDesc: "user registration performed successfully",
	}, nil
}

// VerifyCode verifies code provided by the user
func (fs *FireStore) VerifyCode(ctx context.Context, attrs *domain.RegistrationAttrs) (bool, error) {
	doc := fmt.Sprintf("musers/%s", attrs.MobileNumber)
	muser := fs.Doc(doc)

	err := fs.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(muser)
		if err != nil {
			return err
		}

		code, err := doc.DataAt("verification_code")
		if err != nil {
			return err
		}

		pin, err := doc.DataAt("pin")
		if err != nil {
			return err
		}

		verified, err := doc.DataAt("first_verified_on")
		if err != nil {
			return err
		}

		blocked, err := doc.DataAt("blocked")
		if err != nil {
			return err
		}

		ok := code == attrs.VerificationCode && pin == attrs.Pin && blocked == false
		if ok {
			if verified == nil {
				return tx.Update(muser, []firestore.Update{
					{Path: "verified", Value: true},
					{Path: "first_verified_on", Value: time.Now()},
					{Path: "last_verified_on", Value: time.Now()},
				})
			} else {
				return tx.Update(muser, []firestore.Update{
					{Path: "verified", Value: true},
					{Path: "last_verified_on", Value: time.Now()},
				})
			}
		}
		return errors.New("incorrect code or status")
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
