package auth

import (
	"context"
	"testing"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, email, passwordHash, name string) (*models.User, error) {
	args := m.Called(ctx, email, passwordHash, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// TestRegister_Success tests successful user registration
func TestRegister_Success(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, "test-secret", 24)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	expectedUser := &models.User{
		ID:    uuid.New(),
		Email: req.Email,
		Name:  req.Name,
	}

	mockRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, ErrUserNotFound)
	mockRepo.On("CreateUser", ctx, req.Email, mock.AnythingOfType("string"), req.Name).Return(expectedUser, nil)

	// ACT
	user, token, err := service.Register(ctx, req)

	// ASSERT
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	mockRepo.AssertExpectations(t)
}

// TestRegister_EmailAlreadyExists tests registration with existing email
func TestRegister_EmailAlreadyExists(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, "test-secret", 24)

	req := &models.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	existingUser := &models.User{
		ID:    uuid.New(),
		Email: req.Email,
		Name:  "Existing User",
	}

	mockRepo.On("GetUserByEmail", ctx, req.Email).Return(existingUser, nil)

	// ACT
	user, token, err := service.Register(ctx, req)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, ErrEmailAlreadyExists, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestRegister_InvalidPassword tests registration with short password
func TestRegister_InvalidPassword(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, "test-secret", 24)

	req := &models.RegisterRequest{
		Email:    "test@example.com",
		Password: "short", // Less than 8 characters
		Name:     "Test User",
	}

	// ACT
	user, token, err := service.Register(ctx, req)

	// ASSERT
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockRepo.AssertNotCalled(t, "CreateUser")
}

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, "test-secret", 24)

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Hash the password for comparison
	hashedPassword, err := hashPassword(req.Password)
	require.NoError(t, err)

	existingUser := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         "Test User",
	}

	mockRepo.On("GetUserByEmail", ctx, req.Email).Return(existingUser, nil)

	// ACT
	user, token, err := service.Login(ctx, req)

	// ASSERT
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, req.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

// TestLogin_UserNotFound tests login with non-existent email
func TestLogin_UserNotFound(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, "test-secret", 24)

	req := &models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, ErrUserNotFound)

	// ACT
	user, token, err := service.Login(ctx, req)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

// TestLogin_InvalidPassword tests login with wrong password
func TestLogin_InvalidPassword(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo, "test-secret", 24)

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Hash a different password
	hashedPassword, err := hashPassword("correctpassword")
	require.NoError(t, err)

	existingUser := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         "Test User",
	}

	mockRepo.On("GetUserByEmail", ctx, req.Email).Return(existingUser, nil)

	// ACT
	user, token, err := service.Login(ctx, req)

	// ASSERT
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredentials, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}
