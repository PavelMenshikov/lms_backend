package mocks

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/teacher_dashboard/repository"
)

type TeacherDashboardRepositoryMock struct {
	Report *domain.TeacherMonthlyReport
	Err    error
}

var _ repository.TeacherDashboardRepository = (*TeacherDashboardRepositoryMock)(nil)

func NewTeacherDashboardRepositoryMock() *TeacherDashboardRepositoryMock {
	return &TeacherDashboardRepositoryMock{}
}

func (r *TeacherDashboardRepositoryMock) GetMonthlyReport(ctx context.Context, teacherID string, year, month int) (*domain.TeacherMonthlyReport, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	if r.Report != nil {
		return r.Report, nil
	}
	return &domain.TeacherMonthlyReport{
		TeacherID:          teacherID,
		Year:               year,
		Month:              month,
		TotalLessons:       10,
		SubstitutionsCount: 2,
		ReplacedCount:      1,
		AvgRating:          4.5,
		TotalStudents:      25,
		AttendanceAvg:      85.5,
	}, nil
}
