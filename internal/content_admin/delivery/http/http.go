package http

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"lms_backend/internal/content_admin/usecase"
	"lms_backend/internal/domain"
)

type ContentAdminHandler struct {
	uc *usecase.ContentAdminUseCase
}

func NewContentAdminHandler(uc *usecase.ContentAdminUseCase) *ContentAdminHandler {
	return &ContentAdminHandler{uc: uc}
}

type CreateFullUserRequest struct {
	FullName        string               `json:"full_name" example:"Иван Иванов"`
	Email           string               `json:"email" example:"student@test.kz"`
	Password        string               `json:"password" example:"secret123"`
	Role            domain.Role          `json:"role" example:"student"`
	Phone           string               `json:"phone" example:"+79998887766"`
	City            string               `json:"city" example:"Алматы"`
	SchoolName      string               `json:"school_name" example:"Школа №123"`
	Language        string               `json:"language" example:"ru"`
	Gender          string               `json:"gender" example:"male"`
	BirthDateStr    string               `json:"birth_date" example:"2000-01-01"`
	Whatsapp        string               `json:"whatsapp" example:"https://wa.me/..."`
	Telegram        string               `json:"telegram" example:"https://t.me/..."`
	ExperienceYears int                  `json:"experience_years" example:"5"`
	CourseID        string               `json:"course_id"`
	StreamID        string               `json:"stream_id"`
	GroupID         string               `json:"group_id"`
	Parents         []usecase.ParentInfo `json:"parents"`
}

type EnrollRequest struct {
	UserID   string `json:"user_id" example:"a0000000-0000-0000-0000-000000000001"`
	CourseID string `json:"course_id" example:"c1111111-1111-1111-1111-111111111111"`
}

type CreateTestRequest struct {
	CourseID     string `json:"course_id"`
	LessonNumber int    `json:"lesson_number"`
	LessonID     string `json:"lesson_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PassingScore int    `json:"passing_score"`
}

type CreateProjectRequest struct {
	CourseID     string `json:"course_id" example:"uuid"`
	LessonNumber int    `json:"lesson_number" example:"1"`
	LessonID     string `json:"lesson_id" example:"l2222222-2222-2222-2222-222222222222"`
	Title        string `json:"title" example:"Финальный проект"`
	Description  string `json:"description" example:"Разработка API на Go"`
	MaxScore     int    `json:"max_score" example:"100"`
}

type CreateModuleRequest struct {
	CourseID    string `json:"course_id" example:"c1111111-1111-1111-1111-111111111111"`
	Title       string `json:"title" example:"Основы синтаксиса"`
	Description string `json:"description" example:"Типы данных, переменные, циклы"`
	OrderNum    int    `json:"order_num" example:"1"`
}

type CreateStreamRequest struct {
	CourseID     string `json:"course_id" example:"c1111111-1111-1111-1111-111111111111"`
	Title        string `json:"title" example:"Поток Сентябрь 2024"`
	StartDateStr string `json:"start_date" example:"2024-09-01"`
}

type CreateGroupRequest struct {
	StreamID  string `json:"stream_id" example:"s3333333-3333-3333-3333-333333333333"`
	CuratorID string `json:"curator_id" example:"u4444444-4444-4444-4444-444444444444"`
	TeacherID string `json:"teacher_id" example:"u5555555-5555-5555-5555-555555555555"`
	Title     string `json:"title" example:"Группа А-1"`
}

// CreateCourse godoc
// @Summary ADMIN: Создание нового курса
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Название курса"
// @Param description formData string false "Описание курса"
// @Param is_main formData boolean false "Флаг основного курса (true/false)"
// @Param image_file formData file false "Изображение обложки курса"
// @Success 200 {object} map[string]string "id"
// @Router /admin/courses [post]
func (h *ContentAdminHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "File upload size exceeded limit.", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	isMain := r.FormValue("is_main") == "true"

	var fileHeader *multipart.FileHeader
	file, header, err := r.FormFile("image_file")
	if err == nil {
		file.Close()
		fileHeader = header
	}

	input := usecase.CreateCourseInput{
		Title:       title,
		Description: description,
		IsMain:      isMain,
		FileHeader:  fileHeader,
	}

	courseID, err := h.uc.CreateCourse(r.Context(), input)
	if err != nil {
		log.Printf("ERROR creating course: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": courseID})
}

// UploadMedia godoc
// @Summary ADMIN: Загрузка медиа-файла
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]string "url"
// @Router /admin/media/upload [post]
func (h *ContentAdminHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 10 << 20
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "File upload size exceeded limit.", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, err := h.uc.UploadMedia(r.Context(), header)
	if err != nil {
		log.Printf("Error uploading media: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

// UpdateCourseSettings godoc
// @Summary ADMIN: Обновление настроек курса
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} map[string]string "status"
// @Router /admin/courses/{id}/settings [put]
func (h *ContentAdminHandler) UpdateCourseSettings(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 10 << 20
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "File upload size exceeded limit.", http.StatusBadRequest)
		return
	}

	courseID := chi.URLParam(r, "id")
	parseBool := func(key string) bool { return r.FormValue(key) == "true" || r.FormValue(key) == "on" }

	var fileHeader *multipart.FileHeader
	if f, header, err := r.FormFile("cover_image"); err == nil {
		f.Close()
		fileHeader = header
	}

	input := usecase.UpdateCourseSettingsInput{
		CourseID:            courseID,
		Title:               r.FormValue("title"),
		Description:         r.FormValue("description"),
		IsMain:              parseBool("is_main"),
		Status:              domain.CourseStatus(r.FormValue("status")),
		HasHomework:         parseBool("has_homework"),
		IsHomeworkMandatory: parseBool("is_homework_mandatory"),
		IsTestMandatory:     parseBool("is_test_mandatory"),
		IsProjectMandatory:  parseBool("is_project_mandatory"),
		IsDiscordMandatory:  parseBool("is_discord_mandatory"),
		IsAntiCopyEnabled:   parseBool("is_anti_copy_enabled"),
		FileHeader:          fileHeader,
	}

	if err := h.uc.UpdateCourseSettings(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// GetAllCourses godoc
// @Summary ADMIN: Список всех курсов
// @Tags Admin-Content
// @Produce json
// @Success 200 {array} domain.Course
// @Router /admin/courses [get]
func (h *ContentAdminHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.uc.GetAllCourses(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourseStructure godoc
// @Summary ADMIN: Содержание курса (Дерево)
// @Tags Admin-Content
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} domain.CourseStructure
// @Router /admin/courses/{id}/structure [get]
func (h *ContentAdminHandler) GetCourseStructure(w http.ResponseWriter, r *http.Request) {
	structure, err := h.uc.GetCourseStructure(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(structure)
}

// CreateModule godoc
// @Summary ADMIN: Создание модуля
// @Tags Admin-Content
// @Accept json
// @Produce json
// @Param request body CreateModuleRequest true "Данные модуля"
// @Success 200 {object} map[string]string "id"
// @Router /admin/modules [post]
func (h *ContentAdminHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	var req CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecase.CreateModuleInput{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		OrderNum:    req.OrderNum,
	}

	id, err := h.uc.CreateModule(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// DeleteModule godoc
// @Summary ADMIN: Удаление модуля
// @Tags Admin-Content
// @Param id path string true "Module ID"
// @Success 200
// @Router /admin/modules/{id} [delete]
func (h *ContentAdminHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.DeleteModule(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateLesson godoc
// @Summary ADMIN: Добавление урока
// @Tags Admin-Content
// @Accept multipart/form-data
// @Param course_id formData string true "ID курса"
// @Param module_id formData string false "ID модуля (опционально)"
// @Param teacher_id formData string false "ID преподавателя (опционально)"
// @Param title formData string true "Название урока"
// @Param order_num formData int true "Порядковый номер"
// @Success 200 {object} map[string]string "id"
// @Router /admin/lessons [post]
func (h *ContentAdminHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	const MAX_VIDEO_SIZE = 500 << 20
	if err := r.ParseMultipartForm(MAX_VIDEO_SIZE); err != nil {
		http.Error(w, "File too large.", http.StatusBadRequest)
		return
	}

	orderNum, _ := strconv.Atoi(r.FormValue("order_num"))
	var vH, pH *multipart.FileHeader
	if f, head, err := r.FormFile("video_file"); err == nil {
		f.Close()
		vH = head
	}
	if f, head, err := r.FormFile("presentation_file"); err == nil {
		f.Close()
		pH = head
	}

	input := usecase.CreateLessonInput{
		CourseID:         r.FormValue("course_id"),
		ModuleID:         r.FormValue("module_id"),
		TeacherID:        r.FormValue("teacher_id"),
		Title:            r.FormValue("title"),
		OrderNum:         orderNum,
		ContentText:      r.FormValue("content_text"),
		VideoFile:        vH,
		PresentationFile: pH,
	}

	id, err := h.uc.CreateLesson(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// DeleteLesson godoc
// @Summary ADMIN: Удаление урока
// @Tags Admin-Content
// @Param id path string true "Lesson ID"
// @Success 200
// @Router /admin/lessons/{id} [delete]
func (h *ContentAdminHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.DeleteLesson(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateTest godoc
// @Summary ADMIN: Добавление теста
// @Tags Admin-Content
// @Accept json
// @Produce json
// @Param request body CreateTestRequest true "Данные теста"
// @Success 200 {object} map[string]string "id"
// @Router /admin/tests [post]
func (h *ContentAdminHandler) CreateTest(w http.ResponseWriter, r *http.Request) {
	var req CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := usecase.CreateTestInput{
		CourseID:     req.CourseID,
		LessonNumber: req.LessonNumber,
		LessonID:     req.LessonID,
		Title:        req.Title,
		Description:  req.Description,
		PassingScore: req.PassingScore,
	}
	id, err := h.uc.CreateTest(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// DeleteTest godoc
// @Summary ADMIN: Удаление теста
// @Tags Admin-Content
// @Param id path string true "Test ID"
// @Success 200
// @Router /admin/tests/{id} [delete]
func (h *ContentAdminHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.DeleteTest(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateProject godoc
// @Summary ADMIN: Добавление проекта
// @Tags Admin-Content
// @Accept json
// @Produce json
// @Param request body CreateProjectRequest true "Данные проекта"
// @Success 200 {object} map[string]string "id"
// @Router /admin/projects [post]
func (h *ContentAdminHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := usecase.CreateProjectInput{
		CourseID:     req.CourseID,
		LessonNumber: req.LessonNumber,
		LessonID:     req.LessonID,
		Title:        req.Title,
		Description:  req.Description,
		MaxScore:     req.MaxScore,
	}
	id, err := h.uc.CreateProject(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// DeleteProject godoc
// @Summary ADMIN: Удаление проекта
// @Tags Admin-Content
// @Param id path string true "Project ID"
// @Success 200
// @Router /admin/projects/{id} [delete]
func (h *ContentAdminHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.uc.DeleteProject(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateUser godoc
// @Summary ADMIN: Создание пользователя (Полный профиль + Родители)
// @Description Регистрация сотрудника или ученика. Поддерживает несколько родителей и привязку к курсу/группе.
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param request body CreateFullUserRequest true "Данные пользователя"
// @Success 200 {object} map[string]string "result mapping"
// @Router /admin/users [post]
func (h *ContentAdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateFullUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	birthDate, _ := time.Parse("2006-01-02", req.BirthDateStr)

	input := usecase.ExtendedCreateUserInput{
		FullName:        req.FullName,
		Email:           req.Email,
		Password:        req.Password,
		Role:            req.Role,
		Phone:           req.Phone,
		City:            req.City,
		SchoolName:      req.SchoolName,
		Language:        req.Language,
		Gender:          req.Gender,
		BirthDate:       birthDate,
		Whatsapp:        req.Whatsapp,
		Telegram:        req.Telegram,
		ExperienceYears: req.ExperienceYears,
		CourseID:        req.CourseID,
		StreamID:        req.StreamID,
		GroupID:         req.GroupID,
		Parents:         req.Parents,
	}

	result, err := h.uc.CreateFullUser(r.Context(), input)
	if err != nil {
		log.Printf("ERROR creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetUserInfo godoc
// @Summary ADMIN: Информация о конкретном пользователе (Карточка)
// @Tags Admin-Users
// @Produce json
// @Param id path string true "UserID"
// @Success 200 {object} map[string]interface{}
// @Router /admin/user/{id} [get]
func (h *ContentAdminHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	res, err := h.uc.GetUserInfo(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateUser godoc
// @Summary ADMIN: Изменить данные пользователя
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param id path string true "UserID"
// @Router /admin/users/{id} [put]
func (h *ContentAdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var req CreateFullUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecase.ExtendedCreateUserInput{
		FullName:        req.FullName,
		Email:           req.Email,
		Role:            req.Role,
		Phone:           req.Phone,
		City:            req.City,
		SchoolName:      req.SchoolName,
		Language:        req.Language,
		Gender:          req.Gender,
		Whatsapp:        req.Whatsapp,
		Telegram:        req.Telegram,
		ExperienceYears: req.ExperienceYears,
	}

	if err := h.uc.UpdateUser(r.Context(), userID, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteUser godoc
// @Summary ADMIN: Удалить пользователя
// @Tags Admin-Users
// @Param id path string true "UserID"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id} [delete]
func (h *ContentAdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if err := h.uc.DeleteUser(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}


// GetUsersList godoc
// @Summary ADMIN: Список пользователей (по ролям)
// @Tags Admin-Users
// @Produce json
// @Param role query string false "Фильтр: student, teacher, curator, moderator, admin"
// @Success 200 {array} domain.User
// @Router /admin/users [get]
func (h *ContentAdminHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	users, _ := h.uc.GetUsersList(r.Context(), domain.Role(role))
	json.NewEncoder(w).Encode(users)
}

// GetDetailedStudents godoc
// @Summary ADMIN: Детальный список учеников (Таблица)
// @Tags Admin-Users
// @Produce json
// @Param course_id query string false "Фильтр по курсу"
// @Success 200 {array} domain.StudentTableItem
// @Router /admin/students/detailed [get]
func (h *ContentAdminHandler) GetDetailedStudents(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	list, err := h.uc.GetDetailedStudents(r.Context(), courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GetDetailedTeachers godoc
// @Summary ADMIN: Список учителей (Таблица)
// @Tags Admin-Users
// @Produce json
// @Success 200 {array} domain.TeacherTableItem
// @Router /admin/teachers/detailed [get]
func (h *ContentAdminHandler) GetDetailedTeachers(w http.ResponseWriter, r *http.Request) {
	list, err := h.uc.GetDetailedTeachers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GetDetailedCurators godoc
// @Summary ADMIN: Список кураторов (Таблица)
// @Tags Admin-Users
// @Produce json
// @Success 200 {array} domain.CuratorTableItem
// @Router /admin/curators/detailed [get]
func (h *ContentAdminHandler) GetDetailedCurators(w http.ResponseWriter, r *http.Request) {
	list, err := h.uc.GetDetailedCurators(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GetDetailedModerators godoc
// @Summary ADMIN: Список модераторов (Таблица)
// @Tags Admin-Users
// @Produce json
// @Success 200 {array} domain.ModeratorTableItem
// @Router /admin/moderators/detailed [get]
func (h *ContentAdminHandler) GetDetailedModerators(w http.ResponseWriter, r *http.Request) {
	list, err := h.uc.GetDetailedModerators(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// GetAllUsersTable godoc
// @Summary ADMIN: Показать всех пользователей (Таблица)
// @Tags Admin-Users
// @Produce json
// @Success 200 {array} domain.AllUsersTableItem
// @Router /admin/users/all [get]
func (h *ContentAdminHandler) GetAllUsersTable(w http.ResponseWriter, r *http.Request) {
	list, err := h.uc.GetAllUsersTable(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// EnrollUser godoc
// @Summary ADMIN: Запись на курс
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param request body EnrollRequest true "UserID и CourseID"
// @Success 200 {object} map[string]string
// @Router /admin/enroll [post]
func (h *ContentAdminHandler) EnrollUser(w http.ResponseWriter, r *http.Request) {
	var req EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.uc.EnrollStudent(r.Context(), req.UserID, req.CourseID)
	if err != nil {
		log.Printf("ERROR enrolling user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "enrolled"})
}

// GetCourseStudents godoc
// @Summary ADMIN: Ученики конкретного курса
// @Tags Admin-Users
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {array} domain.AdminStudentProgress
// @Router /admin/courses/{id}/students [get]
func (h *ContentAdminHandler) GetCourseStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.uc.GetCourseStudents(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetCourseStats godoc
// @Summary ADMIN: Статистика курса
// @Tags Admin-Stats
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} domain.AdminCourseStats
// @Router /admin/courses/{id}/stats [get]
func (h *ContentAdminHandler) GetCourseStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.uc.GetCourseStats(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// CreateStream godoc
// @Summary ADMIN: Создать поток
// @Tags Admin-Staff
// @Accept json
// @Produce json
// @Param request body CreateStreamRequest true "Данные потока"
// @Success 200 {object} map[string]string "id"
// @Router /admin/streams [post]
func (h *ContentAdminHandler) CreateStream(w http.ResponseWriter, r *http.Request) {
	var req CreateStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	startDate, _ := time.Parse("2006-01-02", req.StartDateStr)
	id, err := h.uc.CreateStream(r.Context(), usecase.CreateStreamInput{
		CourseID:  req.CourseID,
		Title:     req.Title,
		StartDate: startDate,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetStreams godoc
// @Summary ADMIN: Список потоков курса
// @Tags Admin-Staff
// @Produce json
// @Param course_id query string false "ID курса (опционально)"
// @Success 200 {array} domain.Stream
// @Router /admin/streams [get]
func (h *ContentAdminHandler) GetStreams(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	streams, _ := h.uc.GetStreamsByCourse(r.Context(), courseID)
	json.NewEncoder(w).Encode(streams)
}

// CreateGroup godoc
// @Summary ADMIN: Создать группу
// @Tags Admin-Staff
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "Данные группы"
// @Success 200 {object} map[string]string "id"
// @Router /admin/groups [post]
func (h *ContentAdminHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateGroupInput
	json.NewDecoder(r.Body).Decode(&req)
	id, _ := h.uc.CreateGroup(r.Context(), req)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetGroups godoc
// @Summary ADMIN: Список групп потока
// @Tags Admin-Staff
// @Produce json
// @Param stream_id query string false "ID потока (опционально)"
// @Success 200 {array} domain.Group
// @Router /admin/groups [get]
func (h *ContentAdminHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	streamID := r.URL.Query().Get("stream_id")
	groups, _ := h.uc.GetGroupsByStream(r.Context(), streamID)
	json.NewEncoder(w).Encode(groups)
}