package domain

// User struct google combined scopes profile and email returns
type User struct {
	Email         string `json:"email" firestore:"email"`
	VerifiedEmail bool   `json:"verified_name" firestore:"verified_email"`
	Name          string `json:"name" firestore:"name"`
	GivenName     string `json:"given_name" firestore:"given_name"`
	FamilyName    string `json:"family_name" firestore:"family_name"`
	Picture       string `json:"picture" firestore:"picture"`
	NickName      string `json:"nick_name" firestore:"nick_name"`
	OrgID         string `json:"org_id" firestore:"org_id"`
	UserID        string `json:"user_id" firestore:"user_id"`
}
