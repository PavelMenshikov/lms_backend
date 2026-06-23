package repository

import (
	"context"
	"database/sql"

	"lms_backend/internal/domain"
)

type ContentAdminRepository interface {
	CreateCourse(ctx context.Context, course *domain.Course) (string, error)
	UpdateCourseSettings(ctx context.Context, course *domain.Course) error
	CreateModule(ctx context.Context, module *domain.Module) (string, error)
	DeleteModule(ctx context.Context, id string) error
	CreateLesson(ctx context.Context, lesson *domain.Lesson) (string, error)
	UpdateLesson(ctx context.Context, lesson *domain.Lesson) error
	DeleteLesson(ctx context.Context, id string) error
	AssignTeacherToLesson(ctx context.Context, lessonID, teacherID string) error
	GetLessonIDByOrder(ctx context.Context, courseID string, orderNum int) (string, error)
	GetAllCourses(ctx context.Context) ([]*domain.Course, error)
	GetCourseByID(ctx context.Context, id string) (*domain.Course, error)
	GetModulesByCourseID(ctx context.Context, courseID string) ([]*domain.Module, error)
	GetLessonsByCourseID(ctx context.Context, courseID string) ([]*domain.Lesson, error)
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetTestByID(ctx context.Context, id string) (*domain.Test, error)
	GetProjectByID(ctx context.Context, id string) (*domain.Project, error)
	GetParentsByStudentID(ctx context.Context, studentID string) ([]domain.User, error)
	LinkParentToStudent(ctx context.Context, studentID, parentID string) error
	EnrollStudentExtended(ctx context.Context, userID, courseID, streamID, groupID string) error
	GetCourseIDByStream(ctx context.Context, streamID string) (string, error)
	GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error)
	GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error)
	CreateTest(ctx context.Context, test *domain.Test) (string, error)
	DeleteTest(ctx context.Context, id string) error
	CreateProject(ctx context.Context, project *domain.Project) (string, error)
	DeleteProject(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID string) error
	GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error)
	GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error)
	GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error)
	GetDetailedModeratorList(ctx context.Context) ([]*domain.ModeratorTableItem, error)
	GetAllUsersList(ctx context.Context) ([]*domain.AllUsersTableItem, error)
	CreateStream(ctx context.Context, stream *domain.Stream) (string, error)
	GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error)
	CreateGroup(ctx context.Context, group *domain.Group) (string, error)
	GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error)
	GetStudentEnrollment(ctx context.Context, userID string) (map[string]string, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	UnlinkAllParents(ctx context.Context, studentID string) error
	SetLessonModule(ctx context.Context, lessonID, moduleID string) error
	GetTestsByCourseID(ctx context.Context, courseID string) ([]domain.Test, error)
	GetProjectsByCourseID(ctx context.Context, courseID string) ([]domain.Project, error)
	UnenrollStudent(ctx context.Context, userID, courseID string) error
	LinkTeachersToCourse(ctx context.Context, courseID string, teacherIDs []string) error
	GetLessonByID(ctx context.Context, id string) (*domain.Lesson, error)
	CancelLesson(ctx context.Context, lessonID, reason string) error
	SubstituteTeacher(ctx context.Context, lessonID, teacherID string) error
	EnsureAssignment(ctx context.Context, lessonID, title string) error
}

type ContentAdminRepoImpl struct {
	db *sql.DB
}

var _ ContentAdminRepository = (*ContentAdminRepoImpl)(nil)

func NewContentAdminRepository(db *sql.DB) *ContentAdminRepoImpl {
	return &ContentAdminRepoImpl{db: db}
}
