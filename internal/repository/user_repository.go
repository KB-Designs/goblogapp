package repository

import (
	"blog-app/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

// userRepository implements UserRepository using pgxpool.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// CreateUser inserts a new user into the database.
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsernameOrEmail retrieves a user by username or email.
func (r *userRepository) GetUserByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1 OR email = $1`
	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found or database error: %w", err)
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found or database error: %w", err)
	}
	return user, nil
}
