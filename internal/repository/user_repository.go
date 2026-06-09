package repository

import (
	"context"

	"github.com/ahmed/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

// UserRepository defines the database operations we can perform on the users table.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type postgresUserRepository struct {
	db *sqlx.DB
}

// NewUserRepository returns a new instance of a UserRepository using PostgreSQL.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

// Create inserts a new user record into the database and returns the generated fields (id, created_at, updated_at).
func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowxContext(ctx, query, user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

// GetByID retrieves a user by their UUID.
func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = $1`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address.
func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by their username.
func (r *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE username = $1`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user's details and sets the updated_at timestamp.
func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, password_hash = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at`

	err := r.db.QueryRowxContext(ctx, query, user.Username, user.Email, user.PasswordHash, user.ID).
		Scan(&user.UpdatedAt)
	return err
}

// Delete permanently removes a user from the database.
func (r *postgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users 
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
