package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lms_backend/internal/domain"
)

type TeacherDashboardRepository interface {
	GetMonthlyReport(ctx context.Context, teacherID string, year, month int) (*domain.TeacherMonthlyReport, error)
}

type TeacherDashboardRepoImpl struct {
	db *sql.DB
}

var _ TeacherDashboardRepository = (*TeacherDashboardRepoImpl)(nil)

func NewTeacherDashboardRepository(db *sql.DB) *TeacherDashboardRepoImpl {
	return &TeacherDashboardRepoImpl{db: db}
}

const monthlyReportQuery = `
WITH date_range AS (
    SELECT $2::date AS start_date, $3::date AS end_date
), lesson_counts AS (
    SELECT COUNT(*) AS total_lessons
    FROM lessons, date_range
    WHERE (teacher_id = $1 OR substituted_teacher_id = $1)
      AND lesson_time >= start_date AND lesson_time < end_date
      AND is_cancelled = FALSE
), substitution_counts AS (
    SELECT
        COUNT(*) FILTER (WHERE substitute_teacher_id = $1) AS substitutions_count,
        COUNT(*) FILTER (WHERE original_teacher_id = $1) AS replaced_count
    FROM lesson_substitutions, date_range
    WHERE created_at >= start_date AND created_at < end_date
), rating_data AS (
    SELECT COALESCE(ROUND(AVG(rating), 2), 0) AS avg_rating
    FROM teacher_reviews, date_range
    WHERE teacher_id = $1 AND created_at >= start_date AND created_at < end_date
), student_counts AS (
    SELECT COUNT(DISTINCT uc.user_id) AS total_students
    FROM user_courses uc
    JOIN groups g ON uc.group_id = g.id
    WHERE g.teacher_id = $1
), attendance_data AS (
    SELECT COALESCE(ROUND(AVG(subq.pct), 2), 0) AS attendance_avg FROM (
        SELECT
            COUNT(CASE WHEN ula.status IN ('visited', 'trial') THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) AS pct
        FROM lessons l
        JOIN user_lesson_attendance ula ON ula.lesson_id = l.id, date_range
        WHERE (l.teacher_id = $1 OR l.substituted_teacher_id = $1)
          AND l.lesson_time >= start_date AND l.lesson_time < end_date
          AND l.is_cancelled = FALSE
        GROUP BY ula.user_id
    ) subq
)
SELECT total_lessons, substitutions_count, replaced_count,
       avg_rating, total_students, attendance_avg
FROM lesson_counts, substitution_counts, rating_data, student_counts, attendance_data`

func (r *TeacherDashboardRepoImpl) GetMonthlyReport(ctx context.Context, teacherID string, year, month int) (*domain.TeacherMonthlyReport, error) {
	report := &domain.TeacherMonthlyReport{
		TeacherID: teacherID,
		Year:      year,
		Month:     month,
	}

	monthStart := fmt.Sprintf("%d-%02d-01", year, month)
	monthEnd := fmt.Sprintf("%d-%02d-01", year, month+1)

	err := r.db.QueryRowContext(ctx, monthlyReportQuery, teacherID, monthStart, monthEnd).Scan(
		&report.TotalLessons,
		&report.SubstitutionsCount,
		&report.ReplacedCount,
		&report.AvgRating,
		&report.TotalStudents,
		&report.AttendanceAvg,
	)
	if err != nil {
		return nil, fmt.Errorf("monthly report query: %w", err)
	}

	return report, nil
}
