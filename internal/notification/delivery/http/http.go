package http

import (
	"encoding/json"
	"lms_backend/internal/domain"
	"lms_backend/internal/notification/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type NotificationHandler struct {
	uc usecase.NotificationUseCase
}

func NewNotificationHandler(uc usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

type CreateNotificationRequest struct {
	RecipientID string                   `json:"recipient_id"`
	Title       string                   `json:"title"`
	Content     string                   `json:"content"`
	Type        domain.NotificationType  `json:"type"`
	LinkURL     *string                  `json:"link_url,omitempty"`
}

// CreateNotification godoc
// @Summary Создать уведомление
// @Tags Notifications
// @Param body body CreateNotificationRequest true "Notification data"
// @Success 200 {object} map[string]string
// @Router /api/notifications [post]
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(string)

	err := h.uc.CreateNotification(r.Context(), req.RecipientID, &userID, req.Title, req.Content, req.Type, req.LinkURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification created successfully"})
}

// GetNotifications godoc
// @Summary Получить уведомления пользователя
// @Tags Notifications
// @Success 200 {array} domain.Notification
// @Router /api/notifications [get]
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	notifications, err := h.uc.GetUserNotifications(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// MarkNotificationAsRead godoc
// @Summary Отметить уведомление как прочитанное
// @Tags Notifications
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Router /api/notifications/{notificationId}/read [patch]
func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	notificationID := chi.URLParam(r, "notificationId")

	err := h.uc.MarkAsRead(r.Context(), notificationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification marked as read"})
}
