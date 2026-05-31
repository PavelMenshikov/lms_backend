package usecase

import (
	"context"
	"time"

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

func (uc *DashboardUseCase) GetCuratorDashboard(ctx context.Context, curatorID string) (*domain.CuratorDashboardData, error) {
	groups, err := uc.repo.GetCuratorGroups(ctx, curatorID)
	if err != nil {
		return nil, err
	}
	if groups == nil {
		groups = []domain.Group{}
	}

	attendance, _ := uc.repo.GetCuratorAttendanceStats(ctx, curatorID)
	if attendance == nil {
		attendance = []domain.CuratorGroupAttendance{}
	}

	homework, _ := uc.repo.GetCuratorHomeworkStats(ctx, curatorID)
	if homework == nil {
		homework = []domain.CuratorHomeworkStats{}
	}

	zones, _ := uc.repo.GetCuratorPerformanceZones(ctx, curatorID)

	return &domain.CuratorDashboardData{
		Groups:            groups,
		AttendanceByGroup: attendance,
		HomeworkByGroup:   homework,
		Performance:       zones,
	}, nil
}

func (uc *DashboardUseCase) GetAdminDashboard(ctx context.Context) (*domain.AdminHomeDashboard, error) {
	totalStudents, newStudents, studentsDelta, totalTeachers, activeCourses, err := uc.repo.GetAdminCounters(ctx)
	if err != nil {
		return nil, err
	}

	zones, _ := uc.repo.GetPerformanceStats(ctx)
	hwZones, _ := uc.repo.GetHwPerformanceStats(ctx)
	attZones, _ := uc.repo.GetAttendancePerformanceStats(ctx)
	activity, _ := uc.repo.GetLessonActivity(ctx)

	return &domain.AdminHomeDashboard{
		TotalStudents:         totalStudents,
		StudentsDelta:         studentsDelta,
		NewStudentsMonth:      newStudents,
		TotalTeachers:         totalTeachers,
		ActiveCourses:         activeCourses,
		Performance:           zones,
		HwPerformance:         hwZones,
		AttendancePerformance: attZones,
		LessonActivity:        activity,
		UpdatePeriodMonth:     time.Now().Format("January"),
	}, nil
}
