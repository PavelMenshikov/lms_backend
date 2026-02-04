package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"lms_backend/internal/domain"
	"lms_backend/internal/learning/repository"
	storageService "lms_backend/pkg/storage"
)

type LearningUseCase struct {
	repo      repository.LearningRepository
	s3Storage storageService.ObjectStorage
}

func NewLearningUseCase(repo repository.LearningRepository, s3Storage storageService.ObjectStorage) *LearningUseCase {
	return &LearningUseCase{repo: repo, s3Storage: s3Storage}
}

func (uc *LearningUseCase) GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	return uc.repo.GetMyCourses(ctx, userID)
}

func (uc *LearningUseCase) GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error) {
	return uc.repo.GetCourseContent(ctx, courseID, userID)
}

func (uc *LearningUseCase) GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error) {
	return uc.repo.GetLessonDetail(ctx, lessonID, userID)
}

type SubmitAssignmentInput struct {
	LessonID   string
	UserID     string
	TextAnswer string
	FileHeader *multipart.FileHeader
}

func (uc *LearningUseCase) SubmitAssignment(ctx context.Context, input SubmitAssignmentInput) error {
	assignmentID, err := uc.repo.GetAssignmentIDByLesson(ctx, input.LessonID)
	if err != nil {
		return fmt.Errorf("assignment not found for this lesson: %w", err)
	}

	var fileURL string
	if input.FileHeader != nil {
		file, err := input.FileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		s3Key := fmt.Sprintf("submissions/%s_%s_%s", input.UserID, assignmentID, input.FileHeader.Filename)
		mimeType := input.FileHeader.Header.Get("Content-Type")

		key, err := uc.s3Storage.UploadFile(ctx, file, s3Key, input.FileHeader.Size, mimeType)
		if err != nil {
			return err
		}
		fileURL, _ = uc.s3Storage.GetPublicURL(ctx, key)
	}

	return uc.repo.SaveSubmission(ctx, input.UserID, assignmentID, input.TextAnswer, fileURL)
}

func (uc *LearningUseCase) CompleteLesson(ctx context.Context, lessonID, userID string) error {
	return uc.repo.MarkLessonComplete(ctx, userID, lessonID)
}
