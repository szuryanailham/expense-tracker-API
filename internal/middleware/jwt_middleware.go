package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		userID, err := auth.ParseJWT(parts[1], []byte(env.GetString("JWT_SECRET","")))
		if err != nil {
			if err.Error() == "token expired" {
				http.Error(w, "JWT expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "JWT invalid", http.StatusUnauthorized)
			return
		}

		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusUnauthorized)
			return
		}
		pgUUID := pgtype.UUID{
			Bytes: parsedUUID,
			Valid: true,
		}
		ctx := context.WithValue(r.Context(), userIDKey, pgUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
		})
}

	func GetUserID(ctx context.Context) (pgtype.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(pgtype.UUID)
	return userID, ok
}