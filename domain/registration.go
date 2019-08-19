package domain

// RegistrationAttrs params required for mobile user registration
type RegistrationAttrs struct {
	MobileNumber     string
	VerificationCode string
	Pin              string
}

// RegistrationStatus the status of registation
type RegistrationStatus struct {
	Pin        string `json:"pin"`
	StatusCd   string `json:"status_cd"`
	StatusDesc string `json:"status_desc"`
}
