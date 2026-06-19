package mocks

import (
	"context"
	"errors"
	"lms_backend/internal/domain"
	"lms_backend/internal/freeze/repository"
	"sync"
	"time"
)

type FreezeRepositoryMock struct {
	mu       sync.Mutex
	Requests map[string]*domain.FreezeRequest
	Periods  map[string]*domain.FreezePeriod
	nextID   int
}

var _ repository.FreezeRepository = (*FreezeRepositoryMock)(nil)

func NewFreezeRepositoryMock() *FreezeRepositoryMock {
	return &FreezeRepositoryMock{
		Requests: make(map[string]*domain.FreezeRequest),
		Periods:  make(map[string]*domain.FreezePeriod),
		nextID:   1,
	}
}

func (r *FreezeRepositoryMock) nextIDStr() string {
	id := r.nextID
	r.nextID++
	return string(rune('0'+id%10)) + string(rune('0'+(id/10)%10))
}

func (r *FreezeRepositoryMock) CreateRequest(ctx context.Context, req *domain.FreezeRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := "freeze-req-" + r.nextIDStr()
	req.ID = id
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	r.Requests[id] = req
	return nil
}

func (r *FreezeRepositoryMock) GetRequestByID(ctx context.Context, id string) (*domain.FreezeRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	req, ok := r.Requests[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return req, nil
}

func (r *FreezeRepositoryMock) GetRequestsByStudent(ctx context.Context, studentID string) ([]*domain.FreezeRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.FreezeRequest
	for _, req := range r.Requests {
		if req.StudentID == studentID {
			result = append(result, req)
		}
	}
	return result, nil
}

func (r *FreezeRepositoryMock) GetPendingRequests(ctx context.Context) ([]*domain.FreezeRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.FreezeRequest
	for _, req := range r.Requests {
		if req.Status == domain.FreezeStatusPending {
			result = append(result, req)
		}
	}
	return result, nil
}

func (r *FreezeRepositoryMock) UpdateRequestStatus(ctx context.Context, id string, status domain.FreezeStatus, reviewedBy, reviewComment *string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	req, ok := r.Requests[id]
	if !ok {
		return errors.New("not found")
	}
	req.Status = status
	req.ReviewedBy = reviewedBy
	req.ReviewComment = reviewComment
	return nil
}

func (r *FreezeRepositoryMock) CreatePeriod(ctx context.Context, period *domain.FreezePeriod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := "freeze-period-" + r.nextIDStr()
	period.ID = id
	period.CreatedAt = time.Now()
	period.UpdatedAt = time.Now()
	r.Periods[id] = period
	return nil
}

func (r *FreezeRepositoryMock) GetActivePeriods(ctx context.Context, studentID string) ([]*domain.FreezePeriod, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.FreezePeriod
	for _, p := range r.Periods {
		if p.StudentID == studentID && p.IsActive {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *FreezeRepositoryMock) GetStudentFreezeStatus(ctx context.Context, studentID string) (*domain.FreezePeriod, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.Periods {
		if p.StudentID == studentID && p.IsActive {
			return p, nil
		}
	}
	return nil, nil
}
