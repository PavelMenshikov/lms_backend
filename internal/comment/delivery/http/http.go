package http

import (
	"encoding/json"
	"lms_backend/internal/comment/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	uc usecase.CommentUseCase
}

func NewCommentHandler(uc usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{uc: uc}
}

type CreateCommentRequest struct {
	StudentID       string  `json:"student_id"`
	LessonID        *string `json:"lesson_id,omitempty"`
	RecipientID     *string `json:"recipient_id,omitempty"`
	Content         string  `json:"content"`
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
}

// CreateComment godoc
// @Summary Создать комментарий
// @Tags Comments
// @Param body body CreateCommentRequest true "Comment data"
// @Success 200 {object} map[string]string
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(string)

	err := h.uc.CreateComment(r.Context(), req.StudentID, req.LessonID, userID, req.RecipientID, req.Content, req.ParentCommentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment created successfully"})
}

// GetComments godoc
// @Summary Получить комментарии по ученику
// @Tags Comments
// @Param studentId query string true "Student ID"
// @Success 200 {array} domain.Comment
// @Router /api/comments [get]
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Query().Get("studentId")
	if studentID == "" {
		http.Error(w, "studentId is required", http.StatusBadRequest)
		return
	}

	comments, err := h.uc.GetStudentComments(r.Context(), studentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// MarkCommentAsRead godoc
// @Summary Отметить комментарий как прочитанный
// @Tags Comments
// @Param commentId path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Router /api/comments/{commentId}/read [patch]
func (h *CommentHandler) MarkCommentAsRead(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "commentId")

	err := h.uc.MarkAsRead(r.Context(), commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment marked as read"})
}
