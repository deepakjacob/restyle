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
	RegisterMobileUser(context.Context, *domain.RegistrationAttrs) (*domain.RegistrationStatus, error)
	// RegisterUser(*domain.User) (RegistrationStatus, error)
	VerifyCode(context.Context, *domain.RegistrationAttrs) (bool, error)
	// GenerateCode(string) (bool, error)
}

// RegistrationServiceImpl impl for the registration service
type RegistrationServiceImpl struct {
	FireStoreService FireStoreService
}

// RegisterMobileUser perform user registration using phone number
func (rs *RegistrationServiceImpl) RegisterMobileUser(ctx context.Context, attrs *domain.RegistrationAttrs) (*domain.RegistrationStatus, error) {
	status, err := rs.FireStoreService.RegisterMobileUser(ctx, attrs)
	if err != nil {
		logger.Log.Error("service:registration:firestore", zap.Error(err))
		return nil, errors.New("error while verifying user")
	}
	return status, nil
}

// VerifyCode verifies the code provided by the user
func (rs *RegistrationServiceImpl) VerifyCode(ctx context.Context, attrs *domain.RegistrationAttrs) (bool, error) {
	b, err := rs.FireStoreService.VerifyCode(ctx, attrs)
	if err != nil {
		logger.Log.Error("service:verification:firestore", zap.Error(err))
		return false, errors.New("error while verifying user")
	}
	if b == false {
		return false, errors.New("incorrect code")
	}
	return true, nil
}
