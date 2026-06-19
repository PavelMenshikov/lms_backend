package mocks

import (
	"context"
	"errors"
	"lms_backend/internal/domain"
	"lms_backend/internal/notification/repository"
	"sync"
	"time"
)

type NotificationRepositoryMock struct {
	mu            sync.Mutex
	Notifications map[string]*domain.Notification
	nextID        int
}

var _ repository.NotificationRepository = (*NotificationRepositoryMock)(nil)

func NewNotificationRepositoryMock() *NotificationRepositoryMock {
	return &NotificationRepositoryMock{
		Notifications: make(map[string]*domain.Notification),
		nextID:        1,
	}
}

func (r *NotificationRepositoryMock) nextIDStr() string {
	id := r.nextID
	r.nextID++
	return string(rune('0'+id%10)) + string(rune('0'+(id/10)%10))
}

func (r *NotificationRepositoryMock) Create(ctx context.Context, notif *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := "notif-" + r.nextIDStr()
	notif.ID = id
	notif.CreatedAt = time.Now()
	notif.UpdatedAt = time.Now()
	r.Notifications[id] = notif
	return nil
}

func (r *NotificationRepositoryMock) GetByRecipient(ctx context.Context, recipientID string) ([]*domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Notification
	for _, n := range r.Notifications {
		if n.RecipientID == recipientID {
			result = append(result, n)
		}
	}
	return result, nil
}

func (r *NotificationRepositoryMock) GetUnreadByRecipient(ctx context.Context, recipientID string) ([]*domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Notification
	for _, n := range r.Notifications {
		if n.RecipientID == recipientID && !n.IsRead {
			result = append(result, n)
		}
	}
	return result, nil
}

func (r *NotificationRepositoryMock) MarkAsRead(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	n, ok := r.Notifications[id]
	if !ok {
		return errors.New("not found")
	}
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
	return nil
}

func (r *NotificationRepositoryMock) MarkAllAsRead(ctx context.Context, recipientID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range r.Notifications {
		if n.RecipientID == recipientID {
			n.IsRead = true
			now := time.Now()
			n.ReadAt = &now
		}
	}
	return nil
}
