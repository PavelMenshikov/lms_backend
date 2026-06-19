package usecase_test

import (
	"context"
	"testing"

	"lms_backend/internal/domain"
	"lms_backend/internal/notification/mocks"
	"lms_backend/internal/notification/usecase"
)

func TestNotificationUseCase_CreateAndGet(t *testing.T) {
	repoMock := mocks.NewNotificationRepositoryMock()
	uc := usecase.NewNotificationUseCase(repoMock)
	ctx := context.Background()

	t.Run("CreateNotification", func(t *testing.T) {
		err := uc.CreateNotification(ctx, "user-1", nil, "Hello", "Content", domain.NotificationTypeInfo, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetUserNotifications", func(t *testing.T) {
		notifs, err := uc.GetUserNotifications(ctx, "user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(notifs) != 1 {
			t.Errorf("expected 1 notification, got %d", len(notifs))
		}
	})

	t.Run("GetUnreadNotifications", func(t *testing.T) {
		notifs, err := uc.GetUnreadNotifications(ctx, "user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(notifs) != 1 {
			t.Errorf("expected 1 unread notification, got %d", len(notifs))
		}
	})
}

func TestNotificationUseCase_MarkAsRead(t *testing.T) {
	repoMock := mocks.NewNotificationRepositoryMock()
	uc := usecase.NewNotificationUseCase(repoMock)
	ctx := context.Background()

	uc.CreateNotification(ctx, "user-1", nil, "Test", "Content", domain.NotificationTypeInfo, nil)
	notifs, _ := uc.GetUserNotifications(ctx, "user-1")
	id := notifs[0].ID

	t.Run("MarkAsRead", func(t *testing.T) {
		err := uc.MarkAsRead(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		unread, _ := uc.GetUnreadNotifications(ctx, "user-1")
		if len(unread) != 0 {
			t.Errorf("expected 0 unread after mark, got %d", len(unread))
		}
	})
}

func TestNotificationUseCase_MarkAllAsRead(t *testing.T) {
	repoMock := mocks.NewNotificationRepositoryMock()
	uc := usecase.NewNotificationUseCase(repoMock)
	ctx := context.Background()

	uc.CreateNotification(ctx, "user-1", nil, "Test 1", "Content", domain.NotificationTypeInfo, nil)
	uc.CreateNotification(ctx, "user-1", nil, "Test 2", "Content", domain.NotificationTypeWarning, nil)

	t.Run("MarkAllAsRead", func(t *testing.T) {
		err := uc.MarkAllAsRead(ctx, "user-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		unread, _ := uc.GetUnreadNotifications(ctx, "user-1")
		if len(unread) != 0 {
			t.Errorf("expected 0 unread after mark all, got %d", len(unread))
		}
	})
}
