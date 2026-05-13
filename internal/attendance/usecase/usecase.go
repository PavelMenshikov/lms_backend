package usecase

import (
	"context"
	"lms_backend/internal/attendance/repository"
	"lms_backend/internal/domain"
	"time"
)

type AttendanceUseCase interface {
	GetStudentCalendar(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*domain.AttendanceRecord, error)
	MarkAttendance(ctx context.Context, lessonID, studentID string, status domain.AttendanceStatus, reason, comment *string, markedBy string) error
	UpdateAttendance(ctx context.Context, lessonID, studentID string, status domain.AttendanceStatus, reason, comment *string, updatedBy string) error
	GetLessonAttendance(ctx context.Context, lessonID string) ([]*domain.AttendanceRecord, error)
	GetStudentStats(ctx context.Context, studentID string) (map[string]int, error)
}

type attendanceUseCase struct {
	repo repository.AttendanceRepository
}

func NewAttendanceUseCase(repo repository.AttendanceRepository) AttendanceUseCase {
	return &attendanceUseCase{repo: repo}
}

func (uc *attendanceUseCase) GetStudentCalendar(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*domain.AttendanceRecord, error) {
	return uc.repo.GetByStudent(ctx, studentID, startDate, endDate)
}

func (uc *attendanceUseCase) MarkAttendance(ctx context.Context, lessonID, studentID string, status domain.AttendanceStatus, reason, comment *string, markedBy string) error {
	// Проверяем, существует ли уже запись
	existing, err := uc.repo.GetByLessonAndStudent(ctx, lessonID, studentID)
	if err != nil {
		return err
	}

	if existing != nil {
		// Обновляем существующую запись
		existing.Status = status
		existing.Reason = reason
		existing.Comment = comment
		existing.UpdatedBy = &markedBy
		return uc.repo.Update(ctx, existing)
	}

	// Создаём новую запись
	record := &domain.AttendanceRecord{
		LessonID:  lessonID,
		StudentID: studentID,
		Status:    status,
		Reason:    reason,
		Comment:   comment,
		MarkedBy:  &markedBy,
		UpdatedBy: &markedBy,
	}
	return uc.repo.Create(ctx, record)
}

func (uc *attendanceUseCase) UpdateAttendance(ctx context.Context, lessonID, studentID string, status domain.AttendanceStatus, reason, comment *string, updatedBy string) error {
	record := &domain.AttendanceRecord{
		LessonID:  lessonID,
		StudentID: studentID,
		Status:    status,
		Reason:    reason,
		Comment:   comment,
		UpdatedBy: &updatedBy,
	}
	return uc.repo.Update(ctx, record)
}

func (uc *attendanceUseCase) GetLessonAttendance(ctx context.Context, lessonID string) ([]*domain.AttendanceRecord, error) {
	return uc.repo.GetByLesson(ctx, lessonID)
}

func (uc *attendanceUseCase) GetStudentStats(ctx context.Context, studentID string) (map[string]int, error) {
	return uc.repo.GetStudentStats(ctx, studentID)
}
