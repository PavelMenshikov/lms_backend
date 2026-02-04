package http

import (
	"encoding/json"
	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/profile/usecase"
	"mime/multipart"
	"net/http"
)

type ProfileHandler struct {
	uc *usecase.ProfileUseCase
}

func NewProfileHandler(uc *usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

// GetProfile godoc
// @Summary USER: Мой профиль
// @Description Получить данные текущего авторизованного пользователя.
// @Tags Profile
// @Produce json
// @Success 200 {object} domain.User
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)
	profile, err := h.uc.GetMyProfile(r.Context(), userData.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile godoc
// @Summary USER: Обновить профиль
// @Description Изменение личных данных и загрузка аватара.
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param first_name formData string false "Имя"
// @Param last_name formData string false "Фамилия"
// @Param phone formData string false "Телефон"
// @Param avatar formData file false "Фото профиля"
// @Success 200 {object} map[string]string
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	const MAX_SIZE = 5 << 20
	r.ParseMultipartForm(MAX_SIZE)

	userData := r.Context().Value(authMiddleware.ContextUserDataKey).(*authMiddleware.UserContextData)

	var fileHeader *multipart.FileHeader
	if f, head, err := r.FormFile("avatar"); err == nil {
		f.Close()
		fileHeader = head
	}

	input := usecase.UpdateProfileInput{
		UserID:     userData.UserID,
		FirstName:  r.FormValue("first_name"),
		LastName:   r.FormValue("last_name"),
		Phone:      r.FormValue("phone"),
		City:       r.FormValue("city"),
		Language:   r.FormValue("language"),
		School:     r.FormValue("school_name"),
		Whatsapp:   r.FormValue("whatsapp"),
		Telegram:   r.FormValue("telegram"),
		FileHeader: fileHeader,
	}

	if err := h.uc.UpdateProfile(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "profile updated"})
}
