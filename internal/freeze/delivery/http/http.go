package http

import (
	"encoding/json"
	"lms_backend/internal/freeze/usecase"
	"lms_backend/internal/httperror"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	authMiddleware "lms_backend/internal/auth/delivery/middleware"
)

type FreezeHandler struct {
	uc usecase.FreezeUseCase
}

func NewFreezeHandler(uc usecase.FreezeUseCase) *FreezeHandler {
	return &FreezeHandler{uc: uc}
}

type CreateFreezeRequestReq struct {
	StudentID string `json:"student_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string `json:"reason"`
}

type ReviewFreezeRequestReq struct {
	ReviewComment *string `json:"review_comment,omitempty"`
}

// CreateFreezeRequest godoc
// @Summary Создать запрос на заморозку
// @Tags Freeze
// @Param body body CreateFreezeRequestReq true "Freeze request data"
// @Success 200 {object} map[string]string
// @Router /api/freeze-requests [post]
func (h *FreezeHandler) CreateFreezeRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateFreezeRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format", http.StatusBadRequest)
		return
	}

	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userCtxData.UserID

	err = h.uc.CreateRequest(r.Context(), req.StudentID, userID, startDate, endDate, req.Reason)
	if err != nil {
		httperror.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Freeze request created successfully"})
}

// GetPendingRequests godoc
// @Summary Получить список ожидающих запросов на заморозку
// @Tags Freeze
// @Success 200 {array} domain.FreezeRequest
// @Router /api/freeze-requests [get]
func (h *FreezeHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.uc.GetPendingRequests(r.Context())
	if err != nil {
		httperror.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// ApproveRequest godoc
// @Summary Одобрить запрос на заморозку
// @Tags Freeze
// @Param requestId path string true "Request ID"
// @Param body body ReviewFreezeRequestReq true "Review data"
// @Success 200 {object} map[string]string
// @Router /api/freeze-requests/{requestId}/approve [patch]
func (h *FreezeHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "requestId")

	var req ReviewFreezeRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userCtxData.UserID

	err := h.uc.ApproveRequest(r.Context(), requestID, userID, req.ReviewComment)
	if err != nil {
		httperror.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Freeze request approved"})
}

// RejectRequest godoc
// @Summary Отклонить запрос на заморозку
// @Tags Freeze
// @Param requestId path string true "Request ID"
// @Param body body ReviewFreezeRequestReq true "Review data"
// @Success 200 {object} map[string]string
// @Router /api/freeze-requests/{requestId}/reject [patch]
func (h *FreezeHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	requestID := chi.URLParam(r, "requestId")

	var req ReviewFreezeRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok || userCtxData == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userCtxData.UserID

	err := h.uc.RejectRequest(r.Context(), requestID, userID, req.ReviewComment)
	if err != nil {
		httperror.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Freeze request rejected"})
}

// GetStudentFreezeStatus godoc
// @Summary Получить статус заморозки ученика
// @Tags Freeze
// @Param studentId path string true "Student ID"
// @Success 200 {object} domain.FreezePeriod
// @Router /api/students/{studentId}/freeze-status [get]
func (h *FreezeHandler) GetStudentFreezeStatus(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "studentId")

	status, err := h.uc.GetStudentFreezeStatus(r.Context(), studentID)
	if err != nil {
		httperror.Internal(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
