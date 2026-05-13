package http

import (
	"encoding/json"
	"lms_backend/internal/banner/usecase"
	"lms_backend/internal/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type BannerHandler struct {
	uc usecase.BannerUseCase
}

func NewBannerHandler(uc usecase.BannerUseCase) *BannerHandler {
	return &BannerHandler{uc: uc}
}

type CreateBannerRequest struct {
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	Type        domain.BannerType  `json:"type"`
	IsActive    bool               `json:"is_active"`
	Priority    int                `json:"priority"`
	StartDate   *time.Time         `json:"start_date,omitempty"`
	EndDate     *time.Time         `json:"end_date,omitempty"`
	TargetRoles []string           `json:"target_roles,omitempty"`
}

type UpdateBannerRequest struct {
	Title       string             `json:"title"`
	Content     string             `json:"content"`
	Type        domain.BannerType  `json:"type"`
	IsActive    bool               `json:"is_active"`
	Priority    int                `json:"priority"`
	StartDate   *time.Time         `json:"start_date,omitempty"`
	EndDate     *time.Time         `json:"end_date,omitempty"`
	TargetRoles []string           `json:"target_roles,omitempty"`
}

// GetActiveBanners godoc
// @Summary Получить активные баннеры
// @Tags Banner
// @Success 200 {array} domain.Banner
// @Router /api/banner/active [get]
func (h *BannerHandler) GetActiveBanners(w http.ResponseWriter, r *http.Request) {
	// Получаем роль пользователя из контекста (если есть)
	var role *string
	if userRole, ok := r.Context().Value("user_role").(string); ok {
		role = &userRole
	}

	banners, err := h.uc.GetActiveBanners(r.Context(), role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// CreateBanner godoc
// @Summary Создать баннер (только admin)
// @Tags Banner
// @Param body body CreateBannerRequest true "Banner data"
// @Success 200 {object} map[string]string
// @Router /api/admin/banner [post]
func (h *BannerHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	var req CreateBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(string)

	err := h.uc.CreateBanner(r.Context(), req.Title, req.Content, req.Type, req.IsActive, req.Priority, req.StartDate, req.EndDate, req.TargetRoles, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Banner created successfully"})
}

// UpdateBanner godoc
// @Summary Обновить баннер (только admin)
// @Tags Banner
// @Param bannerId path string true "Banner ID"
// @Param body body UpdateBannerRequest true "Banner data"
// @Success 200 {object} map[string]string
// @Router /api/admin/banner/{bannerId} [patch]
func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	bannerID := chi.URLParam(r, "bannerId")

	var req UpdateBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.uc.UpdateBanner(r.Context(), bannerID, req.Title, req.Content, req.Type, req.IsActive, req.Priority, req.StartDate, req.EndDate, req.TargetRoles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Banner updated successfully"})
}

// DeleteBanner godoc
// @Summary Удалить баннер (только admin)
// @Tags Banner
// @Param bannerId path string true "Banner ID"
// @Success 200 {object} map[string]string
// @Router /api/admin/banner/{bannerId} [delete]
func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerID := chi.URLParam(r, "bannerId")

	err := h.uc.DeleteBanner(r.Context(), bannerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Banner deleted successfully"})
}
