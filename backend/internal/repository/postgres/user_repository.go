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

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context,
	user *models.User) error {
	query := `
			insert into users(
				id , email , password_hash , role , status , phone , is_email_verified,
				is_phone_verified,created_at,updated_at
			)	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		`
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Role, user.Status,
		user.Phone, user.IsEmailVerified, user.IsPhoneVerified,
		user.CreatedAt, user.UpdatedAt,
	)

	return err
}

//special function for helping in scanning the rows
func scanUser(row pgx.Row) (*models.User, error) {
	u := &models.User{}
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Status,
		&u.Phone, &u.ProfilePictureURL, &u.IsEmailVerified, &u.IsPhoneVerified,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) GetByID(ctx context.Context,
	id uuid.UUID) (*models.User,error) {
	query := `select * from users where id=$1 and deleted_at is null`
	row := r.db.QueryRow(ctx, query, id)
	return scanUser(row)
}

func (r *UserRepository) GetByEmail(ctx context.Context , 
	email string)(*models.User , error){
		query := `select * from users where email=$1 and deleted_at is null`;
		row := r.db.QueryRow(ctx , query , email);
		return scanUser(row)
	}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			email = $1,
			phone = $2,
			profile_picture_url = $3,
			status = $4,
			is_email_verified = $5,
			is_phone_verified = $6,
			updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		user.Email, user.Phone, user.ProfilePictureURL,
		user.Status, user.IsEmailVerified, user.IsPhoneVerified,
		user.UpdatedAt, user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, passwordHash, time.Now(), id)
	return err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, now, now, id)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error {
	query := `UPDATE users SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *UserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `UPDATE users SET deleted_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, now, now, id)
	return err
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.db.QueryRow(ctx, query, email).Scan(&count)
	return count > 0, err
}