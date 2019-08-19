package service

import (
	"context"
	"errors"
)

// RegistrationStatus the status of registation
type RegistrationStatus struct {
	MobileNumber     string `'json:"mobile_number"`
	VerificationCode string `json:"verification_code"`
	Status           string `json:"status"`
}

// RegistrationService to register and check validity mobile user registration
type RegistrationService interface {
	// RegisterPhone(string, string) (RegistrationStatus, error)
	// RegisterUser(*domain.User) (RegistrationStatus, error)
	VerifyCode(context.Context, string, string) (bool, error)
	// GenerateCode(string) (bool, error)
}

// RegistrationServiceImpl impl for the registration service
type RegistrationServiceImpl struct {
	FireStoreService FireStoreService
}

// VerifyCode verifies the code provided by the user
func (rs *RegistrationServiceImpl) VerifyCode(ctx context.Context, mobileNumber string, receivedCode string) (bool, error) {
	b, err := rs.FireStoreService.VerifyCode(ctx, mobileNumber, receivedCode)
	if err != nil {
		return false, err
	}
	if b == false {
		return false, errors.New("incorrect code")
	}
	return true, nil
}
