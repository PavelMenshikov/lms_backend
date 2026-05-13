package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetByStudent(ctx context.Context, studentID string) ([]*domain.Comment, error)
	GetByID(ctx context.Context, id string) (*domain.Comment, error)
	MarkAsRead(ctx context.Context, id string) error
	GetUnreadByRecipient(ctx context.Context, recipientID string) ([]*domain.Comment, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (id, student_id, lesson_id, author_id, recipient_id, content, parent_comment_id)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		comment.StudentID, comment.LessonID, comment.AuthorID, comment.RecipientID,
		comment.Content, comment.ParentCommentID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (r *commentRepository) GetByStudent(ctx context.Context, studentID string) ([]*domain.Comment, error) {
	query := `
		SELECT id, student_id, lesson_id, author_id, recipient_id, content, is_read, read_at, parent_comment_id, created_at, updated_at
		FROM comments
		WHERE student_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(
			&c.ID, &c.StudentID, &c.LessonID, &c.AuthorID, &c.RecipientID,
			&c.Content, &c.IsRead, &c.ReadAt, &c.ParentCommentID, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

func (r *commentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	var c domain.Comment
	query := `
		SELECT id, student_id, lesson_id, author_id, recipient_id, content, is_read, read_at, parent_comment_id, created_at, updated_at
		FROM comments
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.StudentID, &c.LessonID, &c.AuthorID, &c.RecipientID,
		&c.Content, &c.IsRead, &c.ReadAt, &c.ParentCommentID, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *commentRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `
		UPDATE comments
		SET is_read = true, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *commentRepository) GetUnreadByRecipient(ctx context.Context, recipientID string) ([]*domain.Comment, error) {
	query := `
		SELECT id, student_id, lesson_id, author_id, recipient_id, content, is_read, read_at, parent_comment_id, created_at, updated_at
		FROM comments
		WHERE recipient_id = $1 AND is_read = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		err := rows.Scan(
			&c.ID, &c.StudentID, &c.LessonID, &c.AuthorID, &c.RecipientID,
			&c.Content, &c.IsRead, &c.ReadAt, &c.ParentCommentID, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}
