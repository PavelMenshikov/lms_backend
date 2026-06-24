package usecase

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"

	"lms_backend/internal/domain"
	"lms_backend/internal/profile/repository"
	storageService "lms_backend/pkg/storage"
)

type ProfileUseCase struct {
	repo      repository.ProfileRepository
	s3Storage storageService.ObjectStorage
}

func NewProfileUseCase(repo repository.ProfileRepository, s3Storage storageService.ObjectStorage) *ProfileUseCase {
	return &ProfileUseCase{repo: repo, s3Storage: s3Storage}
}

func (uc *ProfileUseCase) GetMyProfile(ctx context.Context, userID string) (*domain.User, error) {
	return uc.repo.GetProfile(ctx, userID)
}

type UpdateProfileInput struct {
	UserID     string
	FirstName  string
	LastName   string
	Phone      string
	City       string
	Language   string
	School     string
	Whatsapp   string
	Telegram   string
	FileHeader *multipart.FileHeader
}

func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, input UpdateProfileInput) error {
	user, err := uc.repo.GetProfile(ctx, input.UserID)
	if err != nil {
		return err
	}

	if input.FileHeader != nil {
		file, err := input.FileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		s3Key := fmt.Sprintf("avatars/%s_%s", input.UserID, input.FileHeader.Filename)
		mimeType := input.FileHeader.Header.Get("Content-Type")
		if mimeType == "" {
			ext := filepath.Ext(input.FileHeader.Filename)
			mimeType = mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}
		key, err := uc.s3Storage.UploadFile(ctx, file, s3Key, input.FileHeader.Size, mimeType)
		if err != nil {
			return err
		}
		user.AvatarURL, _ = uc.s3Storage.GetPublicURL(ctx, key)
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Phone = input.Phone
	user.City = input.City
	user.Language = input.Language
	user.SchoolName = input.School
	user.Whatsapp = input.Whatsapp
	user.Telegram = input.Telegram

	return uc.repo.UpdateProfile(ctx, user)
}
func (uc *ProfileUseCase) UpdateTeacherSchedule(ctx context.Context, userID string, scheduleJSON []byte) error {
	return uc.repo.UpdateTeacherSchedule(ctx, userID, scheduleJSON)
}
