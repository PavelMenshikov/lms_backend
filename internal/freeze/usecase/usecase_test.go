package usecase_test

import (
	"context"
	"testing"
	"time"

	"lms_backend/internal/freeze/mocks"
	"lms_backend/internal/freeze/usecase"
)

func TestFreezeUseCase_CreateRequest(t *testing.T) {
	repoMock := mocks.NewFreezeRepositoryMock()
	uc := usecase.NewFreezeUseCase(repoMock)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "student-1", "curator-1",
			now, now.Add(7*24*time.Hour), "need break")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("EndDateBeforeStartDate", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "student-2", "curator-1",
			now.Add(7*24*time.Hour), now, "invalid dates")
		if err == nil {
			t.Error("expected error for end_date before start_date, got nil")
		}
	})
}

func TestFreezeUseCase_ApproveRequest(t *testing.T) {
	repoMock := mocks.NewFreezeRepositoryMock()
	uc := usecase.NewFreezeUseCase(repoMock)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "student-1", "curator-1",
			now, now.Add(7*24*time.Hour), "need break")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, _ := uc.GetPendingRequests(ctx)
		if len(reqs) != 1 {
			t.Fatalf("expected 1 pending request, got %d", len(reqs))
		}

		comment := "approved"
		err = uc.ApproveRequest(ctx, reqs[0].ID, "admin-1", &comment)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		status, err := uc.GetStudentFreezeStatus(ctx, "student-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status == nil {
			t.Fatal("expected active freeze period, got nil")
		}
	})

	t.Run("NotPending", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "student-2", "curator-1",
			now, now.Add(7*24*time.Hour), "need break")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, _ := uc.GetPendingRequests(ctx)
		comment := "approved"

		err = uc.ApproveRequest(ctx, reqs[0].ID, "admin-1", &comment)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = uc.ApproveRequest(ctx, reqs[0].ID, "admin-1", &comment)
		if err == nil {
			t.Error("expected error for already approved request, got nil")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		err := uc.ApproveRequest(ctx, "nonexistent", "admin-1", nil)
		if err == nil {
			t.Error("expected error for nonexistent request, got nil")
		}
	})
}

func TestFreezeUseCase_RejectRequest(t *testing.T) {
	repoMock := mocks.NewFreezeRepositoryMock()
	uc := usecase.NewFreezeUseCase(repoMock)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "student-1", "curator-1",
			now, now.Add(7*24*time.Hour), "need break")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, _ := uc.GetPendingRequests(ctx)
		comment := "not eligible"
		err = uc.RejectRequest(ctx, reqs[0].ID, "admin-1", &comment)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFreezeUseCase_GetStudentRequests(t *testing.T) {
	repoMock := mocks.NewFreezeRepositoryMock()
	uc := usecase.NewFreezeUseCase(repoMock)
	ctx := context.Background()
	now := time.Now()

	uc.CreateRequest(ctx, "student-1", "curator-1", now, now.Add(7*24*time.Hour), "reason 1")
	uc.CreateRequest(ctx, "student-1", "curator-1", now, now.Add(14*24*time.Hour), "reason 2")

	reqs, err := uc.GetStudentRequests(ctx, "student-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reqs) != 2 {
		t.Errorf("expected 2 requests, got %d", len(reqs))
	}
}
