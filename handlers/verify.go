package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/service"
)

//Verification regsitering and verifying users
type Verification struct {
	RegistrationService service.RegistrationService
}

// Handle handles the sending of text message
func (rs *Verification) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobileNumber := r.Form["mobile_number"]
	mobileCode := r.Form["mobile_code"]
	mobilePin := r.Form["mobile_pin"]
	attrs := &domain.RegistrationAttrs{
		MobileNumber:     mobileNumber[0],
		VerificationCode: mobileCode[0],
		Pin:              mobilePin[0],
	}
	rStatus, err := rs.RegistrationService.VerifyCode(r.Context(), attrs)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rjson, err := json.Marshal(rStatus)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rjson)
}
