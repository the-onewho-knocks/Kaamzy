package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type UserStatus string

const (
	RoleWorker   UserRole = "worker"
	RoleBusiness UserRole = "business"
	RoleAdmin    UserRole = "admin"
)

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
	StatusBanned    UserStatus = "banned"
)

type User struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	Email            string     `json:"email" db:"email"`
	PasswordHash     string     `json:"-" db:"password_hash"`
	Role             UserRole   `json:"role" db:"role"`
	Status           UserStatus `json:"status" db:"status"`
	IsEmailVerified  bool       `json:"is_email_verified" db:"is_email_verified"`
	IsPhoneVerified  bool       `json:"is_phone_verified" db:"is_phone_verified"`
	Phone            string     `json:"phone,omitempty" db:"phone"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
	LastLoginAt      *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}