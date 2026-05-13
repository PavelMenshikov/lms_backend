package usecase

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/statistics/repository"
)

type StatisticsUseCase interface {
	GetStudentStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error)
	RefreshStudentStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error)
}

type statisticsUseCase struct {
	repo repository.StatisticsRepository
}

func NewStatisticsUseCase(repo repository.StatisticsRepository) StatisticsUseCase {
	return &statisticsUseCase{repo: repo}
}

func (uc *statisticsUseCase) GetStudentStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
	return uc.repo.GetByStudent(ctx, studentID)
}

func (uc *statisticsUseCase) RefreshStudentStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
	return uc.repo.RecalculateStatistics(ctx, studentID)
}
