package usecase_test

import (
	"context"
	"testing"
	"time"

	"lms_backend/internal/content_admin/mocks"
	"lms_backend/internal/content_admin/usecase"
	"lms_backend/internal/domain"
	s3Mocks "lms_backend/pkg/storage/mocks"
)

func TestCreateCourse(t *testing.T) {
	repoMock := mocks.NewContentAdminRepoMock()
	s3Mock := s3Mocks.NewS3StorageMock()
	uc := usecase.NewContentAdminUseCase(repoMock, s3Mock)

	ctx := context.Background()

	input := usecase.CreateCourseInput{
		Title:       "Golang Basic",
		Description: "Learn Go",
		IsMain:      true,
	}

	id, err := uc.CreateCourse(ctx, input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id == "" {
		t.Error("Expected ID, got empty string")
	}

	if len(repoMock.CreatedCourses) != 1 {
		t.Error("Course was not saved to repository")
	}
}

func TestCreateFullUser_StudentWithParent(t *testing.T) {
	repoMock := mocks.NewContentAdminRepoMock()
	s3Mock := s3Mocks.NewS3StorageMock()
	uc := usecase.NewContentAdminUseCase(repoMock, s3Mock)

	ctx := context.Background()

	input := usecase.ExtendedCreateUserInput{
		FirstName:       "Маленький",
		LastName:        "Бобби",
		Email:           "bobby@school.com",
		Password:        "12345",
		Role:            domain.RoleStudent,
		Phone:           "+79990000000",
		BirthDate:       time.Now(),
		ParentFirstName: "Мама",
		ParentLastName:  "Роберта",
		ParentPhone:     "+78881112233",
		ParentEmail:     "mom@gmail.com",
	}

	result, err := uc.CreateFullUser(ctx, input)
	if err != nil {
		t.Fatalf("CreateFullUser failed: %v", err)
	}

	studentID := result["user_id"]
	parentID := result["parent_id"]

	if studentID == "" {
		t.Error("Student ID missing")
	}
	if parentID == "" {
		t.Error("Parent ID missing")
	}

	if len(repoMock.CreatedUsers) != 2 {
		t.Errorf("Expected 2 users created, got %d", len(repoMock.CreatedUsers))
	}

	linkedParent, exists := repoMock.LinkedParents[studentID]
	if !exists {
		t.Error("Link between Student and Parent NOT created")
	}
	if linkedParent != parentID {
		t.Error("Linked wrong parent ID")
	}
}
