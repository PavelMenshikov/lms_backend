package mocks

import (
	"context"
	"lms_backend/internal/profile/repository"
	"lms_backend/internal/domain"
)

var _ repository.ProfileRepository = (*ProfileRepoMock)(nil)

type ProfileRepoMock struct {
	GetProfileFunc          func(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfileFunc       func(ctx context.Context, user *domain.User) error
	UpdateTeacherScheduleFunc func(ctx context.Context, userID string, scheduleJSON []byte) error
}

func NewProfileRepoMock() *ProfileRepoMock {
	return &ProfileRepoMock{}
}

func (m *ProfileRepoMock) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return m.GetProfileFunc(ctx, userID)
}

func (m *ProfileRepoMock) UpdateProfile(ctx context.Context, user *domain.User) error {
	return m.UpdateProfileFunc(ctx, user)
}

func (m *ProfileRepoMock) UpdateTeacherSchedule(ctx context.Context, userID string, scheduleJSON []byte) error {
	return m.UpdateTeacherScheduleFunc(ctx, userID, scheduleJSON)
}
