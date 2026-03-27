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

type WorkerRepository struct {
	db *pgxpool.Pool
}

func NewWorkerRepository(db *pgxpool.Pool) *WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) Create(ctx context.Context,
	worker *models.Worker) error {
	query := `
		insert into workers(
			id, user_id, first_name, last_name, bio, skills, hourly_rate, currency,
			city, state, country, latitude, longitude, availability_status,
			experience_years, average_rating, total_reviews, total_jobs_done,
			is_verified, resume_url, created_at, updated_at
		) Values(
		 $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22
		)
	`

	worker.ID = uuid.New()
	worker.CreatedAt = time.Now()
	worker.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		worker.ID, worker.UserID, worker.FirstName, worker.LastName,
		worker.Bio, worker.Skills, worker.HourlyRate, worker.Currency,
		worker.City, worker.State, worker.Country,
		worker.Latitude, worker.Longitude,
		worker.AvailabilityStatus, worker.ExperienceYears,
		worker.AverageRating, worker.TotalReviews, worker.TotalJobsDone,
		worker.IsVerified, worker.ResumeURL,
		worker.CreatedAt, worker.UpdatedAt,
	)
	return err

}

func (r *WorkerRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Worker, error) {
	query := `SELECT * FROM workers WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanWorker(row)
}

func (r *WorkerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Worker, error) {
	query := `SELECT * FROM workers WHERE user_id = $1`
	row := r.db.QueryRow(ctx, query, userID)
	return scanWorker(row)
}

func (r *WorkerRepository) Update(ctx context.Context, worker *models.Worker) error {
	query := `
		UPDATE workers SET
			first_name = $1, last_name = $2, bio = $3, skills = $4,
			hourly_rate = $5, currency = $6, city = $7, state = $8,
			country = $9, latitude = $10, longitude = $11,
			availability_status = $12, experience_years = $13,
			resume_url = $14, updated_at = $15
		WHERE id = $16
	`
	worker.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		worker.FirstName, worker.LastName, worker.Bio, worker.Skills,
		worker.HourlyRate, worker.Currency, worker.City, worker.State,
		worker.Country, worker.Latitude, worker.Longitude,
		worker.AvailabilityStatus, worker.ExperienceYears,
		worker.ResumeURL, worker.UpdatedAt, worker.ID,
	)
	return err
}

func (r *WorkerRepository) UpdateRating(ctx context.Context, id uuid.UUID, avg float64, total int) error {
	query := `UPDATE workers SET average_rating = $1, total_reviews = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, avg, total, time.Now(), id)
	return err
}

func (r *WorkerRepository) IncrementJobsDone(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE workers SET total_jobs_done = total_jobs_done + 1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *WorkerRepository) List(ctx context.Context, filter dto.WorkerFilterRequest) ([]*models.Worker, int, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	i := 1

	if len(filter.Skills) > 0 {
		conditions = append(conditions, fmt.Sprintf("skills && $%d", i))
		args = append(args, filter.Skills)
		i++
	}
	if filter.MinHourlyRate > 0 {
		conditions = append(conditions, fmt.Sprintf("hourly_rate >= $%d", i))
		args = append(args, filter.MinHourlyRate)
		i++
	}
	if filter.MaxHourlyRate > 0 {
		conditions = append(conditions, fmt.Sprintf("hourly_rate <= $%d", i))
		args = append(args, filter.MaxHourlyRate)
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
	if filter.Availability != "" {
		conditions = append(conditions, fmt.Sprintf("availability_status = $%d", i))
		args = append(args, filter.Availability)
		i++
	}
	if filter.MinRating > 0 {
		conditions = append(conditions, fmt.Sprintf("average_rating >= $%d", i))
		args = append(args, filter.MinRating)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM workers %s", where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	listQuery := fmt.Sprintf(
		"SELECT * FROM workers %s ORDER BY average_rating DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var workers []*models.Worker
	for rows.Next() {
		w, err := scanWorkerRows(rows)
		if err != nil {
			return nil, 0, err
		}
		workers = append(workers, w)
	}
	return workers, total, rows.Err()
}

func scanWorker(row pgx.Row) (*models.Worker, error) {
	w := &models.Worker{}
	err := row.Scan(
		&w.ID, &w.UserID, &w.FirstName, &w.LastName, &w.Bio, &w.Skills,
		&w.HourlyRate, &w.Currency, &w.City, &w.State, &w.Country,
		&w.Latitude, &w.Longitude, &w.AvailabilityStatus, &w.ExperienceYears,
		&w.AverageRating, &w.TotalReviews, &w.TotalJobsDone,
		&w.IsVerified, &w.ResumeURL, &w.CreatedAt, &w.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return w, err
}

func scanWorkerRows(rows pgx.Rows) (*models.Worker, error) {
	w := &models.Worker{}
	err := rows.Scan(
		&w.ID, &w.UserID, &w.FirstName, &w.LastName, &w.Bio, &w.Skills,
		&w.HourlyRate, &w.Currency, &w.City, &w.State, &w.Country,
		&w.Latitude, &w.Longitude, &w.AvailabilityStatus, &w.ExperienceYears,
		&w.AverageRating, &w.TotalReviews, &w.TotalJobsDone,
		&w.IsVerified, &w.ResumeURL, &w.CreatedAt, &w.UpdatedAt,
	)
	return w, err
}
