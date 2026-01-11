package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/photobooth/backend/config"
	"github.com/photobooth/backend/internal/api/handlers"
	"github.com/photobooth/backend/internal/api/middleware"
)

func SetupRouter() *gin.Engine {
	cfg := config.GetConfig()

	// Initialize Gin router
	router := gin.Default()

	// Setup middleware
	// router.Use(middleware.SetupCORS(cfg.AllowedOrigins))
	router.Use(gin.Recovery())
	// router.Use(middleware.Logging())

	// Health routes
	health := router.Group("/health")
	{
		health.GET("", handlers.Health)
		health.GET("/live", handlers.HealthLive)
		health.GET("/ready", handlers.HealthReady)
	}

 	api := router.Group(cfg.APIPrefix)
    {
        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", handlers.Register)
            auth.POST("/login", handlers.Login)
            auth.POST("/refresh", handlers.Refresh)
            auth.POST("/logout", handlers.Logout)
            auth.GET("/me", handlers.Me)
        }

        // Protected routes (add JWT middleware)
        protected := api.Group("")
        protected.Use(middleware.JWTAuth())
        {
            // User routes
            user := protected.Group("/users")
            {
                user.GET("/profile", handlers.GetProfile)
                user.PATCH("/profile", handlers.UpdateProfile)
                user.GET("/stats", handlers.GetStats)
            }

            // Photos routes
            photos := protected.Group("/photos")
            {
                photos.POST("/upload", handlers.UploadPhoto)
                photos.GET("/list", handlers.ListPhotos)
                photos.GET("/get", handlers.GetPhoto)
                photos.DELETE("/delete", handlers.DeletePhoto)
            }
        }
    }

	return router
}
