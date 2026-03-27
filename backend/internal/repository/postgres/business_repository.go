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

type BusinessRepository struct {
	db *pgxpool.Pool
}

func NewBusinessRepository(db *pgxpool.Pool) *BusinessRepository {
	return &BusinessRepository{db: db}
}

func (r *BusinessRepository) Create(ctx context.Context, business *models.Business) error {
	query := `
		INSERT INTO businesses (
			id, user_id, business_name, description, industry, business_size,
			website_url, logo_url, city, state, country, latitude, longitude,
			is_verified, average_rating, total_reviews, total_jobs_posted,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19
		)
	`
	business.ID = uuid.New()
	business.CreatedAt = time.Now()
	business.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		business.ID, business.UserID, business.BusinessName, business.Description,
		business.Industry, business.BusinessSize, business.WebsiteURL, business.LogoURL,
		business.City, business.State, business.Country,
		business.Latitude, business.Longitude,
		business.IsVerified, business.AverageRating,
		business.TotalReviews, business.TotalJobsPosted,
		business.CreatedAt, business.UpdatedAt,
	)
	return err
}

func (r *BusinessRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Business, error) {
	query := `SELECT * FROM businesses WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanBusiness(row)
}

func (r *BusinessRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Business, error) {
	query := `SELECT * FROM businesses WHERE user_id = $1`
	row := r.db.QueryRow(ctx, query, userID)
	return scanBusiness(row)
}

func (r *BusinessRepository) Update(ctx context.Context, business *models.Business) error {
	query := `
		UPDATE businesses SET
			business_name = $1, description = $2, industry = $3,
			business_size = $4, website_url = $5, logo_url = $6,
			city = $7, state = $8, country = $9,
			latitude = $10, longitude = $11, updated_at = $12
		WHERE id = $13
	`
	business.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		business.BusinessName, business.Description, business.Industry,
		business.BusinessSize, business.WebsiteURL, business.LogoURL,
		business.City, business.State, business.Country,
		business.Latitude, business.Longitude,
		business.UpdatedAt, business.ID,
	)
	return err
}

func (r *BusinessRepository) UpdateRating(ctx context.Context, id uuid.UUID, avg float64, total int) error {
	query := `UPDATE businesses SET average_rating = $1, total_reviews = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, avg, total, time.Now(), id)
	return err
}

func (r *BusinessRepository) IncrementJobsPosted(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE businesses SET total_jobs_posted = total_jobs_posted + 1, updated_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func scanBusiness(row pgx.Row) (*models.Business, error) {
	b := &models.Business{}
	err := row.Scan(
		&b.ID, &b.UserID, &b.BusinessName, &b.Description,
		&b.Industry, &b.BusinessSize, &b.WebsiteURL, &b.LogoURL,
		&b.City, &b.State, &b.Country, &b.Latitude, &b.Longitude,
		&b.IsVerified, &b.AverageRating, &b.TotalReviews, &b.TotalJobsPosted,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return b, err
}
