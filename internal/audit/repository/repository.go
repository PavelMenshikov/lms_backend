package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type AuditRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	GetByEntity(ctx context.Context, entityType, entityID string) ([]*domain.AuditLog, error)
	GetByUser(ctx context.Context, userID string, limit int) ([]*domain.AuditLog, error)
	GetRecent(ctx context.Context, limit int) ([]*domain.AuditLog, error)
}

type auditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		log.UserID, log.Action, log.EntityType, log.EntityID,
		log.OldValues, log.NewValues, log.IPAddress, log.UserAgent,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *auditRepository) GetByEntity(ctx context.Context, entityType, entityID string) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.EntityType, &log.EntityID,
			&log.OldValues, &log.NewValues, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

func (r *auditRepository) GetByUser(ctx context.Context, userID string, limit int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.EntityType, &log.EntityID,
			&log.OldValues, &log.NewValues, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, nil
}

func (r *auditRepository) GetRecent(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent, created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.EntityType, &log.EntityID,
			&log.OldValues, &log.NewValues, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	return logs, nil
}
