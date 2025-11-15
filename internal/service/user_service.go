package service

import (
	"blog-app/internal/models"
	"blog-app/internal/repository"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return UserService{repo}
}

func (s UserService) Register(username, email, password string) error {

	// Check if email already exists
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user object
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPass),
	}

	// Save to DB
	return s.repo.CreateUser(user)
}
