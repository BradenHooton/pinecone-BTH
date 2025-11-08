package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken_Success(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	expiryHours := 24

	// ACT
	token, err := GenerateToken(userID, email, secret, expiryHours)

	// ASSERT
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_ValidToken(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	expiryHours := 24

	token, err := GenerateToken(userID, email, secret, expiryHours)
	require.NoError(t, err)

	// ACT
	claims, err := ValidateToken(token, secret)

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestValidateToken_InvalidSecret(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"
	wrongSecret := "wrong-secret"
	userID := uuid.New()
	email := "test@example.com"
	expiryHours := 24

	token, err := GenerateToken(userID, email, secret, expiryHours)
	require.NoError(t, err)

	// ACT
	_, err = ValidateToken(token, wrongSecret)

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	expiryHours := -1 // Already expired

	token, err := GenerateToken(userID, email, secret, expiryHours)
	require.NoError(t, err)

	// ACT
	_, err = ValidateToken(token, secret)

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestValidateToken_MalformedToken(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"
	malformedToken := "not.a.valid.jwt.token"

	// ACT
	_, err := ValidateToken(malformedToken, secret)

	// ASSERT
	assert.Error(t, err)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"

	// ACT
	_, err := ValidateToken("", secret)

	// ASSERT
	assert.Error(t, err)
}

func TestClaims_ExpiresAt(t *testing.T) {
	// ARRANGE
	secret := "test-secret-key"
	userID := uuid.New()
	email := "test@example.com"
	expiryHours := 24

	token, err := GenerateToken(userID, email, secret, expiryHours)
	require.NoError(t, err)

	// ACT
	claims, err := ValidateToken(token, secret)
	require.NoError(t, err)

	// ASSERT
	// Token should expire approximately 24 hours from now
	expectedExpiry := time.Now().Add(time.Duration(expiryHours) * time.Hour)
	actualExpiry := claims.ExpiresAt.Time

	// Allow 1 minute tolerance for test execution time
	assert.WithinDuration(t, expectedExpiry, actualExpiry, 1*time.Minute)
}
