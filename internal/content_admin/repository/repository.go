package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lms_backend/internal/domain"
)

type ContentAdminRepository interface {
	CreateCourse(ctx context.Context, course *domain.Course) (string, error)
	UpdateCourseSettings(ctx context.Context, course *domain.Course) error
	CreateModule(ctx context.Context, module *domain.Module) (string, error)
	CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error)
	GetAllCourses(ctx context.Context) ([]*domain.Course, error)
	GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
	GetModulesByCourseID(ctx context.Context, courseID string) ([]*domain.Module, error)
	GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error)
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error)
	LinkParentToStudent(ctx context.Context, studentID, parentID string) error
	EnrollStudent(ctx context.Context, userID, courseID string) error
	GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error)
	GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error)
	CreateTest(ctx context.Context, test *domain.Test) (string, error)
	CreateProject(ctx context.Context, project *domain.Project) (string, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID string) error
	GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error)
	CreateStream(ctx context.Context, stream *domain.Stream) (string, error)
	GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error)
	CreateGroup(ctx context.Context, group *domain.Group) (string, error)
	GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error)
	GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error)
	GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error)
}

type ContentAdminRepoImpl struct {
	db *sql.DB
}

var _ ContentAdminRepository = (*ContentAdminRepoImpl)(nil)

func NewContentAdminRepository(db *sql.DB) *ContentAdminRepoImpl {
	return &ContentAdminRepoImpl{db: db}
}

func (r *ContentAdminRepoImpl) CreateCourse(ctx context.Context, course *domain.Course) (string, error) {
	var newID string
	query := `
        INSERT INTO courses (title, description, is_main, image_url, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		course.Title,
		course.Description,
		course.IsMain,
		course.ImageURL,
		course.Status,
	).Scan(&newID)

	if err != nil {
		return "", fmt.Errorf("failed to create course: %w", err)
	}
	course.ID = newID
	return newID, nil
}

func (r *ContentAdminRepoImpl) UpdateCourseSettings(ctx context.Context, course *domain.Course) error {
	query := `
		UPDATE courses 
		SET has_homework = $1,
			is_homework_mandatory = $2,
			is_test_mandatory = $3,
			is_project_mandatory = $4,
			is_discord_mandatory = $5,
			is_anti_copy_enabled = $6,
			status = $7,
			title = $8,
			description = $9,
			image_url = $10,
			is_main = $11
		WHERE id = $12
	`
	_, err := r.db.ExecContext(ctx, query,
		course.HasHomework,
		course.IsHomeworkMandatory,
		course.IsTestMandatory,
		course.IsProjectMandatory,
		course.IsDiscordMandatory,
		course.IsAntiCopyEnabled,
		course.Status,
		course.Title,
		course.Description,
		course.ImageURL,
		course.IsMain,
		course.ID,
	)
	return err
}

func (r *ContentAdminRepoImpl) CreateModule(ctx context.Context, module *domain.Module) (string, error) {
	var newID string
	query := `
		INSERT INTO modules (course_id, title, description, order_num)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, module.CourseID, module.Title, module.Description, module.OrderNum).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("failed to create module: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error) {
	var newID string
	query := `
		INSERT INTO lessons (
			module_id, teacher_id, title, lesson_time, order_num, 
			video_url, presentation_url, content_text, is_published
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		lesson.ModuleID,
		lesson.TeacherID,
		lesson.Title,
		lesson.LessonTime,
		lesson.OrderNum,
		lesson.VideoURL,
		lesson.PresentationURL,
		lesson.ContentText,
		lesson.IsPublished,
	).Scan(&newID)

	if err != nil {
		return "", fmt.Errorf("failed to create lesson: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	query := `SELECT id, title, description, is_main, image_url, status, created_at FROM courses ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
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
	query := `
		SELECT id, title, description, is_main, image_url, status, created_at,
		       has_homework, is_homework_mandatory, is_test_mandatory, 
		       is_project_mandatory, is_discord_mandatory, is_anti_copy_enabled
		FROM courses WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Title, &c.Description, &c.IsMain, &c.ImageURL, &c.Status, &c.CreatedAt,
		&c.HasHomework, &c.IsHomeworkMandatory, &c.IsTestMandatory,
		&c.IsProjectMandatory, &c.IsDiscordMandatory, &c.IsAntiCopyEnabled,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ContentAdminRepoImpl) GetModulesByCourseID(ctx context.Context, courseID string) ([]*domain.Module, error) {
	query := `SELECT id, course_id, title, order_num, description FROM modules WHERE course_id = $1 ORDER BY order_num ASC`
	rows, err := r.db.QueryContext(ctx, query, courseID)
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

func (r *ContentAdminRepoImpl) GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error) {
	query := `
		SELECT l.id, l.module_id, l.teacher_id, l.title, l.lesson_time, l.duration_min, 
		       l.order_num, l.is_published, l.video_url, l.presentation_url, l.content_text
		FROM lessons l
		JOIN modules m ON l.module_id = m.id
		WHERE m.course_id = $1
		ORDER BY l.order_num ASC
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*domain.Lesson
	for rows.Next() {
		l := &domain.Lesson{}
		if err := rows.Scan(&l.ID, &l.ModuleID, &l.TeacherID, &l.Title, &l.LessonTime,
			&l.DurationMin, &l.OrderNum, &l.IsPublished, &l.VideoURL, &l.PresentationURL, &l.ContentText); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *ContentAdminRepoImpl) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	var newID string
	query := `
        INSERT INTO users (
            first_name, last_name, email, password_hash, role,
            phone, city, language, gender, birth_date, 
            school_name, experience_years, whatsapp_link, telegram_link, avatar_url
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query,
		user.FirstName, user.LastName, user.Email, user.Password, user.Role,
		user.Phone, user.City, user.Language, user.Gender, user.BirthDate,
		user.SchoolName, user.ExperienceYears, user.Whatsapp, user.Telegram, user.AvatarURL,
	).Scan(&newID)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	query := `
        SELECT id, first_name, last_name, email, role, created_at, 
               COALESCE(phone, ''), COALESCE(city, ''), COALESCE(school_name, ''), COALESCE(avatar_url, '')
        FROM users WHERE 1=1 
    `
	args := []interface{}{}
	argID := 1

	if filter.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argID)
		args = append(args, filter.Role)
		argID++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.CreatedAt,
			&u.Phone, &u.City, &u.SchoolName, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *ContentAdminRepoImpl) LinkParentToStudent(ctx context.Context, studentID, parentID string) error {
	query := `
        INSERT INTO child_parent_link (child_id, parent_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `
	_, err := r.db.ExecContext(ctx, query, studentID, parentID)
	return err
}

func (r *ContentAdminRepoImpl) EnrollStudent(ctx context.Context, userID, courseID string) error {
	query := `
        INSERT INTO user_courses (user_id, course_id, progress_percent)
        VALUES ($1, $2, 0)
        ON CONFLICT (user_id, course_id) DO NOTHING
    `
	_, err := r.db.ExecContext(ctx, query, userID, courseID)
	if err != nil {
		return fmt.Errorf("failed to enroll student: %w", err)
	}
	return nil
}

func (r *ContentAdminRepoImpl) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	query := `
		SELECT 
			u.id, u.first_name, u.last_name, uc.progress_percent,
			(SELECT COUNT(*) FROM user_lesson_attendance WHERE user_id = u.id AND is_attended = true) as lessons_attended,
			(SELECT COUNT(*) FROM user_assignments_submission WHERE user_id = u.id AND status = 'accepted') as homeworks_done
		FROM users u
		JOIN user_courses uc ON u.id = uc.user_id
		WHERE uc.course_id = $1 AND u.role = 'student'
	`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []*domain.AdminStudentProgress
	for rows.Next() {
		s := &domain.AdminStudentProgress{}
		if err := rows.Scan(&s.UserID, &s.FirstName, &s.LastName, &s.ProgressPercent, &s.LessonsAttended, &s.HomeworksDone); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *ContentAdminRepoImpl) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	stats := &domain.AdminCourseStats{SuccessRateBreakdown: make(map[string]int)}
	query := `
		SELECT 
			COUNT(user_id) as total_students,
			COALESCE(AVG(progress_percent), 0) as avg_progress
		FROM user_courses WHERE course_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&stats.TotalStudents, &stats.AverageScore)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *ContentAdminRepoImpl) CreateTest(ctx context.Context, test *domain.Test) (string, error) {
	var newID string
	query := `INSERT INTO tests (lesson_id, title, description, passing_score) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, test.LessonID, test.Title, test.Description, test.PassingScore).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("failed to create test: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) CreateProject(ctx context.Context, project *domain.Project) (string, error) {
	var newID string
	query := `INSERT INTO projects (lesson_id, title, description, max_score) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, project.LessonID, project.Title, project.Description, project.MaxScore).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) UpdateUser(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET
			first_name = $1, last_name = $2, email = $3, role = $4,
			phone = $5, city = $6, language = $7, gender = $8,
			school_name = $9, experience_years = $10, 
			whatsapp_link = $11, telegram_link = $12
		WHERE id = $13
	`
	_, err := r.db.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Role,
		u.Phone, u.City, u.Language, u.Gender,
		u.SchoolName, u.ExperienceYears,
		u.Whatsapp, u.Telegram, u.ID,
	)
	return err
}

func (r *ContentAdminRepoImpl) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *ContentAdminRepoImpl) GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error) {
	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.email, u.role, u.created_at,
			COALESCE(u.phone, ''), COALESCE(u.city, ''), COALESCE(u.school_name, ''), COALESCE(u.avatar_url, ''),
			COALESCE(uc.status, 'active'), 
			COALESCE(c.title, ''), 
			COALESCE(g.title, ''), 
			COALESCE(cur.first_name || ' ' || cur.last_name, ''),
			COALESCE(t.first_name || ' ' || t.last_name, ''),
			COALESCE(s.title, ''),
			COALESCE(uc.progress_percent, 0),
			COALESCE((SELECT phone FROM users pu JOIN child_parent_link cpl ON pu.id = cpl.parent_id WHERE cpl.child_id = u.id LIMIT 1), '') as parent_phone
		FROM users u
		LEFT JOIN user_courses uc ON u.id = uc.user_id
		LEFT JOIN courses c ON uc.course_id = c.id
		LEFT JOIN groups g ON uc.group_id = g.id
		LEFT JOIN users cur ON g.curator_id = cur.id
		LEFT JOIN users t ON g.teacher_id = t.id
		LEFT JOIN streams s ON uc.stream_id = s.id
		WHERE u.role = 'student'
	`
	args := []interface{}{}
	argID := 1
	if filter.CourseID != "" {
		query += fmt.Sprintf(" AND uc.course_id = $%d", argID)
		args = append(args, filter.CourseID)
		argID++
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.StudentTableItem
	for rows.Next() {
		item := &domain.StudentTableItem{}
		err := rows.Scan(
			&item.ID, &item.FirstName, &item.LastName, &item.Email, &item.Role, &item.CreatedAt,
			&item.Phone, &item.City, &item.SchoolName, &item.AvatarURL, &item.Status, &item.CourseTitle,
			&item.GroupTitle, &item.CuratorName, &item.TeacherName, &item.StreamTitle,
			&item.ProgressPercent, &item.ParentPhone,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) CreateStream(ctx context.Context, stream *domain.Stream) (string, error) {
	var newID string
	query := `INSERT INTO streams (course_id, title, start_date) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, stream.CourseID, stream.Title, stream.StartDate).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("failed to create stream: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error) {
	query := `SELECT id, course_id, title, start_date, created_at FROM streams WHERE course_id = $1`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var streams []*domain.Stream
	for rows.Next() {
		s := &domain.Stream{}
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.StartDate, &s.CreatedAt); err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}
	return streams, nil
}

func (r *ContentAdminRepoImpl) CreateGroup(ctx context.Context, group *domain.Group) (string, error) {
	var newID string
	query := `INSERT INTO groups (stream_id, curator_id, teacher_id, title) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, group.StreamID, group.CuratorID, group.TeacherID, group.Title).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("failed to create group: %w", err)
	}
	return newID, nil
}

func (r *ContentAdminRepoImpl) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	query := `SELECT id, stream_id, curator_id, teacher_id, title, created_at FROM groups WHERE stream_id = $1`
	rows, err := r.db.QueryContext(ctx, query, streamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*domain.Group
	for rows.Next() {
		g := &domain.Group{}
		if err := rows.Scan(&g.ID, &g.StreamID, &g.CuratorID, &g.TeacherID, &g.Title, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *ContentAdminRepoImpl) GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error) {
	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.email, u.role, u.created_at,
			COALESCE(u.phone, ''), COALESCE(u.city, ''), COALESCE(u.avatar_url, ''),
			COALESCE(u.experience_years, 0), COALESCE(u.language, 'ru'),
			COALESCE((SELECT STRING_AGG(title, ', ') FROM groups WHERE teacher_id = u.id), '') as group_titles
		FROM users u
		WHERE u.role = 'teacher'
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.TeacherTableItem
	for rows.Next() {
		item := &domain.TeacherTableItem{}
		err := rows.Scan(
			&item.ID, &item.FirstName, &item.LastName, &item.Email, &item.Role, &item.CreatedAt,
			&item.Phone, &item.City, &item.AvatarURL, &item.ExperienceYears, &item.Language, &item.GroupTitles,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error) {
	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.email, u.role, u.created_at,
			COALESCE(u.phone, ''), 
			COALESCE((SELECT STRING_AGG(title, ', ') FROM groups WHERE curator_id = u.id), '') as group_titles
		FROM users u
		WHERE u.role = 'curator'
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.CuratorTableItem
	for rows.Next() {
		item := &domain.CuratorTableItem{}
		err := rows.Scan(
			&item.ID, &item.FirstName, &item.LastName, &item.Email, &item.Role, &item.CreatedAt,
			&item.Phone, &item.GroupTitles,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}
