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
			COALESCE(uas.status, 'not_started') as assignment_status
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
	query := `SELECT COUNT(*) FROM user_courses WHERE user_id = $1 AND is_active = true`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *DashboardRepositoryImpl) GetAttendancePercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	stats := &domain.StatisticSummary{}
	query := `
		SELECT 
			ROUND(COUNT(CASE WHEN is_attended = true THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 2) as percentage
		FROM user_lesson_attendance 
		WHERE user_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&stats.Percentage)
	if err == sql.ErrNoRows {
		return &domain.StatisticSummary{Percentage: 0}, nil
	}
	return stats, err
}

func (r *DashboardRepositoryImpl) GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	stats := &domain.StatisticSummary{Breakdown: make(map[string]int)}
	query := `
		SELECT 
			ROUND(COUNT(CASE WHEN status = 'accepted' THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 2) as percentage
		FROM user_assignments_submission
		WHERE user_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&stats.Percentage)
	if err == sql.ErrNoRows {
		return &domain.StatisticSummary{Percentage: 0}, nil
	}
	return stats, err
}

func (r *DashboardRepositoryImpl) GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error) {
	query := `
		SELECT 
			l.lesson_time, 
			u.first_name || ' ' || u.last_name as teacher_name,
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