package usecase

import (
	"context"

	"lms_backend/internal/dashboard/repository"
	"lms_backend/internal/domain"
)

type DashboardUseCase struct {
	repo repository.UserDataRepository
}

func NewDashboardUseCase(repo repository.UserDataRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

func (uc *DashboardUseCase) GetUserHomeData(ctx context.Context, user *domain.User) (*domain.HomeDashboard, error) {
	lastLesson, _ := uc.repo.GetLastLessonData(ctx, user.ID)
	coursesCount, _ := uc.repo.GetActiveCoursesCount(ctx, user.ID)
	attendance, _ := uc.repo.GetAttendancePercentage(ctx, user.ID)
	assignments, _ := uc.repo.GetAssignmentsCompletionPercentage(ctx, user.ID)
	upcoming, _ := uc.repo.GetUpcomingLessons(ctx, user.ID)

	return &domain.HomeDashboard{
		UserRole:           user.Role,
		User:               user,
		LastLessonData:     lastLesson,
		ActiveCoursesCount: coursesCount,
		AttendanceStats:    attendance,
		AssignmentStats:    assignments,
		UpcomingLessons:    upcoming,
	}, nil
}

func (uc *DashboardUseCase) GetAdminDashboard(ctx context.Context) (*domain.AdminHomeDashboard, error) {
	totalStudents, newStudents, totalTeachers, activeCourses, err := uc.repo.GetAdminCounters(ctx)
	if err != nil {
		return nil, err
	}

	zones, _ := uc.repo.GetPerformanceStats(ctx)
	activity, _ := uc.repo.GetLessonActivity(ctx)

	return &domain.AdminHomeDashboard{
		TotalStudents:    totalStudents,
		NewStudentsMonth: newStudents,
		TotalTeachers:    totalTeachers,
		ActiveCourses:    activeCourses,
		Performance:      zones,
		LessonActivity:   activity,
		StudentsDelta:    5.4,
	}, nil
}
