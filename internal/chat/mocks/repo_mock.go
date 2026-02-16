package mocks

import (
	"context"
	"lms_backend/internal/chat/repository"
	"lms_backend/internal/domain"
)

type ChatRepoMock struct {
	Messages []*domain.ChatMessage
}

func NewChatRepoMock() *ChatRepoMock {
	return &ChatRepoMock{
		Messages: make([]*domain.ChatMessage, 0),
	}
}

var _ repository.ChatRepository = (*ChatRepoMock)(nil)

func (m *ChatRepoMock) SaveMessage(ctx context.Context, msg *domain.ChatMessage) error {
	msg.ID = "msg-uuid"
	m.Messages = append(m.Messages, msg)
	return nil
}

func (m *ChatRepoMock) GetHistory(ctx context.Context, moduleID, studentID string, limit, offset int) ([]*domain.ChatMessage, error) {
	return m.Messages, nil
}

func (m *ChatRepoMock) MarkAsRead(ctx context.Context, moduleID, studentID, readerID string) error {
	return nil
}
