package usecase_test

import (
	"context"
	"errors"
	"testing"

	"lms_backend/internal/domain"
	"lms_backend/internal/statistics/mocks"
	"lms_backend/internal/statistics/usecase"
)

func TestGetStudentStatistics(t *testing.T) {
	repo := mocks.NewStatisticsRepoMock()
	uc := usecase.NewStatisticsUseCase(repo)

	repo.GetByStudentFunc = func(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
		if studentID == "" {
			return nil, errors.New("not found")
		}
		return &domain.StudentStatistics{
			StudentID:       studentID,
			TotalLessons:    20,
			AttendedLessons: 15,
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		stats, err := uc.GetStudentStatistics(context.Background(), "s1")
		if err != nil {
			t.Fatal(err)
		}
		if stats.TotalLessons != 20 || stats.AttendedLessons != 15 {
			t.Error("wrong stats values")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := uc.GetStudentStatistics(context.Background(), "")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestRefreshStudentStatistics(t *testing.T) {
	repo := mocks.NewStatisticsRepoMock()
	uc := usecase.NewStatisticsUseCase(repo)

	repo.RecalculateStatisticsFunc = func(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
		if studentID == "fail" {
			return nil, errors.New("recalc failed")
		}
		return &domain.StudentStatistics{
			StudentID:       studentID,
			TotalLessons:    30,
			AttendedLessons: 25,
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		stats, err := uc.RefreshStudentStatistics(context.Background(), "s1")
		if err != nil {
			t.Fatal(err)
		}
		if stats.TotalLessons != 30 {
			t.Error("expected 30 total lessons")
		}
	})

	t.Run("recalc error", func(t *testing.T) {
		_, err := uc.RefreshStudentStatistics(context.Background(), "fail")
		if err == nil {
			t.Error("expected error")
		}
	})
}
