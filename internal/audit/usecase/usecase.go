package usecase

import (
	"context"
	"encoding/json"
	"lms_backend/internal/audit/repository"
	"lms_backend/internal/domain"
)

type AuditUseCase interface {
	LogAction(ctx context.Context, userID *string, action, entityType, entityID string, oldValues, newValues interface{}, ipAddress, userAgent *string) error
	GetEntityHistory(ctx context.Context, entityType, entityID string) ([]*domain.AuditLog, error)
	GetUserActivity(ctx context.Context, userID string, limit int) ([]*domain.AuditLog, error)
	GetRecentActivity(ctx context.Context, limit int) ([]*domain.AuditLog, error)
}

type auditUseCase struct {
	repo repository.AuditRepository
}

func NewAuditUseCase(repo repository.AuditRepository) AuditUseCase {
	return &auditUseCase{repo: repo}
}

func (uc *auditUseCase) LogAction(ctx context.Context, userID *string, action, entityType, entityID string, oldValues, newValues interface{}, ipAddress, userAgent *string) error {
	var oldJSON, newJSON *string

	if oldValues != nil {
		data, err := json.Marshal(oldValues)
		if err == nil {
			str := string(data)
			oldJSON = &str
		}
	}

	if newValues != nil {
		data, err := json.Marshal(newValues)
		if err == nil {
			str := string(data)
			newJSON = &str
		}
	}

	log := &domain.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValues:  oldJSON,
		NewValues:  newJSON,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	return uc.repo.Create(ctx, log)
}

func (uc *auditUseCase) GetEntityHistory(ctx context.Context, entityType, entityID string) ([]*domain.AuditLog, error) {
	return uc.repo.GetByEntity(ctx, entityType, entityID)
}

func (uc *auditUseCase) GetUserActivity(ctx context.Context, userID string, limit int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return uc.repo.GetByUser(ctx, userID, limit)
}

func (uc *auditUseCase) GetRecentActivity(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return uc.repo.GetRecent(ctx, limit)
}
