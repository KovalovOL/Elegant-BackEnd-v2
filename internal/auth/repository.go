package auth

import (
	"context"
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

func (r *Repository) CreateRefToken(ctx context.Context, token CreateRefreshToken) (int, error) {
	query := `
	INSERT INTO refresh_tokens (user_id, refresh_token_hash, user_agent, ip, expire_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING session_id
	`
	var id int
	err := r.db.QueryRow(ctx, query,
		&token.UserID,
		&token.RefreshTokenHash,
		&token.UserAgent,
		&token.IP,
		&token.ExpireAt,
		&token.CreatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) DeleteRefToken(ctx context.Context, session_id int) error {
	query := `
	DELETE FROM refresh_tokens
	WHERE session_id = $1
	`
	_, err := r.db.Exec(ctx, query, session_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateAuthProvider(ctx context.Context, provider Provider) (string, error) {
	query := `
	INSERT INTO auth_providers (provider, provider_user_id, user_id, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING provider_user_id
	`

	var provUserID string
	err := r.db.QueryRow(ctx, query,
		&provider.Provider,
		&provider.ProviderUserID,
		&provider.UserID,
		&provider.CreatedAt,
	).Scan(&provUserID)
	if err != nil {
		return "", err
	}
	return provUserID, nil
}

func (r *Repository) DeleteAuthProvider(ctx context.Context, provider string, providerUserID string) error {
	query := `
	DELETE FROM provider_user_id
	WHERE provider = $1 and provider_user_id = $2
	`
	_, err := r.db.Exec(ctx, query, provider, providerUserID)
	if err != nil {
		return err
	}
	return nil
}