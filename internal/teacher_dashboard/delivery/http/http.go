package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/domain"
	"lms_backend/internal/teacher_dashboard/usecase"
)

type TeacherDashboardHandler struct {
	uc *usecase.TeacherDashboardUseCase
}

func NewTeacherDashboardHandler(uc *usecase.TeacherDashboardUseCase) *TeacherDashboardHandler {
	return &TeacherDashboardHandler{uc: uc}
}

// GetTeacherMonthlyReport godoc
// @Summary TEACHER: Месячный отчёт
// @Tags Teacher
// @Param year query int false "Year (default: current)"
// @Param month query int false "Month 1-12 (default: current)"
// @Success 200 {object} domain.TeacherMonthlyReport
// @Router /teacher/monthly-report [get]
func (h *TeacherDashboardHandler) GetTeacherMonthlyReport(w http.ResponseWriter, r *http.Request) {
	userData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userData.Role != domain.RoleTeacher {
		http.Error(w, "Forbidden: Only for teachers", http.StatusForbidden)
		return
	}

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	year := time.Now().Year()
	month := int(time.Now().Month())

	if yearStr != "" {
		if v, err := strconv.Atoi(yearStr); err == nil {
			year = v
		}
	}
	if monthStr != "" {
		if v, err := strconv.Atoi(monthStr); err == nil && v >= 1 && v <= 12 {
			month = v
		}
	}

	report, err := h.uc.GetMonthlyReport(r.Context(), userData.UserID, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
