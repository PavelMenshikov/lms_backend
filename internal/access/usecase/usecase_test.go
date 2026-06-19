package usecase_test

import (
	"context"
	"testing"

	"lms_backend/internal/access/mocks"
	"lms_backend/internal/access/usecase"
	"lms_backend/internal/domain"
)

func TestAccessUseCase_CreateRequest(t *testing.T) {
	repoMock := mocks.NewAccessRepositoryMock()
	uc := usecase.NewAccessUseCase(repoMock)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "user-1", "course", "course-1", "need access")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, err := uc.GetUserRequests(ctx, "user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(reqs) != 1 {
			t.Fatalf("expected 1 request, got %d", len(reqs))
		}
		if reqs[0].Status != domain.AccessRequestStatusPending {
			t.Errorf("expected PENDING, got %s", reqs[0].Status)
		}
	})
}

func TestAccessUseCase_ApproveRequest(t *testing.T) {
	repoMock := mocks.NewAccessRepositoryMock()
	uc := usecase.NewAccessUseCase(repoMock)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "user-1", "course", "course-1", "need access")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, _ := uc.GetUserRequests(ctx, "user-1")
		comment := "approved"
		err = uc.ApproveRequest(ctx, reqs[0].ID, "admin-1", &comment)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("NotPending", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "user-2", "course", "course-2", "need access")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, _ := uc.GetUserRequests(ctx, "user-2")
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
		comment := "approved"
		err := uc.ApproveRequest(ctx, "nonexistent", "admin-1", &comment)
		if err == nil {
			t.Error("expected error for nonexistent request, got nil")
		}
	})
}

func TestAccessUseCase_RejectRequest(t *testing.T) {
	repoMock := mocks.NewAccessRepositoryMock()
	uc := usecase.NewAccessUseCase(repoMock)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		err := uc.CreateRequest(ctx, "user-1", "course", "course-1", "need access")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reqs, _ := uc.GetUserRequests(ctx, "user-1")
		comment := "rejected because"
		err = uc.RejectRequest(ctx, reqs[0].ID, "admin-1", &comment)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("NotPending", func(t *testing.T) {
		err := uc.RejectRequest(ctx, "nonexistent", "admin-1", nil)
		if err == nil {
			t.Error("expected error for nonexistent request, got nil")
		}
	})
}

func TestAccessUseCase_GetPendingRequests(t *testing.T) {
	repoMock := mocks.NewAccessRepositoryMock()
	uc := usecase.NewAccessUseCase(repoMock)
	ctx := context.Background()

	uc.CreateRequest(ctx, "user-1", "course", "course-1", "reason 1")
	uc.CreateRequest(ctx, "user-2", "course", "course-2", "reason 2")

	reqs, err := uc.GetPendingRequests(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reqs) != 2 {
		t.Errorf("expected 2 pending requests, got %d", len(reqs))
	}
}
