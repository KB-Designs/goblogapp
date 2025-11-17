package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// DBPool is our global connection pool.
var DBPool *pgxpool.Pool

// InitDB initializes the PostgreSQL connection pool.
func InitDB() {
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

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	DBPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	err = DBPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL!")
}

// CloseDB closes the PostgreSQL connection pool.
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection pool closed.")
	}
}
