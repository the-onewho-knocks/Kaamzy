package postgres

import (
	"context"
	"errors"
	"kaamzy/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateConversation(ctx context.Context, conv *models.Conversation) error {
	query := `
		INSERT INTO conversations (id, participant_a, participant_b, job_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`
	conv.ID = uuid.New()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		conv.ID, conv.ParticipantA, conv.ParticipantB,
		conv.JobID, conv.CreatedAt, conv.UpdatedAt,
	)
	return err
}

func (r *MessageRepository) GetConversation(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	query := `SELECT * FROM conversations WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanConversation(row)
}

func (r *MessageRepository) GetConversationByParticipants(ctx context.Context, userA, userB uuid.UUID) (*models.Conversation, error) {
	query := `
		SELECT * FROM conversations
		WHERE (participant_a = $1 AND participant_b = $2)
		   OR (participant_a = $2 AND participant_b = $1)
		LIMIT 1
	`
	row := r.db.QueryRow(ctx, query, userA, userB)
	return scanConversation(row)
}

func (r *MessageRepository) GetConversationsByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*models.Conversation, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM conversations WHERE participant_a = $1 OR participant_b = $1`
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := `
		SELECT * FROM conversations
		WHERE participant_a = $1 OR participant_b = $1
		ORDER BY last_message_at DESC NULLS LAST
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, listQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var convs []*models.Conversation
	for rows.Next() {
		c := &models.Conversation{}
		if err := rows.Scan(
			&c.ID, &c.ParticipantA, &c.ParticipantB, &c.JobID,
			&c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		convs = append(convs, c)
	}
	return convs, total, rows.Err()
}

func (r *MessageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	query := `
		INSERT INTO messages (
			id, conversation_id, sender_id, content, attachment_url,
			status, is_deleted, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		msg.ID, msg.ConversationID, msg.SenderID, msg.Content,
		msg.AttachmentURL, msg.Status, msg.IsDeleted,
		msg.CreatedAt, msg.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return r.updateConversationLastMessage(ctx, msg.ConversationID, msg.CreatedAt)
}

func (r *MessageRepository) GetMessagesByConversation(ctx context.Context, convID uuid.UUID, page, pageSize int) ([]*models.Message, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM messages WHERE conversation_id = $1 AND is_deleted = false`
	if err := r.db.QueryRow(ctx, countQuery, convID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := `
		SELECT * FROM messages
		WHERE conversation_id = $1 AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, listQuery, convID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		m := &models.Message{}
		if err := rows.Scan(
			&m.ID, &m.ConversationID, &m.SenderID, &m.Content,
			&m.AttachmentURL, &m.Status, &m.IsDeleted,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		messages = append(messages, m)
	}
	return messages, total, rows.Err()
}

func (r *MessageRepository) MarkMessagesAsRead(ctx context.Context, convID, userID uuid.UUID) error {
	query := `
		UPDATE messages SET status = 'read', updated_at = $1
		WHERE conversation_id = $2
		AND sender_id != $3
		AND status != 'read'
	`
	_, err := r.db.Exec(ctx, query, time.Now(), convID, userID)
	return err
}

func (r *MessageRepository) SoftDeleteMessage(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE messages SET is_deleted = true, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *MessageRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM messages m
		JOIN conversations c ON m.conversation_id = c.id
		WHERE (c.participant_a = $1 OR c.participant_b = $1)
		AND m.sender_id != $1
		AND m.status != 'read'
		AND m.is_deleted = false
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *MessageRepository) updateConversationLastMessage(ctx context.Context, convID uuid.UUID, t time.Time) error {
	query := `UPDATE conversations SET last_message_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, t, t, convID)
	return err
}

func scanConversation(row pgx.Row) (*models.Conversation, error) {
	c := &models.Conversation{}
	err := row.Scan(
		&c.ID, &c.ParticipantA, &c.ParticipantB, &c.JobID,
		&c.LastMessageAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return c, err
}
