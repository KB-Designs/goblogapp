package main

import (
	"fmt"
	"log"
	"net/http"

	"blog-app/internal/config"
	"blog-app/internal/router"
)

func main() {
	// Connect to DB
	config.ConnectDB()

	// Register routes
	router.RegisterRoutes()

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
