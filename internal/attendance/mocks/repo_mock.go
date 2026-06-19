package mocks

import (
	"context"
	"lms_backend/internal/attendance/repository"
	"lms_backend/internal/domain"
	"sync"
	"time"
)

type AttendanceRepositoryMock struct {
	mu      sync.Mutex
	Records map[string]*domain.AttendanceRecord
	nextID  int
}

var _ repository.AttendanceRepository = (*AttendanceRepositoryMock)(nil)

func NewAttendanceRepositoryMock() *AttendanceRepositoryMock {
	return &AttendanceRepositoryMock{
		Records: make(map[string]*domain.AttendanceRecord),
		nextID:  1,
	}
}

func (r *AttendanceRepositoryMock) GetByLessonAndStudent(ctx context.Context, lessonID, studentID string) (*domain.AttendanceRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rec := range r.Records {
		if rec.LessonID == lessonID && rec.StudentID == studentID {
			return rec, nil
		}
	}
	return nil, nil
}

func (r *AttendanceRepositoryMock) GetByStudent(ctx context.Context, studentID string, startDate, endDate time.Time) ([]*domain.AttendanceRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.AttendanceRecord
	for _, rec := range r.Records {
		if rec.StudentID == studentID {
			result = append(result, rec)
		}
	}
	return result, nil
}

func (r *AttendanceRepositoryMock) GetByLesson(ctx context.Context, lessonID string) ([]*domain.AttendanceRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []*domain.AttendanceRecord
	for _, rec := range r.Records {
		if rec.LessonID == lessonID {
			result = append(result, rec)
		}
	}
	return result, nil
}

func (r *AttendanceRepositoryMock) Create(ctx context.Context, record *domain.AttendanceRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	record.ID = "att-" + r.nextIDStr()
	r.Records[record.ID] = record
	return nil
}

func (r *AttendanceRepositoryMock) nextIDStr() string {
	id := r.nextID
	r.nextID++
	return string(rune('0'+id%10)) + string(rune('0'+(id/10)%10))
}

func (r *AttendanceRepositoryMock) Update(ctx context.Context, record *domain.AttendanceRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rec := range r.Records {
		if rec.LessonID == record.LessonID && rec.StudentID == record.StudentID {
			rec.Status = record.Status
			rec.Reason = record.Reason
			rec.Comment = record.Comment
			rec.UpdatedBy = record.UpdatedBy
			return nil
		}
	}
	return nil
}

func (r *AttendanceRepositoryMock) GetStudentStats(ctx context.Context, studentID string) (map[string]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stats := map[string]int{"attended": 0, "absent_excused": 0, "absent_unexcused": 0, "freeze": 0}
	for _, rec := range r.Records {
		if rec.StudentID == studentID {
			switch rec.Status {
			case domain.AttendanceStatusAttended:
				stats["attended"]++
			case domain.AttendanceStatusAbsentExcused:
				stats["absent_excused"]++
			case domain.AttendanceStatusAbsentUnexcused:
				stats["absent_unexcused"]++
			case domain.AttendanceStatusFreeze:
				stats["freeze"]++
			}
		}
	}
	return stats, nil
}
