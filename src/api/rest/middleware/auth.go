package middleware

import (
	"context"
	"net/http"
	"nullableocean-postupashki/src/usecases"
	"strings"
)

var (
	AuthorizationHeader = "Authorization"
	AuthorizationPrefix = "Bearer "
)

func CommiterAuth(accessHeader, accessToken string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(accessHeader)

			if token == "" || token != accessToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Auth(sessiongService usecases.Session) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := strings.CutPrefix(r.Header.Get("Authorization"), AuthorizationPrefix)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			session, err := sessiongService.Get(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
