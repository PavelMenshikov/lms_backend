package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"lms_backend/internal/domain"
)

func (r *ContentAdminRepoImpl) CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error) {
	var newID string
	var tid sql.NullString
	if lesson.TeacherID != "" {
		tid = sql.NullString{String: lesson.TeacherID, Valid: true}
	}

	contentJSON, err := json.Marshal(lesson.Content)
	if err != nil {
		contentJSON = []byte("[]")
	}

	query := `INSERT INTO lessons (course_id, module_id, teacher_id, title, lesson_time, order_num, video_url, presentation_url, content_text, content, is_published)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

	err = r.db.QueryRowContext(ctx, query,
		lesson.CourseID,
		lesson.ModuleID,
		tid,
		lesson.Title,
		lesson.LessonTime,
		lesson.OrderNum,
		lesson.VideoURL,
		lesson.PresentationURL,
		lesson.ContentText,
		contentJSON,
		lesson.IsPublished,
	).Scan(&newID)

	return newID, err
}

func (r *ContentAdminRepoImpl) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	contentJSON, err := json.Marshal(lesson.Content)
	if err != nil {
		contentJSON = []byte("[]")
	}

	query := `
		UPDATE lessons 
		SET title = $1, order_num = $2, video_url = $3, presentation_url = $4, 
		    content_text = $5, content = $6, is_published = $7, module_id = $8, teacher_id = $9
		WHERE id = $10`

	var tid sql.NullString
	if lesson.TeacherID != "" {
		tid = sql.NullString{String: lesson.TeacherID, Valid: true}
	}

	_, err = r.db.ExecContext(ctx, query,
		lesson.Title, lesson.OrderNum, lesson.VideoURL, lesson.PresentationURL,
		lesson.ContentText, contentJSON, lesson.IsPublished, lesson.ModuleID, tid, lesson.ID,
	)
	return err
}

func (r *ContentAdminRepoImpl) DeleteLesson(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM lessons WHERE id = $1", id)
	return err
}

func (r *ContentAdminRepoImpl) GetLessonByID(ctx context.Context, id string) (*domain.Lesson, error) {
	l := &domain.Lesson{}
	var mid, tid sql.NullString
	var contentRaw []byte
	query := `SELECT id, course_id, module_id, teacher_id, title, lesson_time, duration_min, order_num, is_published, video_url, presentation_url, content_text, content 
              FROM lessons WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID, &l.CourseID, &mid, &tid, &l.Title, &l.LessonTime, &l.DurationMin, &l.OrderNum,
		&l.IsPublished, &l.VideoURL, &l.PresentationURL, &l.ContentText, &contentRaw,
	)
	if err != nil {
		return nil, err
	}
	if mid.Valid {
		s := mid.String
		l.ModuleID = &s
	}
	if tid.Valid {
		l.TeacherID = tid.String
	}
	if len(contentRaw) > 0 {
		_ = json.Unmarshal(contentRaw, &l.Content)
	}
	return l, nil
}

func (r *ContentAdminRepoImpl) GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, course_id, module_id, teacher_id, title, lesson_time, duration_min, order_num, is_published, video_url, presentation_url, content_text, content FROM lessons WHERE course_id = $1 ORDER BY order_num ASC", courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lessons []*domain.Lesson
	for rows.Next() {
		l := &domain.Lesson{}
		var mid, tid sql.NullString
		var contentRaw []byte

		err := rows.Scan(&l.ID, &l.CourseID, &mid, &tid, &l.Title, &l.LessonTime, &l.DurationMin, &l.OrderNum, &l.IsPublished, &l.VideoURL, &l.PresentationURL, &l.ContentText, &contentRaw)
		if err != nil {
			return nil, err
		}

		if mid.Valid {
			s := mid.String
			l.ModuleID = &s
		}
		if tid.Valid {
			l.TeacherID = tid.String
		}

		if len(contentRaw) > 0 {
			_ = json.Unmarshal(contentRaw, &l.Content)
		} else {
			l.Content = []domain.ContentBlock{}
		}

		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *ContentAdminRepoImpl) GetLessonIDByOrder(ctx context.Context, courseID string, orderNum int) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM lessons WHERE course_id = $1 AND order_num = $2", courseID, orderNum).Scan(&id)
	return id, err
}

func (r *ContentAdminRepoImpl) AssignTeacherToLesson(ctx context.Context, lessonID, teacherID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE lessons SET teacher_id = $1 WHERE id = $2", teacherID, lessonID)
	return err
}

func (r *ContentAdminRepoImpl) SetLessonModule(ctx context.Context, lessonID, moduleID string) error {
	query := `UPDATE lessons SET module_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, moduleID, lessonID)
	return err
}

func (r *ContentAdminRepoImpl) CancelLesson(ctx context.Context, lessonID, reason string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE lessons SET is_cancelled = TRUE, cancelled_at = NOW(), cancellation_reason = $1 WHERE id = $2`,
		reason, lessonID,
	)
	return err
}

func (r *ContentAdminRepoImpl) SubstituteTeacher(ctx context.Context, lessonID, teacherID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var originalTeacherID sql.NullString
	err = tx.QueryRowContext(ctx,
		`SELECT teacher_id FROM lessons WHERE id = $1 FOR UPDATE`, lessonID,
	).Scan(&originalTeacherID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE lessons SET substituted_teacher_id = $1 WHERE id = $2`,
		teacherID, lessonID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO lesson_substitutions (lesson_id, original_teacher_id, substitute_teacher_id) VALUES ($1, $2, $3)`,
		lessonID, originalTeacherID, teacherID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ContentAdminRepoImpl) CreateTest(ctx context.Context, test *domain.Test) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO tests (lesson_id, title, description, passing_score) VALUES ($1, $2, $3, $4) RETURNING id", test.LessonID, test.Title, test.Description, test.PassingScore).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) DeleteTest(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tests WHERE id = $1", id)
	return err
}

func (r *ContentAdminRepoImpl) GetTestByID(ctx context.Context, id string) (*domain.Test, error) {
	t := &domain.Test{}
	err := r.db.QueryRowContext(ctx, "SELECT id, lesson_id, title, description, passing_score, created_at FROM tests WHERE id = $1", id).
		Scan(&t.ID, &t.LessonID, &t.Title, &t.Description, &t.PassingScore, &t.CreatedAt)
	return t, err
}

func (r *ContentAdminRepoImpl) GetTestsByCourseID(ctx context.Context, courseID string) ([]domain.Test, error) {
	query := `SELECT t.id, t.lesson_id, t.title, t.description, t.passing_score, t.created_at 
			  FROM tests t 
			  JOIN lessons l ON t.lesson_id = l.id 
			  WHERE l.course_id = $1`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tests []domain.Test
	for rows.Next() {
		var t domain.Test
		var lid sql.NullString
		err := rows.Scan(&t.ID, &lid, &t.Title, &t.Description, &t.PassingScore, &t.CreatedAt)
		if err == nil {
			if lid.Valid {
				s := lid.String
				t.LessonID = &s
			}
			tests = append(tests, t)
		}
	}
	return tests, nil
}

func (r *ContentAdminRepoImpl) CreateProject(ctx context.Context, project *domain.Project) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO projects (lesson_id, title, description, max_score) VALUES ($1, $2, $3, $4) RETURNING id", project.LessonID, project.Title, project.Description, project.MaxScore).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) DeleteProject(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}

func (r *ContentAdminRepoImpl) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	p := &domain.Project{}
	err := r.db.QueryRowContext(ctx, "SELECT id, lesson_id, title, description, max_score, created_at FROM projects WHERE id = $1", id).
		Scan(&p.ID, &p.LessonID, &p.Title, &p.Description, &p.MaxScore, &p.CreatedAt)
	return p, err
}

func (r *ContentAdminRepoImpl) GetProjectsByCourseID(ctx context.Context, courseID string) ([]domain.Project, error) {
	query := `SELECT p.id, p.lesson_id, p.title, p.description, p.max_score, p.created_at 
			  FROM projects p 
			  JOIN lessons l ON p.lesson_id = l.id 
			  WHERE l.course_id = $1`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var lid sql.NullString
		err := rows.Scan(&p.ID, &lid, &p.Title, &p.Description, &p.MaxScore, &p.CreatedAt)
		if err == nil {
			if lid.Valid {
				s := lid.String
				p.LessonID = &s
			}
			projects = append(projects, p)
		}
	}
	return projects, nil
}
