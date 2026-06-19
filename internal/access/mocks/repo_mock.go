package mocks

import (
	"context"
	"errors"
	"lms_backend/internal/access/repository"
	"lms_backend/internal/domain"
	"sync"
)

type AccessRepositoryMock struct {
	mu       sync.Mutex
	Requests map[string]*domain.AccessRequest
	nextID   int
}

var _ repository.AccessRepository = (*AccessRepositoryMock)(nil)

func NewAccessRepositoryMock() *AccessRepositoryMock {
	return &AccessRepositoryMock{
		Requests: make(map[string]*domain.AccessRequest),
		nextID:   1,
	}
}

func (r *AccessRepositoryMock) CreateRequest(ctx context.Context, req *domain.AccessRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := "req-" + r.itoa()
	req.ID = id
	req.Status = domain.AccessRequestStatusPending
	r.Requests[id] = req
	return nil
}

func (r *AccessRepositoryMock) itoa() string {
	id := r.nextID
	r.nextID++
	return string(rune('0'+id%10)) + string(rune('0'+(id/10)%10))
}

func (r *AccessRepositoryMock) GetRequestByID(ctx context.Context, id string) (*domain.AccessRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	req, ok := r.Requests[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return req, nil
}

func (r *AccessRepositoryMock) GetPendingRequests(ctx context.Context) ([]*domain.AccessRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.AccessRequest
	for _, req := range r.Requests {
		if req.Status == domain.AccessRequestStatusPending {
			result = append(result, req)
		}
	}
	return result, nil
}

func (r *AccessRepositoryMock) GetUserRequests(ctx context.Context, userID string) ([]*domain.AccessRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.AccessRequest
	for _, req := range r.Requests {
		if req.UserID == userID {
			result = append(result, req)
		}
	}
	return result, nil
}

func (r *AccessRepositoryMock) UpdateRequestStatus(ctx context.Context, id string, status domain.AccessRequestStatus, reviewedBy, reviewComment *string) error {
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
