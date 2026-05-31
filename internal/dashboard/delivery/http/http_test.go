package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/dashboard/repository"
	"lms_backend/internal/dashboard/usecase"
	"lms_backend/internal/domain"
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
	return args.Get(0).(*domain.StatisticSummary), args.Error(1)
}
func (m *mockDashboardRepo) GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.StatisticSummary), args.Error(1)
}
func (m *mockDashboardRepo) GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.UpcomingLesson), args.Error(1)
}
func (m *mockDashboardRepo) GetAdminCounters(ctx context.Context) (totalStudents, newStudents int, studentsDelta float64, totalTeachers, activeCourses int, err error) {
	args := m.Called(ctx)
	return args.Int(0), args.Int(1), args.Get(2).(float64), args.Int(3), args.Int(4), args.Error(5)
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

func TestGetCuratorDashboard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)
		handler := NewDashboardHandler(uc)

		expectedGroups := []domain.Group{{ID: "g1", Title: "Group A"}}
		expectedAttendance := []domain.CuratorGroupAttendance{{
			GroupID: "g1", GroupTitle: "Group A",
			StudentCount: 15, AvgAttendance: 80,
		}}
		expectedHomework := []domain.CuratorHomeworkStats{{
			GroupID: "g1", GroupTitle: "Group A",
			AvgCompletion: 70, TotalSubmitted: 7, TotalAccepted: 5,
		}}
		expectedZones := domain.PerformanceZones{
			Green: 50, Yellow: 30, Red: 20,
		}

		mockRepo.On("GetCuratorGroups", mock.Anything, "curator-1").Return(expectedGroups, nil).Once()
		mockRepo.On("GetCuratorAttendanceStats", mock.Anything, "curator-1").Return(expectedAttendance, nil).Once()
		mockRepo.On("GetCuratorHomeworkStats", mock.Anything, "curator-1").Return(expectedHomework, nil).Once()
		mockRepo.On("GetCuratorPerformanceZones", mock.Anything, "curator-1").Return(expectedZones, nil).Once()

		req := httptest.NewRequest("GET", "/admin/curator/dashboard", nil)
		req = req.WithContext(context.WithValue(req.Context(),
			authMiddleware.ContextUserDataKey,
			&authMiddleware.UserContextData{UserID: "curator-1", Role: "curator"}))

		w := httptest.NewRecorder()
		handler.GetCuratorDashboard(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp domain.CuratorDashboardData
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Len(t, resp.Groups, 1)
		assert.Equal(t, "Group A", resp.Groups[0].Title)
		assert.Len(t, resp.AttendanceByGroup, 1)
		assert.Equal(t, 80.0, resp.AttendanceByGroup[0].AvgAttendance)
		assert.Len(t, resp.HomeworkByGroup, 1)
		assert.Equal(t, 7, resp.HomeworkByGroup[0].TotalSubmitted)
		assert.Equal(t, 50, resp.Performance.Green)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockRepo := new(mockDashboardRepo)
		uc := usecase.NewDashboardUseCase(mockRepo)
		handler := NewDashboardHandler(uc)

		mockRepo.On("GetCuratorGroups", mock.Anything, "curator-2").Return(nil, assert.AnError).Once()
		mockRepo.On("GetCuratorAttendanceStats", mock.Anything, "curator-2").Return([]domain.CuratorGroupAttendance{}, nil).Maybe()
		mockRepo.On("GetCuratorHomeworkStats", mock.Anything, "curator-2").Return([]domain.CuratorHomeworkStats{}, nil).Maybe()
		mockRepo.On("GetCuratorPerformanceZones", mock.Anything, "curator-2").Return(domain.PerformanceZones{}, nil).Maybe()

		req := httptest.NewRequest("GET", "/admin/curator/dashboard", nil)
		req = req.WithContext(context.WithValue(req.Context(),
			authMiddleware.ContextUserDataKey,
			&authMiddleware.UserContextData{UserID: "curator-2", Role: "curator"}))

		w := httptest.NewRecorder()
		handler.GetCuratorDashboard(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
