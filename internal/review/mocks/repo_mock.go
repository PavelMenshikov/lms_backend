package mocks

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/review/repository"
	"sync"
)

type ReviewRepositoryMock struct {
	mu          sync.Mutex
	Submissions map[string]*domain.SubmissionRecord
	nextID      int
}

var _ repository.ReviewRepository = (*ReviewRepositoryMock)(nil)

func NewReviewRepositoryMock() *ReviewRepositoryMock {
	return &ReviewRepositoryMock{
		Submissions: make(map[string]*domain.SubmissionRecord),
		nextID:      1,
	}
}

func (r *ReviewRepositoryMock) GetPendingSubmissions(ctx context.Context, staffID string, role string, studentID string) ([]*domain.SubmissionRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.SubmissionRecord
	for _, s := range r.Submissions {
		if s.Status == "pending" {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *ReviewRepositoryMock) EvaluateSubmission(ctx context.Context, submissionID string, grade int, comment string, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s, ok := r.Submissions[submissionID]
	if !ok {
		s = &domain.SubmissionRecord{ID: submissionID}
		r.Submissions[submissionID] = s
	}
	s.Grade = grade
	s.TeacherComment = comment
	s.Status = status
	return nil
}

func (r *ReviewRepositoryMock) UpdateUserCourseProgress(ctx context.Context, userID, courseID string) error {
	return nil
}
