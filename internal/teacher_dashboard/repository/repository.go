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

	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM lessons
		WHERE (teacher_id = $1 OR substituted_teacher_id = $1)
		  AND lesson_time >= $2::date AND lesson_time < $3::date
		  AND is_cancelled = FALSE
	`, teacherID, monthStart, monthEnd).Scan(&report.TotalLessons)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM lesson_substitutions
		WHERE substitute_teacher_id = $1
		  AND created_at >= $2::date AND created_at < $3::date
	`, teacherID, monthStart, monthEnd).Scan(&report.SubstitutionsCount)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM lesson_substitutions
		WHERE original_teacher_id = $1
		  AND created_at >= $2::date AND created_at < $3::date
	`, teacherID, monthStart, monthEnd).Scan(&report.ReplacedCount)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(ROUND(AVG(rating), 2), 0) FROM teacher_reviews
		WHERE teacher_id = $1
		  AND created_at >= $2::date AND created_at < $3::date
	`, teacherID, monthStart, monthEnd).Scan(&report.AvgRating)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT uc.user_id)
		FROM user_courses uc
		JOIN groups g ON uc.group_id = g.id
		WHERE g.teacher_id = $1
	`, teacherID).Scan(&report.TotalStudents)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(ROUND(AVG(subq.pct), 2), 0) FROM (
			SELECT
				COUNT(CASE WHEN ula.status IN ('visited', 'trial') THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0) as pct
			FROM lessons l
			JOIN user_lesson_attendance ula ON ula.lesson_id = l.id
			WHERE (l.teacher_id = $1 OR l.substituted_teacher_id = $1)
			  AND l.lesson_time >= $2::date AND l.lesson_time < $3::date
			  AND l.is_cancelled = FALSE
			GROUP BY ula.user_id
		) subq
	`, teacherID, monthStart, monthEnd).Scan(&report.AttendanceAvg)

	return report, nil
}
