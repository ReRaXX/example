package main

import (
	"log"

	"user-api/internal/config"
	"user-api/internal/controller"
	"user-api/internal/entity"
	"user-api/internal/middleware"
	"user-api/internal/repository"
	"user-api/internal/route"
	"user-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto migrate database schema
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWTAccessSecret, cfg.JWTRefreshSecret)

	// Initialize controllers
	userController := controller.NewUserController(userService)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(cfg.JWTAccessSecret)

	// Initialize router
	router := route.NewRouter(userController, authMiddleware)

	// Setup routes
	engine := router.SetupRoutes()

	// Set Gin mode
	if cfg.Port == "8080" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Health check available at: http://localhost:%s/health", cfg.Port)
	log.Printf("API documentation:")
	log.Printf("  Register: POST http://localhost:%s/api/auth/register", cfg.Port)
	log.Printf("  Login:    POST http://localhost:%s/api/auth/login", cfg.Port)
	log.Printf("  Profile:  GET  http://localhost:%s/api/users/me", cfg.Port)

	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDatabase initializes the database connection
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Create custom enum types if they don't exist
	enumMigrations := []string{
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN CREATE TYPE user_role AS ENUM ('admin', 'user'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN CREATE TYPE user_status AS ENUM ('active', 'blocked'); END IF; END $$;",
	}

	for _, migration := range enumMigrations {
		if err := db.Exec(migration).Error; err != nil {
			return err
		}
	}

	// Auto migrate all entities
	err := db.AutoMigrate(
		&entity.User{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
