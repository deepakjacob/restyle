package domain

//UserToken  contents of the jwt web token
type UserToken struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
}
