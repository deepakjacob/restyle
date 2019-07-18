package oauth

// GoogleUser struct google combined scopes profile and email returns
type GoogleUser struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_name"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	NickName      string `json:"nick_name"`
	OrgID         string `json:"org_id"`
	UserID        string `json:"user_id"`
}
