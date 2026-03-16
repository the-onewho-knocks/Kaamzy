package dto

type CreateRatingRequest struct {
	JobID      string `json:"job_id" validate:"required,uuid"`
	RevieweeID string `json:"reviewee_id" validate:"required,uuid"`
	Score      int    `json:"score" validate:"required,min=1,max=5"`
	Title      string `json:"title" validate:"omitempty,max=100"`
	Comment    string `json:"comment" validate:"omitempty,max=1000"`
	IsPublic   bool   `json:"is_public"`
}

type RatingResponse struct {
	ID         string  `json:"id"`
	RatingType string  `json:"rating_type"`
	ReviewerID string  `json:"reviewer_id"`
	RevieweeID string  `json:"reviewee_id"`
	JobID      string  `json:"job_id"`
	Score      int     `json:"score"`
	Title      string  `json:"title,omitempty"`
	Comment    string  `json:"comment,omitempty"`
	IsPublic   bool    `json:"is_public"`
	CreatedAt  string  `json:"created_at"`
}

type RatingFilterRequest struct {
	RevieweeID string `form:"reviewee_id" validate:"omitempty,uuid"`
	RatingType string `form:"rating_type" validate:"omitempty,oneof=worker_to_job business_to_worker"`
	MinScore   int    `form:"min_score" validate:"omitempty,min=1,max=5"`
	MaxScore   int    `form:"max_score" validate:"omitempty,min=1,max=5"`
	Page       int    `form:"page" validate:"gte=1"`
	PageSize   int    `form:"page_size" validate:"gte=1,lte=100"`
}