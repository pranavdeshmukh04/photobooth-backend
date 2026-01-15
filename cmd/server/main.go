package main

import (
	"fmt"
	"log"

	"github.com/photobooth/backend/config"
	"github.com/photobooth/backend/internal/api/handlers"
	"github.com/photobooth/backend/internal/api/routes"
	"github.com/photobooth/backend/internal/repositories"
	"github.com/photobooth/backend/internal/services"
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

	// Initialize repositories
	scope := database.Bucket.Scope("_default")
	usersCollection := scope.Collection("users")
	userRepo := repositories.NewUserRepository(usersCollection, database.Cluster)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	log.Println("Services initialized successfully")

	// Initialize router
	router := routes.SetupRouter(authHandler)
	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on %s", address)

	// Run server
	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
