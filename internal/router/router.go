package router

import (
	"blog-app/internal/handler"
	"net/http"
)

// NewRouter sets up all application routes.
func NewRouter(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy!"))
	})
	mux.HandleFunc("POST /register", userHandler.RegisterUser)
	mux.HandleFunc("POST /login", userHandler.LoginUser)

	return mux
}
