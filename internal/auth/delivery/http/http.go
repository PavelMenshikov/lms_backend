package http

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"lms_backend/internal/auth/usecase"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// РЕГИСТРАЦИЯ ОТМЕНЕНА: оставляем роут для старых Swagger-клиентов, но он выдаст ошибку 403.
// Register godoc
// @Summary Регистрация (ЗАБЛОКИРОВАНО)
// @Description Данный функционал отключен, пользователи создаются только Администратором.
// @Tags Аутентификация
// @Produce  json
// @Failure 403 {object} map[string]string "error: Public registration is disabled."
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("ATTEMPT TO USE DEPRECATED /auth/register ROUTE. ACCESS DENIED.")
	http.Error(w, "Public registration is disabled. Users must be created by Admin.", http.StatusForbidden)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Вход в систему (Log in student/admin)
// @Description Ввод email и пароля, возвращает HTTP-Only Cookie.
// @Tags Аутентификация
// @Accept  json
// @Produce  json
// @Param   request body LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "message: Login successful, user: {...}"
// @Failure 401 {object} map[string]string "error: Неправильный логин/пароль"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token := "mock-jwt-token-for-user-" + user.Email

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	response := map[string]interface{}{
		"message": "Login successful. Token stored in HTTP-only cookie.",
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
