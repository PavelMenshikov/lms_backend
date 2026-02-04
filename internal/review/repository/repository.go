package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type ReviewRepository interface {
	GetPendingSubmissions(ctx context.Context) ([]*domain.SubmissionRecord, error)
	EvaluateSubmission(ctx context.Context, submissionID string, grade int, comment string, status string) error
	UpdateUserCourseProgress(ctx context.Context, userID, courseID string) error
}

type ReviewRepoImpl struct {
	db *sql.DB
}

var _ ReviewRepository = (*ReviewRepoImpl)(nil)

func NewReviewRepository(db *sql.DB) *ReviewRepoImpl {
	return &ReviewRepoImpl{db: db}
}

func (r *ReviewRepoImpl) GetPendingSubmissions(ctx context.Context) ([]*domain.SubmissionRecord, error) {
	query := `
		SELECT 
			uas.assignment_id,
			uas.user_id,
			u.first_name || ' ' || u.last_name as student_name,
			l.title as lesson_title,
			c.title as course_title,
			uas.submission_text,
			uas.submission_link,
			uas.submitted_at
		FROM user_assignments_submission uas
		JOIN users u ON uas.user_id = u.id
		JOIN assignments a ON uas.assignment_id = a.id
		JOIN lessons l ON a.lesson_id = l.id
		JOIN modules m ON l.module_id = m.id
		JOIN courses c ON m.course_id = c.id
		WHERE uas.status = 'pending_check'
		ORDER BY uas.submitted_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.SubmissionRecord
	for rows.Next() {
		rec := &domain.SubmissionRecord{}
		if err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.StudentName, &rec.LessonTitle,
			&rec.CourseTitle, &rec.Text, &rec.Link, &rec.SubmittedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r *ReviewRepoImpl) EvaluateSubmission(ctx context.Context, submissionID string, grade int, comment string, status string) error {
	query := `
		UPDATE user_assignments_submission 
		SET grade = $1, status = $2, teacher_comment = $3
		WHERE assignment_id = $4
	`
	_, err := r.db.ExecContext(ctx, query, grade, status, comment, submissionID)
	return err
}

func (r *ReviewRepoImpl) UpdateUserCourseProgress(ctx context.Context, userID, courseID string) error {
	query := `
		UPDATE user_courses 
		SET progress_percent = (
			SELECT 
				COUNT(CASE WHEN uas.status = 'accepted' THEN 1 END) * 100 / COUNT(a.id)
			FROM assignments a
			JOIN lessons l ON a.lesson_id = l.id
			JOIN modules m ON l.module_id = m.id
			LEFT JOIN user_assignments_submission uas ON a.id = uas.assignment_id AND uas.user_id = $1
			WHERE m.course_id = $2
		)
		WHERE user_id = $1 AND course_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, courseID)
	return err
}
