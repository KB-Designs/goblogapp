package handler

import (
	"blog-app/internal/models"
	"blog-app/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// RegisterUser handles new user registration.
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.RegisterUser(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
		} else {
			log.Printf("Error registering user: %v", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginUser handles user login and issues JWT tokens.
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UsernameOrEmail == "" || req.Password == "" {
		http.Error(w, "Username/email and password are required", http.StatusBadRequest)
		return
	}

	tokens, err := h.userService.LoginUser(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized) // 401 Unauthorized
		} else {
			log.Printf("Error logging in user: %v", err)
			http.Error(w, "Failed to login", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}
