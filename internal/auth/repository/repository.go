package repository

import (
	"context"
	"lms_backend/internal/domain"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, u *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
