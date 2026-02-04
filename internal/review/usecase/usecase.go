package usecase

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/review/repository"
)

type ReviewUseCase struct {
	repo repository.ReviewRepository
}

func NewReviewUseCase(repo repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{repo: repo}
}

func (uc *ReviewUseCase) GetPendingList(ctx context.Context) ([]*domain.SubmissionRecord, error) {
	return uc.repo.GetPendingSubmissions(ctx)
}

type EvaluateInput struct {
	SubmissionID string
	Grade        int
	Comment      string
	IsAccepted   bool
}

func (uc *ReviewUseCase) Evaluate(ctx context.Context, input EvaluateInput) error {
	status := "rejected"
	if input.IsAccepted {
		status = "accepted"
	}

	err := uc.repo.EvaluateSubmission(ctx, input.SubmissionID, input.Grade, input.Comment, status)
	if err != nil {
		return err
	}

	return nil
}
