package usecase_test

import (
	"context"
	"testing"

	"lms_backend/internal/comment/mocks"
	"lms_backend/internal/comment/usecase"
)

func TestCommentUseCase_CreateAndGet(t *testing.T) {
	repoMock := mocks.NewCommentRepositoryMock()
	uc := usecase.NewCommentUseCase(repoMock)
	ctx := context.Background()

	t.Run("CreateAndGetByStudent", func(t *testing.T) {
		err := uc.CreateComment(ctx, "student-1", nil, "author-1", nil, "Great job!", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		comments, err := uc.GetStudentComments(ctx, "student-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(comments) != 1 {
			t.Errorf("expected 1 comment, got %d", len(comments))
		}
		if comments[0].Content != "Great job!" {
			t.Errorf("expected 'Great job!', got %s", comments[0].Content)
		}
	})

	t.Run("CreateWithRecipient", func(t *testing.T) {
		recipientID := "teacher-1"
		err := uc.CreateComment(ctx, "student-1", nil, "author-1", &recipientID, "Please check", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCommentUseCase_MarkAsRead(t *testing.T) {
	repoMock := mocks.NewCommentRepositoryMock()
	uc := usecase.NewCommentUseCase(repoMock)
	ctx := context.Background()

	recipientID := "teacher-1"
	err := uc.CreateComment(ctx, "student-1", nil, "student-1", &recipientID, "Please check", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("UnreadBeforeMark", func(t *testing.T) {
		unread, err := uc.GetUnreadComments(ctx, "teacher-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(unread) != 1 {
			t.Errorf("expected 1 unread comment, got %d", len(unread))
		}
	})

	t.Run("MarkAsRead", func(t *testing.T) {
		comments, _ := uc.GetUnreadComments(ctx, "teacher-1")
		err := uc.MarkAsRead(ctx, comments[0].ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		unread, _ := uc.GetUnreadComments(ctx, "teacher-1")
		if len(unread) != 0 {
			t.Errorf("expected 0 unread after marking read, got %d", len(unread))
		}
	})
}
