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

type RatingRepository struct {
	db *pgxpool.Pool
}

func NewRatingRepository(db *pgxpool.Pool) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) Create(ctx context.Context,
	rating *models.Rating) error {

	query := `
		INSERT INTO ratings (
			id, rating_type, reviewer_id, reviewee_id, job_id,
			score, title, comment, is_public, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	rating.ID = uuid.New()
	rating.CreatedAt = time.Now()
	rating.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		rating.ID, rating.RatingType, rating.ReviewerID, rating.RevieweeID,
		rating.JobID, rating.Score, rating.Title, rating.Comment,
		rating.IsPublic, rating.CreatedAt, rating.UpdatedAt,
	)
	return err
}

func (r *RatingRepository) GetByID(ctx context.Context,
	id uuid.UUID) (*models.Rating, error) {
	query := `SELECT * FROM ratings WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanRating(row)

}

func (r *RatingRepository) ExistsByJobAndReviewer(ctx context.Context, jobID, reviewerID uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM ratings WHERE job_id = $1 AND reviewer_id = $2`
	err := r.db.QueryRow(ctx, query, jobID, reviewerID).Scan(&count)
	return count > 0, err
}

func (r *RatingRepository) GetAverageForReviewee(ctx context.Context, revieweeID uuid.UUID) (float64, int, error) {
	var avg float64
	var total int
	query := `
		SELECT COALESCE(AVG(score), 0), COUNT(*)
		FROM ratings
		WHERE reviewee_id = $1 AND is_public = true
	`
	err := r.db.QueryRow(ctx, query, revieweeID).Scan(&avg, &total)
	return avg, total, err
}

func (r *RatingRepository) List(ctx context.Context, filter dto.RatingFilterRequest) ([]*models.Rating, int, error) {
	conditions := []string{"is_public = true"}
	args := []interface{}{}
	i := 1

	if filter.RevieweeID != "" {
		conditions = append(conditions, fmt.Sprintf("reviewee_id = $%d", i))
		args = append(args, filter.RevieweeID)
		i++
	}
	if filter.RatingType != "" {
		conditions = append(conditions, fmt.Sprintf("rating_type = $%d", i))
		args = append(args, filter.RatingType)
		i++
	}
	if filter.MinScore > 0 {
		conditions = append(conditions, fmt.Sprintf("score >= $%d", i))
		args = append(args, filter.MinScore)
		i++
	}
	if filter.MaxScore > 0 {
		conditions = append(conditions, fmt.Sprintf("score <= $%d", i))
		args = append(args, filter.MaxScore)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ratings %s", where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	listQuery := fmt.Sprintf(
		"SELECT * FROM ratings %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ratings []*models.Rating
	for rows.Next() {
		rat := &models.Rating{}
		if err := rows.Scan(
			&rat.ID, &rat.RatingType, &rat.ReviewerID, &rat.RevieweeID,
			&rat.JobID, &rat.Score, &rat.Title, &rat.Comment,
			&rat.IsPublic, &rat.CreatedAt, &rat.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		ratings = append(ratings, rat)
	}
	return ratings, total, rows.Err()
}

func scanRating(row pgx.Row) (*models.Rating, error) {
	r := &models.Rating{}
	err := row.Scan(
		&r.ID, &r.RatingType, &r.ReviewerID, &r.RevieweeID,
		&r.JobID, &r.Score, &r.Title, &r.Comment,
		&r.IsPublic, &r.CreatedAt, &r.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return r, err
}
