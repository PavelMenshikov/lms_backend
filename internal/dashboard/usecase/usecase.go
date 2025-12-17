package usecase

import (
	"context"

	"lms_backend/internal/dashboard/repository"
	"lms_backend/internal/domain"
)

type DashboardUseCase struct {
	repo repository.UserDataRepository
}

func NewDashboardUseCase(repo repository.UserDataRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

// Заглушка, чтобы пакет был валидным
func (uc *DashboardUseCase) GetUserHomeData(ctx context.Context, user *domain.User) (*domain.HomeDashboard, error) {

	lastLesson, _ := uc.repo.GetLastLessonData(ctx, user.ID)
	coursesCount, _ := uc.repo.GetActiveCoursesCount(ctx, user.ID)

	dashboardData := &domain.HomeDashboard{
		UserRole:           user.Role,
		User:               user,
		LastLessonData:     lastLesson,
		ActiveCoursesCount: coursesCount,
	}

	return dashboardData, nil
}
