package http

import (
	"encoding/json"
	"lms_backend/internal/attendance/usecase"
	"lms_backend/internal/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type AttendanceHandler struct {
	uc usecase.AttendanceUseCase
}

func NewAttendanceHandler(uc usecase.AttendanceUseCase) *AttendanceHandler {
	return &AttendanceHandler{uc: uc}
}

type MarkAttendanceRequest struct {
	StudentID string                   `json:"student_id"`
	Status    domain.AttendanceStatus  `json:"status"`
	Reason    *string                  `json:"reason,omitempty"`
	Comment   *string                  `json:"comment,omitempty"`
}

type UpdateAttendanceRequest struct {
	Status  domain.AttendanceStatus `json:"status"`
	Reason  *string                 `json:"reason,omitempty"`
	Comment *string                 `json:"comment,omitempty"`
}

// GetStudentCalendar godoc
// @Summary Получить календарь посещаемости ученика
// @Tags Attendance
// @Param studentId path string true "Student ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} domain.AttendanceRecord
// @Router /api/attendance/students/{studentId}/calendar [get]
func (h *AttendanceHandler) GetStudentCalendar(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0) // По умолчанию последний месяц
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now()
	}

	records, err := h.uc.GetStudentCalendar(r.Context(), studentID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// MarkLessonAttendance godoc
// @Summary Отметить посещаемость урока
// @Tags Attendance
// @Param lessonId path string true "Lesson ID"
// @Param body body MarkAttendanceRequest true "Attendance data"
// @Success 200 {object} map[string]string
// @Router /api/attendance/lessons/{lessonId} [patch]
func (h *AttendanceHandler) MarkLessonAttendance(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonId")

	var req MarkAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из контекста (должен быть установлен middleware)
	userID := r.Context().Value("user_id").(string)

	err := h.uc.MarkAttendance(r.Context(), lessonID, req.StudentID, req.Status, req.Reason, req.Comment, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Attendance marked successfully"})
}

// GetStudentStats godoc
// @Summary Получить статистику посещаемости ученика
// @Tags Attendance
// @Param studentId path string true "Student ID"
// @Success 200 {object} map[string]int
// @Router /api/attendance/students/{studentId}/stats [get]
func (h *AttendanceHandler) GetStudentStats(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	stats, err := h.uc.GetStudentStats(r.Context(), studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetLessonAttendance godoc
// @Summary Получить посещаемость урока
// @Tags Attendance
// @Param lessonId path string true "Lesson ID"
// @Success 200 {array} domain.AttendanceRecord
// @Router /api/attendance/lessons/{lessonId} [get]
func (h *AttendanceHandler) GetLessonAttendance(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonId")

	records, err := h.uc.GetLessonAttendance(r.Context(), lessonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}
