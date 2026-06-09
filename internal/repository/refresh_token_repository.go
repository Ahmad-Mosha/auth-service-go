package repository

import (
	"context"

	"github.com/ahmed/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

// RefreshTokenRepository defines the database operations for managing user sessions/refresh tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
}

type postgresRefreshTokenRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenRepository instantiates a new repository implementation.
func NewRefreshTokenRepository(db *sqlx.DB) RefreshTokenRepository {
	return &postgresRefreshTokenRepository{db: db}
}

// Create stores a new refresh token session in the database.
func (r *postgresRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRowxContext(ctx, query, token.UserID, token.Token, token.ExpiresAt).
		Scan(&token.ID, &token.CreatedAt)
	return err
}

// GetByToken retrieves the refresh token details to check expiration or revocation.
func (r *postgresRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked_at 
		FROM refresh_tokens 
		WHERE token = $1`

	var rt domain.RefreshToken
	err := r.db.GetContext(ctx, &rt, query, token)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// Revoke invalidates a refresh token by setting revoked_at to the current timestamp.
func (r *postgresRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `
		UPDATE refresh_tokens 
		SET revoked_at = CURRENT_TIMESTAMP 
		WHERE token = $1 AND revoked_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
