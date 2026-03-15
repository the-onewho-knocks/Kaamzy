package models

import (
	"time"

	"github.com/google/uuid"
)

type RatingType string

const (
	RatingTypeWorkerToJob      RatingType = "worker_to_job"
	RatingTypeBusinessToWorker RatingType = "business_to_worker"
)

type Rating struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	RatingType RatingType `json:"rating_type" db:"rating_type"`
	ReviewerID uuid.UUID  `json:"reviewer_id" db:"reviewer_id"`
	RevieweeID uuid.UUID  `json:"reviewee_id" db:"reviewee_id"`
	JobID      uuid.UUID  `json:"job_id" db:"job_id"`
	Score      int        `json:"score" db:"score"` // 1–5
	Title      string     `json:"title,omitempty" db:"title"`
	Comment    string     `json:"comment,omitempty" db:"comment"`
	IsPublic   bool       `json:"is_public" db:"is_public"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`

	// Joined
	Reviewer *User `json:"reviewer,omitempty" db:"-"`
	Job      *Job  `json:"job,omitempty" db:"-"`
}