package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/schedule/usecase"
)

type ScheduleHandler struct {
	uc *usecase.ScheduleUseCase
}

func NewScheduleHandler(uc *usecase.ScheduleUseCase) *ScheduleHandler {
	return &ScheduleHandler{uc: uc}
}

// GetWeeklySchedule godoc
// @Summary USER: Расписание на неделю
// @Description Возвращает уроки, сгруппированные по дням недели.
// @Tags Schedule
// @Produce json
// @Param date query string false "Дата недели (YYYY-MM-DD), по умолчанию сегодня"
// @Success 200 {object} domain.WeeklySchedule
// @Router /schedule/weekly [get]
func (h *ScheduleHandler) GetWeeklySchedule(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)

	dateStr := r.URL.Query().Get("date")
	targetDate := time.Now()
	if dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			targetDate = d
		}
	}

	schedule, err := h.uc.GetWeeklySchedule(r.Context(), userData.UserID, targetDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

// GetMonthlySchedule godoc
// @Summary USER: Расписание на месяц
// @Description Возвращает список занятий для календарной сетки месяца.
// @Tags Schedule
// @Produce json
// @Param year query int false "Год"
// @Param month query int false "Месяц"
// @Success 200 {object} domain.MonthlySchedule
// @Router /schedule/monthly [get]
func (h *ScheduleHandler) GetMonthlySchedule(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)

	now := time.Now()
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}

	schedule, err := h.uc.GetMonthlySchedule(r.Context(), userData.UserID, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}
