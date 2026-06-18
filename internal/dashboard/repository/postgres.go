package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type DashboardRepositoryImpl struct {
	db *sql.DB
}

var _ UserDataRepository = (*DashboardRepositoryImpl)(nil)

func NewDashboardRepository(db *sql.DB) *DashboardRepositoryImpl {
	return &DashboardRepositoryImpl{db: db}
}

func (r *DashboardRepositoryImpl) GetLastLessonData(ctx context.Context, userID string) (*domain.LastLesson, error) {
	lesson := &domain.LastLesson{}
	query := `
		SELECT 
			c.title, m.title, l.title, l.id,
			COALESCE(uas.status::text, 'not_started') as assignment_status
		FROM user_courses uc
		JOIN courses c ON uc.course_id = c.id
		JOIN modules m ON m.course_id = c.id
		JOIN lessons l ON l.module_id = m.id
		LEFT JOIN assignments a ON a.lesson_id = l.id
		LEFT JOIN user_assignments_submission uas ON uas.assignment_id = a.id AND uas.user_id = $1
		WHERE uc.user_id = $1
		ORDER BY l.lesson_time DESC
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&lesson.CourseTitle, &lesson.ModuleName, &lesson.LessonTitle,
		&lesson.LessonID, &lesson.AssignmentStatus,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return lesson, err
}

func (r *DashboardRepositoryImpl) GetActiveCoursesCount(ctx context.Context, userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_courses WHERE user_id = $1 AND status = 'active'`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *DashboardRepositoryImpl) GetAttendancePercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	stats := &domain.StatisticSummary{}
	query := `
		SELECT 
			COALESCE(ROUND(COUNT(CASE WHEN status IN ('visited', 'trial') THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 2), 0) as percentage
		FROM user_lesson_attendance 
		WHERE user_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&stats.Percentage)
	return stats, err
}

func (r *DashboardRepositoryImpl) GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	stats := &domain.StatisticSummary{Breakdown: make(map[string]int)}
	query := `
		SELECT 
			COALESCE(ROUND(COUNT(CASE WHEN status = 'accepted' THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 2), 0) as percentage
		FROM user_assignments_submission
		WHERE user_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&stats.Percentage)
	return stats, err
}

func (r *DashboardRepositoryImpl) GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error) {
	query := `
		SELECT 
			l.lesson_time, 
			COALESCE(u.first_name || ' ' || u.last_name, '') as teacher_name,
			c.title as course_title
		FROM user_courses uc
		JOIN courses c ON uc.course_id = c.id
		JOIN modules m ON m.course_id = c.id
		JOIN lessons l ON l.module_id = m.id
		JOIN users u ON l.teacher_id = u.id
		WHERE uc.user_id = $1 AND l.lesson_time > NOW()
		ORDER BY l.lesson_time ASC
		LIMIT 5
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []domain.UpcomingLesson
	for rows.Next() {
		var l domain.UpcomingLesson
		if err := rows.Scan(&l.Date, &l.TeacherName, &l.CourseTitle); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *DashboardRepositoryImpl) GetAdminCounters(ctx context.Context) (totalStudents, newStudents int, studentsDelta float64, totalTeachers, activeCourses int, err error) {
	query := `
		WITH stats AS (
			SELECT
				COUNT(*) FILTER (WHERE role = 'student') as total_s,
				COUNT(*) FILTER (WHERE role = 'student' AND created_at >= date_trunc('month', now())) as new_s,
				COUNT(*) FILTER (WHERE role = 'student' AND created_at >= date_trunc('month', now() - interval '1 month') AND created_at < date_trunc('month', now())) as prev_month_s,
				COUNT(*) FILTER (WHERE role = 'teacher') as total_t,
				(SELECT COUNT(*) FROM courses WHERE status = 'active') as active_c
			FROM users
		)
		SELECT 
			total_s, 
			new_s, 
			CASE 
				WHEN prev_month_s = 0 THEN 100 
				ELSE ROUND(((new_s::numeric - prev_month_s::numeric) / prev_month_s::numeric) * 100, 1) 
			END as delta,
			total_t, 
			active_c
		FROM stats
	`
	err = r.db.QueryRowContext(ctx, query).Scan(&totalStudents, &newStudents, &studentsDelta, &totalTeachers, &activeCourses)
	return
}

func (r *DashboardRepositoryImpl) GetAllPerformanceStats(ctx context.Context) (*domain.AllPerformanceStats, error) {
	result := &domain.AllPerformanceStats{}
	query := `
		WITH student_scores AS (
			SELECT user_id, AVG(progress_percent) as score
			FROM user_courses
			GROUP BY user_id
		),
		student_hw_scores AS (
			SELECT uas.user_id,
				COUNT(CASE WHEN uas.status = 'accepted' THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as score
			FROM user_assignments_submission uas
			GROUP BY uas.user_id
		),
		student_att_scores AS (
			SELECT ula.user_id,
				COUNT(CASE WHEN ula.status IN ('visited', 'trial') THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as score
			FROM user_lesson_attendance ula
			GROUP BY ula.user_id
		)
		SELECT
			COALESCE((SELECT COUNT(CASE WHEN score >= 80 THEN 1 END) FROM student_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) FROM student_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score < 50 THEN 1 END) FROM student_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score >= 80 THEN 1 END) FROM student_hw_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) FROM student_hw_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score < 50 THEN 1 END) FROM student_hw_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score >= 80 THEN 1 END) FROM student_att_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) FROM student_att_scores), 0),
			COALESCE((SELECT COUNT(CASE WHEN score < 50 THEN 1 END) FROM student_att_scores), 0)
	`
	err := r.db.QueryRowContext(ctx, query).Scan(
		&result.CourseZones.Green, &result.CourseZones.Yellow, &result.CourseZones.Red,
		&result.HomeworkZones.Green, &result.HomeworkZones.Yellow, &result.HomeworkZones.Red,
		&result.AttendanceZones.Green, &result.AttendanceZones.Yellow, &result.AttendanceZones.Red,
	)
	return result, err
}

func (r *DashboardRepositoryImpl) GetPerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	var zones domain.PerformanceZones
	query := `
		WITH student_scores AS (
			SELECT user_id, AVG(progress_percent) as score
			FROM user_courses
			GROUP BY user_id
		)
		SELECT
			COUNT(CASE WHEN score >= 80 THEN 1 END) as green,
			COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) as yellow,
			COUNT(CASE WHEN score < 50 THEN 1 END) as red
		FROM student_scores
	`
	err := r.db.QueryRowContext(ctx, query).Scan(&zones.Green, &zones.Yellow, &zones.Red)
	return zones, err
}

func (r *DashboardRepositoryImpl) GetHwPerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	var zones domain.PerformanceZones
	query := `
		WITH student_hw_scores AS (
			SELECT uas.user_id,
				COUNT(CASE WHEN uas.status = 'accepted' THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as score
			FROM user_assignments_submission uas
			GROUP BY uas.user_id
		)
		SELECT
			COUNT(CASE WHEN score >= 80 THEN 1 END) as green,
			COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) as yellow,
			COUNT(CASE WHEN score < 50 THEN 1 END) as red
		FROM student_hw_scores
	`
	err := r.db.QueryRowContext(ctx, query).Scan(&zones.Green, &zones.Yellow, &zones.Red)
	return zones, err
}

func (r *DashboardRepositoryImpl) GetAttendancePerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	var zones domain.PerformanceZones
	query := `
		WITH student_att_scores AS (
			SELECT ula.user_id,
				COUNT(CASE WHEN ula.status IN ('visited', 'trial') THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as score
			FROM user_lesson_attendance ula
			GROUP BY ula.user_id
		)
		SELECT
			COUNT(CASE WHEN score >= 80 THEN 1 END) as green,
			COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) as yellow,
			COUNT(CASE WHEN score < 50 THEN 1 END) as red
		FROM student_att_scores
	`
	err := r.db.QueryRowContext(ctx, query).Scan(&zones.Green, &zones.Yellow, &zones.Red)
	return zones, err
}

func (r *DashboardRepositoryImpl) GetCuratorGroups(ctx context.Context, curatorID string) ([]domain.Group, error) {
	query := `
		SELECT id, stream_id, COALESCE(curator_id, ''), COALESCE(teacher_id, ''), title, created_at
		FROM groups
		WHERE curator_id = $1
		ORDER BY title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, curatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []domain.Group
	for rows.Next() {
		var g domain.Group
		if err := rows.Scan(&g.ID, &g.StreamID, &g.CuratorID, &g.TeacherID, &g.Title, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *DashboardRepositoryImpl) GetCuratorAttendanceStats(ctx context.Context, curatorID string) ([]domain.CuratorGroupAttendance, error) {
	query := `
		WITH group_students AS (
			SELECT g.id as group_id, g.title as group_title, uc.user_id
			FROM groups g
			JOIN user_courses uc ON uc.group_id = g.id
			WHERE g.curator_id = $1
		),
		student_attendance AS (
			SELECT user_id,
				COUNT(CASE WHEN status IN ('visited', 'trial') THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as pct
			FROM user_lesson_attendance
			WHERE user_id IN (SELECT user_id FROM group_students)
			GROUP BY user_id
		)
		SELECT
			gs.group_id,
			gs.group_title,
			COUNT(DISTINCT gs.user_id) as student_count,
			COALESCE(ROUND(AVG(sa.pct), 2), 0) as avg_attendance
		FROM group_students gs
		LEFT JOIN student_attendance sa ON gs.user_id = sa.user_id
		GROUP BY gs.group_id, gs.group_title
		ORDER BY gs.group_title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, curatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []domain.CuratorGroupAttendance
	for rows.Next() {
		var s domain.CuratorGroupAttendance
		if err := rows.Scan(&s.GroupID, &s.GroupTitle, &s.StudentCount, &s.AvgAttendance); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *DashboardRepositoryImpl) GetCuratorHomeworkStats(ctx context.Context, curatorID string) ([]domain.CuratorHomeworkStats, error) {
	query := `
		WITH group_homework AS (
			SELECT g.id as group_id, g.title as group_title, uas.user_id, uas.status
			FROM groups g
			JOIN user_courses uc ON uc.group_id = g.id
			JOIN user_assignments_submission uas ON uas.user_id = uc.user_id
			WHERE g.curator_id = $1
		)
		SELECT
			gh.group_id,
			gh.group_title,
			COUNT(*) as total_submitted,
			COUNT(CASE WHEN gh.status = 'accepted' THEN 1 END) as total_accepted,
			COALESCE(ROUND(COUNT(CASE WHEN gh.status = 'accepted' THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 2), 0) as avg_completion
		FROM group_homework gh
		GROUP BY gh.group_id, gh.group_title
		ORDER BY gh.group_title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, curatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []domain.CuratorHomeworkStats
	for rows.Next() {
		var s domain.CuratorHomeworkStats
		if err := rows.Scan(&s.GroupID, &s.GroupTitle, &s.TotalSubmitted, &s.TotalAccepted, &s.AvgCompletion); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *DashboardRepositoryImpl) GetCuratorPerformanceZones(ctx context.Context, curatorID string) (domain.PerformanceZones, error) {
	var zones domain.PerformanceZones
	query := `
		WITH curator_students AS (
			SELECT DISTINCT uc.user_id, AVG(uc.progress_percent) as score
			FROM groups g
			JOIN user_courses uc ON uc.group_id = g.id
			WHERE g.curator_id = $1
			GROUP BY uc.user_id
		)
		SELECT
			COUNT(CASE WHEN score >= 80 THEN 1 END) as green,
			COUNT(CASE WHEN score >= 50 AND score < 80 THEN 1 END) as yellow,
			COUNT(CASE WHEN score < 50 THEN 1 END) as red
		FROM curator_students
	`
	err := r.db.QueryRowContext(ctx, query, curatorID).Scan(&zones.Green, &zones.Yellow, &zones.Red)
	return zones, err
}

func (r *DashboardRepositoryImpl) GetLessonActivity(ctx context.Context) ([]domain.DailyLessonActivity, error) {
	query := `
		SELECT 
			TO_CHAR(lesson_time, 'YYYY-MM-DD') as day,
			COUNT(CASE WHEN duration_min > 60 THEN 1 END) as group_lessons,
			COUNT(CASE WHEN duration_min = 30 THEN 1 END) as trial_lessons,
			COUNT(CASE WHEN duration_min <= 60 AND duration_min > 30 THEN 1 END) as individual_lessons
		FROM lessons
		WHERE lesson_time > date_trunc('month', now())
		GROUP BY day
		ORDER BY day ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []domain.DailyLessonActivity
	for rows.Next() {
		var a domain.DailyLessonActivity
		if err := rows.Scan(&a.Date, &a.Group, &a.Trial, &a.Individual); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}
