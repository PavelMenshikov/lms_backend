package http

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/go-chi/chi/v5"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "submitted"})
}

// CompleteLesson godoc
// @Summary УЧЕНИК: Завершить урок (без ДЗ)
// @Description Отметить урок как пройденный (для уроков без обязательного ДЗ).
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID урока"
// @Success 200 {object} map[string]string
// @Router /lessons/{id}/complete [post]
func (h *LearningHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	lessonID := chi.URLParam(r, "id")
	err := h.uc.CompleteLesson(r.Context(), lessonID, userCtxData.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := usecase.AddReviewInput{
		TeacherID: teacherID,
		StudentID: userCtxData.UserID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	if err := h.uc.AddReview(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
