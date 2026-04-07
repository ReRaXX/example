package controller

import (
	"net/http"

	"user-api/internal/apperror"
	"user-api/internal/dto"
	"user-api/internal/entity"
	"user-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userService service.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account and return authentication tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} apperror.AppError
// @Failure 409 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.ValidationError("Invalid request body", err))
		return
	}

	response, err := c.userService.Register(&req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   response,
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Failure 403 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.ValidationError("Invalid request body", err))
		return
	}

	response, err := c.userService.Login(&req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/auth/refresh [post]
func (c *UserController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.ValidationError("Invalid request body", err))
		return
	}

	response, err := c.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetProfile handles getting current user profile
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/users/me [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User not authenticated", nil))
		return
	}

	response, err := c.userService.GetProfile(userID.(uuid.UUID))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// UpdateProfile handles updating current user profile
// @Summary Update current user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/users/me [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User not authenticated", nil))
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.ValidationError("Invalid request body", err))
		return
	}

	response, err := c.userService.UpdateProfile(userID.(uuid.UUID), &req)
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetAllUsers handles getting all users (admin only)
// @Summary Get all users
// @Description Get all users (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserResponse
// @Failure 401 {object} apperror.AppError
// @Failure 403 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	userRole, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User not authenticated", nil))
		return
	}

	response, err := c.userService.GetAllUsers(userRole.(entity.UserRole))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetUserByID handles getting a user by ID
// @Summary Get user by ID
// @Description Get a user by ID (with authorization checks)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Failure 403 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User not authenticated", nil))
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User role not found", nil))
		return
	}

	targetIDStr := ctx.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.ValidationError("Invalid user ID", err))
		return
	}

	response, err := c.userService.GetUserByID(userID.(uuid.UUID), targetID, userRole.(entity.UserRole))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// ToggleBlock handles toggling user block status
// @Summary Toggle user block status
// @Description Toggle a user's blocked status (admin can block anyone, users can block themselves)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} apperror.AppError
// @Failure 401 {object} apperror.AppError
// @Failure 403 {object} apperror.AppError
// @Failure 404 {object} apperror.AppError
// @Failure 500 {object} apperror.AppError
// @Router /api/users/{id}/toggle-block [patch]
func (c *UserController) ToggleBlock(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User not authenticated", nil))
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, apperror.UnauthorizedError("User role not found", nil))
		return
	}

	targetIDStr := ctx.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apperror.ValidationError("Invalid user ID", err))
		return
	}

	response, err := c.userService.ToggleBlock(userID.(uuid.UUID), targetID, userRole.(entity.UserRole))
	if err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			ctx.JSON(appErr.StatusCode, appErr)
			return
		}
		ctx.JSON(http.StatusInternalServerError, apperror.InternalServerError("Internal server error", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}
