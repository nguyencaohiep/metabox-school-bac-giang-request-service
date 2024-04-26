package middleware

import (
	"context"
	"net/http"
	"strings"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))
			if tokenString == "" {
				next.ServeHTTP(w, r)
				// http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "has_token", true))
			r = r.WithContext(context.WithValue(r.Context(), "token", tokenString))

			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// tokenString := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(r.Header.Get("Authorization")), "Bearer"))
			// if tokenString == "" {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

			// result, err := service_shared.TokenVerify(r.Context(), &service_shared.TokenVerifyCommand{
			// 	Token:      tokenString,
			// 	TypeVerify: model.TypeVerifyBasic,
			// })
			// if err != nil || result == nil {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

			// r = r.WithContext(context.WithValue(r.Context(), "account_id", result.AccountID))
			// r = r.WithContext(context.WithValue(r.Context(), "user_id", result.UserID))
			// r = r.WithContext(context.WithValue(r.Context(), "email", result.Email))
			// r = r.WithContext(context.WithValue(r.Context(), "token", tokenString))
			next.ServeHTTP(w, r)
		})
	}
}
