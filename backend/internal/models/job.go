package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type JobStatus string
type JobType string
type PaymentType string

const (
	JobStatusDraft     JobStatus = "draft"
	JobStatusOpen      JobStatus = "open"
	JobStatusPaused    JobStatus = "paused"
	JobStatusFilled    JobStatus = "filled"
	JobStatusCompleted JobStatus = "completed"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusExpired   JobStatus = "expired"
)

const (
	JobTypeOnsite  JobType = "onsite"
	JobTypeRemote  JobType = "remote"
	JobTypeHybrid  JobType = "hybrid"
)

const (
	PaymentTypeHourly  PaymentType = "hourly"
	PaymentTypeFixed   PaymentType = "fixed"
	PaymentTypeDaily   PaymentType = "daily"
)

type Job struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	BusinessID       uuid.UUID      `json:"business_id" db:"business_id"`
	Title            string         `json:"title" db:"title"`
	Description      string         `json:"description" db:"description"`
	RequiredSkills   pq.StringArray `json:"required_skills" db:"required_skills"`
	JobType          JobType        `json:"job_type" db:"job_type"`
	PaymentType      PaymentType    `json:"payment_type" db:"payment_type"`
	PaymentAmount    float64        `json:"payment_amount" db:"payment_amount"`
	Currency         string         `json:"currency" db:"currency"`
	City             string         `json:"city" db:"city"`
	State            string         `json:"state" db:"state"`
	Country          string         `json:"country" db:"country"`
	Latitude         float64        `json:"latitude" db:"latitude"`
	Longitude        float64        `json:"longitude" db:"longitude"`
	StartDate        time.Time      `json:"start_date" db:"start_date"`
	EndDate          *time.Time     `json:"end_date,omitempty" db:"end_date"`
	ApplicationDeadline *time.Time  `json:"application_deadline,omitempty" db:"application_deadline"`
	MaxApplicants    int            `json:"max_applicants" db:"max_applicants"`
	Status           JobStatus      `json:"status" db:"status"`
	ViewsCount       int            `json:"views_count" db:"views_count"`
	ApplicationCount int            `json:"application_count" db:"application_count"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`

	// Joined
	Business *Business `json:"business,omitempty" db:"-"`
}