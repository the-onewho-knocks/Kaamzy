package dto

type CreateBusinessProfileRequest struct {
	BusinessName string  `json:"business_name" validate:"required,min=2,max=100"`
	Description  string  `json:"description" validate:"omitempty,max=2000"`
	Industry     string  `json:"industry" validate:"required,max=100"`
	BusinessSize string  `json:"business_size" validate:"required,oneof=solo small medium large"`
	WebsiteURL   string  `json:"website_url" validate:"omitempty,url"`
	City         string  `json:"city" validate:"required"`
	State        string  `json:"state" validate:"required"`
	Country      string  `json:"country" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required,latitude"`
	Longitude    float64 `json:"longitude" validate:"required,longitude"`
}

type UpdateBusinessProfileRequest struct {
	BusinessName *string  `json:"business_name" validate:"omitempty,min=2,max=100"`
	Description  *string  `json:"description" validate:"omitempty,max=2000"`
	Industry     *string  `json:"industry" validate:"omitempty,max=100"`
	BusinessSize *string  `json:"business_size" validate:"omitempty,oneof=solo small medium large"`
	WebsiteURL   *string  `json:"website_url" validate:"omitempty,url"`
	City         *string  `json:"city" validate:"omitempty"`
	State        *string  `json:"state" validate:"omitempty"`
	Country      *string  `json:"country" validate:"omitempty"`
	Latitude     *float64 `json:"latitude" validate:"omitempty,latitude"`
	Longitude    *float64 `json:"longitude" validate:"omitempty,longitude"`
}

type BusinessResponse struct {
	ID              string  `json:"id"`
	BusinessName    string  `json:"business_name"`
	Description     string  `json:"description,omitempty"`
	Industry        string  `json:"industry"`
	BusinessSize    string  `json:"business_size"`
	WebsiteURL      string  `json:"website_url,omitempty"`
	LogoURL         string  `json:"logo_url,omitempty"`
	City            string  `json:"city"`
	State           string  `json:"state"`
	Country         string  `json:"country"`
	IsVerified      bool    `json:"is_verified"`
	AverageRating   float64 `json:"average_rating"`
	TotalReviews    int     `json:"total_reviews"`
	TotalJobsPosted int     `json:"total_jobs_posted"`
}