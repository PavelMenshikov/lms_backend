package mocks

import (
	"context"
	"database/sql"
	"errors"

	"lms_backend/internal/auth/repository"
	"lms_backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryMock struct {
	Users map[string]*domain.User
}

func hashPasswordForMock(password string) string {

	b, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NewAuthRepositoryMock() *AuthRepositoryMock {

	testPasswordHash := hashPasswordForMock("password")

	return &AuthRepositoryMock{
		Users: map[string]*domain.User{
			"test@lms.ru": {
				ID:        "00000000-0000-0000-0000-000000000001",
				FirstName: "Тест",
				LastName:  "Юзер",
				Email:     "test@lms.ru",
				Password:  testPasswordHash,
				Role:      domain.RoleStudent,
			},
		},
	}
}

var _ repository.AuthRepository = (*AuthRepositoryMock)(nil)

func (r *AuthRepositoryMock) CreateUser(ctx context.Context, u *domain.User) error {
	if _, ok := r.Users[u.Email]; ok {
		return errors.New("пользователь уже существует")
	}
	r.Users[u.Email] = u
	return nil
}

func (r *AuthRepositoryMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if user, ok := r.Users[email]; ok {
		copiedUser := *user
		return &copiedUser, nil
	}
	return nil, sql.ErrNoRows
}
