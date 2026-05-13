package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type AccessRepository interface {
	CreateRequest(ctx context.Context, req *domain.AccessRequest) error
	GetRequestByID(ctx context.Context, id string) (*domain.AccessRequest, error)
	GetPendingRequests(ctx context.Context) ([]*domain.AccessRequest, error)
	GetUserRequests(ctx context.Context, userID string) ([]*domain.AccessRequest, error)
	UpdateRequestStatus(ctx context.Context, id string, status domain.AccessRequestStatus, reviewedBy, reviewComment *string) error
}

type accessRepository struct {
	db *sql.DB
}

func NewAccessRepository(db *sql.DB) AccessRepository {
	return &accessRepository{db: db}
}

func (r *accessRepository) CreateRequest(ctx context.Context, req *domain.AccessRequest) error {
	query := `
		INSERT INTO access_requests (id, user_id, resource_type, resource_id, reason, status)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		req.UserID, req.ResourceType, req.ResourceID, req.Reason, req.Status,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *accessRepository) GetRequestByID(ctx context.Context, id string) (*domain.AccessRequest, error) {
	var req domain.AccessRequest
	query := `
		SELECT id, user_id, resource_type, resource_id, reason, status,
		       reviewed_by, reviewed_at, review_comment, created_at, updated_at
		FROM access_requests
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&req.ID, &req.UserID, &req.ResourceType, &req.ResourceID, &req.Reason,
		&req.Status, &req.ReviewedBy, &req.ReviewedAt, &req.ReviewComment,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *accessRepository) GetPendingRequests(ctx context.Context) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_type, resource_id, reason, status,
		       reviewed_by, reviewed_at, review_comment, created_at, updated_at
		FROM access_requests
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.AccessRequest
	for rows.Next() {
		var req domain.AccessRequest
		err := rows.Scan(
			&req.ID, &req.UserID, &req.ResourceType, &req.ResourceID, &req.Reason,
			&req.Status, &req.ReviewedBy, &req.ReviewedAt, &req.ReviewComment,
			&req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

func (r *accessRepository) GetUserRequests(ctx context.Context, userID string) ([]*domain.AccessRequest, error) {
	query := `
		SELECT id, user_id, resource_type, resource_id, reason, status,
		       reviewed_by, reviewed_at, review_comment, created_at, updated_at
		FROM access_requests
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.AccessRequest
	for rows.Next() {
		var req domain.AccessRequest
		err := rows.Scan(
			&req.ID, &req.UserID, &req.ResourceType, &req.ResourceID, &req.Reason,
			&req.Status, &req.ReviewedBy, &req.ReviewedAt, &req.ReviewComment,
			&req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

func (r *accessRepository) UpdateRequestStatus(ctx context.Context, id string, status domain.AccessRequestStatus, reviewedBy, reviewComment *string) error {
	query := `
		UPDATE access_requests
		SET status = $1, reviewed_by = $2, reviewed_at = CURRENT_TIMESTAMP, review_comment = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, status, reviewedBy, reviewComment, id)
	return err
}
