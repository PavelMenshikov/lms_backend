package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"lms_backend/internal/domain"
)

func (r *ContentAdminRepoImpl) CreateCourse(ctx context.Context, course *domain.Course) (string, error) {
	var newID string
	query := `INSERT INTO courses (title, description, is_main, image_url, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, course.Title, course.Description, course.IsMain, course.ImageURL, course.Status).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) UpdateCourseSettings(ctx context.Context, course *domain.Course) error {
	query := `UPDATE courses SET has_homework = $1, is_homework_mandatory = $2, is_test_mandatory = $3, is_project_mandatory = $4, is_discord_mandatory = $5, is_anti_copy_enabled = $6, status = $7, title = $8, description = $9, image_url = $10, is_main = $11 WHERE id = $12`
	_, err := r.db.ExecContext(ctx, query, course.HasHomework, course.IsHomeworkMandatory, course.IsTestMandatory, course.IsProjectMandatory, course.IsDiscordMandatory, course.IsAntiCopyEnabled, course.Status, course.Title, course.Description, course.ImageURL, course.IsMain, course.ID)
	return err
}

func (r *ContentAdminRepoImpl) CreateModule(ctx context.Context, module *domain.Module) (string, error) {
	var newID string
	query := `INSERT INTO modules (course_id, title, description, order_num) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, module.CourseID, module.Title, module.Description, module.OrderNum).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) DeleteModule(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM modules WHERE id = $1", id)
	return err
}

func (r *ContentAdminRepoImpl) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description, is_main, image_url, status, created_at FROM courses ORDER BY created_at DESC LIMIT 200")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []*domain.Course
	for rows.Next() {
		c := &domain.Course{}
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.IsMain, &c.ImageURL, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *ContentAdminRepoImpl) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	c := &domain.Course{}
	err := r.db.QueryRowContext(ctx, "SELECT id, title, description, is_main, image_url, status, created_at, has_homework, is_homework_mandatory, is_test_mandatory, is_project_mandatory, is_discord_mandatory, is_anti_copy_enabled FROM courses WHERE id = $1", id).Scan(&c.ID, &c.Title, &c.Description, &c.IsMain, &c.ImageURL, &c.Status, &c.CreatedAt, &c.HasHomework, &c.IsHomeworkMandatory, &c.IsTestMandatory, &c.IsProjectMandatory, &c.IsDiscordMandatory, &c.IsAntiCopyEnabled)
	return c, err
}

func (r *ContentAdminRepoImpl) GetModulesByCourseID(ctx context.Context, courseID string) ([]*domain.Module, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, course_id, title, order_num, description FROM modules WHERE course_id = $1 ORDER BY order_num ASC", courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var modules []*domain.Module
	for rows.Next() {
		m := &domain.Module{}
		if err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.OrderNum, &m.Description); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *ContentAdminRepoImpl) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	query := `
		SELECT u.id, u.first_name || ' ' || u.last_name, uc.progress_percent,
			COALESCE(ula.attended, 0), COALESCE(uas.accepted, 0)
		FROM users u
		JOIN user_courses uc ON u.id = uc.user_id AND uc.course_id = $1
		LEFT JOIN (
			SELECT user_id, COUNT(*) as attended
			FROM user_lesson_attendance
			WHERE status IN ('visited', 'trial')
			GROUP BY user_id
		) ula ON u.id = ula.user_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as accepted
			FROM user_assignments_submission
			WHERE status = 'accepted'
			GROUP BY user_id
		) uas ON u.id = uas.user_id
		WHERE u.role = 'student'
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.AdminStudentProgress
	for rows.Next() {
		s := &domain.AdminStudentProgress{}
		if err := rows.Scan(&s.UserID, &s.FullName, &s.ProgressPercent, &s.LessonsAttended, &s.HomeworksDone); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	stats := &domain.AdminCourseStats{}
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(user_id), COALESCE(AVG(progress_percent), 0) FROM user_courses WHERE course_id = $1", courseID).Scan(&stats.TotalStudents, &stats.AverageScore)
	return stats, err
}

func (r *ContentAdminRepoImpl) GetCourseIDByStream(ctx context.Context, streamID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT course_id FROM streams WHERE id = $1", streamID).Scan(&id)
	return id, err
}

func (r *ContentAdminRepoImpl) EnrollStudentExtended(ctx context.Context, userID, courseID, streamID, groupID string) error {
	query := `INSERT INTO user_courses (user_id, course_id, stream_id, group_id, progress_percent, status)
		VALUES ($1, $2, $3, $4, 0, 'active')
		ON CONFLICT (user_id, course_id) DO UPDATE SET stream_id = EXCLUDED.stream_id, group_id = EXCLUDED.group_id`
	sid := sql.NullString{String: streamID, Valid: streamID != ""}
	gid := sql.NullString{String: groupID, Valid: groupID != ""}
	_, err := r.db.ExecContext(ctx, query, userID, courseID, sid, gid)
	return err
}

func (r *ContentAdminRepoImpl) UnenrollStudent(ctx context.Context, userID, courseID string) error {
	query := `DELETE FROM user_courses WHERE user_id = $1 AND course_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, courseID)
	return err
}

func (r *ContentAdminRepoImpl) LinkTeachersToCourse(ctx context.Context, courseID string, teacherIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM course_teachers WHERE course_id = $1", courseID); err != nil {
		tx.Rollback()
		return err
	}

	if len(teacherIDs) > 0 {
		valueStrs := make([]string, 0, len(teacherIDs))
		args := make([]interface{}, 0, 1+len(teacherIDs))
		args = append(args, courseID)
		for i, tID := range teacherIDs {
			valueStrs = append(valueStrs, fmt.Sprintf("($1, $%d)", i+2))
			args = append(args, tID)
		}
		query := fmt.Sprintf("INSERT INTO course_teachers (course_id, teacher_id) VALUES %s", strings.Join(valueStrs, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
