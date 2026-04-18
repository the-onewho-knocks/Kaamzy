package postgres

import (
	"context"
	"errors"
	"fmt"
	"kaamzy/internal/dto"
	"kaamzy/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepository struct {
	db *pgxpool.Pool
}
//hiii

func NewApplicationRepository(db *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(ctx context.Context, app *models.Application) error {
	query := `
		INSERT INTO applications (
			id, job_id, worker_id, cover_letter, proposed_rate,
			status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	app.ID = uuid.New()
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		app.ID, app.JobID, app.WorkerID, app.CoverLetter,
		app.ProposedRate, app.Status, app.CreatedAt, app.UpdatedAt,
	)
	return err
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	query := `SELECT * FROM applications WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanApplication(row)
}

func (r *ApplicationRepository) GetByJobAndWorker(ctx context.Context, jobID, workerID uuid.UUID) (*models.Application, error) {
	query := `SELECT * FROM applications WHERE job_id = $1 AND worker_id = $2`
	row := r.db.QueryRow(ctx, query, jobID, workerID)
	return scanApplication(row)
}

func (r *ApplicationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ApplicationStatus) error {
	now := time.Now()
	query := `UPDATE applications SET status = $1, responded_at = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, status, now, now, id)
	return err
}

func (r *ApplicationRepository) MarkAsSeen(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE applications SET seen_at = $1, updated_at = $2 WHERE id = $3 AND seen_at IS NULL`
	_, err := r.db.Exec(ctx, query, now, now, id)
	return err
}

func (r *ApplicationRepository) List(ctx context.Context, filter dto.ApplicationFilterRequest) ([]*models.Application, int, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	i := 1

	if filter.JobID != "" {
		conditions = append(conditions, fmt.Sprintf("job_id = $%d", i))
		args = append(args, filter.JobID)
		i++
	}
	if filter.WorkerID != "" {
		conditions = append(conditions, fmt.Sprintf("worker_id = $%d", i))
		args = append(args, filter.WorkerID)
		i++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", i))
		args = append(args, filter.Status)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM applications %s", where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	listQuery := fmt.Sprintf(
		"SELECT * FROM applications %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var apps []*models.Application
	for rows.Next() {
		a, err := scanApplicationRows(rows)
		if err != nil {
			return nil, 0, err
		}
		apps = append(apps, a)
	}
	return apps, total, rows.Err()
}

func scanApplication(row pgx.Row) (*models.Application, error) {
	a := &models.Application{}
	err := row.Scan(
		&a.ID, &a.JobID, &a.WorkerID, &a.CoverLetter,
		&a.ProposedRate, &a.Status, &a.SeenAt, &a.RespondedAt,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return a, err
}

func scanApplicationRows(rows pgx.Rows) (*models.Application, error) {
	a := &models.Application{}
	err := rows.Scan(
		&a.ID, &a.JobID, &a.WorkerID, &a.CoverLetter,
		&a.ProposedRate, &a.Status, &a.SeenAt, &a.RespondedAt,
		&a.CreatedAt, &a.UpdatedAt,
	)
	return a, err
}
