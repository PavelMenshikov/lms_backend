package domain

import "time"

type NotificationType string

const (
	NotificationTypeInfo    NotificationType = "INFO"
	NotificationTypeWarning NotificationType = "WARNING"
	NotificationTypeError   NotificationType = "ERROR"
	NotificationTypeSuccess NotificationType = "SUCCESS"
)

type Notification struct {
	ID          string           `json:"id" db:"id"`
	RecipientID string           `json:"recipient_id" db:"recipient_id"`
	SenderID    *string          `json:"sender_id,omitempty" db:"sender_id"`
	Title       string           `json:"title" db:"title"`
	Content     string           `json:"content" db:"content"`
	Type        NotificationType `json:"type" db:"type"`
	IsRead      bool             `json:"is_read" db:"is_read"`
	ReadAt      *time.Time       `json:"read_at,omitempty" db:"read_at"`
	LinkURL     *string          `json:"link_url,omitempty" db:"link_url"`
	Metadata    *string          `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}
