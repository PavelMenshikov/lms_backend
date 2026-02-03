package http

import (
	"encoding/json"
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
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Role      domain.Role `json:"role"`
}

type EnrollRequest struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

type CreateTestRequest struct {
	LessonID     string `json:"lesson_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PassingScore int    `json:"passing_score"`
}

type CreateProjectRequest struct {
	LessonID    string `json:"lesson_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MaxScore    int    `json:"max_score"`
}

// @Summary ADMIN: Создание нового курса
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Название курса"
// @Param image_file formData file false "Файл изображения"
// @Success 200 {object} map[string]string
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": courseID})
}

// @Summary ADMIN: Обновление настроек курса
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} map[string]string
// @Router /admin/courses/{id}/settings [put]
func (h *ContentAdminHandler) UpdateCourseSettings(w http.ResponseWriter, r *http.Request) {
	const MAX_UPLOAD_SIZE = 10 << 20
	r.ParseMultipartForm(MAX_UPLOAD_SIZE)

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
}

// @Summary ADMIN: Создание модуля
// @Tags Admin-Content
// @Accept json
// @Produce json
// @Router /admin/modules [post]
func (h *ContentAdminHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateModuleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

// @Summary ADMIN: Список всех курсов
// @Tags Admin-Content
// @Produce json
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

// @Summary ADMIN: Создание пользователя
// @Tags Admin-Users
// @Accept json
// @Produce json
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// @Summary ADMIN: Запись на курс
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Router /admin/enroll [post]
func (h *ContentAdminHandler) EnrollUser(w http.ResponseWriter, r *http.Request) {
	var req EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.uc.EnrollStudent(r.Context(), req.UserID, req.CourseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "enrolled"})
}

// @Summary ADMIN: Список учеников курса
// @Tags Admin-Users
// @Produce json
// @Param id path string true "ID курса"
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