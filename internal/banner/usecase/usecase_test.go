package usecase_test

import (
	"context"
	"testing"

	"lms_backend/internal/banner/mocks"
	"lms_backend/internal/banner/usecase"
	"lms_backend/internal/domain"
)

func TestBannerUseCase_CreateAndGet(t *testing.T) {
	repoMock := mocks.NewBannerRepositoryMock()
	uc := usecase.NewBannerUseCase(repoMock)
	ctx := context.Background()

	t.Run("CreateBanner", func(t *testing.T) {
		err := uc.CreateBanner(ctx, "Test Banner", "Content", domain.BannerTypeInfo, true, 1, nil, nil, []string{"student"}, "admin-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("GetAllBanners", func(t *testing.T) {
		banners, err := uc.GetAllBanners(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(banners) != 1 {
			t.Errorf("expected 1 banner, got %d", len(banners))
		}
	})

	t.Run("GetActiveBanners", func(t *testing.T) {
		role := "student"
		banners, err := uc.GetActiveBanners(ctx, &role)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(banners) != 1 {
			t.Errorf("expected 1 active banner, got %d", len(banners))
		}
	})
}

func TestBannerUseCase_UpdateAndDelete(t *testing.T) {
	repoMock := mocks.NewBannerRepositoryMock()
	uc := usecase.NewBannerUseCase(repoMock)
	ctx := context.Background()

	err := uc.CreateBanner(ctx, "Test", "Content", domain.BannerTypeInfo, true, 1, nil, nil, nil, "admin-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	banners, _ := uc.GetAllBanners(ctx)
	id := banners[0].ID

	t.Run("UpdateBanner", func(t *testing.T) {
		err := uc.UpdateBanner(ctx, id, "Updated", "Updated Content", domain.BannerTypeWarning, false, 2, nil, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		updated, _ := uc.GetAllBanners(ctx)
		if updated[0].Title != "Updated" {
			t.Errorf("expected 'Updated', got %s", updated[0].Title)
		}
	})

	t.Run("DeleteBanner", func(t *testing.T) {
		err := uc.DeleteBanner(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		banners, _ := uc.GetAllBanners(ctx)
		if len(banners) != 0 {
			t.Errorf("expected 0 banners after delete, got %d", len(banners))
		}
	})
}
