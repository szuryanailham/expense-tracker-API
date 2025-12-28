package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/szuryanailham/expense-tracker/internal/auth"
	"github.com/szuryanailham/expense-tracker/internal/env"
)

type contextKey string

const userIDKey contextKey = "user_id"

func JWTAuth(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}
		log.Println("JWT_SECRET:", env.GetString("JWT_SECRET",""))

		userID, err := auth.ParseJWT(
			parts[1],
			[]byte(env.GetString("JWT_SECRET","")),
		)

		if err != nil {
			log.Println("JWT PARSE ERROR:", err)
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
		})
}

	func GetUserID(ctx context.Context) (string, bool) {
		userID, ok := ctx.Value(userIDKey).(string)
		return userID, ok
	}