package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/chat/usecase"
	"lms_backend/internal/domain"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	uc *usecase.ChatUseCase
}

func NewChatHandler(uc *usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{uc: uc}
}

// ConnectToChat godoc
// @Summary WS: Подключение к чату модуля
// @Description WebSocket соединение для обмена сообщениями.
// @Tags Chat
// @Param module_id query string true "ID модуля"
// @Param student_id query string true "ID ученика (для персонала)"
// @Router /chat/ws [get]
func (h *ChatHandler) ConnectToChat(w http.ResponseWriter, r *http.Request) {
	userCtxData, ok := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	moduleID := r.URL.Query().Get("module_id")
	studentID := r.URL.Query().Get("student_id")

	if userCtxData.Role == domain.RoleStudent {
		studentID = userCtxData.UserID
	}

	if moduleID == "" || studentID == "" {
		http.Error(w, "Missing params", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	roomID := moduleID + "_" + studentID
	messageChan := make(chan *domain.ChatMessage)

	h.uc.RegisterClient(roomID, messageChan)
	defer h.uc.UnregisterClient(roomID, messageChan)

	go func() {
		for msg := range messageChan {
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}()

	for {
		var incoming struct {
			Text    string `json:"text"`
			FileURL string `json:"file_url"`
		}
		if err := conn.ReadJSON(&incoming); err != nil {
			break
		}

		msg := &domain.ChatMessage{
			ModuleID:    moduleID,
			StudentID:   studentID,
			SenderID:    userCtxData.UserID,
			MessageText: incoming.Text,
			FileURL:     incoming.FileURL,
		}

		if err := h.uc.SendMessage(r.Context(), msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

// GetChatHistory godoc
// @Summary USER: История чата
// @Description Получить список сообщений чата модуля.
// @Tags Chat
// @Produce json
// @Param module_id query string true "ID модуля"
// @Param student_id query string false "ID ученика (для персонала)"
// @Param limit query int false "Количество"
// @Param offset query int false "Смещение"
// @Success 200 {array} domain.ChatMessage
// @Router /chat/history [get]
func (h *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	userCtxData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	moduleID := r.URL.Query().Get("module_id")
	studentID := r.URL.Query().Get("student_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if userCtxData.Role == domain.RoleStudent {
		studentID = userCtxData.UserID
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit == 0 {
		limit = 50
	}

	history, err := h.uc.GetHistory(r.Context(), moduleID, studentID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
