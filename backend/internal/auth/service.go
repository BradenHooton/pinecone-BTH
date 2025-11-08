package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/BradenHooton/pinecone-api/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailAlreadyExists is returned when email is already registered
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidPassword is returned when password doesn't meet requirements
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
)

// Repository defines the interface for user data access
type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash, name string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// Service handles authentication business logic
type Service struct {
	repo          Repository
	jwtSecret     string
	jwtExpiryHours int
}

// NewService creates a new auth service
func NewService(repo Repository, jwtSecret string, jwtExpiryHours int) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register registers a new user
func (s *Service) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, string, error) {
	// Validate password length
	if len(req.Password) < 8 {
		return nil, "", ErrInvalidPassword
	}

	// Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && err != ErrUserNotFound {
		return nil, "", fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, "", ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.repo.CreateUser(ctx, req.Email, passwordHash, req.Name)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := verifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// hashPassword hashes a plain text password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// verifyPassword compares a hashed password with a plain text password
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
