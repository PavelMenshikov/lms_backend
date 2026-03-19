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

func (uc *ReviewUseCase) GetPendingList(ctx context.Context, staffID string, role string) ([]*domain.SubmissionRecord, error) {
	return uc.repo.GetPendingSubmissions(ctx, staffID, role)
}

type EvaluateInput struct {
	SubmissionID string
	Grade        int
	Comment      string
	Status       string 
}

func (uc *ReviewUseCase) Evaluate(ctx context.Context, input EvaluateInput) error {
	return uc.repo.EvaluateSubmission(ctx, input.SubmissionID, input.Grade, input.Comment, input.Status)
}