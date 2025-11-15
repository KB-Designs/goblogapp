package repository

import (
	"blog-app/internal/config"
	"blog-app/internal/models"
	"context"
)

type UserRepository struct{}

func (r UserRepository) CreateUser(user models.User) error {
	query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
    `

	// Use ExecContext instead of Exec
	_, err := config.DB.ExecContext(context.Background(), query,
		user.Username,
		user.Email,
		user.Password,
	)

	return err
}

func (r UserRepository) FindByEmail(email string) (models.User, error) {
	query := `
        SELECT id, username, email, password, created_at 
        FROM users
        WHERE email = $1 LIMIT 1
    `

	var user models.User

	// Use QueryRowContext instead of QueryRow
	err := config.DB.QueryRowContext(context.Background(), query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	return user, err
}

func NewUserRepository() UserRepository {
	return UserRepository{}
}
