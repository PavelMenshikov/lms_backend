package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"golang.org/x/crypto/bcrypt"

	"lms_backend/internal/content_admin/repository"
	"lms_backend/internal/domain"
	storageService "lms_backend/pkg/storage"
)

type ContentAdminUseCase struct {
	repo      repository.ContentAdminRepository
	s3Storage storageService.ObjectStorage
}

func NewContentAdminUseCase(repo repository.ContentAdminRepository, s3Storage storageService.ObjectStorage) *ContentAdminUseCase {
	return &ContentAdminUseCase{repo: repo, s3Storage: s3Storage}
}

type CreateCourseInput struct {
	Title       string
	Description string
	IsMain      bool
	FileHeader  *multipart.FileHeader
}

func (uc *ContentAdminUseCase) CreateCourse(ctx context.Context, input CreateCourseInput) (string, error) {
	if input.Title == "" {
		return "", errors.New("title is required")
	}

	var imageURL string
	if input.FileHeader != nil {
		file, err := input.FileHeader.Open()
		if err != nil {
			return "", fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		s3Key := fmt.Sprintf("course_previews/%s", input.FileHeader.Filename)
		mimeType := input.FileHeader.Header.Get("Content-Type")

		s3KeyAfterUpload, err := uc.s3Storage.UploadFile(ctx, file, s3Key, input.FileHeader.Size, mimeType)
		if err != nil {
			return "", fmt.Errorf("failed to upload image to S3: %w", err)
		}

		imageURL, err = uc.s3Storage.GetPublicURL(ctx, s3KeyAfterUpload)
		if err != nil {
			return "", fmt.Errorf("failed to get public URL for %s: %w", s3KeyAfterUpload, err)
		}
	}

	course := &domain.Course{
		Title:       input.Title,
		Description: input.Description,
		IsMain:      input.IsMain,
		ImageURL:    imageURL,
		Status:      domain.CourseStatusDraft,
	}

	return uc.repo.CreateCourse(ctx, course)
}

type CreateModuleInput struct {
	CourseID    string
	Title       string
	Description string
	OrderNum    int
}

func (uc *ContentAdminUseCase) CreateModule(ctx context.Context, input CreateModuleInput) (string, error) {
	module := &domain.Module{
		CourseID:    input.CourseID,
		Title:       input.Title,
		Description: input.Description,
		OrderNum:    input.OrderNum,
	}
	return uc.repo.CreateModule(ctx, module)
}

type UpdateCourseSettingsInput struct {
	CourseID            string
	Title               string
	Description         string
	IsMain              bool
	Status              domain.CourseStatus
	HasHomework         bool
	IsHomeworkMandatory bool
	IsTestMandatory     bool
	IsProjectMandatory  bool
	IsDiscordMandatory  bool
	IsAntiCopyEnabled   bool
	FileHeader          *multipart.FileHeader
}

func (uc *ContentAdminUseCase) UpdateCourseSettings(ctx context.Context, input UpdateCourseSettingsInput) error {
	var imageURL string

	if input.FileHeader != nil {
		file, err := input.FileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		s3Key := fmt.Sprintf("course_previews/%s_%s", input.CourseID, input.FileHeader.Filename)
		mimeType := input.FileHeader.Header.Get("Content-Type")

		key, err := uc.s3Storage.UploadFile(ctx, file, s3Key, input.FileHeader.Size, mimeType)
		if err != nil {
			return err
		}
		imageURL, _ = uc.s3Storage.GetPublicURL(ctx, key)
	}

	course := &domain.Course{
		ID:                  input.CourseID,
		Title:               input.Title,
		Description:         input.Description,
		IsMain:              input.IsMain,
		Status:              input.Status,
		ImageURL:            imageURL,
		HasHomework:         input.HasHomework,
		IsHomeworkMandatory: input.IsHomeworkMandatory,
		IsTestMandatory:     input.IsTestMandatory,
		IsProjectMandatory:  input.IsProjectMandatory,
		IsDiscordMandatory:  input.IsDiscordMandatory,
		IsAntiCopyEnabled:   input.IsAntiCopyEnabled,
	}

	return uc.repo.UpdateCourseSettings(ctx, course)
}

type CreateLessonInput struct {
	ModuleID         string
	TeacherID        string
	Title            string
	OrderNum         int
	VideoFile        *multipart.FileHeader
	PresentationFile *multipart.FileHeader
	ContentText      string
}

func (uc *ContentAdminUseCase) CreateLesson(ctx context.Context, input CreateLessonInput) (string, error) {
	var videoURL, presentationURL string

	if input.VideoFile != nil {
		file, err := input.VideoFile.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		key, err := uc.s3Storage.UploadFile(ctx, file, "lessons/videos/"+input.VideoFile.Filename, input.VideoFile.Size, input.VideoFile.Header.Get("Content-Type"))
		if err != nil {
			return "", err
		}
		videoURL, _ = uc.s3Storage.GetPublicURL(ctx, key)
	}

	if input.PresentationFile != nil {
		file, err := input.PresentationFile.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		key, err := uc.s3Storage.UploadFile(ctx, file, "lessons/presentations/"+input.PresentationFile.Filename, input.PresentationFile.Size, input.PresentationFile.Header.Get("Content-Type"))
		if err != nil {
			return "", err
		}
		presentationURL, _ = uc.s3Storage.GetPublicURL(ctx, key)
	}

	lesson := &domain.Lesson{
		ModuleID:        input.ModuleID,
		TeacherID:       input.TeacherID,
		Title:           input.Title,
		LessonTime:      time.Now(),
		OrderNum:        input.OrderNum,
		VideoURL:        videoURL,
		PresentationURL: presentationURL,
		ContentText:     input.ContentText,
		IsPublished:     true,
	}

	return uc.repo.CreateLesson(ctx, lesson)
}

func (uc *ContentAdminUseCase) GetAllCourses(ctx context.Context) ([]*domain.Course, error) {
	return uc.repo.GetAllCourses(ctx)
}

func (uc *ContentAdminUseCase) GetCourseByID(ctx context.Context, id string) (*domain.Course, error) {
	return uc.repo.GetCourseByID(ctx, id)
}

func (uc *ContentAdminUseCase) GetCourseStructure(ctx context.Context, courseID string) (*domain.CourseStructure, error) {
	course, err := uc.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	modules, err := uc.repo.GetModulesByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	allLessons, err := uc.repo.GetLessonsByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	lessonsByModule := make(map[string][]*domain.Lesson)
	for _, l := range allLessons {
		lessonsByModule[l.ModuleID] = append(lessonsByModule[l.ModuleID], l)
	}

	var moduleStructures []*domain.ModuleStructure
	for _, m := range modules {
		ms := &domain.ModuleStructure{
			Module:  m,
			Lessons: lessonsByModule[m.ID],
		}
		if ms.Lessons == nil {
			ms.Lessons = []*domain.Lesson{}
		}
		moduleStructures = append(moduleStructures, ms)
	}

	return &domain.CourseStructure{
		Course:  course,
		Modules: moduleStructures,
	}, nil
}

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      domain.Role
}

func (uc *ContentAdminUseCase) CreateUser(ctx context.Context, input CreateUserInput) (string, error) {
	if input.Email == "" || input.Password == "" {
		return "", errors.New("email and password are required")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedBytes),
		Role:      input.Role,
	}

	return uc.repo.CreateUser(ctx, user)
}

func (uc *ContentAdminUseCase) EnrollStudent(ctx context.Context, userID, courseID string) error {
	if userID == "" || courseID == "" {
		return errors.New("user_id and course_id are required")
	}
	return uc.repo.EnrollStudent(ctx, userID, courseID)
}

func (uc *ContentAdminUseCase) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	return uc.repo.GetCourseStudents(ctx, courseID)
}

func (uc *ContentAdminUseCase) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	return uc.repo.GetCourseStats(ctx, courseID)
}