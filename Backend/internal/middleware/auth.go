package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"taskflow/internal/utils"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			fmt.Println("AUTH HEADER:", authHeader)

			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]
			fmt.Println("TOKEN:", tokenStr)

			claims, err := utils.ValidateToken(tokenStr, secret)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			fmt.Println("CLAIMS:", claims, "ERROR:", err)
			ctx := context.WithValue(r.Context(), utils.UserContextKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
