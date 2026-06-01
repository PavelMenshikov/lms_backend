package mocks

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/statistics/repository"
)

var _ repository.StatisticsRepository = (*StatisticsRepoMock)(nil)

type StatisticsRepoMock struct {
	GetByStudentFunc           func(ctx context.Context, studentID string) (*domain.StudentStatistics, error)
	UpdateStatisticsFunc       func(ctx context.Context, studentID string) error
	RecalculateStatisticsFunc  func(ctx context.Context, studentID string) (*domain.StudentStatistics, error)
}

func NewStatisticsRepoMock() *StatisticsRepoMock {
	return &StatisticsRepoMock{}
}

func (m *StatisticsRepoMock) GetByStudent(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
	return m.GetByStudentFunc(ctx, studentID)
}

func (m *StatisticsRepoMock) UpdateStatistics(ctx context.Context, studentID string) error {
	return m.UpdateStatisticsFunc(ctx, studentID)
}

func (m *StatisticsRepoMock) RecalculateStatistics(ctx context.Context, studentID string) (*domain.StudentStatistics, error) {
	return m.RecalculateStatisticsFunc(ctx, studentID)
}
