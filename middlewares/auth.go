package middlewares

import (
	"net/http"

	"github.com/deepakjacob/restyle/logger"
	"go.uber.org/zap"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(oauth2Cookie)
		if err != nil {
			logger.Log.Error("Not authorised - unable to get auth cookie",
				zap.Error(err))
			// http.Error(w, fmt.Sprintf("Not Authorized: %v", err), http.StatusUnauthorized)
			http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
			return
		}
		u, err := verifyUser(cookie)
		if err != nil {
			logger.Log.Error("Not authorised", zap.Error(err))
			// http.Error(w, fmt.Sprintf("Not Authorized: %v", err), http.StatusUnauthorized)
			http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
			return
		}
		ctx := r.Context()
		user, err := service.GetUser(ctx, u.Email)
		if err != nil {
			logger.Log.Error("User not found", zap.Error(err))
			//http.Error(w, "User not found", http.StatusUnauthorized)
			http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
			return
		}
		usrCtx := WithUser(ctx, user)
		next.ServeHTTP(w, r.WithContext(usrCtx))
	})
}
