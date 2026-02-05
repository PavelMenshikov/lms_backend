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

	GetTeachersList(ctx context.Context) ([]*domain.TeacherPublicInfo, error)
	GetTeacherByID(ctx context.Context, id string) (*domain.TeacherPublicInfo, error)
	AddTeacherReview(ctx context.Context, review *domain.TeacherReview) error
	GetTeacherReviews(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error)
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
	err := r.db.QueryRowContext(ctx, queryCourse, courseID).Scan(&course.ID, &course.Title, &course.Description, &course.ImageURL, &course.IsMain)
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
			CASE WHEN ula.is_attended = true OR uas.status = 'accepted' THEN true ELSE false END as is_completed
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
		l.IsLocked = false
		lessonsByModule[moduleID] = append(lessonsByModule[moduleID], &l)
	}
	for _, m := range modules {
		if list, ok := lessonsByModule[m.ID]; ok {
			m.Lessons = list
		} else {
			m.Lessons = []*domain.StudentLessonRef{}
		}
	}
	return &domain.StudentCourseView{Course: course, Modules: modules}, nil
}

func (r *LearningRepoImpl) GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error) {
	lesson := &domain.Lesson{}
	query := `SELECT id, title, video_url, presentation_url, content_text, duration_min FROM lessons WHERE id = $1 AND is_published = true`
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(&lesson.ID, &lesson.Title, &lesson.VideoURL, &lesson.PresentationURL, &lesson.ContentText, &lesson.DurationMin)
	if err != nil {
		return nil, fmt.Errorf("lesson not found: %w", err)
	}
	var isCompleted bool
	queryStatus := `SELECT EXISTS(SELECT 1 FROM user_lesson_attendance WHERE lesson_id = $1 AND user_id = $2 AND is_attended = true)`
	_ = r.db.QueryRowContext(ctx, queryStatus, lessonID, userID).Scan(&isCompleted)
	return &domain.StudentLessonDetail{Lesson: lesson, IsCompleted: isCompleted}, nil
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
		ON CONFLICT (user_id, assignment_id) DO UPDATE SET submission_text = EXCLUDED.submission_text, submission_link = EXCLUDED.submission_link, status = 'pending_check', submitted_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, userID, assignmentID, text, fileURL)
	return err
}

func (r *LearningRepoImpl) MarkLessonComplete(ctx context.Context, userID, lessonID string) error {
	query := `INSERT INTO user_lesson_attendance (user_id, lesson_id, is_attended) VALUES ($1, $2, true) ON CONFLICT (user_id, lesson_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, lessonID)
	return err
}

func (r *LearningRepoImpl) GetTeachersList(ctx context.Context) ([]*domain.TeacherPublicInfo, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar_url, ''), 
		       COALESCE(t.rating, 0.0), COALESCE(u.experience_years, 0)
		FROM users u
		LEFT JOIN teachers t ON u.id = t.id
		WHERE u.role = 'teacher'
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teachers []*domain.TeacherPublicInfo
	for rows.Next() {
		t := &domain.TeacherPublicInfo{}
		if err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.AvatarURL, &t.Rating, &t.ExperienceYears); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, nil
}

func (r *LearningRepoImpl) GetTeacherByID(ctx context.Context, id string) (*domain.TeacherPublicInfo, error) {
	t := &domain.TeacherPublicInfo{}
	query := `
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar_url, ''), 
		       COALESCE(t.rating, 0.0), COALESCE(u.experience_years, 0), COALESCE(t.bio, '')
		FROM users u
		LEFT JOIN teachers t ON u.id = t.id
		WHERE u.id = $1 AND u.role = 'teacher'
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.FirstName, &t.LastName, &t.AvatarURL, &t.Rating, &t.ExperienceYears, &t.Bio)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *LearningRepoImpl) AddTeacherReview(ctx context.Context, rev *domain.TeacherReview) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO teacher_reviews (teacher_id, student_id, rating, comment)
		VALUES ($1, $2, $3, $4)
	`
	if _, err := tx.ExecContext(ctx, query, rev.TeacherID, rev.StudentID, rev.Rating, rev.Comment); err != nil {
		tx.Rollback()
		return err
	}
	updateRating := `
		UPDATE teachers SET rating = (
			SELECT AVG(rating) FROM teacher_reviews WHERE teacher_id = $1
		) WHERE id = $1
	`
	if _, err := tx.ExecContext(ctx, updateRating, rev.TeacherID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *LearningRepoImpl) GetTeacherReviews(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error) {
	query := `
		SELECT tr.id, tr.teacher_id, tr.student_id, u.first_name || ' ' || u.last_name, tr.rating, tr.comment, tr.created_at
		FROM teacher_reviews tr
		JOIN users u ON tr.student_id = u.id
		WHERE tr.teacher_id = $1
		ORDER BY tr.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []*domain.TeacherReview
	for rows.Next() {
		rev := &domain.TeacherReview{}
		if err := rows.Scan(&rev.ID, &rev.TeacherID, &rev.StudentID, &rev.StudentName, &rev.Rating, &rev.Comment, &rev.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}
