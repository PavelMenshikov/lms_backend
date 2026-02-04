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
// @Tags Student-Learning
// @Produce json
// @Success 200 {array} domain.StudentCoursePreview
// @Router /my-courses [get]
func (h *LearningHandler) GetMyCourses(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	courses, err := h.uc.GetMyCourses(r.Context(), userData.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourseContent godoc
// @Summary УЧЕНИК: Страница курса
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID курса"
// @Success 200 {object} domain.StudentCourseView
// @Router /courses/{id} [get]
func (h *LearningHandler) GetCourseContent(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	courseID := chi.URLParam(r, "id")
	view, err := h.uc.GetCourseContent(r.Context(), courseID, userData.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(view)
}

// GetLessonDetail godoc
// @Summary УЧЕНИК: Просмотр урока
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID урока"
// @Success 200 {object} domain.StudentLessonDetail
// @Router /lessons/{id} [get]
func (h *LearningHandler) GetLessonDetail(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	lessonID := chi.URLParam(r, "id")
	lesson, err := h.uc.GetLessonDetail(r.Context(), lessonID, userData.UserID)
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
// @Success 200 {object} map[string]string "status"
// @Router /lessons/{id}/assignment [post]
func (h *LearningHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	const MAX_SIZE = 10 << 20 // 10MB
	r.ParseMultipartForm(MAX_SIZE)

	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	lessonID := chi.URLParam(r, "id")

	var fileHeader *multipart.FileHeader
	if f, h, err := r.FormFile("file"); err == nil {
		f.Close()
		fileHeader = h
	}

	input := usecase.SubmitAssignmentInput{
		LessonID:   lessonID,
		UserID:     userData.UserID,
		TextAnswer: r.FormValue("text_answer"),
		FileHeader: fileHeader,
	}

	err := h.uc.SubmitAssignment(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "submitted"})
}

// CompleteLesson godoc
// @Summary УЧЕНИК: Завершить урок (без ДЗ)
// @Tags Student-Learning
// @Produce json
// @Param id path string true "ID урока"
// @Success 200 {object} map[string]string "status"
// @Router /lessons/{id}/complete [post]
func (h *LearningHandler) CompleteLesson(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	lessonID := chi.URLParam(r, "id")

	err := h.uc.CompleteLesson(r.Context(), lessonID, userData.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
}
