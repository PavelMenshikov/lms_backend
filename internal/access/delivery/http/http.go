package http

import (
	"encoding/json"
	"lms_backend/internal/access/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
	authMiddleware "lms_backend/internal/auth/delivery/middleware"
)

type AccessHandler struct {
	uc usecase.AccessUseCase
}

func NewAccessHandler(uc usecase.AccessUseCase) *AccessHandler {
	return &AccessHandler{uc: uc}
}

type CreateAccessRequestReq struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	Reason       string `json:"reason"`
}

type ReviewAccessRequestReq struct {
	ReviewComment *string `json:"review_comment,omitempty"`
}

// CreateAccessRequest godoc
// @Summary Создать запрос доступа
// @Tags Access
// @Param body body CreateAccessRequestReq true "Access request data"
// @Success 200 {object} map[string]string
// @Router /api/access-requests [post]
func (h *AccessHandler) CreateAccessRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateAccessRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	userID := userCtxData.UserID

	err := h.uc.CreateRequest(r.Context(), userID, req.ResourceType, req.ResourceID, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Access request created successfully"})
}

// GetPendingRequests godoc
// @Summary Получить список ожидающих запросов доступа
// @Tags Access
// @Success 200 {array} domain.AccessRequest
// @Router /api/access-requests [get]
func (h *AccessHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.uc.GetPendingRequests(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// ApproveRequest godoc
// @Summary Одобрить запрос доступа
// @Tags Access
// @Param requestId path string true "Request ID"
// @Param body body ReviewAccessRequestReq true "Review data"
// @Success 200 {object} map[string]string
// @Router /api/access-requests/{requestId}/approve [patch]
func (h *AccessHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "requestId")

	var req ReviewAccessRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	userID := userCtxData.UserID

	err := h.uc.ApproveRequest(r.Context(), requestID, userID, req.ReviewComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Access request approved"})
}

// RejectRequest godoc
// @Summary Отклонить запрос доступа
// @Tags Access
// @Param requestId path string true "Request ID"
// @Param body body ReviewAccessRequestReq true "Review data"
// @Success 200 {object} map[string]string
// @Router /api/access-requests/{requestId}/reject [patch]
func (h *AccessHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "requestId")

	var req ReviewAccessRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	userID := userCtxData.UserID

	err := h.uc.RejectRequest(r.Context(), requestID, userID, req.ReviewComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Access request rejected"})
}
