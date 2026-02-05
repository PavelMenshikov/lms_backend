package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"lms_backend/internal/domain"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

var _ AuthRepository = (*AuthRepositoryImpl)(nil)

func NewAuthRepository(db *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}

	query := `
		SELECT 
			id, first_name, last_name, email, password_hash, role, created_at
		FROM users
		WHERE email = $1;
	`
	var passwordHash string

	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &passwordHash, &u.Role, &u.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	u.Password = passwordHash
	return u, nil
}
