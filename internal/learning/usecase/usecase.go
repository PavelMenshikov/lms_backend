package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"time"

	"lms_backend/internal/domain"
	"lms_backend/internal/learning/repository"
	storageService "lms_backend/pkg/storage"
)

const s3UploadTimeout = 30 * time.Second

func s3Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, s3UploadTimeout)
}

type LearningUseCase struct {
	repo      repository.LearningRepository
	s3Storage storageService.ObjectStorage
}

func NewLearningUseCase(repo repository.LearningRepository, s3Storage storageService.ObjectStorage) *LearningUseCase {
	return &LearningUseCase{repo: repo, s3Storage: s3Storage}
}

func (uc *LearningUseCase) GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	return uc.repo.GetMyCourses(ctx, userID)
}

func (uc *LearningUseCase) GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error) {
	return uc.repo.GetCourseContent(ctx, courseID, userID)
}

func (uc *LearningUseCase) GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error) {
	return uc.repo.GetLessonDetail(ctx, lessonID, userID)
}

type SubmitAssignmentInput struct {
	LessonID    string
	UserID      string
	TextAnswer  string
	FileHeaders []*multipart.FileHeader
}

func (uc *LearningUseCase) SubmitAssignment(ctx context.Context, input SubmitAssignmentInput) error {
	assignmentID, err := uc.repo.GetAssignmentIDByLesson(ctx, input.LessonID)
	if err != nil {
		if err := uc.repo.EnsureAssignment(ctx, input.LessonID, ""); err != nil {
			return fmt.Errorf("failed to ensure assignment: %w", err)
		}
		assignmentID, err = uc.repo.GetAssignmentIDByLesson(ctx, input.LessonID)
		if err != nil {
			return fmt.Errorf("assignment not found after ensure: %w", err)
		}
	}
	var fileURLs []string
	for _, fh := range input.FileHeaders {
		file, err := fh.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		s3Ctx, cancel := s3Context(ctx)
		defer cancel()
		s3Key := fmt.Sprintf("submissions/%s_%s_%s", input.UserID, assignmentID, fh.Filename)
		mimeType := fh.Header.Get("Content-Type")
		if mimeType == "" {
			ext := filepath.Ext(fh.Filename)
			mimeType = mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}
		key, err := uc.s3Storage.UploadFile(s3Ctx, file, s3Key, fh.Size, mimeType)
		if err != nil {
			return err
		}
		url, _ := uc.s3Storage.GetPublicURL(ctx, key)
		fileURLs = append(fileURLs, url)
	}
	return uc.repo.SaveSubmission(ctx, input.UserID, assignmentID, input.TextAnswer, fileURLs)
}

type SetAttendanceInput struct {
	LessonID       string
	UserID         string
	Status         string
	RecordingURL   string
	TeacherComment string
}

func (uc *LearningUseCase) SetLessonAttendance(ctx context.Context, input SetAttendanceInput) error {
	return uc.repo.SetLessonAttendance(ctx, input.UserID, input.LessonID, input.Status, input.RecordingURL, input.TeacherComment)
}

func (uc *LearningUseCase) GetTeachers(ctx context.Context) ([]*domain.TeacherPublicInfo, error) {
	return uc.repo.GetTeachersList(ctx)
}

func (uc *LearningUseCase) GetTeacherDetails(ctx context.Context, id string) (*domain.TeacherPublicInfo, error) {
	teacher, err := uc.repo.GetTeacherByID(ctx, id)
	if err != nil {
		return nil, err
	}
	reviews, _ := uc.repo.GetTeacherReviews(ctx, id)
	teacher.Reviews = reviews
	return teacher, nil
}

type AddReviewInput struct {
	TeacherID string
	StudentID string
	Rating    int
	Comment   string
}

func (uc *LearningUseCase) AddReview(ctx context.Context, input AddReviewInput) error {
	if input.Rating < 1 || input.Rating > 5 {
		return errors.New("invalid rating")
	}
	review := &domain.TeacherReview{
		TeacherID: input.TeacherID,
		StudentID: input.StudentID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}
	return uc.repo.AddTeacherReview(ctx, review)
}
func (uc *LearningUseCase) GetProject(ctx context.Context, projectID string) (*domain.Project, error) {
	p, err := uc.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, nil
	}
	return p, nil
}

func (uc *LearningUseCase) GetTest(ctx context.Context, testID string) (*domain.Test, error) {
	t, err := uc.repo.GetTestByID(ctx, testID)
	if err != nil {
		return nil, nil
	}
	return t, nil
}

func (uc *LearningUseCase) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	return uc.repo.GetAllCourses(ctx)
}
func (uc *LearningUseCase) GetTeacherCertificates(ctx context.Context, teacherID string) ([]*domain.TeacherCertificate, error) {
	return uc.repo.GetTeacherCertificates(ctx, teacherID)
}

func (uc *LearningUseCase) GetTeacherDashboard(ctx context.Context, teacherID string) (*domain.TeacherDashboardData, error) {
	profile, err := uc.repo.GetTeacherByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	reviews, _ := uc.repo.GetTeacherReviews(ctx, teacherID)
	if reviews == nil {
		reviews = []*domain.TeacherReview{}
	}

	courses, _ := uc.repo.GetTeacherCourses(ctx, teacherID)

	substitutions, _ := uc.repo.GetTeacherSubstitutions(ctx, teacherID)
	cancelledLessons, _ := uc.repo.GetTeacherCancelledLessons(ctx, teacherID)
	upcomingLessons, _ := uc.repo.GetTeacherUpcomingLessons(ctx, teacherID)

	if substitutions == nil {
		substitutions = []*domain.Lesson{}
	}
	if cancelledLessons == nil {
		cancelledLessons = []*domain.Lesson{}
	}
	if upcomingLessons == nil {
		upcomingLessons = []*domain.Lesson{}
	}

	return &domain.TeacherDashboardData{
		Profile:          profile,
		AssignedCourses:  courses,
		MyReviews:        reviews,
		Substitutions:    substitutions,
		CancelledLessons: cancelledLessons,
		UpcomingLessons:  upcomingLessons,
	}, nil
}
