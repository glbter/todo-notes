package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"todoNote/internal/model"
	"todoNote/internal/server/http/auth"
	"todoNote/internal/server/http/dto"
	"todoNote/internal/server/http/log"
)

const UserAuthorized = "UserAuthorized"

type Middleware struct {
	auth auth.IAuth
}

func New(auth auth.IAuth) *Middleware {
	return &Middleware{
		auth: auth,
	}
}

func(md *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk, ok := tokenFromRequest(r)
		if !ok {
			writeErrorMessage(w, http.StatusBadRequest, "Malformed token")
			return
		}

		u, err := md.auth.ValidateToken(tk)
		if err != nil {
			writeErrorMessage(w, http.StatusUnauthorized, "Not valid token")
			return
		}

		ctx := contextWithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithUser(ctx context.Context, user model.UserInReq) context.Context {
	return context.WithValue(ctx, UserAuthorized, user)
}

func UserFromContext(r *http.Request, log log.Logger, method string) (model.UserInReq, bool) {
	u, ok := r.Context().Value(UserAuthorized).(model.UserInReq)
	if !ok {
		log.Error(method + " " + "no user in context")
	}

	return u, ok
}

func tokenFromRequest(r *http.Request) (auth.JwtToken, bool) {
	t := r.Header.Get("Authorization")
	splt := strings.Split(t, "Bearer ")
	if len(splt) != 2 {
		return "", false
	}

	return splt[1], true
}


func writeErrorMessage(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(
		dto.Error{
			Message: msg,
			Error: code,
		})

	return
}
