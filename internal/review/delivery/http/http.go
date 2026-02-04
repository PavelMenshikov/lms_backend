package http

import (
	"encoding/json"
	"lms_backend/internal/review/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
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
// @Description Получить список всех работ со статусом pending_check.
// @Tags Staff-Review
// @Produce json
// @Success 200 {array} domain.SubmissionRecord
// @Router /staff/submissions [get]
func (h *ReviewHandler) GetPendingSubmissions(w http.ResponseWriter, r *http.Request) {
	list, err := h.uc.GetPendingList(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// EvaluateSubmission godoc
// @Summary STAFF: Проверить ДЗ
// @Description Выставить оценку, статус и комментарий к работе ученика.
// @Tags Staff-Review
// @Accept json
// @Produce json
// @Param id path string true "ID задания"
// @Param request body EvaluateRequest true "Результат проверки"
// @Success 200 {object} map[string]string
// @Router /staff/submissions/{id}/evaluate [post]
func (h *ReviewHandler) EvaluateSubmission(w http.ResponseWriter, r *http.Request) {
	submissionID := chi.URLParam(r, "id")
	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := usecase.EvaluateInput{
		SubmissionID: submissionID,
		Grade:        req.Grade,
		Comment:      req.Comment,
		IsAccepted:   req.IsAccepted,
	}

	if err := h.uc.Evaluate(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "evaluated"})
}
