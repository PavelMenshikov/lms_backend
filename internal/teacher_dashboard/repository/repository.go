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

func (r *TeacherDashboardRepoImpl) GetMonthlyReport(ctx context.Context, teacherID string, year, month int) (*domain.TeacherMonthlyReport, error) {
	report := &domain.TeacherMonthlyReport{
		TeacherID: teacherID,
		Year:      year,
		Month:     month,
	}

	monthStart := fmt.Sprintf("%d-%02d-01", year, month)
	monthEnd := fmt.Sprintf("%d-%02d-01", year, month+1)

	prevMonth := month - 1
	prevYear := year
	if prevMonth == 0 {
		prevMonth = 12
		prevYear--
	}
	prevStart := fmt.Sprintf("%d-%02d-01", prevYear, prevMonth)
	prevEnd := fmt.Sprintf("%d-%02d-01", prevYear, prevMonth+1)

	query := `
	WITH date_range AS (
		SELECT $2::date AS start_date, $3::date AS end_date
	), prev_date_range AS (
		SELECT $4::date AS prev_start, $5::date AS prev_end
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
	), prev_substitution_counts AS (
		SELECT
			COUNT(*) FILTER (WHERE substitute_teacher_id = $1) AS sub,
			COUNT(*) FILTER (WHERE original_teacher_id = $1) AS rep
		FROM lesson_substitutions, prev_date_range
		WHERE created_at >= prev_start AND created_at < prev_end
	), cancelled_counts AS (
		SELECT COUNT(*) AS total_cancelled
		FROM lessons, date_range
		WHERE teacher_id = $1 AND lesson_time >= start_date AND lesson_time < end_date AND is_cancelled = TRUE
	), prev_cancelled_counts AS (
		SELECT COUNT(*) AS total
		FROM lessons, prev_date_range
		WHERE teacher_id = $1 AND lesson_time >= prev_start AND lesson_time < prev_end AND is_cancelled = TRUE
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
	), homework_avg AS (
		SELECT COALESCE(ROUND(AVG(uas.grade), 2), 0) AS avg_score
		FROM user_assignments_submission uas
		JOIN assignments a ON uas.assignment_id = a.id
		JOIN lessons l ON a.lesson_id = l.id, date_range
		WHERE l.teacher_id = $1 AND uas.submitted_at >= start_date AND uas.submitted_at < end_date
	)
	SELECT 
		lc.total_lessons,
		CASE WHEN EXTRACT(DAY FROM (date_trunc('month', start_date) + interval '1 month' - date_trunc('month', start_date))) > 0
			 THEN ROUND(lc.total_lessons * 7.0 / EXTRACT(DAY FROM (date_trunc('month', start_date) + interval '1 month' - date_trunc('month', start_date))), 1)
			 ELSE 0 END,
		sc.substitutions_count, sc.replaced_count,
		sc.substitutions_count - COALESCE(psc.sub, 0),
		sc.replaced_count - COALESCE(psc.rep, 0),
		rd.avg_rating, stu.total_students, ad.attendance_avg,
		ha.avg_score,
		cc.total_cancelled,
		cc.total_cancelled - COALESCE(pcc.total, 0)
	FROM lesson_counts lc, substitution_counts sc, prev_substitution_counts psc,
	     cancelled_counts cc, prev_cancelled_counts pcc, rating_data rd, student_counts stu,
	     attendance_data ad, homework_avg ha
	`

	err := r.db.QueryRowContext(ctx, query, teacherID, monthStart, monthEnd, prevStart, prevEnd).Scan(
		&report.TotalLessons,
		&report.LessonsPerWeek,
		&report.SubstitutionsCount,
		&report.ReplacedCount,
		&report.SubstitutionsDelta,
		&report.ReplacedDelta,
		&report.AvgRating,
		&report.TotalStudents,
		&report.AttendanceAvg,
		&report.AverageHomeworkScore,
		&report.TotalCancelled,
		&report.CancelledDelta,
	)
	if err != nil {
		return nil, fmt.Errorf("monthly report query: %w", err)
	}

	return report, nil
}
