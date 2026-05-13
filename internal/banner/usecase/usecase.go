package usecase

import (
	"context"
	"lms_backend/internal/banner/repository"
	"lms_backend/internal/domain"
	"time"
)

type BannerUseCase interface {
	CreateBanner(ctx context.Context, title, content string, bannerType domain.BannerType, isActive bool, priority int, startDate, endDate *time.Time, targetRoles []string, createdBy string) error
	GetActiveBanners(ctx context.Context, role *string) ([]*domain.Banner, error)
	GetAllBanners(ctx context.Context) ([]*domain.Banner, error)
	UpdateBanner(ctx context.Context, id, title, content string, bannerType domain.BannerType, isActive bool, priority int, startDate, endDate *time.Time, targetRoles []string) error
	DeleteBanner(ctx context.Context, id string) error
}

type bannerUseCase struct {
	repo repository.BannerRepository
}

func NewBannerUseCase(repo repository.BannerRepository) BannerUseCase {
	return &bannerUseCase{repo: repo}
}

func (uc *bannerUseCase) CreateBanner(ctx context.Context, title, content string, bannerType domain.BannerType, isActive bool, priority int, startDate, endDate *time.Time, targetRoles []string, createdBy string) error {
	banner := &domain.Banner{
		Title:       title,
		Content:     content,
		Type:        bannerType,
		IsActive:    isActive,
		Priority:    priority,
		StartDate:   startDate,
		EndDate:     endDate,
		TargetRoles: targetRoles,
		CreatedBy:   createdBy,
	}
	return uc.repo.Create(ctx, banner)
}

func (uc *bannerUseCase) GetActiveBanners(ctx context.Context, role *string) ([]*domain.Banner, error) {
	return uc.repo.GetActive(ctx, role)
}

func (uc *bannerUseCase) GetAllBanners(ctx context.Context) ([]*domain.Banner, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *bannerUseCase) UpdateBanner(ctx context.Context, id, title, content string, bannerType domain.BannerType, isActive bool, priority int, startDate, endDate *time.Time, targetRoles []string) error {
	banner := &domain.Banner{
		ID:          id,
		Title:       title,
		Content:     content,
		Type:        bannerType,
		IsActive:    isActive,
		Priority:    priority,
		StartDate:   startDate,
		EndDate:     endDate,
		TargetRoles: targetRoles,
	}
	return uc.repo.Update(ctx, banner)
}

func (uc *bannerUseCase) DeleteBanner(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
