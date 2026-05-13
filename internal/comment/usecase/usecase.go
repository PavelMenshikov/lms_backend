package usecase

import (
	"context"
	"lms_backend/internal/comment/repository"
	"lms_backend/internal/domain"
)

type CommentUseCase interface {
	CreateComment(ctx context.Context, studentID string, lessonID *string, authorID string, recipientID *string, content string, parentCommentID *string) error
	GetStudentComments(ctx context.Context, studentID string) ([]*domain.Comment, error)
	MarkAsRead(ctx context.Context, commentID string) error
	GetUnreadComments(ctx context.Context, recipientID string) ([]*domain.Comment, error)
}

type commentUseCase struct {
	repo repository.CommentRepository
}

func NewCommentUseCase(repo repository.CommentRepository) CommentUseCase {
	return &commentUseCase{repo: repo}
}

func (uc *commentUseCase) CreateComment(ctx context.Context, studentID string, lessonID *string, authorID string, recipientID *string, content string, parentCommentID *string) error {
	comment := &domain.Comment{
		StudentID:       studentID,
		LessonID:        lessonID,
		AuthorID:        authorID,
		RecipientID:     recipientID,
		Content:         content,
		ParentCommentID: parentCommentID,
	}
	return uc.repo.Create(ctx, comment)
}

func (uc *commentUseCase) GetStudentComments(ctx context.Context, studentID string) ([]*domain.Comment, error) {
	return uc.repo.GetByStudent(ctx, studentID)
}

func (uc *commentUseCase) MarkAsRead(ctx context.Context, commentID string) error {
	return uc.repo.MarkAsRead(ctx, commentID)
}

func (uc *commentUseCase) GetUnreadComments(ctx context.Context, recipientID string) ([]*domain.Comment, error) {
	return uc.repo.GetUnreadByRecipient(ctx, recipientID)
}
