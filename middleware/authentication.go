package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/konaro/line-notify-service/jwtutil"
)

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		split := strings.Split(auth, "Bearer ")

		// check token exists
		if len(split) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := split[1]

		if claims, ok := jwtutil.VerifyToken(token); ok {
			ctx := context.WithValue(r.Context(), "account", claims.Account)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
