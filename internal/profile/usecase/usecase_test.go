package usecase_test

import (
	"context"
	"errors"
	"testing"

	"lms_backend/internal/domain"
	"lms_backend/internal/profile/mocks"
	"lms_backend/internal/profile/usecase"
	pkgMocks "lms_backend/pkg/storage/mocks"
)

func TestGetMyProfile(t *testing.T) {
	repo := mocks.NewProfileRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewProfileUseCase(repo, s3)

	repo.GetProfileFunc = func(ctx context.Context, userID string) (*domain.User, error) {
		if userID == "" {
			return nil, errors.New("not found")
		}
		return &domain.User{
			ID: userID, FirstName: "John", LastName: "Doe",
			Email: "john@test.com", Role: domain.RoleStudent,
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		user, err := uc.GetMyProfile(context.Background(), "u1")
		if err != nil {
			t.Fatal(err)
		}
		if user.FirstName != "John" {
			t.Error("expected 'John'")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := uc.GetMyProfile(context.Background(), "")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestUpdateTeacherSchedule(t *testing.T) {
	repo := mocks.NewProfileRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewProfileUseCase(repo, s3)

	scheduleJSON := []byte(`{"monday":{"start":"09:00","end":"18:00"}}`)
	var savedSchedule []byte

	repo.UpdateTeacherScheduleFunc = func(ctx context.Context, userID string, scheduleJSON []byte) error {
		if userID == "fail" {
			return errors.New("save failed")
		}
		savedSchedule = scheduleJSON
		return nil
	}

	t.Run("success", func(t *testing.T) {
		savedSchedule = nil
		err := uc.UpdateTeacherSchedule(context.Background(), "t1", scheduleJSON)
		if err != nil {
			t.Fatal(err)
		}
		if savedSchedule == nil {
			t.Error("schedule not saved")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		err := uc.UpdateTeacherSchedule(context.Background(), "fail", scheduleJSON)
		if err == nil {
			t.Error("expected error")
		}
	})
}
