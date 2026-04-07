package dto

import (
	"time"

	"github.com/google/uuid"
	"user-api/internal/entity"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	FullName  string `json:"fullName" binding:"required" validate:"min=2,max=255"`
	BirthDate string `json:"birthDate" binding:"required" validate:"datetime=2006-01-02"`
	Email     string `json:"email" binding:"required,email" validate:"max=255"`
	Password  string `json:"password" binding:"required" validate:"min=6,max=255"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"max=255"`
	Password string `json:"password" binding:"required" validate:"min=6,max=255"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

// UserResponse represents the user data in responses (without password)
type UserResponse struct {
	ID        uuid.UUID          `json:"id"`
	FullName  string             `json:"fullName"`
	BirthDate string             `json:"birthDate"`
	Email     string             `json:"email"`
	Role      entity.UserRole    `json:"role"`
	Status    entity.UserStatus  `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	FullName  *string `json:"fullName,omitempty" binding:"omitempty" validate:"omitempty,min=2,max=255"`
	BirthDate *string `json:"birthDate,omitempty" binding:"omitempty" validate:"omitempty,datetime=2006-01-02"`
}

// ToUserResponse converts a User entity to UserResponse
func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		BirthDate: user.BirthDate.Format("2006-01-02"),
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponses converts a slice of User entities to UserResponses
func ToUserResponses(users []*entity.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
