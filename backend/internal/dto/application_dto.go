package dto

type CreateApplicationRequest struct {
	JobID        string  `json:"job_id" validate:"required,uuid"`
	CoverLetter  string  `json:"cover_letter" validate:"omitempty,max=2000"`
	ProposedRate float64 `json:"proposed_rate" validate:"omitempty,gt=0"`
}

type UpdateApplicationStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=reviewed shortlisted accepted rejected"`
}

type ApplicationResponse struct {
	ID           string          `json:"id"`
	JobID        string          `json:"job_id"`
	WorkerID     string          `json:"worker_id"`
	CoverLetter  string          `json:"cover_letter,omitempty"`
	ProposedRate float64         `json:"proposed_rate,omitempty"`
	Status       string          `json:"status"`
	CreatedAt    string          `json:"created_at"`
	Job          *JobResponse    `json:"job,omitempty"`
	Worker       *WorkerResponse `json:"worker,omitempty"`
}

type ApplicationFilterRequest struct {
	JobID    string `form:"job_id" validate:"omitempty,uuid"`
	WorkerID string `form:"worker_id" validate:"omitempty,uuid"`
	Status   string `form:"status" validate:"omitempty,oneof=pending reviewed shortlisted accepted rejected withdrawn"`
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=100"`
}