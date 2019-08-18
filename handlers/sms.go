package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deepakjacob/restyle/service"
	"github.com/deepakjacob/restyle/util"
)

//Sms for sending text messages
type Sms struct {
	SmsService *service.SmsService
	// TODO for personalization
	// CustomizationService CustomizationService
}

// Handle handles the sending of text message
func (s *Sms) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobileNum := r.Form["mobile_num"]
	msg := fmt.Sprintf("Restyle - your code is %d", util.RandNum())
	status, err := s.SmsService.Send(context.Background(), mobileNum[0], msg)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	smsStatus, err := json.Marshal(status)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(smsStatus)
}
