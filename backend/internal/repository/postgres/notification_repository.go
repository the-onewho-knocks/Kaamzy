package postgres

import (
	"context"
	"kaamzy/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, body, payload, is_read, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	n.ID = uuid.New()
	n.CreatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		n.ID, n.UserID, n.Type, n.Title, n.Body, n.Payload, n.IsRead, n.CreatedAt,
	)
	return err
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*models.Notification, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := `
		SELECT * FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, listQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		n := &models.Notification{}
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body,
			&n.Payload, &n.IsRead, &n.ReadAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}
	return notifications, total, rows.Err()
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID uuid.UUID) error {
	now := time.Now()
	query := `UPDATE notifications SET is_read = true, read_at = $1 WHERE id = $2 AND user_id = $3`
	_, err := r.db.Exec(ctx, query, now, id, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	query := `UPDATE notifications SET is_read = true, read_at = $1 WHERE user_id = $2 AND is_read = false`
	_, err := r.db.Exec(ctx, query, now, userID)
	return err
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) DeleteOld(ctx context.Context, before time.Time) error {
	query := `DELETE FROM notifications WHERE created_at < $1 AND is_read = true`
	_, err := r.db.Exec(ctx, query, before)
	return err
}
