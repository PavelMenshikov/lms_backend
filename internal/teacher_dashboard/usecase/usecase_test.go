package usecase_test

import (
	"context"
	"errors"
	"testing"

	"lms_backend/internal/teacher_dashboard/mocks"
	"lms_backend/internal/teacher_dashboard/usecase"
)

func TestTeacherDashboardUseCase_GetMonthlyReport(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		repoMock := mocks.NewTeacherDashboardRepositoryMock()
		uc := usecase.NewTeacherDashboardUseCase(repoMock)

		report, err := uc.GetMonthlyReport(ctx, "teacher-1", 2026, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if report.TeacherID != "teacher-1" {
			t.Errorf("expected teacher-1, got %s", report.TeacherID)
		}
		if report.Year != 2026 || report.Month != 6 {
			t.Errorf("expected 2026/6, got %d/%d", report.Year, report.Month)
		}
		if report.TotalLessons != 10 {
			t.Errorf("expected 10 lessons, got %d", report.TotalLessons)
		}
	})

	t.Run("ErrorFromRepo", func(t *testing.T) {
		repoMock := mocks.NewTeacherDashboardRepositoryMock()
		repoMock.Err = errors.New("db error")
		uc := usecase.NewTeacherDashboardUseCase(repoMock)

		_, err := uc.GetMonthlyReport(ctx, "teacher-1", 2026, 6)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
