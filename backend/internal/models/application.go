package models

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	ApplicationStatusPending    ApplicationStatus = "pending"
	ApplicationStatusReviewed   ApplicationStatus = "reviewed"
	ApplicationStatusShortlisted ApplicationStatus = "shortlisted"
	ApplicationStatusAccepted   ApplicationStatus = "accepted"
	ApplicationStatusRejected   ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn  ApplicationStatus = "withdrawn"
)

type Application struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	JobID       uuid.UUID         `json:"job_id" db:"job_id"`
	WorkerID    uuid.UUID         `json:"worker_id" db:"worker_id"`
	CoverLetter string            `json:"cover_letter,omitempty" db:"cover_letter"`
	ProposedRate float64          `json:"proposed_rate,omitempty" db:"proposed_rate"`
	Status      ApplicationStatus `json:"status" db:"status"`
	SeenAt      *time.Time        `json:"seen_at,omitempty" db:"seen_at"`
	RespondedAt *time.Time        `json:"responded_at,omitempty" db:"responded_at"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`

	// Joined
	Job    *Job    `json:"job,omitempty" db:"-"`
	Worker *Worker `json:"worker,omitempty" db:"-"`
}