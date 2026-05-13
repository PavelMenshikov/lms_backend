package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type GroupRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Group, error)
	UpdateGroup(ctx context.Context, groupID, name string, teacherID *string) error
	AddStudentToGroup(ctx context.Context, groupID, studentID string) error
	RemoveStudentFromGroup(ctx context.Context, groupID, studentID string) error
	ChangeStudentGroup(ctx context.Context, studentID, newGroupID string) error
	ChangeTeacherGroup(ctx context.Context, teacherID, newGroupID string) error
	GetGroupStudents(ctx context.Context, groupID string) ([]string, error)
}

type groupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) GetByID(ctx context.Context, id string) (*domain.Group, error) {
	var group domain.Group
	query := `SELECT id, title, stream_id, created_at FROM groups WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&group.ID, &group.Title, &group.StreamID, &group.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) UpdateGroup(ctx context.Context, groupID, name string, teacherID *string) error {
	query := `UPDATE groups SET title = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, name, groupID)
	return err
}

func (r *groupRepository) AddStudentToGroup(ctx context.Context, groupID, studentID string) error {
	// Проверяем, что ученик не в другой группе этого же курса
	query := `
		INSERT INTO user_enrollments (user_id, course_id, enrolled_at)
		SELECT $1, c.id, CURRENT_TIMESTAMP
		FROM groups g
		JOIN streams s ON g.stream_id = s.id
		JOIN courses c ON s.course_id = c.id
		WHERE g.id = $2
		ON CONFLICT (user_id, course_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, studentID, groupID)
	if err != nil {
		return err
	}

	// Обновляем group_id в lessons для этого ученика
	updateQuery := `
		UPDATE lessons
		SET group_id = $1
		WHERE id IN (
			SELECT l.id FROM lessons l
			JOIN courses c ON l.course_id = c.id
			JOIN streams s ON s.course_id = c.id
			JOIN groups g ON g.stream_id = s.id
			WHERE g.id = $1
		)
	`
	_, err = r.db.ExecContext(ctx, updateQuery, groupID)
	return err
}

func (r *groupRepository) RemoveStudentFromGroup(ctx context.Context, groupID, studentID string) error {
	// Удаляем enrollment для курса этой группы
	query := `
		DELETE FROM user_enrollments
		WHERE user_id = $1
		AND course_id IN (
			SELECT c.id FROM courses c
			JOIN streams s ON s.course_id = c.id
			JOIN groups g ON g.stream_id = s.id
			WHERE g.id = $2
		)
	`
	_, err := r.db.ExecContext(ctx, query, studentID, groupID)
	return err
}

func (r *groupRepository) ChangeStudentGroup(ctx context.Context, studentID, newGroupID string) error {
	// Сначала удаляем из текущей группы
	deleteQuery := `
		DELETE FROM user_enrollments
		WHERE user_id = $1
		AND course_id IN (
			SELECT c.id FROM courses c
			JOIN streams s ON s.course_id = c.id
			JOIN groups g ON g.stream_id = s.id
			WHERE g.id = $2
		)
	`
	_, err := r.db.ExecContext(ctx, deleteQuery, studentID, newGroupID)
	if err != nil {
		return err
	}

	// Добавляем в новую группу
	return r.AddStudentToGroup(ctx, newGroupID, studentID)
}

func (r *groupRepository) ChangeTeacherGroup(ctx context.Context, teacherID, newGroupID string) error {
	// Обновляем teacher_id в lessons для новой группы
	query := `
		UPDATE lessons
		SET teacher_id = $1
		WHERE group_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, teacherID, newGroupID)
	return err
}

func (r *groupRepository) GetGroupStudents(ctx context.Context, groupID string) ([]string, error) {
	query := `
		SELECT DISTINCT ue.user_id
		FROM user_enrollments ue
		JOIN courses c ON ue.course_id = c.id
		JOIN streams s ON s.course_id = c.id
		JOIN groups g ON g.stream_id = s.id
		WHERE g.id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentIDs []string
	for rows.Next() {
		var studentID string
		if err := rows.Scan(&studentID); err != nil {
			return nil, err
		}
		studentIDs = append(studentIDs, studentID)
	}
	return studentIDs, nil
}
