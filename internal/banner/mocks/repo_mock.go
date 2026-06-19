package mocks

import (
	"context"
	"errors"
	"lms_backend/internal/banner/repository"
	"lms_backend/internal/domain"
	"sync"
	"time"
)

type BannerRepositoryMock struct {
	mu      sync.Mutex
	Banners map[string]*domain.Banner
	nextID  int
}

var _ repository.BannerRepository = (*BannerRepositoryMock)(nil)

func NewBannerRepositoryMock() *BannerRepositoryMock {
	return &BannerRepositoryMock{
		Banners: make(map[string]*domain.Banner),
		nextID:  1,
	}
}

func (r *BannerRepositoryMock) Create(ctx context.Context, banner *domain.Banner) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := "banner-" + r.nextIDStr()
	banner.ID = id
	banner.CreatedAt = time.Now()
	banner.UpdatedAt = time.Now()
	r.Banners[id] = banner
	return nil
}

func (r *BannerRepositoryMock) nextIDStr() string {
	id := r.nextID
	r.nextID++
	return string(rune('0'+id%10)) + string(rune('0'+(id/10)%10))
}

func (r *BannerRepositoryMock) GetByID(ctx context.Context, id string) (*domain.Banner, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.Banners[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (r *BannerRepositoryMock) GetActive(ctx context.Context, role *string) ([]*domain.Banner, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Banner
	for _, b := range r.Banners {
		if b.IsActive {
			result = append(result, b)
		}
	}
	return result, nil
}

func (r *BannerRepositoryMock) GetAll(ctx context.Context) ([]*domain.Banner, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.Banner
	for _, b := range r.Banners {
		result = append(result, b)
	}
	return result, nil
}

func (r *BannerRepositoryMock) Update(ctx context.Context, banner *domain.Banner) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Banners[banner.ID] = banner
	return nil
}

func (r *BannerRepositoryMock) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Banners, id)
	return nil
}
