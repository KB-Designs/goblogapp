package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // File source for migrations
	"github.com/joho/godotenv"
)

// RunMigrations applies database migrations.
func RunMigrations() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming environment variables are set.")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required database environment variables (DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)")
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New(
		"file:///home/brian/blog-app/migrations", // Path to our migration files
		databaseURL)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	// Apply all up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not apply migrations: %v", err)
	}

	log.Println("Database migrations applied successfully!")
}
