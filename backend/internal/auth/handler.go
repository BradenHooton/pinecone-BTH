package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
)

// ServiceInterface defines the interface for auth service
type ServiceInterface interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, string, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// Handler handles HTTP requests for authentication
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new auth handler
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// HandleRegister handles user registration
func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, token, err := h.service.Register(r.Context(), &req)
	if err != nil {
		switch err {
		case ErrEmailAlreadyExists:
			respondWithError(w, http.StatusConflict, "Email already exists")
		case ErrInvalidPassword:
			respondWithError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		default:
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	// Set JWT cookie
	setJWTCookie(w, token)

	// Return user data
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"user": user.ToResponse(),
	})
}

// HandleLogin handles user login
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, token, err := h.service.Login(r.Context(), &req)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		default:
			respondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	// Set JWT cookie
	setJWTCookie(w, token)

	// Return user data
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user": user.ToResponse(),
	})
}

// HandleLogout handles user logout
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear JWT cookie
	clearJWTCookie(w)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Logged out successfully",
	})
}

// setJWTCookie sets the JWT token in an HTTP-only cookie
func setJWTCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 24 hours
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// clearJWTCookie clears the JWT token cookie
func clearJWTCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	})
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
		"meta": map[string]interface{}{
			"timestamp": time.Now().UTC(),
		},
	})
}
