package main

import (
	"fmt"
	"log"

	"github.com/photobooth/backend/config"
	"github.com/photobooth/backend/internal/api/routes"
	"github.com/photobooth/backend/pkg/database"
)

func main() {
	fmt.Println("PhotoBooth Backend Server")
	log.Println("Server starting...")

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize configuration
	cfg := config.GetConfig()

	// Initialize database connection
	if err := database.Connect(
		cfg.CouchbaseURL, 
		cfg.CouchbaseUsername, 
		cfg.CouchbasePassword, 
		cfg.CouchbaseBucket, 
		cfg.CouchbaseTimeout,
	); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connected successfully")
	
	// Initialize router
	router := routes.SetupRouter()

	// Start server
	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on %s", address)

	// Run server
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
