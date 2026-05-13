package http

import (
	"encoding/json"
	"lms_backend/internal/statistics/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StatisticsHandler struct {
	uc usecase.StatisticsUseCase
}

func NewStatisticsHandler(uc usecase.StatisticsUseCase) *StatisticsHandler {
	return &StatisticsHandler{uc: uc}
}

// GetStudentStatistics godoc
// @Summary Получить статистику ученика
// @Tags Statistics
// @Param studentId path string true "Student ID"
// @Success 200 {object} domain.StudentStatistics
// @Router /api/statistics/students/{studentId} [get]
func (h *StatisticsHandler) GetStudentStatistics(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	stats, err := h.uc.GetStudentStatistics(r.Context(), studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// RefreshStudentStatistics godoc
// @Summary Пересчитать статистику ученика
// @Tags Statistics
// @Param studentId path string true "Student ID"
// @Success 200 {object} domain.StudentStatistics
// @Router /api/statistics/students/{studentId}/refresh [post]
func (h *StatisticsHandler) RefreshStudentStatistics(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	stats, err := h.uc.RefreshStudentStatistics(r.Context(), studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
