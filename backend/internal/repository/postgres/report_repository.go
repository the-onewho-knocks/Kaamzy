package postgres

import (
	"context"
	"errors"
	"fmt"
	"kaamzy/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(ctx context.Context, report *models.Report) error {
	query := `
		INSERT INTO reports (
			id, reporter_id, target_type, target_id, reason,
			description, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	report.ID = uuid.New()
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		report.ID, report.ReporterID, report.TargetType, report.TargetID,
		report.Reason, report.Description, report.Status,
		report.CreatedAt, report.UpdatedAt,
	)
	return err
}

func (r *ReportRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Report, error) {
	query := `SELECT * FROM reports WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanReport(row)
}

func (r *ReportRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.ReportStatus, reviewerID uuid.UUID, note string) error {
	now := time.Now()
	query := `
		UPDATE reports SET
			status = $1, reviewed_by = $2, review_note = $3,
			reviewed_at = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, status, reviewerID, note, now, now, id)
	return err
}

func (r *ReportRepository) List(ctx context.Context, status models.ReportStatus, page, pageSize int) ([]*models.Report, int, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	i := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", i))
		args = append(args, status)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM reports %s", where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf(
		"SELECT * FROM reports %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reports []*models.Report
	for rows.Next() {
		rep, err := scanReportRows(rows)
		if err != nil {
			return nil, 0, err
		}
		reports = append(reports, rep)
	}
	return reports, total, rows.Err()
}

func (r *ReportRepository) ExistsByReporterAndTarget(ctx context.Context, reporterID, targetID uuid.UUID, targetType models.ReportTargetType) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM reports WHERE reporter_id = $1 AND target_id = $2 AND target_type = $3`
	err := r.db.QueryRow(ctx, query, reporterID, targetID, targetType).Scan(&count)
	return count > 0, err
}

func scanReport(row pgx.Row) (*models.Report, error) {
	r := &models.Report{}
	err := row.Scan(
		&r.ID, &r.ReporterID, &r.TargetType, &r.TargetID,
		&r.Reason, &r.Description, &r.Status,
		&r.ReviewedBy, &r.ReviewNote, &r.ReviewedAt,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return r, err
}

func scanReportRows(rows pgx.Rows) (*models.Report, error) {
	r := &models.Report{}
	err := rows.Scan(
		&r.ID, &r.ReporterID, &r.TargetType, &r.TargetID,
		&r.Reason, &r.Description, &r.Status,
		&r.ReviewedBy, &r.ReviewNote, &r.ReviewedAt,
		&r.CreatedAt, &r.UpdatedAt,
	)
	return r, err
}
