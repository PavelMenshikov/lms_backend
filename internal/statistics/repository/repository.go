package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type StatisticsRepository interface {
	GetByStudent(ctx context.Context, studentID string) (*domain.StudentStatistics, error)
	UpdateStatistics(ctx context.Context, studentID string) error
	RecalculateStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error)
}

type statisticsRepository struct {
	db *sql.DB
}

func NewStatisticsRepository(db *sql.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) GetByStudent(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
	var stats domain.StudentStatistics
	query := `
		SELECT id, student_id, total_lessons, attended_lessons, absent_excused, absent_unexcused,
		       freeze_days, remaining_lessons, remaining_excused, last_attendance_date,
		       current_freeze_end_date, updated_at, created_at
		FROM student_statistics
		WHERE student_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&stats.ID, &stats.StudentID, &stats.TotalLessons, &stats.AttendedLessons,
		&stats.AbsentExcused, &stats.AbsentUnexcused, &stats.FreezeDays,
		&stats.RemainingLessons, &stats.RemainingExcused, &stats.LastAttendanceDate,
		&stats.CurrentFreezeEndDate, &stats.UpdatedAt, &stats.CreatedAt,
	)
	if err == sql.ErrNoRows {
		// Если статистики нет, создаём её
		return r.RecalculateStatistics(ctx, studentID)
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *statisticsRepository) UpdateStatistics(ctx context.Context, studentID string) error {
	query := `
		INSERT INTO student_statistics (student_id, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP)
		ON CONFLICT (student_id)
		DO UPDATE SET updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.ExecContext(ctx, query, studentID)
	return err
}

func (r *statisticsRepository) RecalculateStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
	// Пересчитываем статистику на основе attendance_records
	query := `
		WITH stats AS (
			SELECT
				COUNT(*) as total_lessons,
				COUNT(*) FILTER (WHERE status = 'ATTENDED') as attended_lessons,
				COUNT(*) FILTER (WHERE status = 'ABSENT_EXCUSED') as absent_excused,
				COUNT(*) FILTER (WHERE status = 'ABSENT_UNEXCUSED') as absent_unexcused,
				MAX(marked_at)::date as last_attendance_date
			FROM attendance_records
			WHERE student_id = $1
		),
		freeze_info AS (
			SELECT
				COALESCE(SUM(end_date - start_date + 1), 0) as freeze_days,
				MAX(end_date) FILTER (WHERE is_active = true AND end_date >= CURRENT_DATE) as current_freeze_end_date
			FROM freeze_periods
			WHERE student_id = $1
		)
		INSERT INTO student_statistics (
			id, student_id, total_lessons, attended_lessons, absent_excused, absent_unexcused,
			freeze_days, remaining_lessons, remaining_excused, last_attendance_date, current_freeze_end_date
		)
		SELECT
			gen_random_uuid(),
			$1,
			COALESCE(s.total_lessons, 0),
			COALESCE(s.attended_lessons, 0),
			COALESCE(s.absent_excused, 0),
			COALESCE(s.absent_unexcused, 0),
			COALESCE(f.freeze_days, 0),
			0, -- remaining_lessons (нужна бизнес-логика)
			0, -- remaining_excused (нужна бизнес-логика)
			s.last_attendance_date,
			f.current_freeze_end_date
		FROM stats s, freeze_info f
		ON CONFLICT (student_id)
		DO UPDATE SET
			total_lessons = EXCLUDED.total_lessons,
			attended_lessons = EXCLUDED.attended_lessons,
			absent_excused = EXCLUDED.absent_excused,
			absent_unexcused = EXCLUDED.absent_unexcused,
			freeze_days = EXCLUDED.freeze_days,
			last_attendance_date = EXCLUDED.last_attendance_date,
			current_freeze_end_date = EXCLUDED.current_freeze_end_date,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, student_id, total_lessons, attended_lessons, absent_excused, absent_unexcused,
		          freeze_days, remaining_lessons, remaining_excused, last_attendance_date,
		          current_freeze_end_date, updated_at, created_at
	`

	var stats domain.StudentStatistics
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&stats.ID, &stats.StudentID, &stats.TotalLessons, &stats.AttendedLessons,
		&stats.AbsentExcused, &stats.AbsentUnexcused, &stats.FreezeDays,
		&stats.RemainingLessons, &stats.RemainingExcused, &stats.LastAttendanceDate,
		&stats.CurrentFreezeEndDate, &stats.UpdatedAt, &stats.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
