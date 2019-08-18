package service

import (
	"context"
	"os"

	"github.com/deepakjacob/restyle/config"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/sms"
	"go.uber.org/zap"
)

//SmsStatus status of the sms text sending operation
type SmsStatus struct {
	MobileNumber string `json:"mobile_number"`
}

// SmsService implementation for the SmsService
type SmsService struct {
	client *sms.Twilio
}

// SmsService return firestore client connection
func NewSmsService(ctx context.Context) (*SmsService, error) {
	env, err := config.Getenv(ctx)
	if err != nil {
		logger.Log.Error("error in pulling env vars details for twilio", zap.Error(err))
		return nil, err
	}
	client := sms.NewTwilioClient(env.TwilioAccountSid, env.TwilioAuthToken)
	return &SmsService{client}, nil
}

// Send text message
func (s SmsService) Send(ctx context.Context, mobileNum string, msg string) (*SmsStatus, error) {
	_, _, err := s.client.SendSMS(os.Getenv("TWILIO_NUMBER"), mobileNum, msg, "", "")
	if err != nil {
		return nil, err
	}

	return &SmsStatus{MobileNumber: mobileNum}, nil
}
