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

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (
			id, business_id, title, description, required_skills, job_type,
			payment_type, payment_amount, currency, city, state, country,
			latitude, longitude, start_date, end_date, application_deadline,
			max_applicants, status, views_count, application_count,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,
			$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23
		)
	`
	job.ID = uuid.New()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		job.ID, job.BusinessID, job.Title, job.Description,
		job.RequiredSkills, job.JobType, job.PaymentType, job.PaymentAmount,
		job.Currency, job.City, job.State, job.Country,
		job.Latitude, job.Longitude, job.StartDate, job.EndDate,
		job.ApplicationDeadline, job.MaxApplicants, job.Status,
		job.ViewsCount, job.ApplicationCount,
		job.CreatedAt, job.UpdatedAt,
	)
	return err
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	query := `SELECT * FROM jobs WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)
	return scanJob(row)
}

func (r *JobRepository) Update(ctx context.Context, job *models.Job) error {
	query := `
		UPDATE jobs SET
			title = $1, description = $2, required_skills = $3,
			job_type = $4, payment_type = $5, payment_amount = $6,
			currency = $7, city = $8, state = $9, country = $10,
			latitude = $11, longitude = $12, start_date = $13,
			end_date = $14, application_deadline = $15,
			max_applicants = $16, status = $17, updated_at = $18
		WHERE id = $19 AND deleted_at IS NULL
	`
	job.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		job.Title, job.Description, job.RequiredSkills,
		job.JobType, job.PaymentType, job.PaymentAmount,
		job.Currency, job.City, job.State, job.Country,
		job.Latitude, job.Longitude, job.StartDate,
		job.EndDate, job.ApplicationDeadline,
		job.MaxApplicants, job.Status, job.UpdatedAt, job.ID,
	)
	return err
}

func (r *JobRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE jobs SET deleted_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, now, now, id)
	return err
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.JobStatus) error {
	query := `UPDATE jobs SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *JobRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE jobs SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *JobRepository) IncrementApplicationCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE jobs SET application_count = application_count + 1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *JobRepository) GetByBusinessID(ctx context.Context, businessID uuid.UUID) ([]*models.Job, error) {
	query := `SELECT * FROM jobs WHERE business_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobRows(rows)
}

func (r *JobRepository) GetExpiredJobs(ctx context.Context) ([]*models.Job, error) {
	query := `
		SELECT * FROM jobs
		WHERE status = 'open'
		AND application_deadline < NOW()
		AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobRows(rows)
}

func (r *JobRepository) List(ctx context.Context, filter dto.JobFilterRequest) ([]*models.Job, int, error) {
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	i := 1

	if len(filter.Skills) > 0 {
		conditions = append(conditions, fmt.Sprintf("required_skills && $%d", i))
		args = append(args, filter.Skills)
		i++
	}
	if filter.JobType != "" {
		conditions = append(conditions, fmt.Sprintf("job_type = $%d", i))
		args = append(args, filter.JobType)
		i++
	}
	if filter.PaymentType != "" {
		conditions = append(conditions, fmt.Sprintf("payment_type = $%d", i))
		args = append(args, filter.PaymentType)
		i++
	}
	if filter.MinPayment > 0 {
		conditions = append(conditions, fmt.Sprintf("payment_amount >= $%d", i))
		args = append(args, filter.MinPayment)
		i++
	}
	if filter.MaxPayment > 0 {
		conditions = append(conditions, fmt.Sprintf("payment_amount <= $%d", i))
		args = append(args, filter.MaxPayment)
		i++
	}
	if filter.City != "" {
		conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", i))
		args = append(args, "%"+filter.City+"%")
		i++
	}
	if filter.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country = $%d", i))
		args = append(args, filter.Country)
		i++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", i))
		args = append(args, filter.Status)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs %s", where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	listQuery := fmt.Sprintf(
		"SELECT * FROM jobs %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	jobs, err := scanJobRows(rows)
	return jobs, total, err
}

func scanJob(row pgx.Row) (*models.Job, error) {
	j := &models.Job{}
	err := row.Scan(
		&j.ID, &j.BusinessID, &j.Title, &j.Description, &j.RequiredSkills,
		&j.JobType, &j.PaymentType, &j.PaymentAmount, &j.Currency,
		&j.City, &j.State, &j.Country, &j.Latitude, &j.Longitude,
		&j.StartDate, &j.EndDate, &j.ApplicationDeadline,
		&j.MaxApplicants, &j.Status, &j.ViewsCount, &j.ApplicationCount,
		&j.CreatedAt, &j.UpdatedAt, &j.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return j, err
}

func scanJobRows(rows pgx.Rows) ([]*models.Job, error) {
	var jobs []*models.Job
	for rows.Next() {
		j := &models.Job{}
		err := rows.Scan(
			&j.ID, &j.BusinessID, &j.Title, &j.Description, &j.RequiredSkills,
			&j.JobType, &j.PaymentType, &j.PaymentAmount, &j.Currency,
			&j.City, &j.State, &j.Country, &j.Latitude, &j.Longitude,
			&j.StartDate, &j.EndDate, &j.ApplicationDeadline,
			&j.MaxApplicants, &j.Status, &j.ViewsCount, &j.ApplicationCount,
			&j.CreatedAt, &j.UpdatedAt, &j.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}
