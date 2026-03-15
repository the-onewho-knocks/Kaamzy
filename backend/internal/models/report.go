package models

import (
	"time"

	"github.com/google/uuid"
)

type ReportReason string
type ReportStatus string
type ReportTargetType string

const (
	ReportReasonSpam        ReportReason = "spam"
	ReportReasonFakeProfile ReportReason = "fake_profile"
	ReportReasonAbuse       ReportReason = "abuse"
	ReportReasonFraud       ReportReason = "fraud"
	ReportReasonInappropriate ReportReason = "inappropriate_content"
	ReportReasonOther       ReportReason = "other"
)

const (
	ReportStatusPending    ReportStatus = "pending"
	ReportStatusReviewing  ReportStatus = "reviewing"
	ReportStatusResolved   ReportStatus = "resolved"
	ReportStatusDismissed  ReportStatus = "dismissed"
)

const (
	ReportTargetUser    ReportTargetType = "user"
	ReportTargetJob     ReportTargetType = "job"
	ReportTargetMessage ReportTargetType = "message"
)

type Report struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	ReporterID   uuid.UUID        `json:"reporter_id" db:"reporter_id"`
	TargetType   ReportTargetType `json:"target_type" db:"target_type"`
	TargetID     uuid.UUID        `json:"target_id" db:"target_id"`
	Reason       ReportReason     `json:"reason" db:"reason"`
	Description  string           `json:"description,omitempty" db:"description"`
	Status       ReportStatus     `json:"status" db:"status"`
	ReviewedBy   *uuid.UUID       `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewNote   string           `json:"review_note,omitempty" db:"review_note"`
	ReviewedAt   *time.Time       `json:"reviewed_at,omitempty" db:"reviewed_at"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" db:"updated_at"`

	// Joined
	Reporter *User `json:"reporter,omitempty" db:"-"`
}