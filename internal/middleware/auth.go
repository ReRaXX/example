package middleware

import (
	"net/http"
	"strings"

	"user-api/internal/apperror"
	"user-api/internal/entity"
	"user-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("Missing authorization header", nil))
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("Invalid authorization header format", nil))
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("Missing token", nil))
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &service.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("Invalid or expired token", err))
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("Invalid token", nil))
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*service.JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("Invalid token claims", nil))
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("userID", claims.ID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// AdminMiddleware creates a middleware that requires admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User not authenticated", nil))
			c.Abort()
			return
		}

		role, ok := userRole.(entity.UserRole)
		if !ok || role != entity.UserRoleAdmin {
			c.JSON(http.StatusForbidden, apperror.ForbiddenError("Admin access required", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ErrorHandlerMiddleware handles application errors globally
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// If it's an AppError, return it with proper status code
			if appErr, ok := err.(*apperror.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
				return
			}

			// For other errors, return internal server error
			c.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
			return
		}
	}
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
