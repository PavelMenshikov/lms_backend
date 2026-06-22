package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"lms_backend/internal/domain"
	"lms_backend/internal/schedule/mocks"
	"lms_backend/internal/schedule/usecase"
)

func TestGetWeeklySchedule(t *testing.T) {
	repo := mocks.NewScheduleRepoMock()
	uc := usecase.NewScheduleUseCase(repo)

	monday := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC) // Monday

	repo.GetStudentLessonsInRangeFunc = func(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
		if userID == "" {
			return nil, errors.New("unauthorized")
		}
		return []domain.ScheduleLesson{
			{ID: "l1", Title: "Go Basics", StartTime: monday.Add(10 * time.Hour)},
			{ID: "l2", Title: "API Design", StartTime: monday.AddDate(0, 0, 2).Add(14 * time.Hour)},
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		sched, err := uc.GetWeeklySchedule(context.Background(), "user-1", monday)
		if err != nil {
			t.Fatal(err)
		}
		if len(sched.Days) != 2 {
			t.Errorf("expected 2 days, got %d", len(sched.Days))
		}
		if sched.StartDate.Weekday() != time.Monday {
			t.Error("week should start on Monday")
		}
	})

	t.Run("repo error returns empty schedule", func(t *testing.T) {
		repo.GetStudentLessonsInRangeFunc = func(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
			return nil, errors.New("db error")
		}
		sched, err := uc.GetWeeklySchedule(context.Background(), "user-1", monday)
		if err != nil {
			t.Fatalf("expected graceful degradation, got error: %v", err)
		}
		if len(sched.Days) != 0 {
			t.Error("expected empty days on error")
		}
	})

	t.Run("empty schedule", func(t *testing.T) {
		repo.GetStudentLessonsInRangeFunc = func(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
			return []domain.ScheduleLesson{}, nil
		}
		sched, err := uc.GetWeeklySchedule(context.Background(), "user-1", monday)
		if err != nil {
			t.Fatal(err)
		}
		if len(sched.Days) != 0 {
			t.Error("expected empty days")
		}
	})
}

func TestGetMonthlySchedule(t *testing.T) {
	repo := mocks.NewScheduleRepoMock()
	uc := usecase.NewScheduleUseCase(repo)

	repo.GetStudentLessonsInRangeFunc = func(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
		if userID == "fail" {
			return nil, errors.New("db error")
		}
		return []domain.ScheduleLesson{
			{ID: "l1", Title: "Lesson 1", StartTime: time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC)},
			{ID: "l2", Title: "Lesson 2", StartTime: time.Date(2026, 6, 15, 14, 0, 0, 0, time.UTC)},
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		sched, err := uc.GetMonthlySchedule(context.Background(), "user-1", 2026, 6)
		if err != nil {
			t.Fatal(err)
		}
		if sched.Month != 6 || sched.Year != 2026 {
			t.Error("wrong month/year")
		}
		if len(sched.Days) != 2 {
			t.Errorf("expected 2 days, got %d", len(sched.Days))
		}
	})

	t.Run("repo error returns empty schedule", func(t *testing.T) {
		sched, err := uc.GetMonthlySchedule(context.Background(), "fail", 2026, 6)
		if err != nil {
			t.Fatalf("expected graceful degradation, got error: %v", err)
		}
		if len(sched.Days) != 0 {
			t.Error("expected empty days on error")
		}
	})

	t.Run("empty month", func(t *testing.T) {
		repo.GetStudentLessonsInRangeFunc = func(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
			return []domain.ScheduleLesson{}, nil
		}
		sched, err := uc.GetMonthlySchedule(context.Background(), "user-2", 2026, 7)
		if err != nil {
			t.Fatal(err)
		}
		if len(sched.Days) != 0 {
			t.Error("expected empty days")
		}
	})
}
