package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"lms_backend/internal/auth/usecase"
	"lms_backend/internal/httperror"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Register godoc
// @Summary Регистрация (ЗАБЛОКИРОВАНО)
// @Description Данный функционал отключен, пользователи создаются только Администратором.
// @Tags Аутентификация
// @Produce  json
// @Failure 403 {object} map[string]string "error: Public registration is disabled."
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	slog.Info("ATTEMPT TO USE DEPRECATED /auth/register ROUTE. ACCESS DENIED.")
	http.Error(w, "Public registration is disabled. Users must be created by Admin.", http.StatusForbidden)
}

type LoginRequest struct {
	Email    string `json:"email" example:"admin@capedu.kz"`
	Password string `json:"password" example:"capedu123"`
}

// Login godoc
// @Summary Вход в систему
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
		httperror.BadRequest(w, err)
		return
	}

	user, err := h.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httperror.Unauthorized(w, err)
		return
	}

	token := user.ID + ":" + string(user.Role)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		MaxAge:   int(24 * 7 * time.Hour / time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	// -------------------------------------------------------

	response := map[string]interface{}{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout godoc
// @Summary Выход из системы
// @Description Удаляет HTTP-Only Cookie.
// @Tags Аутентификация
// @Produce json
// @Success 200 {object} map[string]string "message: Logged out successfully"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

type ForgotPasswordRequest struct {
	Email string `json:"email" example:"user@capedu.kz"`
}

// ForgotPassword godoc
// @Summary Запрос сброса пароля
// @Description Отправляет email для сброса пароля (заглушка).
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Заглушка: в реальном проекте здесь отправка email
	slog.Info("password reset requested", slog.String("email", req.Email))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If the email exists, a password reset link has been sent.",
	})
}

type ResetPasswordRequest struct {
	Token    string `json:"token" example:"reset-token"`
	Password string `json:"password" example:"newpassword123"`
}

// ResetPassword godoc
// @Summary Сброс пароля
// @Description Сбрасывает пароль по токену (заглушка).
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Токен и новый пароль"
// @Success 200 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	if req.Token == "" || req.Password == "" {
		http.Error(w, "Token and password are required", http.StatusBadRequest)
		return
	}

	// Заглушка: в реальном проекте здесь проверка токена и смена пароля
	slog.Info("password reset with token", slog.String("token", req.Token))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been reset successfully.",
	})
}
