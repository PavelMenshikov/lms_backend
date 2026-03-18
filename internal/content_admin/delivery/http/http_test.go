package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"lms_backend/internal/domain"
	"lms_backend/internal/content_admin/usecase"
)

type MockContentAdminUseCase struct {
	mock.Mock
}

func (m *MockContentAdminUseCase) GetLesson(ctx context.Context, id string) (*domain.Lesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lesson), args.Error(1)
}

func TestGetLessonHandler(t *testing.T) {
	mockUC := new(MockContentAdminUseCase)
	handler := &ContentAdminHandler{uc: mockUC}

	expectedLesson := &domain.Lesson{
		ID:    "test-uuid",
		Title: "Test Lesson",
	}

	mockUC.On("GetLesson", mock.Anything, "test-uuid").Return(expectedLesson, nil)

	req := httptest.NewRequest("GET", "/admin/lessons/test-uuid", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "test-uuid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler.GetLesson(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response domain.Lesson
	json.NewDecoder(w.Body).Decode(&response)
	assert.Equal(t, "Test Lesson", response.Title)
}