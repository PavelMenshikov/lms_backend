package http

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

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

type CreateUserRequest struct {
	FirstName string      `json:"first_name" example:"Иван"`
	LastName  string      `json:"last_name" example:"Иванов"`
	Email     string      `json:"email" example:"student@test.kz"`
	Password  string      `json:"password" example:"secret123"`
	Role      domain.Role `json:"role" example:"student"`
}

type EnrollRequest struct {
	UserID   string `json:"user_id" example:"a0000000-0000-0000-0000-000000000001"`
	CourseID string `json:"course_id" example:"c1111111-1111-1111-1111-111111111111"`
}

type CreateTestRequest struct {
	LessonID     string `json:"lesson_id" example:"l2222222-2222-2222-2222-222222222222"`
	Title        string `json:"title" example:"Итоговый тест по модулю 1"`
	Description  string `json:"description" example:"Тест на проверку базовых знаний Go"`
	PassingScore int    `json:"passing_score" example:"70"`
}

type CreateProjectRequest struct {
	LessonID    string `json:"lesson_id" example:"l2222222-2222-2222-2222-222222222222"`
	Title       string `json:"title" example:"Финальный проект"`
	Description string `json:"description" example:"Разработка API на Go"`
	MaxScore    int    `json:"max_score" example:"100"`
}

type CreateModuleRequest struct {
	CourseID    string `json:"course_id" example:"c1111111-1111-1111-1111-111111111111"`
	Title       string `json:"title" example:"Основы синтаксиса"`
	Description string `json:"description" example:"Типы данных, переменные, циклы"`
	OrderNum    int    `json:"order_num" example:"1"`
}

// CreateCourse godoc
// @Summary ADMIN: Создание нового курса
// @Description Создает карточку курса с загрузкой изображения.
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Название курса"
// @Param description formData string false "Описание курса"
// @Param is_main formData boolean false "Флаг основного курса (true/false)"
// @Param image_file formData file false "Изображение обложки курса"
// @Success 200 {object} map[string]string "id"
// @Failure 400 {object} map[string]string "error"
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

// UpdateCourseSettings godoc
// @Summary ADMIN: Обновление настроек курса
// @Description Позволяет изменить параметры курса, обложку и статус.
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID курса"
// @Param title formData string false "Новое название"
// @Param description formData string false "Новое описание"
// @Param is_main formData boolean false "Тип курса"
// @Param status formData string false "Статус: draft, active, archived"
// @Param has_homework formData boolean false "Наличие ДЗ"
// @Param is_homework_mandatory formData boolean false "Обязательность ДЗ"
// @Param is_test_mandatory formData boolean false "Обязательность тестов"
// @Param is_project_mandatory formData boolean false "Обязательность проектов"
// @Param is_discord_mandatory formData boolean false "Обязательность Discord"
// @Param is_anti_copy_enabled formData boolean false "Запрет копирования"
// @Param cover_image formData file false "Новая обложка"
// @Success 200 {object} map[string]string "status"
// @Failure 500 {object} map[string]string "error"
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// GetAllCourses godoc
// @Summary ADMIN: Список всех курсов
// @Description Возвращает полный список курсов со всеми метаданными.
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
// @Description Возвращает иерархическую структуру: Курс -> Модули -> Уроки.
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
// @Description Добавляет модуль в учебный план курса.
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

// CreateLesson godoc
// @Summary ADMIN: Добавление урока
// @Description Создает урок с загрузкой видео и презентации в S3.
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param module_id formData string true "ID модуля"
// @Param teacher_id formData string true "ID преподавателя"
// @Param title formData string true "Название урока"
// @Param order_num formData int true "Порядковый номер"
// @Param content_text formData string false "HTML/JSON контент урока"
// @Param video_file formData file false "Видео файл (до 500MB)"
// @Param presentation_file formData file false "Файл презентации"
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

// CreateTest godoc
// @Summary ADMIN: Добавление теста
// @Description Создает тест и привязывает его к уроку.
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

// CreateProject godoc
// @Summary ADMIN: Добавление проекта
// @Description Создает финальный проект и привязывает его к уроку.
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
		LessonID:    req.LessonID,
		Title:       req.Title,
		Description: req.Description,
		MaxScore:    req.MaxScore,
	}
	id, err := h.uc.CreateProject(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// CreateUser godoc
// @Summary ADMIN: Создание пользователя
// @Description Создает аккаунт для ученика, учителя, родителя или куратора.
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Данные пользователя"
// @Success 200 {object} map[string]string "id"
// @Router /admin/users [post]
func (h *ContentAdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := usecase.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
	}
	id, err := h.uc.CreateUser(r.Context(), input)
	if err != nil {
		log.Printf("ERROR creating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// EnrollUser godoc
// @Summary ADMIN: Запись на курс
// @Description Привязывает ученика к конкретному курсу.
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param request body EnrollRequest true "Связка UserID и CourseID"
// @Success 200 {object} map[string]string "status"
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
// @Summary ADMIN: Список учеников курса
// @Description Возвращает список всех учеников с их текущим прогрессом на курсе.
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
// @Description Возвращает агрегированные данные: кол-во учеников, средний балл и т.д.
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