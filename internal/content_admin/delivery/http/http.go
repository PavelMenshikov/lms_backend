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
	FirstName       string      `json:"first_name" example:"Иван"`
	LastName        string      `json:"last_name" example:"Иванов"`
	Email           string      `json:"email" example:"student@test.kz"`
	Password        string      `json:"password" example:"secret123"`
	Role            domain.Role `json:"role" example:"student"`
	Phone           string      `json:"phone" example:"+79998887766"`
	City            string      `json:"city" example:"Алматы"`
	Language        string      `json:"language" example:"ru"`
	Gender          string      `json:"gender" example:"male"`
	BirthDateStr    string      `json:"birth_date" example:"2000-01-01"`
	Whatsapp        string      `json:"whatsapp" example:"https://wa.me/..."`
	Telegram        string      `json:"telegram" example:"https://t.me/..."`
	ParentFirstName string      `json:"parent_first_name" example:"Мама"`
	ParentLastName  string      `json:"parent_last_name" example:"Иванова"`
	ParentPhone     string      `json:"parent_phone" example:"+70001112233"`
	ParentEmail     string      `json:"parent_email" example:"mom@test.kz"`
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

// UploadMedia godoc
// @Summary ADMIN: Загрузка медиа-файла
// @Description Загружает произвольный файл в S3 (для использования в редакторе контента).
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} map[string]string "url"
// @Failure 400 {object} map[string]string "error"
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
// @Description Редактирование флагов (ДЗ, Тесты, Discord), статусов и обложки.
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID курса"
// @Param title formData string false "Название"
// @Param description formData string false "Описание"
// @Param is_main formData boolean false "Основной"
// @Param status formData string false "Статус"
// @Param cover_image formData file false "Новая обложка"
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
// @Description Получение краткого списка всех существующих курсов.
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
// @Description Возвращает иерархию Модули -> Уроки для конкретного курса.
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
// @Description Добавление нового модуля в учебный план курса.
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
// @Description Создает урок с опциональной загрузкой видео и презентации.
// @Tags Admin-Content
// @Accept multipart/form-data
// @Produce json
// @Param module_id formData string true "ID модуля"
// @Param teacher_id formData string true "ID преподавателя (лектора)"
// @Param title formData string true "Название урока"
// @Param order_num formData int true "Порядковый номер в модуле"
// @Param content_text formData string false "Markdown/HTML конспект"
// @Param video_file formData file false "Видео файл"
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
// @Description Создает проверочный тест к уроку.
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
// @Description Создает курсовой проект к уроку/модулю.
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
// @Summary ADMIN: Создание пользователя (Полный профиль + Родитель)
// @Description Регистрация сотрудника или ученика. Если роль student, автоматически создает связанного родителя.
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
		Whatsapp:        req.Whatsapp,
		Telegram:        req.Telegram,
		ParentFirstName: req.ParentFirstName,
		ParentLastName:  req.ParentLastName,
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

// UpdateUser godoc
// @Summary ADMIN: Изменить данные пользователя
// @Description Обновление анкетных данных любого пользователя.
// @Tags Admin-Users
// @Accept json
// @Produce json
// @Param id path string true "UserID"
// @Param request body CreateFullUserRequest true "Новые данные"
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
		Whatsapp:  req.Whatsapp,
		Telegram:  req.Telegram,
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
// @Description Полное удаление аккаунта из системы.
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
// @Description Получение таблицы всех сотрудников или студентов.
// @Tags Admin-Users
// @Produce json
// @Param role query string false "Фильтр: student, teacher, curator, moderator, admin"
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

// GetDetailedStudents godoc
// @Summary ADMIN: Детальный список учеников
// @Description Список учеников с расширенной инфо: курс, группа, куратор, прогресс.
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

// EnrollUser godoc
// @Summary ADMIN: Запись на курс
// @Description Привязка существующего ученика к курсу.
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
// @Description Краткий список учеников с их прогрессом.
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
// @Description Агрегированные метрики по курсу (кол-во, успеваемость).
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
// @Description Создание нового потока для курса.
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
// @Param course_id query string true "ID курса"
// @Success 200 {array} domain.Stream
// @Router /admin/streams [get]
func (h *ContentAdminHandler) GetStreams(w http.ResponseWriter, r *http.Request) {
	courseID := r.URL.Query().Get("course_id")
	streams, err := h.uc.GetStreamsByCourse(r.Context(), courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(streams)
}

// CreateGroup godoc
// @Summary ADMIN: Создать группу
// @Description Создание группы внутри потока, назначение куратора и учителя.
// @Tags Admin-Staff
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "Данные группы"
// @Success 200 {object} map[string]string "id"
// @Router /admin/groups [post]
func (h *ContentAdminHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.uc.CreateGroup(r.Context(), usecase.CreateGroupInput(req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetGroups godoc
// @Summary ADMIN: Список групп потока
// @Tags Admin-Staff
// @Produce json
// @Param stream_id query string true "ID потока"
// @Success 200 {array} domain.Group
// @Router /admin/groups [get]
func (h *ContentAdminHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	streamID := r.URL.Query().Get("stream_id")
	groups, err := h.uc.GetGroupsByStream(r.Context(), streamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
