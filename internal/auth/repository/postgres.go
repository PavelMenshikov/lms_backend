package repository

import (
	"context"
	"database/sql"

	"lms_backend/internal/domain"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

var _ AuthRepository = (*AuthRepositoryImpl)(nil)

func NewAuthRepository(db *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.Role,
	).Scan(&u.ID, &u.CreatedAt)
}

func (r *AuthRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT id, first_name, last_name, email, password, role, created_at, password 
		FROM users
		WHERE email = $1;
	`

	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.Password)

	if err != nil {
		return nil, err
	}

	return u, nil
}
