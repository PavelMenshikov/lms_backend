package usecase_test

import (
	"context"
	"testing"
	"time"

	"lms_backend/internal/attendance/mocks"
	"lms_backend/internal/attendance/usecase"
	"lms_backend/internal/domain"
)

func ptr(s string) *string { return &s }

func TestAttendanceUseCase_MarkAttendance(t *testing.T) {
	repoMock := mocks.NewAttendanceRepositoryMock()
	uc := usecase.NewAttendanceUseCase(repoMock)
	ctx := context.Background()

	t.Run("CreateNew", func(t *testing.T) {
		err := uc.MarkAttendance(ctx, "lesson-1", "student-1",
			domain.AttendanceStatusAttended, nil, nil, "teacher-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		records, err := uc.GetLessonAttendance(ctx, "lesson-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(records))
		}
		if records[0].Status != domain.AttendanceStatusAttended {
			t.Errorf("expected ATTENDED, got %s", records[0].Status)
		}
	})

	t.Run("UpdateExisting", func(t *testing.T) {
		reason := "sick"
		err := uc.MarkAttendance(ctx, "lesson-1", "student-1",
			domain.AttendanceStatusAbsentExcused, &reason, nil, "teacher-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		records, _ := uc.GetLessonAttendance(ctx, "lesson-1")
		if len(records) != 1 {
			t.Fatalf("expected 1 record, got %d", len(records))
		}
		if records[0].Status != domain.AttendanceStatusAbsentExcused {
			t.Errorf("expected ABSENT_EXCUSED, got %s", records[0].Status)
		}
		if records[0].Reason == nil || *records[0].Reason != "sick" {
			t.Errorf("expected reason 'sick', got %v", records[0].Reason)
		}
	})
}

func TestAttendanceUseCase_GetLessonAttendance(t *testing.T) {
	repoMock := mocks.NewAttendanceRepositoryMock()
	uc := usecase.NewAttendanceUseCase(repoMock)
	ctx := context.Background()

	uc.MarkAttendance(ctx, "lesson-1", "student-1", domain.AttendanceStatusAttended, nil, nil, "teacher-1")
	uc.MarkAttendance(ctx, "lesson-1", "student-2", domain.AttendanceStatusAbsentExcused, ptr("sick"), nil, "teacher-1")

	records, err := uc.GetLessonAttendance(ctx, "lesson-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

func TestAttendanceUseCase_GetStudentStats(t *testing.T) {
	repoMock := mocks.NewAttendanceRepositoryMock()
	uc := usecase.NewAttendanceUseCase(repoMock)
	ctx := context.Background()

	uc.MarkAttendance(ctx, "lesson-1", "student-1", domain.AttendanceStatusAttended, nil, nil, "teacher-1")
	uc.MarkAttendance(ctx, "lesson-2", "student-1", domain.AttendanceStatusAttended, nil, nil, "teacher-1")
	uc.MarkAttendance(ctx, "lesson-3", "student-1", domain.AttendanceStatusAbsentExcused, ptr("sick"), nil, "teacher-1")

	stats, err := uc.GetStudentStats(ctx, "student-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats["attended"] != 2 {
		t.Errorf("expected 2 attended, got %d", stats["attended"])
	}
	if stats["absent_excused"] != 1 {
		t.Errorf("expected 1 absent_excused, got %d", stats["absent_excused"])
	}
}

func TestAttendanceUseCase_GetStudentCalendar(t *testing.T) {
	repoMock := mocks.NewAttendanceRepositoryMock()
	uc := usecase.NewAttendanceUseCase(repoMock)
	ctx := context.Background()

	now := time.Now()
	uc.MarkAttendance(ctx, "lesson-1", "student-1", domain.AttendanceStatusAttended, nil, nil, "teacher-1")

	records, err := uc.GetStudentCalendar(ctx, "student-1", now.Add(-24*time.Hour), now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}
}
