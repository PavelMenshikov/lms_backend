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
	DeleteModule(ctx context.Context, id string) error
	CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error)
	DeleteLesson(ctx context.Context, id string) error
	AssignTeacherToLesson(ctx context.Context, lessonID, teacherID string) error
	GetLessonIDByOrder(ctx context.Context, courseID string, orderNum int) (string, error)
	GetAllCourses(ctx context.Context) ([]*domain.Course, error)
	GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
	GetModulesByCourseID(ctx context.Context, courseID string) ([]*domain.Module, error)
	GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error)
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetParentsByStudentID(ctx context.Context, studentID string) ([]domain.User, error)
	LinkParentToStudent(ctx context.Context, studentID, parentID string) error
	EnrollStudentExtended(ctx context.Context, userID, courseID, streamID, groupID string) error
	GetCourseIDByStream(ctx context.Context, streamID string) (string, error)
	GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error)
	GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error)
	CreateTest(ctx context.Context, test *domain.Test) (string, error)
	DeleteTest(ctx context.Context, id string) error
	CreateProject(ctx context.Context, project *domain.Project) (string, error)
	DeleteProject(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID string) error
	GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error)
	GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error)
	GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error)
	GetDetailedModeratorList(ctx context.Context) ([]*domain.ModeratorTableItem, error)
	GetAllUsersList(ctx context.Context) ([]*domain.AllUsersTableItem, error)
	CreateStream(ctx context.Context, stream *domain.Stream) (string, error)
	GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error)
	CreateGroup(ctx context.Context, group *domain.Group) (string, error)
	GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error)
	GetStudentEnrollment(ctx context.Context, userID string) (map[string]string, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
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

func (r *ContentAdminRepoImpl) CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error) {
	var newID string
	var tid sql.NullString
	if lesson.TeacherID != "" {
		tid = sql.NullString{String: lesson.TeacherID, Valid: true}
	}
	query := `INSERT INTO lessons (course_id, module_id, teacher_id, title, lesson_time, order_num, video_url, presentation_url, content_text, is_published)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, lesson.CourseID, lesson.ModuleID, tid, lesson.Title, lesson.LessonTime, lesson.OrderNum, lesson.VideoURL, lesson.PresentationURL, lesson.ContentText, lesson.IsPublished).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) DeleteLesson(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM lessons WHERE id = $1", id)
	return err
}

func (r *ContentAdminRepoImpl) AssignTeacherToLesson(ctx context.Context, lessonID, teacherID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE lessons SET teacher_id = $1 WHERE id = $2", teacherID, lessonID)
	return err
}

func (r *ContentAdminRepoImpl) GetLessonIDByOrder(ctx context.Context, courseID string, orderNum int) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM lessons WHERE course_id = $1 AND order_num = $2", courseID, orderNum).Scan(&id)
	return id, err
}

func (r *ContentAdminRepoImpl) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description, is_main, image_url, status, created_at FROM courses ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []*domain.Course
	for rows.Next() {
		c := &domain.Course{}
		rows.Scan(&c.ID, &c.Title, &c.Description, &c.IsMain, &c.ImageURL, &c.Status, &c.CreatedAt)
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
		rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.OrderNum, &m.Description)
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *ContentAdminRepoImpl) GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, course_id, module_id, teacher_id, title, lesson_time, duration_min, order_num, is_published, video_url, presentation_url, content_text FROM lessons WHERE course_id = $1 ORDER BY order_num ASC", courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lessons []*domain.Lesson
	for rows.Next() {
		l := &domain.Lesson{}
		var mid, tid sql.NullString
		rows.Scan(&l.ID, &l.CourseID, &mid, &tid, &l.Title, &l.LessonTime, &l.DurationMin, &l.OrderNum, &l.IsPublished, &l.VideoURL, &l.PresentationURL, &l.ContentText)
		if mid.Valid {
			s := mid.String
			l.ModuleID = &s
		}
		if tid.Valid {
			l.TeacherID = tid.String
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *ContentAdminRepoImpl) CreateUser(ctx context.Context, u *domain.User) (string, error) {
	var newID string
	query := `INSERT INTO users (first_name, last_name, email, password_hash, role, phone, city, language, gender, birth_date, school_name, experience_years, whatsapp_link, telegram_link, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role, u.Phone, u.City, u.Language, u.Gender, u.BirthDate, u.SchoolName, u.ExperienceYears, u.Whatsapp, u.Telegram, u.AvatarURL).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	query := `
        SELECT id, first_name, last_name, first_name || ' ' || last_name as full_name, email, role, created_at, 
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
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.FullName, &u.Email, &u.Role, &u.CreatedAt,
			&u.Phone, &u.City, &u.SchoolName, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *ContentAdminRepoImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT 
			id, first_name, last_name, first_name || ' ' || last_name as full_name, 
			email, role, created_at, COALESCE(phone, ''), COALESCE(city, ''), 
			COALESCE(school_name, ''), COALESCE(language, ''), COALESCE(gender, ''), 
			COALESCE(birth_date, NOW()), COALESCE(experience_years, 0), 
			COALESCE(whatsapp_link, ''), COALESCE(telegram_link, ''), COALESCE(avatar_url, '') 
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.FullName, &u.Email, &u.Role, &u.CreatedAt,
		&u.Phone, &u.City, &u.SchoolName, &u.Language, &u.Gender, &u.BirthDate,
		&u.ExperienceYears, &u.Whatsapp, &u.Telegram, &u.AvatarURL,
	)
	return u, err
}

func (r *ContentAdminRepoImpl) GetParentsByStudentID(ctx context.Context, studentID string) ([]domain.User, error) {
	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.first_name || ' ' || u.last_name as full_name, 
			u.email, u.phone 
		FROM users u 
		JOIN child_parent_link cpl ON u.id = cpl.parent_id 
		WHERE cpl.child_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parents []domain.User
	for rows.Next() {
		u := domain.User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.FullName, &u.Email, &u.Phone); err != nil {
			return nil, err
		}
		parents = append(parents, u)
	}
	return parents, nil
}

func (r *ContentAdminRepoImpl) LinkParentToStudent(ctx context.Context, studentID, parentID string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO child_parent_link (child_id, parent_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", studentID, parentID)
	return err
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

func (r *ContentAdminRepoImpl) GetCourseIDByStream(ctx context.Context, streamID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT course_id FROM streams WHERE id = $1", streamID).Scan(&id)
	return id, err
}

func (r *ContentAdminRepoImpl) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	query := `SELECT u.id, u.first_name || ' ' || u.last_name as full_name, uc.progress_percent, (SELECT COUNT(*) FROM user_lesson_attendance WHERE user_id = u.id AND is_attended = true), (SELECT COUNT(*) FROM user_assignments_submission WHERE user_id = u.id AND status = 'accepted') FROM users u JOIN user_courses uc ON u.id = uc.user_id WHERE uc.course_id = $1 AND u.role = 'student'`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.AdminStudentProgress
	for rows.Next() {
		s := &domain.AdminStudentProgress{}
		rows.Scan(&s.UserID, &s.FullName, &s.ProgressPercent, &s.LessonsAttended, &s.HomeworksDone)
		list = append(list, s)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	stats := &domain.AdminCourseStats{}
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(user_id), COALESCE(AVG(progress_percent), 0) FROM user_courses WHERE course_id = $1", courseID).Scan(&stats.TotalStudents, &stats.AverageScore)
	return stats, err
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

func (r *ContentAdminRepoImpl) CreateProject(ctx context.Context, project *domain.Project) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO projects (lesson_id, title, description, max_score) VALUES ($1, $2, $3, $4) RETURNING id", project.LessonID, project.Title, project.Description, project.MaxScore).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) DeleteProject(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}

func (r *ContentAdminRepoImpl) UpdateUser(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET
			first_name = $1, last_name = $2, email = $3, role = $4,
			phone = $5, city = $6, language = $7, gender = $8,
			school_name = $9, experience_years = $10, 
			whatsapp_link = $11, telegram_link = $12,
			birth_date = $13
		WHERE id = $14
	`
	_, err := r.db.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Role,
		u.Phone, u.City, u.Language, u.Gender,
		u.SchoolName, u.ExperienceYears,
		u.Whatsapp, u.Telegram, u.BirthDate,
		u.ID,
	)
	return err
}

func (r *ContentAdminRepoImpl) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}

func (r *ContentAdminRepoImpl) GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error) {
	query := `
		SELECT 
			u.id, COALESCE(u.avatar_url, ''), u.first_name || ' ' || u.last_name, u.created_at, u.gender,
			COALESCE(EXTRACT(YEAR FROM AGE(u.birth_date)), 0) as age,
			COALESCE(uc.status, 'inactive'), COALESCE(c.title, ''), COALESCE(g.title, ''),
			COALESCE(cur.first_name || ' ' || cur.last_name, ''),
			COALESCE(teach.first_name || ' ' || teach.last_name, ''),
			COALESCE(s.title, ''), COALESCE(uc.progress_percent, 0),
			COALESCE((SELECT phone FROM users pu JOIN child_parent_link cpl ON pu.id = cpl.parent_id WHERE cpl.child_id = u.id LIMIT 1), ''),
			COALESCE(u.city, ''), COALESCE(u.school_name, ''), COALESCE(u.language, ''),
			COALESCE(u.phone, ''), u.email
		FROM users u
		LEFT JOIN user_courses uc ON u.id = uc.user_id
		LEFT JOIN courses c ON uc.course_id = c.id
		LEFT JOIN groups g ON uc.group_id = g.id
		LEFT JOIN users cur ON g.curator_id = cur.id
		LEFT JOIN users teach ON g.teacher_id = teach.id
		LEFT JOIN streams s ON uc.stream_id = s.id
		WHERE u.role = 'student'`

	args := []interface{}{}
	if filter.CourseID != "" {
		query += " AND uc.course_id = $1"
		args = append(args, filter.CourseID)
	}
	query += " ORDER BY u.created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.StudentTableItem
	for rows.Next() {
		item := &domain.StudentTableItem{}
		rows.Scan(&item.ID, &item.Photo, &item.FullName, &item.CreatedAt, &item.Gender, &item.Age, &item.Status, &item.Course, &item.Group, &item.Curator, &item.Teacher, &item.Stream, &item.Performance, &item.ParentPhone, &item.City, &item.School, &item.Language, &item.Phone, &item.Email)
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error) {
	query := `SELECT u.id, COALESCE(u.avatar_url, ''), u.first_name || ' ' || u.last_name, u.created_at, u.gender, u.phone, COALESCE((SELECT STRING_AGG(title, ', ') FROM groups WHERE teacher_id = u.id), ''), COALESCE(u.city, ''), u.email, COALESCE(u.experience_years, 0), COALESCE(u.language, '') FROM users u WHERE u.role = 'teacher' ORDER BY u.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.TeacherTableItem
	for rows.Next() {
		item := &domain.TeacherTableItem{}
		rows.Scan(&item.ID, &item.Photo, &item.FullName, &item.CreatedAt, &item.Gender, &item.Phone, &item.Groups, &item.City, &item.Email, &item.ExperienceYears, &item.Language)
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error) {
	query := `SELECT u.id, u.first_name || ' ' || u.last_name, u.created_at, COALESCE((SELECT STRING_AGG(title, ', ') FROM groups WHERE curator_id = u.id), ''), u.phone, u.email FROM users u WHERE u.role = 'curator' ORDER BY u.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.CuratorTableItem
	for rows.Next() {
		item := &domain.CuratorTableItem{}
		rows.Scan(&item.ID, &item.FullName, &item.CreatedAt, &item.Groups, &item.Phone, &item.Email)
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedModeratorList(ctx context.Context) ([]*domain.ModeratorTableItem, error) {
	query := `SELECT id, first_name || ' ' || last_name, created_at, phone, email FROM users WHERE role = 'moderator' ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.ModeratorTableItem
	for rows.Next() {
		item := &domain.ModeratorTableItem{}
		rows.Scan(&item.ID, &item.FullName, &item.CreatedAt, &item.Phone, &item.Email)
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetAllUsersList(ctx context.Context) ([]*domain.AllUsersTableItem, error) {
	query := `SELECT id, COALESCE(avatar_url, ''), first_name || ' ' || last_name, role, created_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.AllUsersTableItem
	for rows.Next() {
		item := &domain.AllUsersTableItem{}
		rows.Scan(&item.ID, &item.Photo, &item.FullName, &item.Role, &item.CreatedAt)
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) CreateStream(ctx context.Context, s *domain.Stream) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO streams (course_id, title, start_date) VALUES ($1, $2, $3) RETURNING id", s.CourseID, s.Title, s.StartDate).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error) {
	query := `SELECT id, course_id, title, start_date, created_at FROM streams`
	args := []interface{}{}
	if courseID != "" {
		query += " WHERE course_id = $1"
		args = append(args, courseID)
	}
	query += " ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var streams []*domain.Stream
	for rows.Next() {
		s := &domain.Stream{}
		rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.StartDate, &s.CreatedAt)
		streams = append(streams, s)
	}
	return streams, nil
}

func (r *ContentAdminRepoImpl) CreateGroup(ctx context.Context, g *domain.Group) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO groups (stream_id, curator_id, teacher_id, title) VALUES ($1, $2, $3, $4) RETURNING id", g.StreamID, g.CuratorID, g.TeacherID, g.Title).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	query := `SELECT id, stream_id, curator_id, teacher_id, title, created_at FROM groups`
	args := []interface{}{}
	if streamID != "" {
		query += " WHERE stream_id = $1"
		args = append(args, streamID)
	}
	query += " ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*domain.Group
	for rows.Next() {
		g := &domain.Group{}
		rows.Scan(&g.ID, &g.StreamID, &g.CuratorID, &g.TeacherID, &g.Title, &g.CreatedAt)
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *ContentAdminRepoImpl) GetStudentEnrollment(ctx context.Context, userID string) (map[string]string, error) {
	var cid, sid, gid sql.NullString
	query := `SELECT course_id, stream_id, group_id FROM user_courses WHERE user_id = $1 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&cid, &sid, &gid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	res := make(map[string]string)
	if cid.Valid { res["course_id"] = cid.String }
	if sid.Valid { res["stream_id"] = sid.String }
	if gid.Valid { res["group_id"] = gid.String }

	return res, nil
}

func (r *ContentAdminRepoImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, first_name, last_name, email FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
	return u, err
}