package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type FreezeRepository interface {
	CreateRequest(ctx context.Context, req *domain.FreezeRequest) error
	GetRequestByID(ctx context.Context, id string) (*domain.FreezeRequest, error)
	GetRequestsByStudent(ctx context.Context, studentID string) ([]*domain.FreezeRequest, error)
	GetPendingRequests(ctx context.Context) ([]*domain.FreezeRequest, error)
	UpdateRequestStatus(ctx context.Context, id string, status domain.FreezeStatus, reviewedBy, reviewComment *string) error
	CreatePeriod(ctx context.Context, period *domain.FreezePeriod) error
	GetActivePeriods(ctx context.Context, studentID string) ([]*domain.FreezePeriod, error)
	GetStudentFreezeStatus(ctx context.Context, studentID string) (*domain.FreezePeriod, error)
}

type freezeRepository struct {
	db *sql.DB
}

func NewFreezeRepository(db *sql.DB) FreezeRepository {
	return &freezeRepository{db: db}
}

func (r *freezeRepository) CreateRequest(ctx context.Context, req *domain.FreezeRequest) error {
	query := `
		INSERT INTO freeze_requests (id, student_id, requested_by, start_date, end_date, reason, status)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		req.StudentID, req.RequestedBy, req.StartDate, req.EndDate, req.Reason, req.Status,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *freezeRepository) GetRequestByID(ctx context.Context, id string) (*domain.FreezeRequest, error) {
	var req domain.FreezeRequest
	query := `
		SELECT id, student_id, requested_by, start_date, end_date, reason, status,
		       reviewed_by, reviewed_at, review_comment, created_at, updated_at
		FROM freeze_requests
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.StudentID, &req.RequestedBy, &req.StartDate, &req.EndDate,
		&req.Reason, &req.Status, &req.ReviewedBy, &req.ReviewedAt,
		&req.ReviewComment, &req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *freezeRepository) GetRequestsByStudent(ctx context.Context, studentID string) ([]*domain.FreezeRequest, error) {
	query := `
		SELECT id, student_id, requested_by, start_date, end_date, reason, status,
		       reviewed_by, reviewed_at, review_comment, created_at, updated_at
		FROM freeze_requests
		WHERE student_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.FreezeRequest
	for rows.Next() {
		var req domain.FreezeRequest
		err := rows.Scan(
			&req.ID, &req.StudentID, &req.RequestedBy, &req.StartDate, &req.EndDate,
			&req.Reason, &req.Status, &req.ReviewedBy, &req.ReviewedAt,
			&req.ReviewComment, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

func (r *freezeRepository) GetPendingRequests(ctx context.Context) ([]*domain.FreezeRequest, error) {
	query := `
		SELECT id, student_id, requested_by, start_date, end_date, reason, status,
		       reviewed_by, reviewed_at, review_comment, created_at, updated_at
		FROM freeze_requests
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.FreezeRequest
	for rows.Next() {
		var req domain.FreezeRequest
		err := rows.Scan(
			&req.ID, &req.StudentID, &req.RequestedBy, &req.StartDate, &req.EndDate,
			&req.Reason, &req.Status, &req.ReviewedBy, &req.ReviewedAt,
			&req.ReviewComment, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

func (r *freezeRepository) UpdateRequestStatus(ctx context.Context, id string, status domain.FreezeStatus, reviewedBy, reviewComment *string) error {
	query := `
		UPDATE freeze_requests
		SET status = $1, reviewed_by = $2, reviewed_at = CURRENT_TIMESTAMP, review_comment = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, status, reviewedBy, reviewComment, id)
	return err
}

func (r *freezeRepository) CreatePeriod(ctx context.Context, period *domain.FreezePeriod) error {
	query := `
		INSERT INTO freeze_periods (id, student_id, freeze_request_id, start_date, end_date, is_active, created_by)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		period.StudentID, period.FreezeRequestID, period.StartDate, period.EndDate,
		period.IsActive, period.CreatedBy,
	).Scan(&period.ID, &period.CreatedAt, &period.UpdatedAt)
}

func (r *freezeRepository) GetActivePeriods(ctx context.Context, studentID string) ([]*domain.FreezePeriod, error) {
	query := `
		SELECT id, student_id, freeze_request_id, start_date, end_date, is_active, created_by, created_at, updated_at
		FROM freeze_periods
		WHERE student_id = $1 AND is_active = true
		ORDER BY start_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []*domain.FreezePeriod
	for rows.Next() {
		var period domain.FreezePeriod
		err := rows.Scan(
			&period.ID, &period.StudentID, &period.FreezeRequestID, &period.StartDate,
			&period.EndDate, &period.IsActive, &period.CreatedBy, &period.CreatedAt, &period.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		periods = append(periods, &period)
	}
	return periods, nil
}

func (r *freezeRepository) GetStudentFreezeStatus(ctx context.Context, studentID string) (*domain.FreezePeriod, error) {
	var period domain.FreezePeriod
	query := `
		SELECT id, student_id, freeze_request_id, start_date, end_date, is_active, created_by, created_at, updated_at
		FROM freeze_periods
		WHERE student_id = $1 AND is_active = true AND end_date >= CURRENT_DATE
		ORDER BY end_date DESC
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(
		&period.ID, &period.StudentID, &period.FreezeRequestID, &period.StartDate,
		&period.EndDate, &period.IsActive, &period.CreatedBy, &period.CreatedAt, &period.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &period, nil
}
