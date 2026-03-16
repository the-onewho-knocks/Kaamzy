package dto

type SendMessageRequest struct {
	RecipientID   string `json:"recipient_id" validate:"required,uuid"`
	Content       string `json:"content" validate:"required,min=1,max=2000"`
	AttachmentURL string `json:"attachment_url" validate:"omitempty,url"`
	JobID         string `json:"job_id" validate:"omitempty,uuid"`
}

type ConversationResponse struct {
	ID            string           `json:"id"`
	ParticipantA  string           `json:"participant_a"`
	ParticipantB  string           `json:"participant_b"`
	JobID         string           `json:"job_id,omitempty"`
	LastMessage   *MessageResponse `json:"last_message,omitempty"`
	LastMessageAt string           `json:"last_message_at,omitempty"`
	UnreadCount   int              `json:"unread_count"`
}

type MessageResponse struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	AttachmentURL  string `json:"attachment_url,omitempty"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}

type GetMessagesRequest struct {
	ConversationID string `form:"conversation_id" validate:"required,uuid"`
	Page           int    `form:"page" validate:"gte=1"`
	PageSize       int    `form:"page_size" validate:"gte=1,lte=100"`
}