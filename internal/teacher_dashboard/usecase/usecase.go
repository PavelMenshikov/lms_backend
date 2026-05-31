package usecase

import (
	"context"

	"lms_backend/internal/domain"
	"lms_backend/internal/teacher_dashboard/repository"
)

type TeacherDashboardUseCase struct {
	repo repository.TeacherDashboardRepository
}

func NewTeacherDashboardUseCase(repo repository.TeacherDashboardRepository) *TeacherDashboardUseCase {
	return &TeacherDashboardUseCase{repo: repo}
}

func (uc *TeacherDashboardUseCase) GetMonthlyReport(ctx context.Context, teacherID string, year, month int) (*domain.TeacherMonthlyReport, error) {
	return uc.repo.GetMonthlyReport(ctx, teacherID, year, month)
}
