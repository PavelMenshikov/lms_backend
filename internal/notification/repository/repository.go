package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *domain.Notification) error
	GetByRecipient(ctx context.Context, recipientID string) ([]*domain.Notification, error)
	GetUnreadByRecipient(ctx context.Context, recipientID string) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, recipientID string) error
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notif *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, recipient_id, sender_id, title, content, type, link_url)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		notif.RecipientID, notif.SenderID, notif.Title, notif.Content, notif.Type, notif.LinkURL,
	).Scan(&notif.ID, &notif.CreatedAt, &notif.UpdatedAt)
}

func (r *notificationRepository) GetByRecipient(ctx context.Context, recipientID string) ([]*domain.Notification, error) {
	query := `
		SELECT id, recipient_id, sender_id, title, content, type, is_read, read_at, link_url, metadata, created_at, updated_at
		FROM notifications
		WHERE recipient_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(
			&n.ID, &n.RecipientID, &n.SenderID, &n.Title, &n.Content,
			&n.Type, &n.IsRead, &n.ReadAt, &n.LinkURL, &n.Metadata,
			&n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

func (r *notificationRepository) GetUnreadByRecipient(ctx context.Context, recipientID string) ([]*domain.Notification, error) {
	query := `
		SELECT id, recipient_id, sender_id, title, content, type, is_read, read_at, link_url, metadata, created_at, updated_at
		FROM notifications
		WHERE recipient_id = $1 AND is_read = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(
			&n.ID, &n.RecipientID, &n.SenderID, &n.Title, &n.Content,
			&n.Type, &n.IsRead, &n.ReadAt, &n.LinkURL, &n.Metadata,
			&n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, recipientID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE recipient_id = $1 AND is_read = false
	`
	_, err := r.db.ExecContext(ctx, query, recipientID)
	return err
}
