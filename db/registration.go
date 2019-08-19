package db

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
)

// VerifyCode verifies code provided by the user
func (fs *FireStore) VerifyCode(ctx context.Context, mobileNumber string, verificationCode string) (bool, error) {
	doc := fmt.Sprintf("musers/%s", mobileNumber)
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

		verificationStatus, err := doc.DataAt("verified")
		if err != nil {
			return err
		}

		ok := code == verificationCode && verificationStatus == false
		if ok {
			return tx.Update(muser, []firestore.Update{{Path: "verified", Value: true}})
		}
		return errors.New("incorrect code or status")
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
