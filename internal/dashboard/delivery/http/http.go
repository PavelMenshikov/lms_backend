package http

import (
	"encoding/json"
	"net/http"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/dashboard/usecase"
	"lms_backend/internal/domain"
)

type DashboardHandler struct {
	uc *usecase.DashboardUseCase
}

func NewDashboardHandler(uc *usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

// GetUserHome godoc
// @Summary УЧЕНИК: Главная страница
// @Description Возвращает данные для дашборда ученика или родителя.
// @Tags Dashboard
// @Produce json
// @Success 200 {object} domain.HomeDashboard
// @Router /dashboard/home [get]
func (h *DashboardHandler) GetUserHome(w http.ResponseWriter, r *http.Request) {
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)

	user := &domain.User{
		ID:   userCtxData.UserID,
		Role: userCtxData.Role,
	}

	data, err := h.uc.GetUserHomeData(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetAdminDashboard godoc
// @Summary ADMIN: Статистика главной страницы
// @Description Возвращает агрегированные данные для админ-панели (счетчики, графики, успеваемость).
// @Tags Admin-Stats
// @Produce json
// @Success 200 {object} domain.AdminHomeDashboard
// @Router /admin/dashboard/stats [get]
func (h *DashboardHandler) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.uc.GetAdminDashboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
