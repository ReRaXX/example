package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusBlocked UserStatus = "blocked"
)

// User represents the user entity in the database
type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FullName  string     `gorm:"type:varchar(255);not null" json:"fullName"`
	BirthDate time.Time  `gorm:"type:date;not null" json:"birthDate"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Role      UserRole   `gorm:"type:user_role;default:'user'" json:"role"`
	Status    UserStatus `gorm:"type:user_status;default:'active'" json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
