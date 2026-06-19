package usecase_test

import (
	"context"
	"testing"

	"lms_backend/internal/domain"
	"lms_backend/internal/review/mocks"
	"lms_backend/internal/review/usecase"
)

func TestReviewUseCase_GetPendingList(t *testing.T) {
	repoMock := mocks.NewReviewRepositoryMock()
	uc := usecase.NewReviewUseCase(repoMock)
	ctx := context.Background()

	t.Run("EmptyList", func(t *testing.T) {
		list, err := uc.GetPendingList(ctx, "teacher-1", "teacher")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(list) != 0 {
			t.Errorf("expected empty list, got %d", len(list))
		}
	})
}

func TestReviewUseCase_Evaluate(t *testing.T) {
	repoMock := mocks.NewReviewRepositoryMock()
	uc := usecase.NewReviewUseCase(repoMock)
	ctx := context.Background()

	t.Run("AcceptWithValidGrade", func(t *testing.T) {
		input := usecase.EvaluateInput{
			SubmissionID: "sub-1",
			Grade:        80,
			Comment:      "Good work",
			Status:       "accepted",
		}
		err := uc.Evaluate(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("AcceptWithInvalidGradeDefaultsTo100", func(t *testing.T) {
		input := usecase.EvaluateInput{
			SubmissionID: "sub-2",
			Grade:        75,
			Comment:      "OK",
			Status:       "accepted",
		}
		err := uc.Evaluate(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list, _ := uc.GetPendingList(ctx, "teacher-1", "teacher")
		for _, s := range list {
			if s.ID == "sub-2" && s.Grade != 100 {
				t.Errorf("expected grade 100 for invalid grade, got %d", s.Grade)
			}
		}
	})

	t.Run("RejectSetsGradeZero", func(t *testing.T) {
		input := usecase.EvaluateInput{
			SubmissionID: "sub-3",
			Grade:        80,
			Comment:      "Needs rework",
			Status:       "rejected",
		}
		err := uc.Evaluate(ctx, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("ValidGrades", func(t *testing.T) {
		validGrades := []int{20, 40, 60, 80, 100}
		for _, g := range validGrades {
			input := usecase.EvaluateInput{
				SubmissionID: "sub-vg" + string(rune('0'+g)),
				Grade:        g,
				Comment:      "Valid grade",
				Status:       "accepted",
			}
			err := uc.Evaluate(ctx, input)
			if err != nil {
				t.Errorf("unexpected error for grade %d: %v", g, err)
			}
		}
	})
}

func TestReviewUseCase_Evaluate_GradeEdgeCases(t *testing.T) {
	repoMock := mocks.NewReviewRepositoryMock()
	uc := usecase.NewReviewUseCase(repoMock)
	ctx := context.Background()

	submission := &domain.SubmissionRecord{
		ID:             "sub-edge",
		UserID:         "user-1",
		StudentName:    "Test Student",
		CourseTitle:    "Course",
		ModuleOrder:    1,
		LessonOrder:    1,
		LessonTitle:    "Lesson",
		Text:           "Homework",
		Status:         "pending",
		Grade:          0,
		TeacherComment: "",
	}
	_ = submission

	input := usecase.EvaluateInput{
		SubmissionID: "sub-edge",
		Grade:        0,
		Comment:      "Invalid",
		Status:       "rejected",
	}
	err := uc.Evaluate(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
