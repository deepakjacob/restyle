package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/oauth"
)

// UserService get user from firestore
type UserService interface {
	Find(context.Context, string) (*domain.User, error)
}

type provider interface {
	RedirectURL(string) string
	GoogleUser(context.Context, string) (*oauth.GoogleUser, error)
}

// OAuth2 provides auth services
type OAuth2 struct {
	Provider    provider
	UserService UserService
	RandStr     func() string
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
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	gUser, err := o.Provider.GoogleUser(r.Context(), code)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	_, err = o.UserService.Find(r.Context(), gUser.Email)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// check user in the system before forward to the app
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}
