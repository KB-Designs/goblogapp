package main

import (
	"blog-app/internal/database" // Import our new package
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. Run Database Migrations
	database.RunMigrations()

	// 2. Our simple health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API is healthy!")
	})

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
