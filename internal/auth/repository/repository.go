package repository

import (
	"context"
	"lms_backend/internal/domain"
)

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
