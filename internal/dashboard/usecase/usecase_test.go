package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"lms_backend/internal/dashboard/repository"
	"lms_backend/internal/dashboard/usecase"
	"lms_backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDashboardRepo struct {
	mock.Mock
}

var _ repository.UserDataRepository = (*mockDashboardRepo)(nil)

func (m *mockDashboardRepo) GetLastLessonData(ctx context.Context, userID string) (*domain.LastLesson, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LastLesson), args.Error(1)
}

func (m *mockDashboardRepo) GetActiveCoursesCount(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *mockDashboardRepo) GetAttendancePercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatisticSummary), args.Error(1)
}

func (m *mockDashboardRepo) GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatisticSummary), args.Error(1)
}

func (m *mockDashboardRepo) GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UpcomingLesson), args.Error(1)
}

func (m *mockDashboardRepo) GetAdminCounters(ctx context.Context) (totalStudents, newStudents int, studentsDelta float64, totalTeachers, activeCourses int, err error) {
	args := m.Called(ctx)
	return args.Int(0), args.Int(1), args.Get(2).(float64), args.Int(3), args.Int(4), args.Error(5)
}

func (m *mockDashboardRepo) GetAllPerformanceStats(ctx context.Context) (*domain.AllPerformanceStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AllPerformanceStats), args.Error(1)
}

func (m *mockDashboardRepo) GetPerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.PerformanceZones), args.Error(1)
}

func (m *mockDashboardRepo) GetHwPerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.PerformanceZones), args.Error(1)
}

func (m *mockDashboardRepo) GetAttendancePerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.PerformanceZones), args.Error(1)
}

func (m *mockDashboardRepo) GetLessonActivity(ctx context.Context) ([]domain.DailyLessonActivity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DailyLessonActivity), args.Error(1)
}

func (m *mockDashboardRepo) GetCuratorGroups(ctx context.Context, curatorID string) ([]domain.Group, error) {
	args := m.Called(ctx, curatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Group), args.Error(1)
}

func (m *mockDashboardRepo) GetCuratorAttendanceStats(ctx context.Context, curatorID string) ([]domain.CuratorGroupAttendance, error) {
	args := m.Called(ctx, curatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CuratorGroupAttendance), args.Error(1)
}

func (m *mockDashboardRepo) GetCuratorHomeworkStats(ctx context.Context, curatorID string) ([]domain.CuratorHomeworkStats, error) {
	args := m.Called(ctx, curatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CuratorHomeworkStats), args.Error(1)
}

func (m *mockDashboardRepo) GetCuratorPerformanceZones(ctx context.Context, curatorID string) (domain.PerformanceZones, error) {
	args := m.Called(ctx, curatorID)
	return args.Get(0).(domain.PerformanceZones), args.Error(1)
}

func TestDashboardUseCase_GetUserHomeData(t *testing.T) {
	user := &domain.User{
		ID:    "user-1",
		Role:  domain.RoleStudent,
		Email: "student@test.ru",
	}
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		lastLesson := &domain.LastLesson{
			CourseTitle:      "Go Basics",
			ModuleName:       "Module 1",
			LessonTitle:      "Variables",
			AssignmentStatus: "pending",
			LessonID:         "lesson-1",
			HomeworkID:       "hw-1",
		}
		attendance := &domain.StatisticSummary{Percentage: 75.0, Delta: 5.0}
		assignment := &domain.StatisticSummary{Percentage: 60.0, Delta: -10.0}
		upcoming := []domain.UpcomingLesson{
			{Date: now, CourseTitle: "Go Basics", TeacherName: "John Doe"},
		}

		mockRepo.On("GetLastLessonData", mock.Anything, "user-1").Return(lastLesson, nil).Once()
		mockRepo.On("GetActiveCoursesCount", mock.Anything, "user-1").Return(3, nil).Once()
		mockRepo.On("GetAttendancePercentage", mock.Anything, "user-1").Return(attendance, nil).Once()
		mockRepo.On("GetAssignmentsCompletionPercentage", mock.Anything, "user-1").Return(assignment, nil).Once()
		mockRepo.On("GetUpcomingLessons", mock.Anything, "user-1").Return(upcoming, nil).Once()

		result, err := uc.GetUserHomeData(context.Background(), user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.RoleStudent, result.UserRole)
		assert.Equal(t, 3, result.ActiveCoursesCount)
		assert.Equal(t, 75.0, result.AttendanceStats.Percentage)
		assert.Equal(t, 60.0, result.AssignmentStats.Percentage)
		assert.Len(t, result.UpcomingLessons, 1)
		assert.Equal(t, "Go Basics", result.LastLessonData.CourseTitle)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error propagated", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		expectedErr := errors.New("db connection failed")
		mockRepo.On("GetLastLessonData", mock.Anything, "user-1").Return(nil, expectedErr).Once()
		mockRepo.On("GetActiveCoursesCount", mock.Anything, "user-1").Return(0, nil).Maybe()
		mockRepo.On("GetAttendancePercentage", mock.Anything, "user-1").Return(nil, nil).Maybe()
		mockRepo.On("GetAssignmentsCompletionPercentage", mock.Anything, "user-1").Return(nil, nil).Maybe()
		mockRepo.On("GetUpcomingLessons", mock.Anything, "user-1").Return([]domain.UpcomingLesson{}, nil).Maybe()

		result, err := uc.GetUserHomeData(context.Background(), user)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("partial data with nil last lesson", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		mockRepo.On("GetLastLessonData", mock.Anything, "user-1").Return(nil, nil).Once()
		mockRepo.On("GetActiveCoursesCount", mock.Anything, "user-1").Return(0, nil).Once()
		mockRepo.On("GetAttendancePercentage", mock.Anything, "user-1").Return(&domain.StatisticSummary{Percentage: 0}, nil).Once()
		mockRepo.On("GetAssignmentsCompletionPercentage", mock.Anything, "user-1").Return(&domain.StatisticSummary{Percentage: 0}, nil).Once()
		mockRepo.On("GetUpcomingLessons", mock.Anything, "user-1").Return([]domain.UpcomingLesson{}, nil).Once()

		result, err := uc.GetUserHomeData(context.Background(), user)
		assert.NoError(t, err)
		assert.Nil(t, result.LastLessonData)
		assert.Equal(t, 0, result.ActiveCoursesCount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("admin role", func(t *testing.T) {
		admin := &domain.User{ID: "admin-1", Role: domain.RoleAdmin}
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		mockRepo.On("GetLastLessonData", mock.Anything, "admin-1").Return(nil, nil).Once()
		mockRepo.On("GetActiveCoursesCount", mock.Anything, "admin-1").Return(0, nil).Once()
		mockRepo.On("GetAttendancePercentage", mock.Anything, "admin-1").Return(&domain.StatisticSummary{Percentage: 0}, nil).Once()
		mockRepo.On("GetAssignmentsCompletionPercentage", mock.Anything, "admin-1").Return(&domain.StatisticSummary{Percentage: 0}, nil).Once()
		mockRepo.On("GetUpcomingLessons", mock.Anything, "admin-1").Return([]domain.UpcomingLesson{}, nil).Once()

		result, err := uc.GetUserHomeData(context.Background(), admin)
		assert.NoError(t, err)
		assert.Equal(t, domain.RoleAdmin, result.UserRole)
		mockRepo.AssertExpectations(t)
	})
}

func TestDashboardUseCase_GetCuratorDashboard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		groups := []domain.Group{
			{ID: "group-1", Title: "Group A"},
			{ID: "group-2", Title: "Group B"},
		}
		attendance := []domain.CuratorGroupAttendance{
			{GroupID: "group-1", GroupTitle: "Group A", StudentCount: 10, AvgAttendance: 85.0},
		}
		homework := []domain.CuratorHomeworkStats{
			{GroupID: "group-1", GroupTitle: "Group A", AvgCompletion: 70.0, TotalSubmitted: 8, TotalAccepted: 5},
		}
		zones := domain.PerformanceZones{Green: 5, Yellow: 3, Red: 2}

		mockRepo.On("GetCuratorGroups", mock.Anything, "curator-1").Return(groups, nil).Once()
		mockRepo.On("GetCuratorAttendanceStats", mock.Anything, "curator-1").Return(attendance, nil).Once()
		mockRepo.On("GetCuratorHomeworkStats", mock.Anything, "curator-1").Return(homework, nil).Once()
		mockRepo.On("GetCuratorPerformanceZones", mock.Anything, "curator-1").Return(zones, nil).Once()

		result, err := uc.GetCuratorDashboard(context.Background(), "curator-1")

		assert.NoError(t, err)
		assert.Len(t, result.Groups, 2)
		assert.Len(t, result.AttendanceByGroup, 1)
		assert.Equal(t, 85.0, result.AttendanceByGroup[0].AvgAttendance)
		assert.Len(t, result.HomeworkByGroup, 1)
		assert.Equal(t, 5, result.Performance.Green)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		mockRepo.On("GetCuratorGroups", mock.Anything, "curator-2").Return(nil, errors.New("not found")).Once()
		mockRepo.On("GetCuratorAttendanceStats", mock.Anything, "curator-2").Return([]domain.CuratorGroupAttendance{}, nil).Maybe()
		mockRepo.On("GetCuratorHomeworkStats", mock.Anything, "curator-2").Return([]domain.CuratorHomeworkStats{}, nil).Maybe()
		mockRepo.On("GetCuratorPerformanceZones", mock.Anything, "curator-2").Return(domain.PerformanceZones{}, nil).Maybe()

		result, err := uc.GetCuratorDashboard(context.Background(), "curator-2")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDashboardUseCase_GetAdminDashboard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		stats := &domain.AllPerformanceStats{
			CourseZones:     domain.PerformanceZones{Green: 10, Yellow: 5, Red: 3},
			HomeworkZones:   domain.PerformanceZones{Green: 8, Yellow: 4, Red: 6},
			AttendanceZones: domain.PerformanceZones{Green: 12, Yellow: 3, Red: 3},
		}
		activity := []domain.DailyLessonActivity{
			{Date: "2026-06-01", Group: 5, Trial: 2, Individual: 1},
		}

		mockRepo.On("GetAdminCounters", mock.Anything).Return(100, 10, 15.5, 20, 8, nil).Once()
		mockRepo.On("GetAllPerformanceStats", mock.Anything).Return(stats, nil).Once()
		mockRepo.On("GetLessonActivity", mock.Anything).Return(activity, nil).Once()

		result, err := uc.GetAdminDashboard(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 100, result.TotalStudents)
		assert.Equal(t, 10, result.NewStudentsMonth)
		assert.Equal(t, 15.5, result.StudentsDelta)
		assert.Equal(t, 20, result.TotalTeachers)
		assert.Equal(t, 8, result.ActiveCourses)
		assert.Equal(t, 10, result.Performance.Green)
		assert.Len(t, result.LessonActivity, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)

		mockRepo.On("GetAdminCounters", mock.Anything).Return(0, 0, 0.0, 0, 0, errors.New("db error")).Once()
		mockRepo.On("GetAllPerformanceStats", mock.Anything).Return(nil, nil).Maybe()
		mockRepo.On("GetLessonActivity", mock.Anything).Return([]domain.DailyLessonActivity{}, nil).Maybe()

		result, err := uc.GetAdminDashboard(context.Background())
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
