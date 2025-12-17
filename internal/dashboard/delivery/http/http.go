package http

import (
	"lms_backend/internal/dashboard/usecase"
	"net/http"
)

type DashboardHandler struct {
	uc *usecase.DashboardUseCase
}

func NewDashboardHandler(uc *usecase.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

// @Summary Get user home dashboard
// @Description Returns the home dashboard for the authenticated user
// @Tags dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /dashboard/home [get]
func (h *DashboardHandler) GetUserHome(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dashboard Data (TODO)"))
}
