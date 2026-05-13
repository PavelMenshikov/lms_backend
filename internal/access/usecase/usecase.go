package usecase

import (
	"context"
	"errors"
	"lms_backend/internal/access/repository"
	"lms_backend/internal/domain"
)

type AccessUseCase interface {
	CreateRequest(ctx context.Context, userID, resourceType, resourceID, reason string) error
	GetPendingRequests(ctx context.Context) ([]*domain.AccessRequest, error)
	ApproveRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error
	RejectRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error
	GetUserRequests(ctx context.Context, userID string) ([]*domain.AccessRequest, error)
}

type accessUseCase struct {
	repo repository.AccessRepository
}

func NewAccessUseCase(repo repository.AccessRepository) AccessUseCase {
	return &accessUseCase{repo: repo}
}

func (uc *accessUseCase) CreateRequest(ctx context.Context, userID, resourceType, resourceID, reason string) error {
	req := &domain.AccessRequest{
		UserID:       userID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Reason:       reason,
		Status:       domain.AccessRequestStatusPending,
	}
	return uc.repo.CreateRequest(ctx, req)
}

func (uc *accessUseCase) GetPendingRequests(ctx context.Context) ([]*domain.AccessRequest, error) {
	return uc.repo.GetPendingRequests(ctx)
}

func (uc *accessUseCase) ApproveRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error {
	req, err := uc.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != domain.AccessRequestStatusPending {
		return errors.New("request is not pending")
	}

	return uc.repo.UpdateRequestStatus(ctx, requestID, domain.AccessRequestStatusApproved, &reviewedBy, reviewComment)
}

func (uc *accessUseCase) RejectRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error {
	req, err := uc.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != domain.AccessRequestStatusPending {
		return errors.New("request is not pending")
	}

	return uc.repo.UpdateRequestStatus(ctx, requestID, domain.AccessRequestStatusRejected, &reviewedBy, reviewComment)
}

func (uc *accessUseCase) GetUserRequests(ctx context.Context, userID string) ([]*domain.AccessRequest, error) {
	return uc.repo.GetUserRequests(ctx, userID)
}
