package http

import (
	"encoding/json"
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
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
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
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
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
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
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
// @Param file formData file false "Файл с заданием"
// @Success 200 {object} map[string]string
// @Router /lessons/{id}/assignment [post]
func (h *LearningHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	const MAX_SIZE = 10 << 20
	if err := r.ParseMultipartForm(MAX_SIZE); err != nil {
		httperror.BadRequest(w, err)
		return
	}
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	lessonID := chi.URLParam(r, "id")
	var fileHeader *multipart.FileHeader
	if f, header, err := r.FormFile("file"); err == nil {
		f.Close()
		fileHeader = header
	}
	input := usecase.SubmitAssignmentInput{
		LessonID:   lessonID,
		UserID:     userCtxData.UserID,
		TextAnswer: r.FormValue("text_answer"),
		FileHeader: fileHeader,
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
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
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
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
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
		http.Error(w, "Test not found", http.StatusNotFound)
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
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proj)
}

// GetTeacherDashboard godoc
// @Summary ТИЧЕР: Дашборд ЛК
// @Tags Teacher-Dashboard
// @Produce json
// @Success 200 {object} domain.TeacherDashboardData
// @Router /teacher/profile [get]
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
