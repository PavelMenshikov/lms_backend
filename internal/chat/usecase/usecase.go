package usecase

import (
	"context"
	"lms_backend/internal/chat/repository"
	"lms_backend/internal/domain"
	"sync"
)

type ChatUseCase struct {
	repo repository.ChatRepository
	// Храним активные подключения в памяти: map[RoomID][]Channels
	mu      sync.RWMutex
	clients map[string]map[chan *domain.ChatMessage]bool
}

func NewChatUseCase(repo repository.ChatRepository) *ChatUseCase {
	return &ChatUseCase{
		repo:    repo,
		clients: make(map[string]map[chan *domain.ChatMessage]bool),
	}
}

func (uc *ChatUseCase) SendMessage(ctx context.Context, msg *domain.ChatMessage) error {
	if err := uc.repo.SaveMessage(ctx, msg); err != nil {
		return err
	}

	roomID := msg.ModuleID + "_" + msg.StudentID

	uc.mu.RLock()
	if channels, ok := uc.clients[roomID]; ok {
		for ch := range channels {
			ch <- msg
		}
	}
	uc.mu.RUnlock()

	return nil
}

func (uc *ChatUseCase) GetHistory(ctx context.Context, moduleID, studentID string, limit, offset int) ([]*domain.ChatMessage, error) {
	return uc.repo.GetHistory(ctx, moduleID, studentID, limit, offset)
}

func (uc *ChatUseCase) RegisterClient(roomID string, ch chan *domain.ChatMessage) {
	uc.mu.Lock()
	if uc.clients[roomID] == nil {
		uc.clients[roomID] = make(map[chan *domain.ChatMessage]bool)
	}
	uc.clients[roomID][ch] = true
	uc.mu.Unlock()
}

func (uc *ChatUseCase) UnregisterClient(roomID string, ch chan *domain.ChatMessage) {
	uc.mu.Lock()
	if clients, ok := uc.clients[roomID]; ok {
		delete(clients, ch)
		if len(clients) == 0 {
			delete(uc.clients, roomID)
		}
	}
	uc.mu.Unlock()
}
