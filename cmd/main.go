package main

import (
	"context" // Import context for graceful shutdown
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog-app/internal/config"
	"blog-app/internal/database"
	"blog-app/internal/handler"
	"blog-app/internal/repository"
	"blog-app/internal/router"
	"blog-app/internal/service"
)

func main() {
	// 0. Load Configuration
	cfg := config.LoadConfig()

	// 1. Run Database Migrations
	database.RunMigrations()

	// 2. Initialize Database Connection Pool
	database.InitDB()
	defer database.CloseDB() // Ensure DB connection is closed on exit

	// 3. Setup Dependencies (Repository, Service, Handler)
	userRepo := repository.NewUserRepository(database.DBPool)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)

	// 4. Setup Router
	appRouter := router.NewRouter(userHandler)

	// 5. Start Server with Graceful Shutdown
	server := &http.Server{
		Addr:         ":8080",
		Handler:      appRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Goroutine to start the server
	go func() {
		fmt.Println("Server starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful Shutdown
	// Listen for OS signals to gracefully shutdown the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM) // Capture Ctrl+C and Kubernetes termination
	<-stop                                             // Block until a signal is received

	log.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
