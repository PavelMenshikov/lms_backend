package usecase

import (
	"context"
	"errors"
	"lms_backend/internal/groups/repository"
)

type GroupUseCase interface {
	UpdateGroup(ctx context.Context, groupID, name string, teacherID *string) error
	AddStudentToGroup(ctx context.Context, groupID, studentID string) error
	RemoveStudentFromGroup(ctx context.Context, groupID, studentID string) error
	ChangeStudentGroup(ctx context.Context, studentID, newGroupID string) error
	ChangeTeacherGroup(ctx context.Context, teacherID, newGroupID string) error
}

type groupUseCase struct {
	repo repository.GroupRepository
}

func NewGroupUseCase(repo repository.GroupRepository) GroupUseCase {
	return &groupUseCase{repo: repo}
}

func (uc *groupUseCase) UpdateGroup(ctx context.Context, groupID, name string, teacherID *string) error {
	if name == "" {
		return errors.New("group name cannot be empty")
	}

	// Проверяем, что группа существует
	_, err := uc.repo.GetByID(ctx, groupID)
	if err != nil {
		return errors.New("group not found")
	}

	return uc.repo.UpdateGroup(ctx, groupID, name, teacherID)
}

func (uc *groupUseCase) AddStudentToGroup(ctx context.Context, groupID, studentID string) error {
	// Проверяем, что группа существует
	_, err := uc.repo.GetByID(ctx, groupID)
	if err != nil {
		return errors.New("group not found")
	}

	return uc.repo.AddStudentToGroup(ctx, groupID, studentID)
}

func (uc *groupUseCase) RemoveStudentFromGroup(ctx context.Context, groupID, studentID string) error {
	// Проверяем, что группа существует
	_, err := uc.repo.GetByID(ctx, groupID)
	if err != nil {
		return errors.New("group not found")
	}

	return uc.repo.RemoveStudentFromGroup(ctx, groupID, studentID)
}

func (uc *groupUseCase) ChangeStudentGroup(ctx context.Context, studentID, newGroupID string) error {
	// Проверяем, что новая группа существует
	_, err := uc.repo.GetByID(ctx, newGroupID)
	if err != nil {
		return errors.New("new group not found")
	}

	return uc.repo.ChangeStudentGroup(ctx, studentID, newGroupID)
}

func (uc *groupUseCase) ChangeTeacherGroup(ctx context.Context, teacherID, newGroupID string) error {
	// Проверяем, что новая группа существует
	_, err := uc.repo.GetByID(ctx, newGroupID)
	if err != nil {
		return errors.New("new group not found")
	}

	return uc.repo.ChangeTeacherGroup(ctx, teacherID, newGroupID)
}
