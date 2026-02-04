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
	FirstName    string      `json:"first_name" example:"Иван"`
	LastName     string      `json:"last_name" example:"Иванов"`
	Email        string      `json:"email" example:"student@test.kz"`
	Password     string      `json:"password" example:"secret123"`
	Role         domain.Role `json:"role" example:"student"`
	Phone        string      `json:"phone" example:"+79998887766"`
	City         string      `json:"city" example:"Алматы"`
	Language     string      `json:"language" example:"ru"`
	Gender       string      `json:"gender" example:"male"`
	BirthDateStr string      `json:"birth_date" example:"2000-01-01"`
	ParentName   string      `json:"parent_name" example:"Мама Иванова"`
	ParentPhone  string      `json:"parent_phone" example:"+7000..."`
	ParentEmail  string      `json:"parent_email" example:"mom@test.kz"`
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

// @Summary ADMIN: Загрузка медиа-файла
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл"
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

// @Summary ADMIN: Добавление урока
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param module_id formData string true "ID модуля"
// @Param teacher_id formData string true "ID преподавателя"
// @Param title formData string true "Название урока"
// @Param order_num formData int true "Номер"
// @Param content_text formData string false "HTML/JSON контент урока"
// @Param video_file formData file false "Видео (до 500MB)"
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

// @Summary ADMIN: Создание пользователя (Полный профиль + Родитель)
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param request body CreateFullUserRequest true "Данные пользователя"
// @Success 200 {object} map[string]string "id, role, parent_info"
// @Router /admin/users [post]
func (h *ContentAdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateFullUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	birthDate, _ := time.Parse("2006-01-02", req.BirthDateStr)

	input := usecase.ExtendedCreateUserInput{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Password:        req.Password,
		Role:            req.Role,
		Phone:           req.Phone,
		City:            req.City,
		Language:        req.Language,
		Gender:          req.Gender,
		BirthDate:       birthDate,
		ParentFirstName: req.ParentName,
		ParentPhone:     req.ParentPhone,
		ParentEmail:     req.ParentEmail,
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

// @Summary ADMIN: Изменить пользователя
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param id path string true "UserID"
// @Param request body CreateFullUserRequest true "Данные"
// @Success 200 {object} map[string]string
// @Router /admin/users/{id} [put]
func (h *ContentAdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var req CreateFullUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecase.ExtendedCreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
		Phone:     req.Phone,
		City:      req.City,
		Language:  req.Language,
		Gender:    req.Gender,
	}

	if err := h.uc.UpdateUser(r.Context(), userID, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// @Summary ADMIN: Запись на курс
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param request body EnrollRequest true "UserID и CourseID"
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

// @Summary ADMIN: Список пользователей
// @Description Получение таблицы пользователей с фильтрацией по роли.
// @Tags Admin-Users
// @Produce json
// @Param role query string false "Роль (student, teacher...)"
// @Success 200 {array} domain.User
// @Router /admin/users [get]
func (h *ContentAdminHandler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	users, err := h.uc.GetUsersList(r.Context(), domain.Role(role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @Summary ADMIN: Список учеников курса
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
