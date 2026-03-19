package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/domain"
	"lms_backend/internal/review/usecase"
)

type ReviewHandler struct {
	uc *usecase.ReviewUseCase
}

func NewReviewHandler(uc *usecase.ReviewUseCase) *ReviewHandler {
	return &ReviewHandler{uc: uc}
}

type EvaluateRequest struct {
	Grade      int    `json:"grade"`
	Comment    string `json:"comment"`
	IsAccepted bool   `json:"is_accepted"`
}

// GetPendingSubmissions godoc
// @Summary STAFF: Список ДЗ на проверку
// @Tags Staff-Review
// @Produce json
// @Success 200 {array} domain.SubmissionRecord
// @Router /staff/submissions [get]
func (h *ReviewHandler) GetPendingSubmissions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)

	list, err := h.uc.GetPendingList(r.Context(), userCtx.UserID, string(userCtx.Role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// EvaluateSubmission godoc
// @Summary STAFF: Проверить ДЗ
// @Tags Staff-Review
// @Accept json
// @Produce json
// @Param id path string true "ID задания"
// @Param request body EvaluateRequest true "Результат проверки"
// @Success 200 {object} map[string]string
// @Router /staff/submissions/{id}/evaluate [post]
func (h *ReviewHandler) EvaluateSubmission(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	
	if userCtx.Role == domain.RoleCurator {
		http.Error(w, "Forbidden: Curators cannot evaluate homework", http.StatusForbidden)
		return
	}

	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	submissionID := chi.URLParam(r, "id")

	status := "on_revision"
	if req.IsAccepted {
		status = "accepted"
	}

	input := usecase.EvaluateInput{
		SubmissionID: submissionID,
		Grade:        req.Grade,
		Comment:      req.Comment,
		Status:       status,
	}

	if err := h.uc.Evaluate(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "evaluated"})
}