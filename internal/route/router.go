package route

import (
	"user-api/internal/controller"
	"user-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router sets up all application routes
type Router struct {
	userController *controller.UserController
	authMiddleware gin.HandlerFunc
}

// NewRouter creates a new router instance
func NewRouter(
	userController *controller.UserController,
	authMiddleware gin.HandlerFunc,
) *Router {
	return &Router{
		userController: userController,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() *gin.Engine {
	// Create Gin router
	router := gin.New()

	// Add global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.ErrorHandlerMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "User API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userController.Register)
			auth.POST("/login", r.userController.Login)
			auth.POST("/refresh", r.userController.RefreshToken)
		}

		// User routes (auth required)
		users := v1.Group("/users")
		users.Use(r.authMiddleware)
		{
			// Current user routes
			users.GET("/me", r.userController.GetProfile)
			users.PUT("/me", r.userController.UpdateProfile)

			// Admin-only routes
			users.GET("", r.userController.GetAllUsers) // AdminMiddleware is applied in the service layer

			// User management routes (with authorization checks in service layer)
			users.GET("/:id", r.userController.GetUserByID)
			users.PATCH("/:id/toggle-block", r.userController.ToggleBlock)
		}
	}

	// Legacy API routes for backward compatibility
	api := router.Group("/api")
	{
		// Authentication routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.userController.Register)
			auth.POST("/login", r.userController.Login)
			auth.POST("/refresh", r.userController.RefreshToken)
		}

		// User routes (auth required)
		users := api.Group("/users")
		users.Use(r.authMiddleware)
		{
			// Current user routes
			users.GET("/me", r.userController.GetProfile)
			users.PUT("/me", r.userController.UpdateProfile)

			// Admin-only routes
			users.GET("", r.userController.GetAllUsers)

			// User management routes
			users.GET("/:id", r.userController.GetUserByID)
			users.PATCH("/:id/toggle-block", r.userController.ToggleBlock)
		}
	}

	return router
}
