package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

func (uc *ContentAdminUseCase) UploadMedia(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", errors.New("no file provided")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	s3Key := fmt.Sprintf("editor_content/%d_%s", time.Now().Unix(), fileHeader.Filename)
	mimeType := fileHeader.Header.Get("Content-Type")

	key, err := uc.s3Storage.UploadFile(ctx, file, s3Key, fileHeader.Size, mimeType)
	if err != nil {
		return "", fmt.Errorf("failed to upload to s3: %w", err)
	}

	return uc.s3Storage.GetPublicURL(ctx, key)
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

type CreateCourseInput struct {
	Title       string
	Description string
	IsMain      bool
	FileHeader  *multipart.FileHeader
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

func (uc *ContentAdminUseCase) CreateModule(ctx context.Context, input CreateModuleInput) (string, error) {
	module := &domain.Module{
		CourseID:    input.CourseID,
		Title:       input.Title,
		Description: input.Description,
		OrderNum:    input.OrderNum,
	}
	return uc.repo.CreateModule(ctx, module)
}

type CreateModuleInput struct {
	CourseID    string
	Title       string
	Description string
	OrderNum    int
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

type CreateLessonInput struct {
	ModuleID         string
	TeacherID        string
	Title            string
	OrderNum         int
	VideoFile        *multipart.FileHeader
	PresentationFile *multipart.FileHeader
	ContentText      string
}

func (uc *ContentAdminUseCase) CreateFullUser(ctx context.Context, input ExtendedCreateUserInput) (map[string]string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Email:           input.Email,
		Password:        string(hashedPass),
		Role:            input.Role,
		Phone:           input.Phone,
		City:            input.City,
		Language:        input.Language,
		Gender:          input.Gender,
		BirthDate:       input.BirthDate,
		ExperienceYears: input.ExperienceYears,
		Whatsapp:        input.Whatsapp,
		Telegram:        input.Telegram,
	}

	userID, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	result := map[string]string{
		"user_id": userID,
		"role":    string(input.Role),
	}

	if input.Role == domain.RoleStudent && input.ParentPhone != "" {
		parentPassRaw := generateSecurePassword()
		parentHash, _ := bcrypt.GenerateFromPassword([]byte(parentPassRaw), 12)

		parentEmail := input.ParentEmail
		if parentEmail == "" {
			parentEmail = fmt.Sprintf("p_%s", input.Email)
		}

		parent := &domain.User{
			FirstName: input.ParentFirstName,
			LastName:  input.ParentLastName,
			Email:     parentEmail,
			Phone:     input.ParentPhone,
			Password:  string(parentHash),
			Role:      domain.RoleParent,
			City:      input.City,
		}

		parentID, err := uc.repo.CreateUser(ctx, parent)
		if err == nil {
			_ = uc.repo.LinkParentToStudent(ctx, userID, parentID)
			result["parent_id"] = parentID
			result["parent_email"] = parentEmail
			result["parent_password"] = parentPassRaw
		}
	}

	return result, nil
}

type ExtendedCreateUserInput struct {
	FirstName       string
	LastName        string
	Email           string
	Role            domain.Role
	Password        string
	Phone           string
	City            string
	Language        string
	Gender          string
	BirthDate       time.Time
	ExperienceYears int
	Whatsapp        string
	Telegram        string
	ParentFirstName string
	ParentLastName  string
	ParentPhone     string
	ParentEmail     string
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

func (uc *ContentAdminUseCase) GetUsersList(ctx context.Context, role domain.Role) ([]*domain.User, error) {
	filter := domain.UserFilter{
		Role:   role,
		Limit:  100,
		Offset: 0,
	}
	return uc.repo.GetUsers(ctx, filter)
}

func (uc *ContentAdminUseCase) GetDetailedStudents(ctx context.Context, courseID string) ([]*domain.StudentTableItem, error) {
	filter := domain.UserFilter{CourseID: courseID}
	return uc.repo.GetDetailedStudentList(ctx, filter)
}

func (uc *ContentAdminUseCase) GetDetailedTeachers(ctx context.Context) ([]*domain.TeacherTableItem, error) {
	return uc.repo.GetDetailedTeacherList(ctx)
}

func (uc *ContentAdminUseCase) GetDetailedCurators(ctx context.Context) ([]*domain.CuratorTableItem, error) {
	return uc.repo.GetDetailedCuratorList(ctx)
}

func (uc *ContentAdminUseCase) UpdateUser(ctx context.Context, userID string, input ExtendedCreateUserInput) error {
	user := &domain.User{
		ID:              userID,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Email:           input.Email,
		Role:            input.Role,
		Phone:           input.Phone,
		City:            input.City,
		Language:        input.Language,
		Gender:          input.Gender,
		ExperienceYears: input.ExperienceYears,
		Whatsapp:        input.Whatsapp,
		Telegram:        input.Telegram,
	}
	return uc.repo.UpdateUser(ctx, user)
}

func (uc *ContentAdminUseCase) DeleteUser(ctx context.Context, userID string) error {
	return uc.repo.DeleteUser(ctx, userID)
}

func (uc *ContentAdminUseCase) CreateTest(ctx context.Context, input CreateTestInput) (string, error) {
	test := &domain.Test{
		LessonID:     input.LessonID,
		Title:        input.Title,
		Description:  input.Description,
		PassingScore: input.PassingScore,
	}
	return uc.repo.CreateTest(ctx, test)
}

type CreateTestInput struct {
	LessonID     string
	Title        string
	Description  string
	PassingScore int
}

func (uc *ContentAdminUseCase) CreateProject(ctx context.Context, input CreateProjectInput) (string, error) {
	project := &domain.Project{
		LessonID:    input.LessonID,
		Title:       input.Title,
		Description: input.Description,
		MaxScore:    input.MaxScore,
	}
	return uc.repo.CreateProject(ctx, project)
}

type CreateProjectInput struct {
	LessonID    string
	Title       string
	Description string
	MaxScore    int
}

func (uc *ContentAdminUseCase) CreateStream(ctx context.Context, input CreateStreamInput) (string, error) {
	stream := &domain.Stream{
		CourseID:  input.CourseID,
		Title:     input.Title,
		StartDate: input.StartDate,
	}
	return uc.repo.CreateStream(ctx, stream)
}

type CreateStreamInput struct {
	CourseID  string
	Title     string
	StartDate time.Time
}

func (uc *ContentAdminUseCase) GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error) {
	return uc.repo.GetStreamsByCourse(ctx, courseID)
}

func (uc *ContentAdminUseCase) CreateGroup(ctx context.Context, input CreateGroupInput) (string, error) {
	group := &domain.Group{
		StreamID:  input.StreamID,
		CuratorID: input.CuratorID,
		TeacherID: input.TeacherID,
		Title:     input.Title,
	}
	return uc.repo.CreateGroup(ctx, group)
}

type CreateGroupInput struct {
	StreamID  string
	CuratorID string
	TeacherID string
	Title     string
}

func (uc *ContentAdminUseCase) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	return uc.repo.GetGroupsByStream(ctx, streamID)
}

func generateSecurePassword() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
