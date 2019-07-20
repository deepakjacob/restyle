package mocks

import (
	"github.com/deepakjacob/restyle/domain"
)

type Signer struct {
	SignEncryptCall struct {
		Receives struct {
			User *domain.UserToken
		}
		Returns struct {
			SignedValue string
			Err         error
		}
	}
	DecryptCall struct {
		Receives struct {
			SignedValue string
		}
		Returns struct {
			User *domain.UserToken
			Err  error
		}
	}
}

func (s *Signer) SignEncrypt(user *domain.UserToken) (string, error) {
	s.SignEncryptCall.Receives.User = user
	return s.SignEncryptCall.Returns.SignedValue, s.SignEncryptCall.Returns.Err
}

func (s *Signer) Decrypt(value string) (*domain.UserToken, error) {
	s.DecryptCall.Receives.SignedValue = value
	return s.DecryptCall.Returns.User, s.DecryptCall.Returns.Err
}
