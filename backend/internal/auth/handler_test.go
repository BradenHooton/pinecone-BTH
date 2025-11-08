package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the service
type MockService struct {
	mock.Mock
}

func (m *MockService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// TestHandleRegister_Success tests successful registration
func TestHandleRegister_Success(t *testing.T) {
	// ARRANGE
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: reqBody.Email,
		Name:  reqBody.Name,
	}

	mockService.On("Register", mock.Anything, reqBody).Return(expectedUser, "test-token", nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// ACT
	handler.HandleRegister(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check response body
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	data := response["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.Equal(t, reqBody.Email, user["email"])
	assert.Equal(t, reqBody.Name, user["name"])

	// Check cookie
	cookies := rr.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "jwt_token", cookies[0].Name)
	assert.Equal(t, "test-token", cookies[0].Value)
	assert.True(t, cookies[0].HttpOnly)

	mockService.AssertExpectations(t)
}

// TestHandleRegister_EmailAlreadyExists tests registration with existing email
func TestHandleRegister_EmailAlreadyExists(t *testing.T) {
	// ARRANGE
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := &models.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	mockService.On("Register", mock.Anything, reqBody).Return(nil, "", ErrEmailAlreadyExists)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// ACT
	handler.HandleRegister(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusConflict, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	errorData := response["error"].(map[string]interface{})
	assert.Contains(t, errorData["message"], "email already exists")

	mockService.AssertExpectations(t)
}

// TestHandleRegister_InvalidJSON tests registration with invalid JSON
func TestHandleRegister_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockService := new(MockService)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// ACT
	handler.HandleRegister(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockService.AssertNotCalled(t, "Register")
}

// TestHandleLogin_Success tests successful login
func TestHandleLogin_Success(t *testing.T) {
	// ARRANGE
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: reqBody.Email,
		Name:  "Test User",
	}

	mockService.On("Login", mock.Anything, reqBody).Return(expectedUser, "test-token", nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// ACT
	handler.HandleLogin(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check response body
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	data := response["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.Equal(t, reqBody.Email, user["email"])

	// Check cookie
	cookies := rr.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "jwt_token", cookies[0].Name)
	assert.Equal(t, "test-token", cookies[0].Value)

	mockService.AssertExpectations(t)
}

// TestHandleLogin_InvalidCredentials tests login with wrong credentials
func TestHandleLogin_InvalidCredentials(t *testing.T) {
	// ARRANGE
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockService.On("Login", mock.Anything, reqBody).Return(nil, "", ErrInvalidCredentials)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// ACT
	handler.HandleLogin(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Contains(t, response, "error")
	errorData := response["error"].(map[string]interface{})
	assert.Contains(t, errorData["message"], "invalid credentials")

	mockService.AssertExpectations(t)
}

// TestHandleLogout_Success tests successful logout
func TestHandleLogout_Success(t *testing.T) {
	// ARRANGE
	handler := NewHandler(nil) // No service needed for logout

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	rr := httptest.NewRecorder()

	// ACT
	handler.HandleLogout(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that cookie is cleared
	cookies := rr.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "jwt_token", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.Equal(t, -1, cookies[0].MaxAge)
}
