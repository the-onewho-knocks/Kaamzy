package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AvailabilityStatus string

const (
	AvailabilityAvailable    AvailabilityStatus = "available"
	AvailabilityUnavailable  AvailabilityStatus = "unavailable"
	AvailabilityPartly       AvailabilityStatus = "partly_available"
)

type Worker struct {
	ID                uuid.UUID          `json:"id" db:"id"`
	UserID            uuid.UUID          `json:"user_id" db:"user_id"`
	FirstName         string             `json:"first_name" db:"first_name"`
	LastName          string             `json:"last_name" db:"last_name"`
	Bio               string             `json:"bio,omitempty" db:"bio"`
	Skills            pq.StringArray     `json:"skills" db:"skills"`
	HourlyRate        float64            `json:"hourly_rate" db:"hourly_rate"`
	Currency          string             `json:"currency" db:"currency"`
	City              string             `json:"city" db:"city"`
	State             string             `json:"state" db:"state"`
	Country           string             `json:"country" db:"country"`
	Latitude          float64            `json:"latitude" db:"latitude"`
	Longitude         float64            `json:"longitude" db:"longitude"`
	AvailabilityStatus AvailabilityStatus `json:"availability_status" db:"availability_status"`
	ExperienceYears   int                `json:"experience_years" db:"experience_years"`
	AverageRating     float64            `json:"average_rating" db:"average_rating"`
	TotalReviews      int                `json:"total_reviews" db:"total_reviews"`
	TotalJobsDone     int                `json:"total_jobs_done" db:"total_jobs_done"`
	IsVerified        bool               `json:"is_verified" db:"is_verified"`
	ResumeURL         string             `json:"resume_url,omitempty" db:"resume_url"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" db:"updated_at"`

	// Joined
	User *User `json:"user,omitempty" db:"-"`
}