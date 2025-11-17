package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"blog-app/internal/config"
	"blog-app/internal/models"
	"blog-app/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("username or email already exists")
	ErrInvalidCredentials = errors.New("invalid username/email or password")
	ErrTokenInvalid       = errors.New("invalid or expired token")
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	RegisterUser(ctx context.Context, req *models.UserRegisterRequest) (*models.User, error)
	LoginUser(ctx context.Context, req *models.UserLoginRequest) (*models.AuthTokens, error)
	GenerateTokens(userID, username string) (*models.AuthTokens, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

// userService implements UserService.
type userService struct {
	userRepo repository.UserRepository
	config   *config.AppConfig
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository, cfg *config.AppConfig) UserService {
	return &userService{userRepo: userRepo, config: cfg}
}

// RegisterUser handles the registration of a new user.
func (s *userService) RegisterUser(ctx context.Context, req *models.UserRegisterRequest) (*models.User, error) {
	// Check if username or email already exists
	existingUser, err := s.userRepo.GetUserByUsernameOrEmail(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}
	existingUser, err = s.userRepo.GetUserByUsernameOrEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// Store user in the database
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %w", err)
	}

	return user, nil
}

// LoginUser handles user login and generates JWT tokens.
func (s *userService) LoginUser(ctx context.Context, req *models.UserLoginRequest) (*models.AuthTokens, error) {
	user, err := s.userRepo.GetUserByUsernameOrEmail(ctx, req.UsernameOrEmail)
	if err != nil {
		return nil, ErrInvalidCredentials // User not found
	}

	// Compare provided password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials // Password mismatch
	}

	// Generate and return tokens
	return s.GenerateTokens(user.ID, user.Username)
}

// GenerateTokens creates access and refresh JWT tokens.
func (s *userService) GenerateTokens(userID, username string) (*models.AuthTokens, error) {
	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(s.config.AccessTokenExp).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh Token (for now, just a longer-lived token with same claims)
	// In a real app, refresh tokens would ideally be stored in DB and used for one-time exchange
	refreshTokenClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(s.config.RefreshTokenExp).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &models.AuthTokens{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

// ValidateToken validates a JWT token and returns its claims.
func (s *userService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
