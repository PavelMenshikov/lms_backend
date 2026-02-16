package usecase_test

import (
	"context"
	"testing"

	"lms_backend/internal/chat/mocks"
	"lms_backend/internal/chat/usecase"
	"lms_backend/internal/domain"
)

func TestSendMessage(t *testing.T) {
	repoMock := mocks.NewChatRepoMock()
	uc := usecase.NewChatUseCase(repoMock)

	ctx := context.Background()

	msg := &domain.ChatMessage{
		ModuleID:    "mod-1",
		StudentID:   "stud-1",
		SenderID:    "stud-1",
		MessageText: "Hello Teacher!",
	}

	err := uc.SendMessage(ctx, msg)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if len(repoMock.Messages) != 1 {
		t.Error("Message not saved to repo")
	}

	if repoMock.Messages[0].MessageText != "Hello Teacher!" {
		t.Error("Message text mismatch")
	}
}
