package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleModerator  AdminRole = "moderator"
	AdminRoleSupport    AdminRole = "support"
)

type Admin struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	AdminRole   AdminRole `json:"admin_role" db:"admin_role"`
	Permissions []string  `json:"permissions" db:"permissions"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Joined
	User *User `json:"user,omitempty" db:"-"`
}

type AuditLog struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	AdminID    uuid.UUID              `json:"admin_id" db:"admin_id"`
	Action     string                 `json:"action" db:"action"`
	TargetType string                 `json:"target_type" db:"target_type"`
	TargetID   uuid.UUID              `json:"target_id" db:"target_id"`
	Details    map[string]interface{} `json:"details,omitempty" db:"details"`
	IPAddress  string                 `json:"ip_address" db:"ip_address"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
}