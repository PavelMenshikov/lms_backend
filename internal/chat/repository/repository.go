package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *domain.ChatMessage) error
	GetHistory(ctx context.Context, moduleID, studentID string, limit, offset int) ([]*domain.ChatMessage, error)
	MarkAsRead(ctx context.Context, moduleID, studentID, readerID string) error
}

type ChatRepoImpl struct {
	db *sql.DB
}

var _ ChatRepository = (*ChatRepoImpl)(nil)

func NewChatRepository(db *sql.DB) *ChatRepoImpl {
	return &ChatRepoImpl{db: db}
}

func (r *ChatRepoImpl) SaveMessage(ctx context.Context, msg *domain.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (module_id, student_id, sender_id, message_text, file_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		msg.ModuleID, msg.StudentID, msg.SenderID, msg.MessageText, msg.FileURL,
	).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *ChatRepoImpl) GetHistory(ctx context.Context, moduleID, studentID string, limit, offset int) ([]*domain.ChatMessage, error) {
	query := `
		SELECT 
			m.id, m.module_id, m.student_id, m.sender_id, 
			u.first_name || ' ' || u.last_name as sender_name,
			u.role as sender_role,
			m.message_text, m.file_url, m.is_read, m.created_at
		FROM chat_messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.module_id = $1 AND m.student_id = $2
		ORDER BY m.created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, query, moduleID, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.ChatMessage
	for rows.Next() {
		msg := &domain.ChatMessage{}
		err := rows.Scan(
			&msg.ID, &msg.ModuleID, &msg.StudentID, &msg.SenderID,
			&msg.SenderName, &msg.SenderRole, &msg.MessageText,
			&msg.FileURL, &msg.IsRead, &msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *ChatRepoImpl) MarkAsRead(ctx context.Context, moduleID, studentID, readerID string) error {
	query := `
		UPDATE chat_messages 
		SET is_read = true 
		WHERE module_id = $1 AND student_id = $2 AND sender_id != $3 AND is_read = false
	`
	_, err := r.db.ExecContext(ctx, query, moduleID, studentID, readerID)
	return err
}
