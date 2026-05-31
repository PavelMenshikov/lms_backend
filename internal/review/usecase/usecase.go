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
	if input.Status == "accepted" {
		validGrades := map[int]bool{20: true, 40: true, 60: true, 80: true, 100: true}
		if !validGrades[input.Grade] {
			input.Grade = 100
		}
	} else {
		input.Grade = 0
	}

	return uc.repo.EvaluateSubmission(ctx, input.SubmissionID, input.Grade, input.Comment, input.Status)
}
