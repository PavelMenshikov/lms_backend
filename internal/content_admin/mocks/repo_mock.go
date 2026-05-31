package mocks

import (
	"context"
	"lms_backend/internal/content_admin/repository"
	"lms_backend/internal/domain"
)

type ContentAdminRepoMock struct {
	CreatedUsers   map[string]*domain.User
	CreatedCourses map[string]*domain.Course
	LinkedParents  map[string]string
}

func NewContentAdminRepoMock() *ContentAdminRepoMock {
	return &ContentAdminRepoMock{
		CreatedUsers:   make(map[string]*domain.User),
		CreatedCourses: make(map[string]*domain.Course),
		LinkedParents:  make(map[string]string),
	}
}

var _ repository.ContentAdminRepository = (*ContentAdminRepoMock)(nil)

func (m *ContentAdminRepoMock) CreateCourse(ctx context.Context, course *domain.Course) (string, error) {
	id := "id"
	m.CreatedCourses[id] = course
	return id, nil
}
func (m *ContentAdminRepoMock) UpdateCourseSettings(ctx context.Context, course *domain.Course) error {
	return nil
}
func (m *ContentAdminRepoMock) CreateModule(ctx context.Context, module *domain.Module) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) DeleteModule(ctx context.Context, id string) error { return nil }
func (m *ContentAdminRepoMock) CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) DeleteLesson(ctx context.Context, id string) error { return nil }
func (m *ContentAdminRepoMock) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	return nil
}
func (m *ContentAdminRepoMock) AssignTeacherToLesson(ctx context.Context, lessonID, teacherID string) error {
	return nil
}
func (m *ContentAdminRepoMock) GetLessonIDByOrder(ctx context.Context, courseID string, orderNum int) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetModulesByCourseID(ctx context.Context, courseID string) ([]*domain.Module, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetParentsByStudentID(ctx context.Context, studentID string) ([]domain.User, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) LinkParentToStudent(ctx context.Context, studentID, parentID string) error {
	return nil
}
func (m *ContentAdminRepoMock) EnrollStudentExtended(ctx context.Context, userID, courseID, streamID, groupID string) error {
	return nil
}
func (m *ContentAdminRepoMock) GetCourseIDByStream(ctx context.Context, streamID string) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) CreateTest(ctx context.Context, test *domain.Test) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) DeleteTest(ctx context.Context, id string) error { return nil }
func (m *ContentAdminRepoMock) CreateProject(ctx context.Context, project *domain.Project) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) DeleteProject(ctx context.Context, id string) error      { return nil }
func (m *ContentAdminRepoMock) UpdateUser(ctx context.Context, user *domain.User) error { return nil }
func (m *ContentAdminRepoMock) DeleteUser(ctx context.Context, userID string) error     { return nil }
func (m *ContentAdminRepoMock) GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetDetailedModeratorList(ctx context.Context) ([]*domain.ModeratorTableItem, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetAllUsersList(ctx context.Context) ([]*domain.AllUsersTableItem, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) CreateStream(ctx context.Context, stream *domain.Stream) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) CreateGroup(ctx context.Context, group *domain.Group) (string, error) {
	return "id", nil
}
func (m *ContentAdminRepoMock) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetStudentEnrollment(ctx context.Context, userID string) (map[string]string, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) UnlinkAllParents(ctx context.Context, studentID string) error {
	return nil
}
func (m *ContentAdminRepoMock) SetLessonModule(ctx context.Context, lessonID, moduleID string) error {
	return nil
}
func (m *ContentAdminRepoMock) GetTestsByCourseID(ctx context.Context, courseID string) ([]domain.Test, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetProjectsByCourseID(ctx context.Context, courseID string) ([]domain.Project, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) UnenrollStudent(ctx context.Context, userID, courseID string) error {
	return nil
}
func (m *ContentAdminRepoMock) LinkTeachersToCourse(ctx context.Context, courseID string, teacherIDs []string) error {
	return nil
}
func (m *ContentAdminRepoMock) GetLessonByID(ctx context.Context, id string) (*domain.Lesson, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetTestByID(ctx context.Context, id string) (*domain.Test, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) GetProjectByID(ctx context.Context, id string) (*domain.Project, error) {
	return nil, nil
}
func (m *ContentAdminRepoMock) CancelLesson(ctx context.Context, lessonID, reason string) error {
	return nil
}
func (m *ContentAdminRepoMock) SubstituteTeacher(ctx context.Context, lessonID, teacherID string) error {
	return nil
}
