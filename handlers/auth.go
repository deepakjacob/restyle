package handlers

import (
	"net/http"
	"time"
)

type provider interface {
	RedirectURL(string) string
}

// OAuth2 provides auth services
type OAuth2 struct {
	Provider provider
	RandStr  func() string
}

// Handle handles /auth
func (o *OAuth2) Handle(w http.ResponseWriter, r *http.Request) {
	state := o.RandStr()
	cookie := http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Value:    state,
		Name:     "state",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	redirectURL := o.Provider.RedirectURL(state)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// HandleCallback handles /auth/callback
func (o *OAuth2) HandleCallback(w http.ResponseWriter, r *http.Request) {
	state, _ := r.Cookie("state")
	urlState := r.URL.Query().Get("state")
	if urlState != state.Value {
		http.Redirect(w, r, "/error?status=403reason=state", http.StatusForbidden)
		return
	}
	code := r.URL.Query().Get("code")
	w.Write([]byte(code))
}
