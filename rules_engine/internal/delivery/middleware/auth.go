package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	authservice "rules-engine/internal/clients/auth_service"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(authClient *authservice.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header format"})
				return
			}

			token := parts[1]
			user, err := authClient.VerifyToken(token)
			if err != nil {
				sendJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*authservice.UserResponse, bool) {
	user, ok := ctx.Value(UserContextKey).(*authservice.UserResponse)
	return user, ok
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
