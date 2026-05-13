package usecase

import (
	"context"
	"errors"
	"lms_backend/internal/domain"
	"lms_backend/internal/freeze/repository"
	"time"
)

type FreezeUseCase interface {
	CreateRequest(ctx context.Context, studentID, requestedBy string, startDate, endDate time.Time, reason string) error
	GetPendingRequests(ctx context.Context) ([]*domain.FreezeRequest, error)
	ApproveRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error
	RejectRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error
	GetStudentFreezeStatus(ctx context.Context, studentID string) (*domain.FreezePeriod, error)
	GetStudentRequests(ctx context.Context, studentID string) ([]*domain.FreezeRequest, error)
}

type freezeUseCase struct {
	repo repository.FreezeRepository
}

func NewFreezeUseCase(repo repository.FreezeRepository) FreezeUseCase {
	return &freezeUseCase{repo: repo}
}

func (uc *freezeUseCase) CreateRequest(ctx context.Context, studentID, requestedBy string, startDate, endDate time.Time, reason string) error {
	if endDate.Before(startDate) {
		return errors.New("end_date must be after start_date")
	}

	req := &domain.FreezeRequest{
		StudentID:   studentID,
		RequestedBy: requestedBy,
		StartDate:   startDate,
		EndDate:     endDate,
		Reason:      reason,
		Status:      domain.FreezeStatusPending,
	}
	return uc.repo.CreateRequest(ctx, req)
}

func (uc *freezeUseCase) GetPendingRequests(ctx context.Context) ([]*domain.FreezeRequest, error) {
	return uc.repo.GetPendingRequests(ctx)
}

func (uc *freezeUseCase) ApproveRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error {
	// Получаем запрос
	req, err := uc.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != domain.FreezeStatusPending {
		return errors.New("request is not pending")
	}

	// Обновляем статус запроса
	err = uc.repo.UpdateRequestStatus(ctx, requestID, domain.FreezeStatusApproved, &reviewedBy, reviewComment)
	if err != nil {
		return err
	}

	// Создаём период заморозки
	period := &domain.FreezePeriod{
		StudentID:       req.StudentID,
		FreezeRequestID: &req.ID,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsActive:        true,
		CreatedBy:       reviewedBy,
	}
	return uc.repo.CreatePeriod(ctx, period)
}

func (uc *freezeUseCase) RejectRequest(ctx context.Context, requestID, reviewedBy string, reviewComment *string) error {
	req, err := uc.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	if req.Status != domain.FreezeStatusPending {
		return errors.New("request is not pending")
	}

	return uc.repo.UpdateRequestStatus(ctx, requestID, domain.FreezeStatusRejected, &reviewedBy, reviewComment)
}

func (uc *freezeUseCase) GetStudentFreezeStatus(ctx context.Context, studentID string) (*domain.FreezePeriod, error) {
	return uc.repo.GetStudentFreezeStatus(ctx, studentID)
}

func (uc *freezeUseCase) GetStudentRequests(ctx context.Context, studentID string) ([]*domain.FreezeRequest, error) {
	return uc.repo.GetRequestsByStudent(ctx, studentID)
}
