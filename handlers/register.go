package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/service"
)

//Registration regsitering and verifying users
type Registration struct {
	RegistrationService service.RegistrationService
}

// Handle handles the sending of text message
func (rs *Registration) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mobileNumber := r.Form["mobile_number"]
	mobileCode := r.Form["mobile_code"]
	attrs := &domain.RegistrationAttrs{
		// TODO: remove array syntax
		MobileNumber:     mobileNumber[0],
		VerificationCode: mobileCode[0],
	}
	rStatus, err := rs.RegistrationService.RegisterMobileUser(r.Context(), attrs)
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
