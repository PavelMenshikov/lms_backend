package http

import (
	"encoding/json"
	"net/http"
	"time"

	"lms_backend/internal/auth/usecase"
	"lms_backend/internal/domain"
)

type AuthHandler struct {
	uc *usecase.AuthUsecase
}

func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type RegisterRequest struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Role      domain.Role `json:"role"`
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает аккаунт ученика или родителя
// @Tags Аутентификация
// @Accept  json
// @Produce  json
// @Param   request body RegisterRequest true "Данные для регистрации"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string "error: Ошибка валидации"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role:      req.Role,
	}

	if err := h.uc.Register(r.Context(), user, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Вход в систему (Log in student)
// @Description Ввод email и пароля, возвращает JWT токен
// @Tags Аутентификация
// @Accept  json
// @Produce  json
// @Param   request body LoginRequest true "Данные для входа"
// @Success 200 {object} map[string]string "token: JWT-токен"
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
		MaxAge:   int(24 * 7 * time.Hour / time.Second),
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
