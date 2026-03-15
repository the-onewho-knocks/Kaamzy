package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeJobMatch          NotificationType = "job_match"
	NotificationTypeApplicationUpdate NotificationType = "application_update"
	NotificationTypeNewMessage        NotificationType = "new_message"
	NotificationTypeRatingReceived    NotificationType = "rating_received"
	NotificationTypeJobExpiringSoon   NotificationType = "job_expiring_soon"
	NotificationTypeAccountUpdate     NotificationType = "account_update"
	NotificationTypeSystemAlert       NotificationType = "system_alert"
)

type Notification struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	UserID     uuid.UUID        `json:"user_id" db:"user_id"`
	Type       NotificationType `json:"type" db:"type"`
	Title      string           `json:"title" db:"title"`
	Body       string           `json:"body" db:"body"`
	Payload    map[string]interface{} `json:"payload,omitempty" db:"payload"`
	IsRead     bool             `json:"is_read" db:"is_read"`
	ReadAt     *time.Time       `json:"read_at,omitempty" db:"read_at"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
}