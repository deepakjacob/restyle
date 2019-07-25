package handlers

import (
	"net/http"
	"time"

	"github.com/deepakjacob/restyle/domain"
	"github.com/deepakjacob/restyle/logger"
	"github.com/deepakjacob/restyle/oauth"
	"github.com/deepakjacob/restyle/service"
	"github.com/deepakjacob/restyle/signer"
	"go.uber.org/zap"
)

// OAuth2 provides auth services
type OAuth2 struct {
	Provider    oauth.Provider
	UserService service.UserService
	Signer      signer.JWTSigner
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

	user, err := o.UserService.Find(r.Context(), gUser.Email)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// TODO: change to the below impl to redirect to register a new user
		// logger.Log.Info("auth:callback::new user", zap.String("email", gUser.Email))
		// http.RedirectHandler("/api/newuser", 307)
		return
	}
	logger.Log.Debug("auth:callback::user",
		zap.String("email", user.Email), zap.String("ID", user.UserID))
	ut := &domain.UserToken{
		UserID: user.UserID,
		Email:  user.Email,
	}
	token, err := o.Signer.SignEncrypt(ut)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Value:    token,
		Name:     "_ut",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/api/", http.StatusTemporaryRedirect)
}

//Middleware function to execute before accessing secure urls
func (o *OAuth2) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("_ut")
		if err != nil {
			http.Error(w,
				http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		u, err := o.Signer.Decrypt(cookie.Value)
		if err != nil {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		user, err := o.UserService.Find(r.Context(), u.Email)
		if err != nil {
			http.Error(w,
				http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		usrCtx := oauth.UserToCtx(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(usrCtx))
	})
}
