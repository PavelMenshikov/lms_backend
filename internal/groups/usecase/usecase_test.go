package usecase

import (
	"context"
	"errors"
	"testing"

	"lms_backend/internal/domain"
	"lms_backend/internal/groups/mocks"
)

func TestUpdateGroup(t *testing.T) {
	mock := mocks.NewGroupRepoMock()
	uc := NewGroupUseCase(mock)

	mock.GetByIDFunc = func(ctx context.Context, id string) (*domain.Group, error) {
		if id == "g1" {
			return &domain.Group{ID: "g1", Title: "Group 1"}, nil
		}
		return nil, errors.New("not found")
	}

	mock.UpdateGroupFunc = func(ctx context.Context, groupID, name string, teacherID *string) error {
		return nil
	}

	t.Run("success", func(t *testing.T) {
		err := uc.UpdateGroup(context.Background(), "g1", "New Name", nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		err := uc.UpdateGroup(context.Background(), "g1", "", nil)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("group not found", func(t *testing.T) {
		err := uc.UpdateGroup(context.Background(), "nonexistent", "Name", nil)
		if err == nil || err.Error() != "group not found" {
			t.Error("expected 'group not found'")
		}
	})
}

func TestAddStudentToGroup(t *testing.T) {
	mock := mocks.NewGroupRepoMock()
	uc := NewGroupUseCase(mock)

	mock.GetByIDFunc = func(ctx context.Context, id string) (*domain.Group, error) {
		if id == "g1" {
			return &domain.Group{ID: "g1"}, nil
		}
		return nil, errors.New("not found")
	}
	mock.AddStudentToGroupFunc = func(ctx context.Context, groupID, studentID string) error {
		return nil
	}

	t.Run("success", func(t *testing.T) {
		err := uc.AddStudentToGroup(context.Background(), "g1", "s1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		err := uc.AddStudentToGroup(context.Background(), "bad", "s1")
		if err == nil || err.Error() != "group not found" {
			t.Error("expected 'group not found'")
		}
	})
}

func TestRemoveStudentFromGroup(t *testing.T) {
	mock := mocks.NewGroupRepoMock()
	uc := NewGroupUseCase(mock)

	mock.GetByIDFunc = func(ctx context.Context, id string) (*domain.Group, error) {
		if id == "g1" {
			return &domain.Group{ID: "g1"}, nil
		}
		return nil, errors.New("not found")
	}
	mock.RemoveStudentFromGroupFunc = func(ctx context.Context, groupID, studentID string) error {
		return nil
	}

	t.Run("success", func(t *testing.T) {
		err := uc.RemoveStudentFromGroup(context.Background(), "g1", "s1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		err := uc.RemoveStudentFromGroup(context.Background(), "bad", "s1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestChangeStudentGroup(t *testing.T) {
	mock := mocks.NewGroupRepoMock()
	uc := NewGroupUseCase(mock)

	mock.GetByIDFunc = func(ctx context.Context, id string) (*domain.Group, error) {
		if id == "g2" {
			return &domain.Group{ID: "g2"}, nil
		}
		return nil, errors.New("not found")
	}
	mock.ChangeStudentGroupFunc = func(ctx context.Context, studentID, newGroupID string) error {
		return nil
	}

	t.Run("success", func(t *testing.T) {
		err := uc.ChangeStudentGroup(context.Background(), "s1", "g2")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		err := uc.ChangeStudentGroup(context.Background(), "s1", "bad")
		if err == nil || err.Error() != "new group not found" {
			t.Error("expected 'new group not found'")
		}
	})
}

func TestChangeTeacherGroup(t *testing.T) {
	mock := mocks.NewGroupRepoMock()
	uc := NewGroupUseCase(mock)

	mock.GetByIDFunc = func(ctx context.Context, id string) (*domain.Group, error) {
		if id == "g2" {
			return &domain.Group{ID: "g2"}, nil
		}
		return nil, errors.New("not found")
	}
	mock.ChangeTeacherGroupFunc = func(ctx context.Context, teacherID, newGroupID string) error {
		return nil
	}

	t.Run("success", func(t *testing.T) {
		err := uc.ChangeTeacherGroup(context.Background(), "t1", "g2")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("group not found", func(t *testing.T) {
		err := uc.ChangeTeacherGroup(context.Background(), "t1", "bad")
		if err == nil || err.Error() != "new group not found" {
			t.Error("expected 'new group not found'")
		}
	})
}
