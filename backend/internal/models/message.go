package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
)

type Conversation struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ParticipantA uuid.UUID  `json:"participant_a" db:"participant_a"`
	ParticipantB uuid.UUID  `json:"participant_b" db:"participant_b"`
	JobID        *uuid.UUID `json:"job_id,omitempty" db:"job_id"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty" db:"last_message_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	// Joined
	LastMessage *Message `json:"last_message,omitempty" db:"-"`
}

type Message struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	ConversationID uuid.UUID     `json:"conversation_id" db:"conversation_id"`
	SenderID       uuid.UUID     `json:"sender_id" db:"sender_id"`
	Content        string        `json:"content" db:"content"`
	AttachmentURL  string        `json:"attachment_url,omitempty" db:"attachment_url"`
	Status         MessageStatus `json:"status" db:"status"`
	IsDeleted      bool          `json:"is_deleted" db:"is_deleted"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`

	// Joined
	Sender *User `json:"sender,omitempty" db:"-"`
}