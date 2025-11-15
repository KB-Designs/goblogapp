package router

import (
	"blog-app/internal/handler"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/health", handler.HealthCheck)
	http.HandleFunc("/register", handler.RegisterUser)

}
