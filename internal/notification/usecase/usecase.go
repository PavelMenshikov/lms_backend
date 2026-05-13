package usecase

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/notification/repository"
)

type NotificationUseCase interface {
	CreateNotification(ctx context.Context, recipientID string, senderID *string, title, content string, notifType domain.NotificationType, linkURL *string) error
	GetUserNotifications(ctx context.Context, recipientID string) ([]*domain.Notification, error)
	GetUnreadNotifications(ctx context.Context, recipientID string) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID string) error
	MarkAllAsRead(ctx context.Context, recipientID string) error
}

type notificationUseCase struct {
	repo repository.NotificationRepository
}

func NewNotificationUseCase(repo repository.NotificationRepository) NotificationUseCase {
	return &notificationUseCase{repo: repo}
}

func (uc *notificationUseCase) CreateNotification(ctx context.Context, recipientID string, senderID *string, title, content string, notifType domain.NotificationType, linkURL *string) error {
	notif := &domain.Notification{
		RecipientID: recipientID,
		SenderID:    senderID,
		Title:       title,
		Content:     content,
		Type:        notifType,
		LinkURL:     linkURL,
	}
	return uc.repo.Create(ctx, notif)
}

func (uc *notificationUseCase) GetUserNotifications(ctx context.Context, recipientID string) ([]*domain.Notification, error) {
	return uc.repo.GetByRecipient(ctx, recipientID)
}

func (uc *notificationUseCase) GetUnreadNotifications(ctx context.Context, recipientID string) ([]*domain.Notification, error) {
	return uc.repo.GetUnreadByRecipient(ctx, recipientID)
}

func (uc *notificationUseCase) MarkAsRead(ctx context.Context, notificationID string) error {
	return uc.repo.MarkAsRead(ctx, notificationID)
}

func (uc *notificationUseCase) MarkAllAsRead(ctx context.Context, recipientID string) error {
	return uc.repo.MarkAllAsRead(ctx, recipientID)
}
