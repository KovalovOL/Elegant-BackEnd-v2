package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
    Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(ctx context.Context, user User) (uuid.UUID, error) {
	query := `
	INSERT INTO users (id, name, email, git_url, likedin_url, bio)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	var userUUID uuid.UUID
	err := r.db.QueryRow(ctx, query,
		user.ID, user.Name, user.Email, user.GitUrl, user.LikedinUrl, user.Bio).Scan(&userUUID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create user: %w", err)
	}
	return userUUID, nil
}

func (r *Repository) GetByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	query := `
	SELECT * FROM users 
	WHERE id = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.GitUrl,
		&user.LikedinUrl,
		&user.Bio,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
	SELECT * FROM users
	WHERE users.email = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.GitUrl,
		&user.LikedinUrl,
		&user.Bio,
	);
	if err != nil {
		return nil, err
	}
	return &user, nil
}