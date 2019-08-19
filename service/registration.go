package service

import (
	"context"
	"errors"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
)

// RegistrationService to register and check validity mobile user registration
type RegistrationService interface {
	RegisterMobileUser(string, string) (*domain.RegistrationStatus, error)
	// RegisterUser(*domain.User) (RegistrationStatus, error)
	VerifyCode(context.Context, string, string) (bool, error)
	// GenerateCode(string) (bool, error)
}

// RegistrationServiceImpl impl for the registration service
type RegistrationServiceImpl struct {
	FireStoreService FireStoreService
}

// RegisterMobileUser perform user registration using phone number
func (rs *RegistrationService) RegisterMobileUser(mobileNumber string, generatedCode string, pin string) (*domain.RegistrationStatus, error) {
	status, err := rs.FireStoreService.RegisterMobileUser(mobileNumber, generatedCode, pin)
	if err != nil {
		logger.Log.Error("service:registration:firestore", zap.Error(err))
		return nil, errors.New("error while verifying user")
	}
	return status, nil
}

// VerifyCode verifies the code provided by the user
func (rs *RegistrationServiceImpl) VerifyCode(ctx context.Context, mobileNumber string, receivedCode string) (bool, error) {
	b, err := rs.FireStoreService.VerifyCode(ctx, mobileNumber, receivedCode)
	if err != nil {
		logger.Log.Error("service:registration:verification:firestore", zap.Error(err))
		return false, errors.New("error while verifying user")
	}
	if b == false {
		return false, errors.New("incorrect code")
	}
	return true, nil
}
