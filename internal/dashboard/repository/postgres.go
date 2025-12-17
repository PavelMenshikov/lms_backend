package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type DashboardRepositoryImpl struct {
	db *sql.DB
}

var _ UserDataRepository = (*DashboardRepositoryImpl)(nil)

func NewDashboardRepository(db *sql.DB) *DashboardRepositoryImpl {
	return &DashboardRepositoryImpl{db: db}
}

func (r *DashboardRepositoryImpl) GetLastLessonData(ctx context.Context, userID string) (*domain.LastLesson, error) {

	return &domain.LastLesson{}, nil
}

func (r *DashboardRepositoryImpl) GetActiveCoursesCount(ctx context.Context, userID string) (int, error) {

	return 0, nil
}

func (r *DashboardRepositoryImpl) GetAttendancePercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {

	return &domain.StatisticSummary{}, nil
}

func (r *DashboardRepositoryImpl) GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {

	return &domain.StatisticSummary{}, nil
}

func (r *DashboardRepositoryImpl) GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error) {

	return []domain.UpcomingLesson{}, nil
}
