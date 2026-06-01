package mocks

import (
	"context"
	"lms_backend/internal/domain"
	"lms_backend/internal/learning/repository"
)

var _ repository.LearningRepository = (*LearningRepoMock)(nil)

type LearningRepoMock struct {
	GetMyCoursesFunc            func(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error)
	GetCourseContentFunc        func(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error)
	GetLessonDetailFunc         func(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error)
	GetAssignmentIDByLessonFunc func(ctx context.Context, lessonID string) (string, error)
	SaveSubmissionFunc          func(ctx context.Context, userID, assignmentID, text string, files []string) error
	SetLessonAttendanceFunc     func(ctx context.Context, userID, lessonID, status, recordingURL, teacherComment string) error
	GetTeachersListFunc         func(ctx context.Context) ([]*domain.TeacherPublicInfo, error)
	GetTeacherByIDFunc          func(ctx context.Context, id string) (*domain.TeacherPublicInfo, error)
	AddTeacherReviewFunc        func(ctx context.Context, review *domain.TeacherReview) error
	GetTeacherReviewsFunc       func(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error)
	GetTeacherCoursesFunc       func(ctx context.Context, teacherID string) ([]*domain.StudentCoursePreview, error)
	GetTestByIDFunc             func(ctx context.Context, testID string) (*domain.Test, error)
	GetProjectByIDFunc          func(ctx context.Context, projectID string) (*domain.Project, error)
	GetTeacherSubstitutionsFunc    func(ctx context.Context, teacherID string) ([]*domain.Lesson, error)
	GetTeacherUpcomingLessonsFunc  func(ctx context.Context, teacherID string) ([]*domain.Lesson, error)
	GetTeacherCancelledLessonsFunc func(ctx context.Context, teacherID string) ([]*domain.Lesson, error)
}

func NewLearningRepoMock() *LearningRepoMock {
	return &LearningRepoMock{}
}

func (m *LearningRepoMock) GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
	return m.GetMyCoursesFunc(ctx, userID)
}

func (m *LearningRepoMock) GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error) {
	return m.GetCourseContentFunc(ctx, courseID, userID)
}

func (m *LearningRepoMock) GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error) {
	return m.GetLessonDetailFunc(ctx, lessonID, userID)
}

func (m *LearningRepoMock) GetAssignmentIDByLesson(ctx context.Context, lessonID string) (string, error) {
	return m.GetAssignmentIDByLessonFunc(ctx, lessonID)
}

func (m *LearningRepoMock) SaveSubmission(ctx context.Context, userID, assignmentID, text string, files []string) error {
	return m.SaveSubmissionFunc(ctx, userID, assignmentID, text, files)
}

func (m *LearningRepoMock) SetLessonAttendance(ctx context.Context, userID, lessonID, status, recordingURL, teacherComment string) error {
	return m.SetLessonAttendanceFunc(ctx, userID, lessonID, status, recordingURL, teacherComment)
}

func (m *LearningRepoMock) GetTeachersList(ctx context.Context) ([]*domain.TeacherPublicInfo, error) {
	return m.GetTeachersListFunc(ctx)
}

func (m *LearningRepoMock) GetTeacherByID(ctx context.Context, id string) (*domain.TeacherPublicInfo, error) {
	return m.GetTeacherByIDFunc(ctx, id)
}

func (m *LearningRepoMock) AddTeacherReview(ctx context.Context, review *domain.TeacherReview) error {
	return m.AddTeacherReviewFunc(ctx, review)
}

func (m *LearningRepoMock) GetTeacherReviews(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error) {
	return m.GetTeacherReviewsFunc(ctx, teacherID)
}

func (m *LearningRepoMock) GetTeacherCourses(ctx context.Context, teacherID string) ([]*domain.StudentCoursePreview, error) {
	return m.GetTeacherCoursesFunc(ctx, teacherID)
}

func (m *LearningRepoMock) GetTestByID(ctx context.Context, testID string) (*domain.Test, error) {
	return m.GetTestByIDFunc(ctx, testID)
}

func (m *LearningRepoMock) GetProjectByID(ctx context.Context, projectID string) (*domain.Project, error) {
	return m.GetProjectByIDFunc(ctx, projectID)
}

func (m *LearningRepoMock) GetTeacherSubstitutions(ctx context.Context, teacherID string) ([]*domain.Lesson, error) {
	return m.GetTeacherSubstitutionsFunc(ctx, teacherID)
}

func (m *LearningRepoMock) GetTeacherUpcomingLessons(ctx context.Context, teacherID string) ([]*domain.Lesson, error) {
	return m.GetTeacherUpcomingLessonsFunc(ctx, teacherID)
}

func (m *LearningRepoMock) GetTeacherCancelledLessons(ctx context.Context, teacherID string) ([]*domain.Lesson, error) {
	return m.GetTeacherCancelledLessonsFunc(ctx, teacherID)
}
