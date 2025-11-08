package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BradenHooton/pinecone-api/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// ARRANGE
	secret := "test-secret"
	userID := uuid.New()
	email := "test@example.com"

	token, err := jwt.GenerateToken(userID, email, secret, 24)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
	})

	rr := httptest.NewRecorder()

	// Create a test handler that checks the context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxUserID := r.Context().Value(UserIDKey)
		ctxEmail := r.Context().Value(UserEmailKey)

		assert.Equal(t, userID.String(), ctxUserID)
		assert.Equal(t, email, ctxEmail)

		w.WriteHeader(http.StatusOK)
	})

	middleware := Auth(secret)
	handler := middleware(testHandler)

	// ACT
	handler.ServeHTTP(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	// ARRANGE
	secret := "test-secret"

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	middleware := Auth(secret)
	handler := middleware(testHandler)

	// ACT
	handler.ServeHTTP(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "missing authentication token")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// ARRANGE
	secret := "test-secret"

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: "invalid.jwt.token",
	})

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	middleware := Auth(secret)
	handler := middleware(testHandler)

	// ACT
	handler.ServeHTTP(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid token")
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	// ARRANGE
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()
	email := "test@example.com"

	token, err := jwt.GenerateToken(userID, email, correctSecret, 24)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
	})

	rr := httptest.NewRecorder()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	})

	middleware := Auth(wrongSecret)
	handler := middleware(testHandler)

	// ACT
	handler.ServeHTTP(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetUserIDFromContext_Success(t *testing.T) {
	// ARRANGE
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), UserIDKey, userID.String())

	// ACT
	retrievedID, err := GetUserIDFromContext(ctx)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, userID, retrievedID)
}

func TestGetUserIDFromContext_MissingUserID(t *testing.T) {
	// ARRANGE
	ctx := context.Background()

	// ACT
	_, err := GetUserIDFromContext(ctx)

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID not found in context")
}

func TestGetUserIDFromContext_InvalidUUID(t *testing.T) {
	// ARRANGE
	ctx := context.WithValue(context.Background(), UserIDKey, "not-a-valid-uuid")

	// ACT
	_, err := GetUserIDFromContext(ctx)

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID format")
}
