package dto

import "time"

type CreateJobRequest struct {
	Title               string    `json:"title" validate:"required,min=5,max=150"`
	Description         string    `json:"description" validate:"required,min=20,max=5000"`
	RequiredSkills      []string  `json:"required_skills" validate:"required,min=1,dive,min=1,max=50"`
	JobType             string    `json:"job_type" validate:"required,oneof=onsite remote hybrid"`
	PaymentType         string    `json:"payment_type" validate:"required,oneof=hourly fixed daily"`
	PaymentAmount       float64   `json:"payment_amount" validate:"required,gt=0"`
	Currency            string    `json:"currency" validate:"required,len=3"`
	City                string    `json:"city" validate:"required_if=JobType onsite,required_if=JobType hybrid"`
	State               string    `json:"state" validate:"required_if=JobType onsite,required_if=JobType hybrid"`
	Country             string    `json:"country" validate:"required"`
	Latitude            float64   `json:"latitude" validate:"omitempty,latitude"`
	Longitude           float64   `json:"longitude" validate:"omitempty,longitude"`
	StartDate           time.Time `json:"start_date" validate:"required"`
	EndDate             *time.Time `json:"end_date" validate:"omitempty,gtfield=StartDate"`
	ApplicationDeadline *time.Time `json:"application_deadline" validate:"omitempty"`
	MaxApplicants       int       `json:"max_applicants" validate:"omitempty,gte=1"`
}

type UpdateJobRequest struct {
	Title               *string    `json:"title" validate:"omitempty,min=5,max=150"`
	Description         *string    `json:"description" validate:"omitempty,min=20,max=5000"`
	RequiredSkills      []string   `json:"required_skills" validate:"omitempty,min=1,dive,min=1,max=50"`
	JobType             *string    `json:"job_type" validate:"omitempty,oneof=onsite remote hybrid"`
	PaymentType         *string    `json:"payment_type" validate:"omitempty,oneof=hourly fixed daily"`
	PaymentAmount       *float64   `json:"payment_amount" validate:"omitempty,gt=0"`
	Currency            *string    `json:"currency" validate:"omitempty,len=3"`
	City                *string    `json:"city" validate:"omitempty"`
	State               *string    `json:"state" validate:"omitempty"`
	Country             *string    `json:"country" validate:"omitempty"`
	Latitude            *float64   `json:"latitude" validate:"omitempty,latitude"`
	Longitude           *float64   `json:"longitude" validate:"omitempty,longitude"`
	StartDate           *time.Time `json:"start_date" validate:"omitempty"`
	EndDate             *time.Time `json:"end_date" validate:"omitempty"`
	ApplicationDeadline *time.Time `json:"application_deadline" validate:"omitempty"`
	MaxApplicants       *int       `json:"max_applicants" validate:"omitempty,gte=1"`
	Status              *string    `json:"status" validate:"omitempty,oneof=draft open paused cancelled"`
}

type JobFilterRequest struct {
	Skills        []string  `form:"skills"`
	JobType       string    `form:"job_type" validate:"omitempty,oneof=onsite remote hybrid"`
	PaymentType   string    `form:"payment_type" validate:"omitempty,oneof=hourly fixed daily"`
	MinPayment    float64   `form:"min_payment"`
	MaxPayment    float64   `form:"max_payment"`
	City          string    `form:"city"`
	Country       string    `form:"country"`
	RadiusKm      float64   `form:"radius_km"`
	Latitude      float64   `form:"latitude"`
	Longitude     float64   `form:"longitude"`
	Status        string    `form:"status"`
	Page          int       `form:"page" validate:"gte=1"`
	PageSize      int       `form:"page_size" validate:"gte=1,lte=100"`
}

type JobResponse struct {
	ID                  string     `json:"id"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	RequiredSkills      []string   `json:"required_skills"`
	JobType             string     `json:"job_type"`
	PaymentType         string     `json:"payment_type"`
	PaymentAmount       float64    `json:"payment_amount"`
	Currency            string     `json:"currency"`
	City                string     `json:"city"`
	State               string     `json:"state"`
	Country             string     `json:"country"`
	StartDate           time.Time  `json:"start_date"`
	EndDate             *time.Time `json:"end_date,omitempty"`
	ApplicationDeadline *time.Time `json:"application_deadline,omitempty"`
	MaxApplicants       int        `json:"max_applicants"`
	Status              string     `json:"status"`
	ViewsCount          int        `json:"views_count"`
	ApplicationCount    int        `json:"application_count"`
	Business            *BusinessResponse `json:"business,omitempty"`
}