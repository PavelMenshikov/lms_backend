package http

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"lms_backend/internal/content_admin/usecase"
	"lms_backend/internal/domain"
)

type MockContentAdminUseCase struct {
	mock.Mock
}

func (m *MockContentAdminUseCase) CreateCourse(ctx context.Context, input usecase.CreateCourseInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) UploadMedia(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, fh)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) UpdateCourseSettings(ctx context.Context, input usecase.UpdateCourseSettingsInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Course), args.Error(1)
}
func (m *MockContentAdminUseCase) GetCourseStructure(ctx context.Context, courseID string) (*domain.CourseStructure, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).(*domain.CourseStructure), args.Error(1)
}
func (m *MockContentAdminUseCase) CreateModule(ctx context.Context, input usecase.CreateModuleInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) DeleteModule(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) CreateLesson(ctx context.Context, input usecase.CreateLessonInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) DeleteLesson(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) CreateTest(ctx context.Context, input usecase.CreateTestInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) DeleteTest(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) CreateProject(ctx context.Context, input usecase.CreateProjectInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) DeleteProject(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) CreateModulesBulk(ctx context.Context, input []usecase.CreateModuleInput) ([]string, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockContentAdminUseCase) CreateLessonsBulk(ctx context.Context, input []usecase.CreateLessonInput) ([]string, error) {
	args := m.Called(ctx, input)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockContentAdminUseCase) CreateFullCourse(ctx context.Context, input usecase.CreateBulkCourseInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) CreateFullUser(ctx context.Context, input usecase.ExtendedCreateUserInput) (map[string]string, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(map[string]string), args.Error(1)
}
func (m *MockContentAdminUseCase) GetUserInfo(ctx context.Context, userID string) (map[string]interface{}, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
func (m *MockContentAdminUseCase) UpdateUser(ctx context.Context, userID string, input usecase.ExtendedCreateUserInput) error {
	args := m.Called(ctx, userID, input)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) GetUsersList(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.User), args.Error(1)
}
func (m *MockContentAdminUseCase) GetDetailedStudents(ctx context.Context, courseID string) ([]*domain.StudentTableItem, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]*domain.StudentTableItem), args.Error(1)
}
func (m *MockContentAdminUseCase) GetDetailedTeachers(ctx context.Context) ([]*domain.TeacherTableItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.TeacherTableItem), args.Error(1)
}
func (m *MockContentAdminUseCase) GetDetailedCurators(ctx context.Context) ([]*domain.CuratorTableItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.CuratorTableItem), args.Error(1)
}
func (m *MockContentAdminUseCase) GetDetailedModerators(ctx context.Context) ([]*domain.ModeratorTableItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.ModeratorTableItem), args.Error(1)
}
func (m *MockContentAdminUseCase) GetAllUsersTable(ctx context.Context) ([]*domain.AllUsersTableItem, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.AllUsersTableItem), args.Error(1)
}
func (m *MockContentAdminUseCase) EnrollStudent(ctx context.Context, userID, courseID string) error {
	args := m.Called(ctx, userID, courseID)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]*domain.AdminStudentProgress), args.Error(1)
}
func (m *MockContentAdminUseCase) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).(*domain.AdminCourseStats), args.Error(1)
}
func (m *MockContentAdminUseCase) CreateStream(ctx context.Context, input usecase.CreateStreamInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error) {
	args := m.Called(ctx, courseID)
	return args.Get(0).([]*domain.Stream), args.Error(1)
}
func (m *MockContentAdminUseCase) CreateGroup(ctx context.Context, input usecase.CreateGroupInput) (string, error) {
	args := m.Called(ctx, input)
	return args.String(0), args.Error(1)
}
func (m *MockContentAdminUseCase) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	args := m.Called(ctx, streamID)
	return args.Get(0).([]*domain.Group), args.Error(1)
}
func (m *MockContentAdminUseCase) UnenrollStudent(ctx context.Context, userID, courseID string) error {
	args := m.Called(ctx, userID, courseID)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) UpdateLesson(ctx context.Context, lessonID string, input usecase.CreateLessonInput) error {
	args := m.Called(ctx, lessonID, input)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) GetLesson(ctx context.Context, id string) (*domain.Lesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Lesson), args.Error(1)
}
func (m *MockContentAdminUseCase) GetTest(ctx context.Context, id string) (*domain.Test, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Test), args.Error(1)
}
func (m *MockContentAdminUseCase) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Project), args.Error(1)
}
func (m *MockContentAdminUseCase) LinkTeachersToCourse(ctx context.Context, courseID string, teacherIDs []string) error {
	args := m.Called(ctx, courseID, teacherIDs)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) CancelLesson(ctx context.Context, lessonID, reason string) error {
	args := m.Called(ctx, lessonID, reason)
	return args.Error(0)
}
func (m *MockContentAdminUseCase) SubstituteTeacher(ctx context.Context, lessonID, teacherID string) error {
	args := m.Called(ctx, lessonID, teacherID)
	return args.Error(0)
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
