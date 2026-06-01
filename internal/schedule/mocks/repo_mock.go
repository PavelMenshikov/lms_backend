package mocks

import (
	"context"
	"time"
	"lms_backend/internal/domain"
	"lms_backend/internal/schedule/repository"
)

var _ repository.ScheduleRepository = (*ScheduleRepoMock)(nil)

type ScheduleRepoMock struct {
	GetStudentLessonsInRangeFunc func(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error)
}

func NewScheduleRepoMock() *ScheduleRepoMock {
	return &ScheduleRepoMock{}
}

func (m *ScheduleRepoMock) GetStudentLessonsInRange(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
	return m.GetStudentLessonsInRangeFunc(ctx, userID, start, end)
}
