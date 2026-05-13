package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
	"time"
)

type AttendanceRepository interface {
	GetByLessonAndStudent(ctx context.Context, lessonID, studentID string) (*domain.AttendanceRecord, error)
	GetByStudent(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*domain.AttendanceRecord, error)
	GetByLesson(ctx context.Context, lessonID string) ([]*domain.AttendanceRecord, error)
	Create(ctx context.Context, record *domain.AttendanceRecord) error
	Update(ctx context.Context, record *domain.AttendanceRecord) error
	GetStudentStats(ctx context.Context, studentID string) (map[string]int, error)
}

type attendanceRepository struct {
	db *sql.DB
}

func NewAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) GetByLessonAndStudent(ctx context.Context, lessonID, studentID string) (*domain.AttendanceRecord, error) {
	var record domain.AttendanceRecord
	query := `
		SELECT id, lesson_id, student_id, status, reason, comment, marked_by, marked_at, updated_by, updated_at, created_at
		FROM attendance_records
		WHERE lesson_id = $1 AND student_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, lessonID, studentID).Scan(
		&record.ID, &record.LessonID, &record.StudentID, &record.Status,
		&record.Reason, &record.Comment, &record.MarkedBy, &record.MarkedAt,
		&record.UpdatedBy, &record.UpdatedAt, &record.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *attendanceRepository) GetByStudent(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*domain.AttendanceRecord, error) {
	query := `
		SELECT ar.id, ar.lesson_id, ar.student_id, ar.status, ar.reason, ar.comment,
		       ar.marked_by, ar.marked_at, ar.updated_by, ar.updated_at, ar.created_at
		FROM attendance_records ar
		JOIN lessons l ON ar.lesson_id = l.id
		WHERE ar.student_id = $1 AND l.scheduled_at BETWEEN $2 AND $3
		ORDER BY l.scheduled_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.AttendanceRecord
	for rows.Next() {
		var record domain.AttendanceRecord
		err := rows.Scan(
			&record.ID, &record.LessonID, &record.StudentID, &record.Status,
			&record.Reason, &record.Comment, &record.MarkedBy, &record.MarkedAt,
			&record.UpdatedBy, &record.UpdatedAt, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	return records, nil
}

func (r *attendanceRepository) GetByLesson(ctx context.Context, lessonID string) ([]*domain.AttendanceRecord, error) {
	query := `
		SELECT id, lesson_id, student_id, status, reason, comment, marked_by, marked_at, updated_by, updated_at, created_at
		FROM attendance_records
		WHERE lesson_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.AttendanceRecord
	for rows.Next() {
		var record domain.AttendanceRecord
		err := rows.Scan(
			&record.ID, &record.LessonID, &record.StudentID, &record.Status,
			&record.Reason, &record.Comment, &record.MarkedBy, &record.MarkedAt,
			&record.UpdatedBy, &record.UpdatedAt, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	return records, nil
}

func (r *attendanceRepository) Create(ctx context.Context, record *domain.AttendanceRecord) error {
	query := `
		INSERT INTO attendance_records (id, lesson_id, student_id, status, reason, comment, marked_by, marked_at, updated_by, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, $7, CURRENT_TIMESTAMP)
		RETURNING id, marked_at, updated_at, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		record.LessonID, record.StudentID, record.Status, record.Reason,
		record.Comment, record.MarkedBy, record.UpdatedBy,
	).Scan(&record.ID, &record.MarkedAt, &record.UpdatedAt, &record.CreatedAt)
}

func (r *attendanceRepository) Update(ctx context.Context, record *domain.AttendanceRecord) error {
	query := `
		UPDATE attendance_records
		SET status = $1, reason = $2, comment = $3, updated_by = $4, updated_at = CURRENT_TIMESTAMP
		WHERE lesson_id = $5 AND student_id = $6
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		record.Status, record.Reason, record.Comment, record.UpdatedBy,
		record.LessonID, record.StudentID,
	).Scan(&record.UpdatedAt)
}

func (r *attendanceRepository) GetStudentStats(ctx context.Context, studentID string) (map[string]int, error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'ATTENDED') as attended,
			COUNT(*) FILTER (WHERE status = 'ABSENT_EXCUSED') as absent_excused,
			COUNT(*) FILTER (WHERE status = 'ABSENT_UNEXCUSED') as absent_unexcused,
			COUNT(*) FILTER (WHERE status = 'FREEZE') as freeze
		FROM attendance_records
		WHERE student_id = $1
	`
	stats := make(map[string]int)
	var attended, absentExcused, absentUnexcused, freeze int
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&attended, &absentExcused, &absentUnexcused, &freeze)
	if err != nil {
		return nil, err
	}
	stats["attended"] = attended
	stats["absent_excused"] = absentExcused
	stats["absent_unexcused"] = absentUnexcused
	stats["freeze"] = freeze
	return stats, nil
}
