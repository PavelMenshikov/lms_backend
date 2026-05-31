package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"lms_backend/internal/content_admin/usecase"
	"lms_backend/internal/domain"
)

type mockService struct {
	mock.Mock
	ContentAdminService
}

func (m *mockService) CreateCourse(ctx context.Context, input usecase.CreateCourseInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *mockService) CreateLesson(ctx context.Context, input usecase.CreateLessonInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *mockService) GetLesson(ctx context.Context, id string) (*domain.Lesson, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Lesson), args.Error(1)
}
func (m *mockService) UpdateLesson(ctx context.Context, id string, input usecase.CreateLessonInput) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func TestCreateLesson_JSON(t *testing.T) {
	m := new(mockService)
	h := &ContentAdminHandler{uc: m}

	reqBody := CreateLessonRequest{
		Title: "New Lesson",
		Content: []domain.ContentBlock{
			{Type: "text", Content: "Hello"},
		},
	}
	body, _ := json.Marshal(reqBody)

	m.On("CreateLesson", mock.Anything, mock.MatchedBy(func(i usecase.CreateLessonInput) bool {
		return i.Title == "New Lesson"
	})).Return("new-uuid", nil)

	req := httptest.NewRequest("POST", "/admin/lessons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateLesson(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "new-uuid")
	m.AssertExpectations(t)
}
