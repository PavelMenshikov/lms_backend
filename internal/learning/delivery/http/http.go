package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"lms_backend/internal/httperror"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/domain"
	"lms_backend/internal/learning/usecase"
)

type LearningHandler struct {
	uc *usecase.LearningUseCase
}

func NewLearningHandler(uc *usecase.LearningUseCase) *LearningHandler {
	return &LearningHandler{uc: uc}
}

// GetMyCourses godoc
// @Summary УЧЕНИК: Мои курсы
// @Description Получить список курсов, на которые записан текущий пользователь.
// @Tags Student-Learning
// @Produce json
// @Success 200 {array} domain.StudentCoursePreview
// @Router /my-courses [get]
func (h *LearningHandler) GetMyCourses(w http.ResponseWriter, r *http.Request) {
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	courses, err := h.uc.GetMyCourses(r.Context(), userCtxData.UserID)
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourseContent godoc
// @Summary УЧЕНИК: Страница курса
// @Description Получить полную структуру курса (модули, уроки) с отметками о прохождении.
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} domain.StudentCourseView
// @Router /courses/{id} [get]
func (h *LearningHandler) GetCourseContent(w http.ResponseWriter, r *http.Request) {
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	courseID := chi.URLParam(r, "id")
	view, err := h.uc.GetCourseContent(r.Context(), courseID, userCtxData.UserID)
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(view)
}

// GetLessonDetail godoc
// @Summary УЧЕНИК: Просмотр урока
// @Description Получить контент конкретного урока (видео, текст).
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID урока"
// @Success 200 {object} domain.StudentLessonDetail
// @Router /lessons/{id} [get]
func (h *LearningHandler) GetLessonDetail(w http.ResponseWriter, r *http.Request) {
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	lessonID := chi.URLParam(r, "id")
	lesson, err := h.uc.GetLessonDetail(r.Context(), lessonID, userCtxData.UserID)
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lesson)
}

// SubmitAssignment godoc
// @Summary УЧЕНИК: Сдать домашнее задание
// @Description Отправка текстового ответа или файла.
// @Tags Student-Learning
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID урока"
// @Param text_answer formData string false "Текст ответа"
// @Param file formData file false "Файл с заданием (можно несколько)"
// @Success 200 {object} map[string]string
// @Router /lessons/{id}/assignment [post]
func (h *LearningHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	const MAX_SIZE = 50 << 20
	if err := r.ParseMultipartForm(MAX_SIZE); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	lessonID := chi.URLParam(r, "id")

	files := r.MultipartForm.File["file"]
	var fileHeaders []*multipart.FileHeader
	for _, fh := range files {
		fileHeaders = append(fileHeaders, fh)
	}
	input := usecase.SubmitAssignmentInput{
		LessonID:    lessonID,
		UserID:      userCtxData.UserID,
		TextAnswer:  r.FormValue("text_answer"),
		FileHeaders: fileHeaders,
	}
	if err := h.uc.SubmitAssignment(r.Context(), input); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "submitted"})
}

type SetAttendanceRequest struct {
	Status         string `json:"status"`
	RecordingURL   string `json:"recording_url,omitempty"`
	TeacherComment string `json:"teacher_comment,omitempty"`
}

// SetLessonAttendance godoc
// @Summary УЧЕНИК: Отметить посещение урока
// @Description Отметить урок с указанием статуса посещения, ссылки на запись и комментария.
// @Tags Student-Learning
// @Accept json
// @Produce json
// @Param id path string true "ID урока"
// @Param body body SetAttendanceRequest true "Данные посещения"
// @Success 200 {object} map[string]string
// @Router /lessons/{id}/attendance [post]
func (h *LearningHandler) SetLessonAttendance(w http.ResponseWriter, r *http.Request) {
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	lessonID := chi.URLParam(r, "id")

	var req SetAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	input := usecase.SetAttendanceInput{
		LessonID:       lessonID,
		UserID:         userCtxData.UserID,
		Status:         req.Status,
		RecordingURL:   req.RecordingURL,
		TeacherComment: req.TeacherComment,
	}
	if err := h.uc.SetLessonAttendance(r.Context(), input); err != nil {
		httperror.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

// GetTeachers godoc
// @Summary УЧЕНИК: Список преподавателей
// @Description Получить список преподавателей с их рейтингом.
// @Tags Student-Learning
// @Produce json
// @Success 200 {array} domain.TeacherPublicInfo
// @Router /teachers [get]
func (h *LearningHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.uc.GetTeachers(r.Context())
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

// GetTeacherDetails godoc
// @Summary УЧЕНИК: Детали преподавателя
// @Description Получить полную инфо о преподавателе и его отзывы.
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID преподавателя"
// @Success 200 {object} domain.TeacherPublicInfo
// @Router /teachers/{id} [get]
func (h *LearningHandler) GetTeacherDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	teacher, err := h.uc.GetTeacherDetails(r.Context(), id)
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

type ReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// AddReview godoc
// @Summary УЧЕНИК: Оставить отзыв преподавателю
// @Description Добавить оценку и текстовый отзыв.
// @Tags Student-Learning
// @Accept json
// @Produce json
// @Param id path string true "ID преподавателя"
// @Param request body ReviewRequest true "Данные отзыва"
// @Success 201
// @Router /teachers/{id}/reviews [post]
func (h *LearningHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	teacherID := chi.URLParam(r, "id")
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	input := usecase.AddReviewInput{
		TeacherID: teacherID,
		StudentID: userCtxData.UserID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	if err := h.uc.AddReview(r.Context(), input); err != nil {
		httperror.Internal(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetTest godoc
// @Summary USER: Детали теста
// @Tags Student-Learning
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} domain.Test
// @Router /tests/{id} [get]
func (h *LearningHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	test, err := h.uc.GetTest(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(&domain.Test{ID: id})
			return
		}
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}

// GetProject godoc
// @Summary USER: Детали проекта
// @Tags Student-Learning
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Project
// @Router /projects/{id} [get]
func (h *LearningHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	proj, err := h.uc.GetProject(r.Context(), id)
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	if proj == nil {
		proj = &domain.Project{ID: id}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proj)
}

// GetAllCourses godoc
// @Summary USER: Каталог курсов
// @Description Получить список всех активных курсов.
// @Tags Courses
// @Produce json
// @Success 200 {array} domain.Course
// @Router /courses [get]
func (h *LearningHandler) GetAllCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.uc.GetAllCourses(r.Context())
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	if courses == nil {
		courses = []*domain.Course{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// SubmitTest godoc
// @Summary УЧЕНИК: Отправить тест
// @Tags Student-Learning
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} map[string]any
// @Router /tests/{id}/submit [post]
func (h *LearningHandler) SubmitTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "submitted"})
}

// SubmitProject godoc
// @Summary УЧЕНИК: Отправить проект
// @Tags Student-Learning
// @Accept mpfd
// @Produce json
// @Param id path string true "Project ID"
// @Param file formData file false "File attachment"
// @Success 200 {object} map[string]any
// @Router /projects/{id}/submission [post]
func (h *LearningHandler) SubmitProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "submitted"})
}

// GetTeacherDashboard godoc
// @Summary ТИЧЕР: Дашборд ЛК
// @Tags Teacher-Dashboard
// @Produce json
// @Success 200 {object} domain.TeacherDashboardData
// @Router /teacher/profile [get]
// GetTeacherCertificates godoc
// @Summary USER: Сертификаты преподавателя
// @Tags Teacher-Certificates
// @Produce json
// @Success 200 {array} domain.TeacherCertificate
// @Router /teacher/certificates [get]
func (h *LearningHandler) GetTeacherCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]domain.TeacherCertificate{})
}

func (h *LearningHandler) GetTeacherDashboard(w http.ResponseWriter, r *http.Request) {
	userData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userData.Role != domain.RoleTeacher {
		http.Error(w, "Forbidden: Only for teachers", http.StatusForbidden)
		return
	}

	dashboard, err := h.uc.GetTeacherDashboard(r.Context(), userData.UserID)
	if err != nil {
		httperror.Internal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}
