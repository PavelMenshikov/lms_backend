package domain

import "time"

type Comment struct {
	ID              string     `json:"id" db:"id"`
	StudentID       string     `json:"student_id" db:"student_id"`
	LessonID        *string    `json:"lesson_id,omitempty" db:"lesson_id"`
	AuthorID        string     `json:"author_id" db:"author_id"`
	RecipientID     *string    `json:"recipient_id,omitempty" db:"recipient_id"`
	Content         string     `json:"content" db:"content"`
	IsRead          bool       `json:"is_read" db:"is_read"`
	ReadAt          *time.Time `json:"read_at,omitempty" db:"read_at"`
	ParentCommentID *string    `json:"parent_comment_id,omitempty" db:"parent_comment_id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
