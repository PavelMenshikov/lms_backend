package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lms_backend/internal/domain"
)

type LearningRepository interface {
	GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error)
	GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error)
	GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error)

	GetAssignmentIDByLesson(ctx context.Context, lessonID string) (string, error)
	SaveSubmission(ctx context.Context, userID, assignmentID, text, fileURL string) error
	MarkLessonComplete(ctx context.Context, userID, lessonID string) error
}

type LearningRepoImpl struct {
	db *sql.DB
}

var _ LearningRepository = (*LearningRepoImpl)(nil)

func NewLearningRepository(db *sql.DB) *LearningRepoImpl {
	return &LearningRepoImpl{db: db}
}

func (r *LearningRepoImpl) GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
	query := `
		SELECT c.id, c.title, c.description, c.image_url, c.is_main, uc.progress_percent
		FROM courses c
		JOIN user_courses uc ON c.id = uc.course_id
		WHERE uc.user_id = $1 AND c.status = 'active'
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*domain.StudentCoursePreview
	for rows.Next() {
		c := &domain.StudentCoursePreview{}
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.ImageURL, &c.IsMain, &c.ProgressPercent); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *LearningRepoImpl) GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error) {
	course := &domain.Course{}
	queryCourse := `SELECT id, title, description, image_url, is_main FROM courses WHERE id = $1`
	err := r.db.QueryRowContext(ctx, queryCourse, courseID).Scan(
		&course.ID, &course.Title, &course.Description, &course.ImageURL, &course.IsMain,
	)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	queryModules := `SELECT id, title, description, order_num FROM modules WHERE course_id = $1 ORDER BY order_num ASC`
	rowsM, err := r.db.QueryContext(ctx, queryModules, courseID)
	if err != nil {
		return nil, err
	}
	defer rowsM.Close()

	var modules []*domain.StudentModuleView
	for rowsM.Next() {
		m := &domain.StudentModuleView{}
		if err := rowsM.Scan(&m.ID, &m.Title, &m.Description, &m.OrderNum); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}

	queryLessons := `
		SELECT 
			l.id, l.module_id, l.title, l.order_num, l.duration_min,
			CASE 
				WHEN ula.is_attended = true THEN true 
				WHEN uas.status = 'accepted' THEN true 
				ELSE false 
			END as is_completed
		FROM lessons l
		JOIN modules m ON l.module_id = m.id
		LEFT JOIN user_lesson_attendance ula ON l.id = ula.lesson_id AND ula.user_id = $2
		LEFT JOIN assignments a ON l.id = a.lesson_id
		LEFT JOIN user_assignments_submission uas ON a.id = uas.assignment_id AND uas.user_id = $2
		WHERE m.course_id = $1 AND l.is_published = true
		ORDER BY l.order_num ASC
	`
	rowsL, err := r.db.QueryContext(ctx, queryLessons, courseID, userID)
	if err != nil {
		return nil, err
	}
	defer rowsL.Close()

	lessonsByModule := make(map[string][]*domain.StudentLessonRef)
	for rowsL.Next() {
		var l domain.StudentLessonRef
		var moduleID string
		if err := rowsL.Scan(&l.ID, &moduleID, &l.Title, &l.OrderNum, &l.DurationMin, &l.IsCompleted); err != nil {
			return nil, err
		}
		l.IsLocked = false // Пока открываем всё
		lessonsByModule[moduleID] = append(lessonsByModule[moduleID], &l)
	}

	for _, m := range modules {
		if list, ok := lessonsByModule[m.ID]; ok {
			m.Lessons = list
		} else {
			m.Lessons = []*domain.StudentLessonRef{}
		}
	}

	return &domain.StudentCourseView{
		Course:  course,
		Modules: modules,
	}, nil
}

func (r *LearningRepoImpl) GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error) {
	lesson := &domain.Lesson{}
	query := `
		SELECT id, title, video_url, presentation_url, content_text, duration_min 
		FROM lessons WHERE id = $1 AND is_published = true
	`
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(
		&lesson.ID, &lesson.Title, &lesson.VideoURL, &lesson.PresentationURL, &lesson.ContentText, &lesson.DurationMin,
	)
	if err != nil {
		return nil, fmt.Errorf("lesson not found: %w", err)
	}

	var isCompleted bool
	var assignmentStatus sql.NullString
	var teacherComment sql.NullString

	queryStatus := `
		SELECT 
			(ula.is_attended IS TRUE) as is_attended,
			uas.status,
			uas.grade -- можно добавить комментарий в таблицу submission позже, пока grade
		FROM lessons l
		LEFT JOIN user_lesson_attendance ula ON l.id = ula.lesson_id AND ula.user_id = $2
		LEFT JOIN assignments a ON l.id = a.lesson_id
		LEFT JOIN user_assignments_submission uas ON a.id = uas.assignment_id AND uas.user_id = $2
		WHERE l.id = $1
	`

	_ = r.db.QueryRowContext(ctx, queryStatus, lessonID, userID).Scan(&isCompleted, &assignmentStatus, &teacherComment)

	return &domain.StudentLessonDetail{
		Lesson:           lesson,
		IsCompleted:      isCompleted,
		AssignmentStatus: assignmentStatus.String,
	}, nil
}

func (r *LearningRepoImpl) GetAssignmentIDByLesson(ctx context.Context, lessonID string) (string, error) {
	var id string
	query := `SELECT id FROM assignments WHERE lesson_id = $1 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(&id)
	return id, err
}

func (r *LearningRepoImpl) SaveSubmission(ctx context.Context, userID, assignmentID, text, fileURL string) error {
	query := `
		INSERT INTO user_assignments_submission (user_id, assignment_id, submission_text, submission_link, status, submitted_at)
		VALUES ($1, $2, $3, $4, 'pending_check', NOW())
		ON CONFLICT (user_id, assignment_id) 
		DO UPDATE SET 
			submission_text = EXCLUDED.submission_text, 
			submission_link = EXCLUDED.submission_link,
			status = 'pending_check',
			submitted_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, userID, assignmentID, text, fileURL)
	return err
}

func (r *LearningRepoImpl) MarkLessonComplete(ctx context.Context, userID, lessonID string) error {
	query := `
		INSERT INTO user_lesson_attendance (user_id, lesson_id, is_attended)
		VALUES ($1, $2, true)
		ON CONFLICT (user_id, lesson_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, lessonID)
	return err
}
