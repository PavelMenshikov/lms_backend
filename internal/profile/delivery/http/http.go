package http

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	authMiddleware "lms_backend/internal/auth/delivery/middleware"
	"lms_backend/internal/profile/usecase"
)

type ProfileHandler struct {
	uc *usecase.ProfileUseCase
}

func NewProfileHandler(uc *usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

type UpdateProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	Language   string `json:"language"`
	SchoolName string `json:"school_name"`
	Whatsapp   string `json:"whatsapp"`
	Telegram   string `json:"telegram"`
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
// @Param avatar formData file false "Фото профиля"
// @Param first_name formData string false "Имя"
// @Param last_name formData string false "Фамилия"
// @Param phone formData string false "Телефон"
// @Param city formData string false "Населенный пункт"
// @Param language formData string false "Родной язык"
// @Param school_name formData string false "Учебное заведение"
// @Param whatsapp formData string false "WhatsApp ссылка"
// @Param telegram formData string false "Telegram ссылка"
// @Success 200 {object} map[string]string "status"
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	const MAX_SIZE = 5 << 20
	if err := r.ParseMultipartForm(MAX_SIZE); err != nil {
		http.Error(w, "File upload size exceeded limit.", http.StatusBadRequest)
		return
	}

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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "profile updated"})
}
