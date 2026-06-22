package usecase

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"lms_backend/internal/dashboard/repository"
	"lms_backend/internal/domain"
)

type DashboardUseCase struct {
	repo repository.UserDataRepository
}

func NewDashboardUseCase(repo repository.UserDataRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

type homeData struct {
	lastLesson   *domain.LastLesson
	coursesCount int
	attendance   *domain.StatisticSummary
	assignments  *domain.StatisticSummary
	upcoming     []domain.UpcomingLesson
}

func (uc *DashboardUseCase) GetUserHomeData(ctx context.Context, user *domain.User) (*domain.HomeDashboard, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var d homeData
	eg.Go(func() (err error) {
		d.lastLesson, err = uc.repo.GetLastLessonData(egCtx, user.ID)
		return err
	})
	eg.Go(func() (err error) {
		d.coursesCount, err = uc.repo.GetActiveCoursesCount(egCtx, user.ID)
		return err
	})
	eg.Go(func() (err error) {
		d.attendance, err = uc.repo.GetAttendancePercentage(egCtx, user.ID)
		return err
	})
	eg.Go(func() (err error) {
		d.assignments, err = uc.repo.GetAssignmentsCompletionPercentage(egCtx, user.ID)
		return err
	})
	eg.Go(func() (err error) {
		d.upcoming, err = uc.repo.GetUpcomingLessons(egCtx, user.ID)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &domain.HomeDashboard{
		UserRole:           user.Role,
		User:               user,
		LastLessonData:     d.lastLesson,
		ActiveCoursesCount: d.coursesCount,
		AttendanceStats:    d.attendance,
		AssignmentStats:    d.assignments,
		UpcomingLessons:    d.upcoming,
	}, nil
}

type curatorData struct {
	groups     []domain.Group
	attendance []domain.CuratorGroupAttendance
	homework   []domain.CuratorHomeworkStats
	zones      domain.PerformanceZones
}

func (uc *DashboardUseCase) GetCuratorDashboard(ctx context.Context, curatorID string) (*domain.CuratorDashboardData, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var d curatorData
	eg.Go(func() (err error) {
		d.groups, err = uc.repo.GetCuratorGroups(egCtx, curatorID)
		if err == nil && d.groups == nil {
			d.groups = []domain.Group{}
		}
		return err
	})
	eg.Go(func() (err error) {
		d.attendance, err = uc.repo.GetCuratorAttendanceStats(egCtx, curatorID)
		if err == nil && d.attendance == nil {
			d.attendance = []domain.CuratorGroupAttendance{}
		}
		return err
	})
	eg.Go(func() (err error) {
		d.homework, err = uc.repo.GetCuratorHomeworkStats(egCtx, curatorID)
		if err == nil && d.homework == nil {
			d.homework = []domain.CuratorHomeworkStats{}
		}
		return err
	})
	eg.Go(func() (err error) {
		d.zones, err = uc.repo.GetCuratorPerformanceZones(egCtx, curatorID)
		if err != nil {
			d.zones = domain.PerformanceZones{}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &domain.CuratorDashboardData{
		Groups:            d.groups,
		AttendanceByGroup: d.attendance,
		HomeworkByGroup:   d.homework,
		Performance:       d.zones,
	}, nil
}

func (uc *DashboardUseCase) GetAdminDashboard(ctx context.Context) (*domain.AdminHomeDashboard, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	var (
		totalStudents int
		newStudents   int
		studentsDelta float64
		totalTeachers int
		activeCourses int
		allPerf       *domain.AllPerformanceStats
		activity      []domain.DailyLessonActivity
	)

	eg.Go(func() (err error) {
		totalStudents, newStudents, studentsDelta, totalTeachers, activeCourses, err = uc.repo.GetAdminCounters(egCtx)
		return err
	})
	eg.Go(func() (err error) {
		allPerf, err = uc.repo.GetAllPerformanceStats(egCtx)
		return err
	})
	eg.Go(func() (err error) {
		activity, err = uc.repo.GetLessonActivity(egCtx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &domain.AdminHomeDashboard{
		TotalStudents:         totalStudents,
		StudentsDelta:         studentsDelta,
		NewStudentsMonth:      newStudents,
		TotalTeachers:         totalTeachers,
		ActiveCourses:         activeCourses,
		Performance:           allPerf.CourseZones,
		HwPerformance:         allPerf.HomeworkZones,
		AttendancePerformance: allPerf.AttendanceZones,
		LessonActivity:        activity,
		UpdatePeriodMonth:     time.Now().Format("January"),
	}, nil
}
