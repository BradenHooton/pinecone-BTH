package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BradenHooton/pinecone-api/pkg/jwt"
	"github.com/google/uuid"
)

// Context keys for user information
type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
)

// Auth is a middleware that validates JWT tokens from cookies
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get JWT token from cookie
			cookie, err := r.Cookie("jwt_token")
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "missing authentication token")
				return
			}

			// Validate token
			claims, err := jwt.ValidateToken(cookie.Value, jwtSecret)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
				return
			}

			// Parse user ID from string to UUID
			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid user ID in token")
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}

	return userID, nil
}

// GetUserEmailFromContext extracts the user email from the request context
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(UserEmailKey).(string)
	if !ok {
		return "", fmt.Errorf("user email not found in context")
	}

	return email, nil
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
	})
}
