package usecase_test

import (
	"context"
	"lms_backend/internal/auth/mocks"
	"lms_backend/internal/auth/usecase"
	"lms_backend/internal/domain"
	"testing"
)

func TestAuthUsecase_Login(t *testing.T) {
	repoMock := mocks.NewAuthRepositoryMock()
	uc := usecase.NewAuthUsecase(repoMock)

	ctx := context.Background()

	t.Run("Success Login", func(t *testing.T) {
		user, err := uc.Login(ctx, "test@lms.ru", "password")
		if err != nil {
			t.Fatalf("Ожидалась успешная аутентификация, получена ошибка: %v", err)
		}
		if user.Email != "test@lms.ru" {
			t.Errorf("Ожидался user 'test@lms.ru', получен %s", user.Email)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, err := uc.Login(ctx, "test@lms.ru", "wrong_password")
		if err == nil {
			t.Error("Ожидалась ошибка 'invalid credentials', но ошибок нет")
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		_, err := uc.Login(ctx, "unknown@lms.ru", "password")
		if err == nil {
			t.Error("Ожидалась ошибка 'user not found', но ошибок нет")
		}
	})
}

func TestAuthUsecase_Register(t *testing.T) {
	repoMock := mocks.NewAuthRepositoryMock()
	uc := usecase.NewAuthUsecase(repoMock)

	ctx := context.Background()
	newPassword := "newpassword"
	newUser := &domain.User{
		FirstName: "New",
		LastName:  "Guy",
		Email:     "new@lms.ru",
		Role:      domain.RoleParent,
	}

	t.Run("Success Register", func(t *testing.T) {
		err := uc.Register(ctx, newUser, newPassword)
		if err != nil {
			t.Fatalf("Ожидалась успешная регистрация, получена ошибка: %v", err)
		}

		if newUser.Password == newPassword {
			t.Error("Пароль должен быть захеширован, но он не захеширован.")
		}

		if _, ok := repoMock.Users["new@lms.ru"]; !ok {
			t.Error("Пользователь должен быть добавлен в Mock-репозиторий")
		}
	})

	t.Run("User Exists", func(t *testing.T) {
		err := uc.Register(ctx, newUser, newPassword)
		if err == nil {
			t.Error("Ожидалась ошибка 'пользователь уже существует', но ошибок нет")
		}
	})
}
