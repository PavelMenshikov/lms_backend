package usecase

import (
	"context"
	"errors"

	"lms_backend/internal/auth/repository"
	"lms_backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo repository.AuthRepository
}

func NewAuthUsecase(repo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(b), err
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (u *AuthUsecase) Register(ctx context.Context, user *domain.User, password string) error {
	return errors.New("Public registration is disabled.")
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := checkPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
