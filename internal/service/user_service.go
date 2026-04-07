package service

import (
	"strings"
	"time"

	"user-api/internal/apperror"
	"user-api/internal/dto"
	"user-api/internal/entity"
	"user-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	ID    uuid.UUID        `json:"id"`
	Email string           `json:"email"`
	Role  entity.UserRole  `json:"role"`
	jwt.RegisteredClaims
}

// UserService defines the interface for user business logic
type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(refreshToken string) (*dto.AuthResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	GetAllUsers(requesterRole entity.UserRole) ([]dto.UserResponse, error)
	GetUserByID(requesterID, targetID uuid.UUID, requesterRole entity.UserRole) (*dto.UserResponse, error)
	ToggleBlock(requesterID, targetID uuid.UUID, requesterRole entity.UserRole) (*dto.UserResponse, error)
}

// userService implements UserService
type userService struct {
	userRepo  repository.UserRepository
	accessSecret  string
	refreshSecret string
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository, accessSecret, refreshSecret string) UserService {
	return &userService{
		userRepo:      userRepo,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// Register creates a new user and returns authentication tokens
func (s *userService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	email := strings.ToLower(strings.TrimSpace(req.Email))
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, apperror.ConflictError("Email is already registered", nil)
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, apperror.ValidationError("Invalid birth date format", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to hash password", err)
	}

	// Create user
	user := &entity.User{
		FullName:  strings.TrimSpace(req.FullName),
		BirthDate: birthDate,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      entity.UserRoleUser,
		Status:    entity.UserStatusActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, apperror.InternalServerError("Failed to create user", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to generate refresh token", err)
	}

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *userService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	
	// Find user with password
	user, err := s.userRepo.FindByEmailWithPassword(email)
	if err != nil {
		return nil, apperror.UnauthorizedError("Invalid email or password", nil)
	}

	// Check if user is blocked
	if user.Status == entity.UserStatusBlocked {
		return nil, apperror.ForbiddenError("Account is blocked", nil)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperror.UnauthorizedError("Invalid email or password", nil)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to generate refresh token", err)
	}

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *userService) RefreshToken(refreshToken string) (*dto.AuthResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, apperror.UnauthorizedError("Invalid refresh token", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, apperror.UnauthorizedError("Invalid token claims", nil)
	}

	// Get user from database
	user, err := s.userRepo.FindByID(claims.ID)
	if err != nil {
		return nil, apperror.UnauthorizedError("User not found", err)
	}

	if user.Status == entity.UserStatusBlocked {
		return nil, apperror.ForbiddenError("Account is blocked", nil)
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to generate access token", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to generate refresh token", err)
	}

	return &dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetProfile returns the current user's profile
func (s *userService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.NotFoundError("User not found", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// UpdateProfile updates the current user's profile
func (s *userService) UpdateProfile(userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	updates := make(map[string]interface{})

	if req.FullName != nil {
		updates["full_name"] = strings.TrimSpace(*req.FullName)
	}

	if req.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return nil, apperror.ValidationError("Invalid birth date format", err)
		}
		updates["birth_date"] = birthDate
	}

	if len(updates) == 0 {
		// No updates provided
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, apperror.NotFoundError("User not found", err)
		}
		response := dto.ToUserResponse(user)
		return &response, nil
	}

	user, err := s.userRepo.Update(userID, updates)
	if err != nil {
		return nil, apperror.InternalServerError("Failed to update user", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// GetAllUsers returns all users (admin only)
func (s *userService) GetAllUsers(requesterRole entity.UserRole) ([]dto.UserResponse, error) {
	if requesterRole != entity.UserRoleAdmin {
		return nil, apperror.ForbiddenError("Only admins can view all users", nil)
	}

	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, apperror.InternalServerError("Failed to fetch users", err)
	}

	return dto.ToUserResponses(users), nil
}

// GetUserByID returns a user by ID (with authorization checks)
func (s *userService) GetUserByID(requesterID, targetID uuid.UUID, requesterRole entity.UserRole) (*dto.UserResponse, error) {
	// Authorization check: users can only view their own profile unless they're admin
	if requesterRole != entity.UserRoleAdmin && requesterID != targetID {
		return nil, apperror.ForbiddenError("You can only view your own profile", nil)
	}

	user, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return nil, apperror.NotFoundError("User not found", err)
	}

	response := dto.ToUserResponse(user)
	return &response, nil
}

// ToggleBlock toggles a user's blocked status
func (s *userService) ToggleBlock(requesterID, targetID uuid.UUID, requesterRole entity.UserRole) (*dto.UserResponse, error) {
	// Authorization check: only admins can block other users, users can block themselves
	if requesterRole != entity.UserRoleAdmin && requesterID != targetID {
		return nil, apperror.ForbiddenError("You can only block/unblock your own account", nil)
	}

	user, err := s.userRepo.FindByID(targetID)
	if err != nil {
		return nil, apperror.NotFoundError("User not found", err)
	}

	// Toggle status
	newStatus := entity.UserStatusActive
	if user.Status == entity.UserStatusActive {
		newStatus = entity.UserStatusBlocked
	}

	updatedUser, err := s.userRepo.Update(targetID, map[string]interface{}{
		"status": newStatus,
	})
	if err != nil {
		return nil, apperror.InternalServerError("Failed to update user status", err)
	}

	response := dto.ToUserResponse(updatedUser)
	return &response, nil
}

// generateAccessToken generates a new JWT access token
func (s *userService) generateAccessToken(user *entity.User) (string, error) {
	claims := &JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

// generateRefreshToken generates a new JWT refresh token
func (s *userService) generateRefreshToken(user *entity.User) (string, error) {
	claims := &JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}
