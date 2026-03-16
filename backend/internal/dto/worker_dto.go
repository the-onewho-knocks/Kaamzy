package dto

type CreateWorkerProfileRequest struct {
	FirstName       string   `json:"first_name" validate:"required,min=2,max=50"`
	LastName        string   `json:"last_name" validate:"required,min=2,max=50"`
	Bio             string   `json:"bio" validate:"omitempty,max=1000"`
	Skills          []string `json:"skills" validate:"required,min=1,dive,min=1,max=50"`
	HourlyRate      float64  `json:"hourly_rate" validate:"required,gt=0"`
	Currency        string   `json:"currency" validate:"required,len=3"`
	City            string   `json:"city" validate:"required"`
	State           string   `json:"state" validate:"required"`
	Country         string   `json:"country" validate:"required"`
	Latitude        float64  `json:"latitude" validate:"required,latitude"`
	Longitude       float64  `json:"longitude" validate:"required,longitude"`
	ExperienceYears int      `json:"experience_years" validate:"gte=0,lte=60"`
}

type UpdateWorkerProfileRequest struct {
	FirstName          *string  `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName           *string  `json:"last_name" validate:"omitempty,min=2,max=50"`
	Bio                *string  `json:"bio" validate:"omitempty,max=1000"`
	Skills             []string `json:"skills" validate:"omitempty,min=1,dive,min=1,max=50"`
	HourlyRate         *float64 `json:"hourly_rate" validate:"omitempty,gt=0"`
	Currency           *string  `json:"currency" validate:"omitempty,len=3"`
	City               *string  `json:"city" validate:"omitempty"`
	State              *string  `json:"state" validate:"omitempty"`
	Country            *string  `json:"country" validate:"omitempty"`
	Latitude           *float64 `json:"latitude" validate:"omitempty,latitude"`
	Longitude          *float64 `json:"longitude" validate:"omitempty,longitude"`
	ExperienceYears    *int     `json:"experience_years" validate:"omitempty,gte=0,lte=60"`
	AvailabilityStatus *string  `json:"availability_status" validate:"omitempty,oneof=available unavailable partly_available"`
}

type WorkerFilterRequest struct {
	Skills          []string `json:"skills" form:"skills"`
	MinHourlyRate   float64  `json:"min_hourly_rate" form:"min_hourly_rate"`
	MaxHourlyRate   float64  `json:"max_hourly_rate" form:"max_hourly_rate"`
	City            string   `json:"city" form:"city"`
	Country         string   `json:"country" form:"country"`
	Availability    string   `json:"availability" form:"availability"`
	MinRating       float64  `json:"min_rating" form:"min_rating"`
	RadiusKm        float64  `json:"radius_km" form:"radius_km"`
	Latitude        float64  `json:"latitude" form:"latitude"`
	Longitude       float64  `json:"longitude" form:"longitude"`
	Page            int      `json:"page" form:"page" validate:"gte=1"`
	PageSize        int      `json:"page_size" form:"page_size" validate:"gte=1,lte=100"`
}

type WorkerResponse struct {
	ID                 string   `json:"id"`
	FirstName          string   `json:"first_name"`
	LastName           string   `json:"last_name"`
	Bio                string   `json:"bio,omitempty"`
	Skills             []string `json:"skills"`
	HourlyRate         float64  `json:"hourly_rate"`
	Currency           string   `json:"currency"`
	City               string   `json:"city"`
	State              string   `json:"state"`
	Country            string   `json:"country"`
	AvailabilityStatus string   `json:"availability_status"`
	ExperienceYears    int      `json:"experience_years"`
	AverageRating      float64  `json:"average_rating"`
	TotalReviews       int      `json:"total_reviews"`
	TotalJobsDone      int      `json:"total_jobs_done"`
	IsVerified         bool     `json:"is_verified"`
	ProfilePictureURL  string   `json:"profile_picture_url,omitempty"`
}