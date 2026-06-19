package mocks

import (
	"context"
	"errors"
	"lms_backend/internal/comment/repository"
	"lms_backend/internal/domain"
	"sync"
	"time"
)

type CommentRepositoryMock struct {
	mu       sync.Mutex
	Comments map[string]*domain.Comment
	nextID   int
}

var _ repository.CommentRepository = (*CommentRepositoryMock)(nil)

func NewCommentRepositoryMock() *CommentRepositoryMock {
	return &CommentRepositoryMock{
		Comments: make(map[string]*domain.Comment),
		nextID:   1,
	}
}

func (r *CommentRepositoryMock) Create(ctx context.Context, comment *domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := "comment-" + r.nextIDStr()
	comment.ID = id
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	r.Comments[id] = comment
	return nil
}

func (r *CommentRepositoryMock) nextIDStr() string {
	id := r.nextID
	r.nextID++
	return string(rune('0'+id%10)) + string(rune('0'+(id/10)%10))
}

func (r *CommentRepositoryMock) GetByStudent(ctx context.Context, studentID string) ([]*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Comment
	for _, c := range r.Comments {
		if c.StudentID == studentID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *CommentRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.Comments[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (r *CommentRepositoryMock) MarkAsRead(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.Comments[id]
	if !ok {
		return errors.New("not found")
	}
	c.IsRead = true
	now := time.Now()
	c.ReadAt = &now
	return nil
}

func (r *CommentRepositoryMock) GetUnreadByRecipient(ctx context.Context, recipientID string) ([]*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Comment
	for _, c := range r.Comments {
		if c.RecipientID != nil && *c.RecipientID == recipientID && !c.IsRead {
			result = append(result, c)
		}
	}
	return result, nil
}
