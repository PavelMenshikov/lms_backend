package mocks

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/groups/repository"
)

var _ repository.GroupRepository = (*GroupRepoMock)(nil)

type GroupRepoMock struct {
	GetByIDFunc            func(ctx context.Context, id string) (*domain.Group, error)
	UpdateGroupFunc        func(ctx context.Context, groupID, name string, teacherID *string) error
	AddStudentToGroupFunc  func(ctx context.Context, groupID, studentID string) error
	RemoveStudentFromGroupFunc func(ctx context.Context, groupID, studentID string) error
	ChangeStudentGroupFunc func(ctx context.Context, studentID, newGroupID string) error
	ChangeTeacherGroupFunc func(ctx context.Context, teacherID, newGroupID string) error
	GetGroupStudentsFunc   func(ctx context.Context, groupID string) ([]string, error)
}

func NewGroupRepoMock() *GroupRepoMock {
	return &GroupRepoMock{}
}

func (m *GroupRepoMock) GetByID(ctx context.Context, id string) (*domain.Group, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *GroupRepoMock) UpdateGroup(ctx context.Context, groupID, name string, teacherID *string) error {
	return m.UpdateGroupFunc(ctx, groupID, name, teacherID)
}

func (m *GroupRepoMock) AddStudentToGroup(ctx context.Context, groupID, studentID string) error {
	return m.AddStudentToGroupFunc(ctx, groupID, studentID)
}

func (m *GroupRepoMock) RemoveStudentFromGroup(ctx context.Context, groupID, studentID string) error {
	return m.RemoveStudentFromGroupFunc(ctx, groupID, studentID)
}

func (m *GroupRepoMock) ChangeStudentGroup(ctx context.Context, studentID, newGroupID string) error {
	return m.ChangeStudentGroupFunc(ctx, studentID, newGroupID)
}

func (m *GroupRepoMock) ChangeTeacherGroup(ctx context.Context, teacherID, newGroupID string) error {
	return m.ChangeTeacherGroupFunc(ctx, teacherID, newGroupID)
}

func (m *GroupRepoMock) GetGroupStudents(ctx context.Context, groupID string) ([]string, error) {
	return m.GetGroupStudentsFunc(ctx, groupID)
}
