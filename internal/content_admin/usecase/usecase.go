package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
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

func splitName(fullName string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(fullName), " ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return fullName, ""
}

type CreateCourseInput struct {
	Title       string
	Description string
	IsMain      bool
	FileHeader  *multipart.FileHeader
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

type CreateModuleInput struct {
	CourseID    string
	Title       string
	Description string
	OrderNum    int
}

type CreateLessonInput struct {
	CourseID         string
	ModuleID         string
	TeacherID        string
	Title            string
	OrderNum         int
	VideoFile        *multipart.FileHeader
	PresentationFile *multipart.FileHeader
	ContentText      string
}

type ExtendedCreateUserInput struct {
	FullName        string
	Email           string
	Role            domain.Role
	Password        string
	Phone           string
	City            string
	SchoolName      string
	Language        string
	Gender          string
	BirthDate       time.Time
	ExperienceYears int
	Whatsapp        string
	Telegram        string
	CourseID        string
	StreamID        string
	GroupID         string
	Parents         []ParentInfo
}

type ParentInfo struct {
	FullName string
	Phone    string
	Email    string
}

type CreateTestInput struct {
	CourseID     string
	LessonNumber int
	LessonID     string
	Title        string
	Description  string
	PassingScore int
}

type CreateProjectInput struct {
	CourseID     string
	LessonNumber int
	LessonID     string
	Title        string
	Description  string
	MaxScore     int
}

type CreateStreamInput struct {
	CourseID  string
	Title     string
	StartDate time.Time
}

type CreateGroupInput struct {
	StreamID  string
	CuratorID string
	TeacherID string
	Title     string
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
	var rootLessons []*domain.Lesson

	for _, l := range allLessons {
		if l.ModuleID != nil && *l.ModuleID != "" {
			lessonsByModule[*l.ModuleID] = append(lessonsByModule[*l.ModuleID], l)
		} else {
			rootLessons = append(rootLessons, l)
		}
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

	if rootLessons == nil {
		rootLessons = []*domain.Lesson{}
	}

	return &domain.CourseStructure{
		Course:      course,
		Modules:     moduleStructures,
		RootLessons: rootLessons,
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

func (uc *ContentAdminUseCase) DeleteModule(ctx context.Context, id string) error {
	return uc.repo.DeleteModule(ctx, id)
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

	var modID *string
	if input.ModuleID != "" {
		modID = &input.ModuleID
	}

	lesson := &domain.Lesson{
		CourseID:        input.CourseID,
		ModuleID:        modID,
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

func (uc *ContentAdminUseCase) DeleteLesson(ctx context.Context, id string) error {
	return uc.repo.DeleteLesson(ctx, id)
}

func (uc *ContentAdminUseCase) CreateFullUser(ctx context.Context, input ExtendedCreateUserInput) (map[string]string, error) {
	firstName, lastName := splitName(input.FullName)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           input.Email,
		Password:        string(hashedPass),
		Role:            input.Role,
		Phone:           input.Phone,
		City:            input.City,
		SchoolName:      input.SchoolName,
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

	if input.Role == domain.RoleStudent {
		for _, pInfo := range input.Parents {
			pFirst, pLast := splitName(pInfo.FullName)
			parentPassRaw := generateSecurePassword()
			parentHash, _ := bcrypt.GenerateFromPassword([]byte(parentPassRaw), 12)

			pEmail := pInfo.Email
			if pEmail == "" {
				pEmail = fmt.Sprintf("p_%s_%s", userID[:8], pInfo.Phone)
			}

			parent := &domain.User{
				FirstName: pFirst,
				LastName:  pLast,
				Email:     pEmail,
				Phone:     pInfo.Phone,
				Password:  string(parentHash),
				Role:      domain.RoleParent,
				City:      input.City,
			}

			parentID, err := uc.repo.CreateUser(ctx, parent)
			if err == nil {
				_ = uc.repo.LinkParentToStudent(ctx, userID, parentID)
			}
		}

		courseID := input.CourseID
		if courseID == "" && input.StreamID != "" {
			derivedCourseID, err := uc.repo.GetCourseIDByStream(ctx, input.StreamID)
			if err == nil {
				courseID = derivedCourseID
			}
		}

		if courseID != "" || input.StreamID != "" || input.GroupID != "" {
			_ = uc.repo.EnrollStudentExtended(ctx, userID, courseID, input.StreamID, input.GroupID)
		}
	}

	return map[string]string{"user_id": userID}, nil
}

func (uc *ContentAdminUseCase) EnrollStudent(ctx context.Context, userID, courseID string) error {
	return uc.repo.EnrollStudentExtended(ctx, userID, courseID, "", "")
}

func (uc *ContentAdminUseCase) GetCourseStudents(ctx context.Context, courseID string) ([]*domain.AdminStudentProgress, error) {
	return uc.repo.GetCourseStudents(ctx, courseID)
}

func (uc *ContentAdminUseCase) GetCourseStats(ctx context.Context, courseID string) (*domain.AdminCourseStats, error) {
	return uc.repo.GetCourseStats(ctx, courseID)
}

func (uc *ContentAdminUseCase) GetUsersList(ctx context.Context, role domain.Role) ([]*domain.User, error) {
	filter := domain.UserFilter{
		Role: role,
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

func (uc *ContentAdminUseCase) GetDetailedModerators(ctx context.Context) ([]*domain.ModeratorTableItem, error) {
	return uc.repo.GetDetailedModeratorList(ctx)
}


func (uc *ContentAdminUseCase) GetAllUsersTable(ctx context.Context) ([]*domain.AllUsersTableItem, error) {
	return uc.repo.GetAllUsersList(ctx)
}

func (uc *ContentAdminUseCase) GetUserInfo(ctx context.Context, userID string) (map[string]interface{}, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usecase: user not found: %w", err)
	}

	res := map[string]interface{}{
		"user": user,
	}

	if user.Role == domain.RoleStudent {
		parents, err := uc.repo.GetParentsByStudentID(ctx, userID)
		if err == nil {
			res["parents"] = parents
		} else {
			res["parents"] = []domain.User{}
		}
	}

	return res, nil
}


func (uc *ContentAdminUseCase) UpdateUser(ctx context.Context, userID string, input ExtendedCreateUserInput) error {
	firstName, lastName := splitName(input.FullName)
	user := &domain.User{
		ID:              userID,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           input.Email,
		Role:            input.Role,
		Phone:           input.Phone,
		City:            input.City,
		SchoolName:      input.SchoolName,
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
	lessonID := input.LessonID
	if lessonID == "" && input.CourseID != "" && input.LessonNumber > 0 {
		id, err := uc.repo.GetLessonIDByOrder(ctx, input.CourseID, input.LessonNumber)
		if err != nil {
			return "", fmt.Errorf("lesson not found by order: %w", err)
		}
		lessonID = id
	}

	if lessonID == "" {
		return "", errors.New("lesson_id or (course_id + lesson_number) required")
	}

	test := &domain.Test{
		LessonID:     lessonID,
		Title:        input.Title,
		Description:  input.Description,
		PassingScore: input.PassingScore,
	}
	return uc.repo.CreateTest(ctx, test)
}

func (uc *ContentAdminUseCase) DeleteTest(ctx context.Context, id string) error {
	return uc.repo.DeleteTest(ctx, id)
}

func (uc *ContentAdminUseCase) CreateProject(ctx context.Context, input CreateProjectInput) (string, error) {
	lessonID := input.LessonID
	if lessonID == "" && input.CourseID != "" && input.LessonNumber > 0 {
		id, err := uc.repo.GetLessonIDByOrder(ctx, input.CourseID, input.LessonNumber)
		if err != nil {
			return "", fmt.Errorf("lesson not found by number: %w", err)
		}
		lessonID = id
	}

	if lessonID == "" {
		return "", errors.New("either lesson_id or course_id + lesson_number must be provided")
	}

	project := &domain.Project{
		LessonID:    lessonID,
		Title:       input.Title,
		Description: input.Description,
		MaxScore:    input.MaxScore,
	}
	return uc.repo.CreateProject(ctx, project)
}

func (uc *ContentAdminUseCase) DeleteProject(ctx context.Context, id string) error {
	return uc.repo.DeleteProject(ctx, id)
}

func (uc *ContentAdminUseCase) CreateStream(ctx context.Context, input CreateStreamInput) (string, error) {
	stream := &domain.Stream{
		CourseID:  input.CourseID,
		Title:     input.Title,
		StartDate: input.StartDate,
	}
	return uc.repo.CreateStream(ctx, stream)
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

func (uc *ContentAdminUseCase) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	return uc.repo.GetGroupsByStream(ctx, streamID)
}

func generateSecurePassword() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}