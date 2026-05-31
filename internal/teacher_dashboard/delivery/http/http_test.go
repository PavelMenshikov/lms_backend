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
	"lms_backend/internal/domain"
	"lms_backend/internal/teacher_dashboard/repository"
	"lms_backend/internal/teacher_dashboard/usecase"
)

type mockTeacherDashboardRepo struct {
	mock.Mock
}

var _ repository.TeacherDashboardRepository = (*mockTeacherDashboardRepo)(nil)

func (m *mockTeacherDashboardRepo) GetMonthlyReport(ctx context.Context, teacherID string, year, month int) (*domain.TeacherMonthlyReport, error) {
	args := m.Called(ctx, teacherID, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TeacherMonthlyReport), args.Error(1)
}

func TestGetTeacherMonthlyReport(t *testing.T) {
	mockRepo := new(mockTeacherDashboardRepo)
	uc := usecase.NewTeacherDashboardUseCase(mockRepo)
	handler := NewTeacherDashboardHandler(uc)

	t.Run("success without params", func(t *testing.T) {
		expected := &domain.TeacherMonthlyReport{
			TeacherID:         "teacher-1",
			Year:              2026,
			Month:             5,
			TotalLessons:      10,
			SubstitutionsCount: 2,
			ReplacedCount:     1,
			AvgRating:         4.5,
			TotalStudents:     25,
			AttendanceAvg:     85.5,
		}
		mockRepo.On("GetMonthlyReport", mock.Anything, "teacher-1", 2026, 5).
			Return(expected, nil).Once()

		req := httptest.NewRequest("GET", "/teacher/monthly-report", nil)
		req = req.WithContext(context.WithValue(req.Context(),
			authMiddleware.ContextUserDataKey,
			&authMiddleware.UserContextData{UserID: "teacher-1", Role: domain.RoleTeacher}))

		w := httptest.NewRecorder()
		handler.GetTeacherMonthlyReport(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp domain.TeacherMonthlyReport
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, 10, resp.TotalLessons)
		assert.Equal(t, 2, resp.SubstitutionsCount)
		assert.Equal(t, 4.5, resp.AvgRating)
	})

	t.Run("success with params", func(t *testing.T) {
		expected := &domain.TeacherMonthlyReport{
			TeacherID:    "teacher-1",
			Year:         2026,
			Month:        3,
			TotalLessons: 8,
		}
		mockRepo.On("GetMonthlyReport", mock.Anything, "teacher-1", 2026, 3).
			Return(expected, nil).Once()

		req := httptest.NewRequest("GET", "/teacher/monthly-report?year=2026&month=3", nil)
		req = req.WithContext(context.WithValue(req.Context(),
			authMiddleware.ContextUserDataKey,
			&authMiddleware.UserContextData{UserID: "teacher-1", Role: domain.RoleTeacher}))

		w := httptest.NewRecorder()
		handler.GetTeacherMonthlyReport(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp domain.TeacherMonthlyReport
		json.NewDecoder(w.Body).Decode(&resp)
		assert.Equal(t, 8, resp.TotalLessons)
		assert.Equal(t, 3, resp.Month)
		assert.Equal(t, 2026, resp.Year)
	})

	t.Run("forbidden for non-teacher", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/teacher/monthly-report", nil)
		req = req.WithContext(context.WithValue(req.Context(),
			authMiddleware.ContextUserDataKey,
			&authMiddleware.UserContextData{UserID: "student-1", Role: domain.Role("student")}))

		w := httptest.NewRecorder()
		handler.GetTeacherMonthlyReport(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockRepo.On("GetMonthlyReport", mock.Anything, "teacher-1", 2026, 5).
			Return(nil, assert.AnError).Once()

		req := httptest.NewRequest("GET", "/teacher/monthly-report", nil)
		req = req.WithContext(context.WithValue(req.Context(),
			authMiddleware.ContextUserDataKey,
			&authMiddleware.UserContextData{UserID: "teacher-1", Role: domain.RoleTeacher}))

		w := httptest.NewRecorder()
		handler.GetTeacherMonthlyReport(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("missing context data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/teacher/monthly-report", nil)

		w := httptest.NewRecorder()
		handler.GetTeacherMonthlyReport(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
