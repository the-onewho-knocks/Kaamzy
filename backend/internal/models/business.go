package models

import (
	"time"

	"github.com/google/uuid"
)

type BusinessSize string

const (
	BusinessSizeSolo   BusinessSize = "solo"
	BusinessSizeSmall  BusinessSize = "small"
	BusinessSizeMedium BusinessSize = "medium"
	BusinessSizeLarge  BusinessSize = "large"
)

type Business struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	UserID           uuid.UUID    `json:"user_id" db:"user_id"`
	BusinessName     string       `json:"business_name" db:"business_name"`
	Description      string       `json:"description,omitempty" db:"description"`
	Industry         string       `json:"industry" db:"industry"`
	BusinessSize     BusinessSize `json:"business_size" db:"business_size"`
	WebsiteURL       string       `json:"website_url,omitempty" db:"website_url"`
	LogoURL          string       `json:"logo_url,omitempty" db:"logo_url"`
	City             string       `json:"city" db:"city"`
	State            string       `json:"state" db:"state"`
	Country          string       `json:"country" db:"country"`
	Latitude         float64      `json:"latitude" db:"latitude"`
	Longitude        float64      `json:"longitude" db:"longitude"`
	IsVerified       bool         `json:"is_verified" db:"is_verified"`
	AverageRating    float64      `json:"average_rating" db:"average_rating"`
	TotalReviews     int          `json:"total_reviews" db:"total_reviews"`
	TotalJobsPosted  int          `json:"total_jobs_posted" db:"total_jobs_posted"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`

	// Joined
	User *User `json:"user,omitempty" db:"-"`
}