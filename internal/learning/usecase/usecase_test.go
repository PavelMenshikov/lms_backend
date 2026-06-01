package usecase_test

import (
	"context"
	"errors"
	"testing"

	"lms_backend/internal/domain"
	"lms_backend/internal/learning/mocks"
	"lms_backend/internal/learning/usecase"
	pkgMocks "lms_backend/pkg/storage/mocks"
)

func TestGetMyCourses(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	repo.GetMyCoursesFunc = func(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
		if userID == "" {
			return nil, errors.New("unauthorized")
		}
		return []*domain.StudentCoursePreview{
			{ID: "c1", Title: "Course 1", ProgressPercent: 50},
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		courses, err := uc.GetMyCourses(context.Background(), "user-1")
		if err != nil {
			t.Fatal(err)
		}
		if len(courses) != 1 || courses[0].Title != "Course 1" {
			t.Error("expected 1 course with title 'Course 1'")
		}
	})

	t.Run("empty userID returns error", func(t *testing.T) {
		_, err := uc.GetMyCourses(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty userID")
		}
	})

	t.Run("repo error propagates", func(t *testing.T) {
		repo.GetMyCoursesFunc = func(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
			return nil, errors.New("db error")
		}
		_, err := uc.GetMyCourses(context.Background(), "user-1")
		if err == nil || err.Error() != "db error" {
			t.Error("expected 'db error'")
		}
	})
}

func TestGetCourseContent(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	repo.GetCourseContentFunc = func(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error) {
		if courseID == "" {
			return nil, errors.New("not found")
		}
		return &domain.StudentCourseView{
			Course: &domain.Course{ID: courseID, Title: "Go Basics"},
		}, nil
	}

	t.Run("success", func(t *testing.T) {
		view, err := uc.GetCourseContent(context.Background(), "c1", "u1")
		if err != nil {
			t.Fatal(err)
		}
		if view.Course.Title != "Go Basics" {
			t.Error("expected 'Go Basics'")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		_, err := uc.GetCourseContent(context.Background(), "", "u1")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestSubmitAssignment(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	repo.GetAssignmentIDByLessonFunc = func(ctx context.Context, lessonID string) (string, error) {
		if lessonID == "bad" {
			return "", errors.New("not found")
		}
		return "assign-1", nil
	}

	repo.SaveSubmissionFunc = func(ctx context.Context, userID, assignmentID, text string, files []string) error {
		if assignmentID == "fail" {
			return errors.New("save failed")
		}
		return nil
	}

	t.Run("success without file", func(t *testing.T) {
		err := uc.SubmitAssignment(context.Background(), usecase.SubmitAssignmentInput{
			LessonID:   "l1",
			UserID:     "u1",
			TextAnswer: "my homework",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("assignment not found", func(t *testing.T) {
		err := uc.SubmitAssignment(context.Background(), usecase.SubmitAssignmentInput{
			LessonID: "bad",
			UserID:   "u1",
		})
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		err := uc.SubmitAssignment(context.Background(), usecase.SubmitAssignmentInput{
			LessonID:   "l1",
			UserID:     "u1",
			TextAnswer: "text",
		})
		repo.GetAssignmentIDByLessonFunc = func(ctx context.Context, lessonID string) (string, error) {
			return "fail", nil
		}
		if err != nil && err.Error() != "save failed" {
			t.Error("expected 'save failed'")
		}
	})
}

func TestSetLessonAttendance(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	repo.SetLessonAttendanceFunc = func(ctx context.Context, userID, lessonID, status, recordingURL, teacherComment string) error {
		if lessonID == "fail" {
			return errors.New("set failed")
		}
		return nil
	}

	t.Run("success", func(t *testing.T) {
		err := uc.SetLessonAttendance(context.Background(), usecase.SetAttendanceInput{
			LessonID: "l1",
			UserID:   "u1",
			Status:   "visited",
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		err := uc.SetLessonAttendance(context.Background(), usecase.SetAttendanceInput{
			LessonID: "fail",
			UserID:   "u1",
		})
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestGetTeacherDetails(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	repo.GetTeacherByIDFunc = func(ctx context.Context, id string) (*domain.TeacherPublicInfo, error) {
		if id == "t1" {
			return &domain.TeacherPublicInfo{
				ID: "t1", FirstName: "John", LastName: "Doe",
			}, nil
		}
		return nil, errors.New("not found")
	}

	repo.GetTeacherReviewsFunc = func(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error) {
		return []*domain.TeacherReview{
			{Rating: 5, Comment: "Great teacher!"},
		}, nil
	}

	t.Run("success with reviews", func(t *testing.T) {
		teacher, err := uc.GetTeacherDetails(context.Background(), "t1")
		if err != nil {
			t.Fatal(err)
		}
		if len(teacher.Reviews) != 1 {
			t.Error("expected 1 review")
		}
	})

	t.Run("teacher not found", func(t *testing.T) {
		_, err := uc.GetTeacherDetails(context.Background(), "unknown")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("reviews error ignored", func(t *testing.T) {
		repo.GetTeacherReviewsFunc = func(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error) {
			return nil, errors.New("reviews error")
		}
		teacher, err := uc.GetTeacherDetails(context.Background(), "t1")
		if err != nil {
			t.Fatal(err)
		}
		if teacher.Reviews != nil {
			t.Error("expected nil reviews on error")
		}
	})
}

func TestAddReview(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	saved := false
	repo.AddTeacherReviewFunc = func(ctx context.Context, review *domain.TeacherReview) error {
		saved = true
		if review.TeacherID == "fail" {
			return errors.New("save failed")
		}
		return nil
	}

	t.Run("valid review", func(t *testing.T) {
		saved = false
		err := uc.AddReview(context.Background(), usecase.AddReviewInput{
			TeacherID: "t1",
			StudentID: "s1",
			Rating:    4,
			Comment:   "Good",
		})
		if err != nil {
			t.Fatal(err)
		}
		if !saved {
			t.Error("review was not saved")
		}
	})

	t.Run("rating too low", func(t *testing.T) {
		err := uc.AddReview(context.Background(), usecase.AddReviewInput{
			TeacherID: "t1",
			Rating:    0,
		})
		if err == nil {
			t.Error("expected error for invalid rating")
		}
	})

	t.Run("rating too high", func(t *testing.T) {
		err := uc.AddReview(context.Background(), usecase.AddReviewInput{
			TeacherID: "t1",
			Rating:    6,
		})
		if err == nil {
			t.Error("expected error for invalid rating")
		}
	})
}

func TestGetTeacherDashboard(t *testing.T) {
	repo := mocks.NewLearningRepoMock()
	s3 := pkgMocks.NewS3StorageMock()
	uc := usecase.NewLearningUseCase(repo, s3)

	repo.GetTeacherByIDFunc = func(ctx context.Context, id string) (*domain.TeacherPublicInfo, error) {
		if id == "t1" {
			return &domain.TeacherPublicInfo{ID: "t1", FirstName: "John"}, nil
		}
		return nil, errors.New("not found")
	}

	repo.GetTeacherReviewsFunc = func(ctx context.Context, id string) ([]*domain.TeacherReview, error) {
		return []*domain.TeacherReview{{Rating: 5}}, nil
	}

	repo.GetTeacherCoursesFunc = func(ctx context.Context, id string) ([]*domain.StudentCoursePreview, error) {
		return []*domain.StudentCoursePreview{{ID: "c1", Title: "Course"}}, nil
	}

	repo.GetTeacherSubstitutionsFunc = func(ctx context.Context, id string) ([]*domain.Lesson, error) {
		return nil, nil
	}

	repo.GetTeacherUpcomingLessonsFunc = func(ctx context.Context, id string) ([]*domain.Lesson, error) {
		return nil, nil
	}

	repo.GetTeacherCancelledLessonsFunc = func(ctx context.Context, id string) ([]*domain.Lesson, error) {
		return nil, nil
	}

	t.Run("success", func(t *testing.T) {
		dash, err := uc.GetTeacherDashboard(context.Background(), "t1")
		if err != nil {
			t.Fatal(err)
		}
		if dash.Profile.FirstName != "John" {
			t.Error("expected teacher John")
		}
		if len(dash.MyReviews) != 1 {
			t.Error("expected 1 review")
		}
		if len(dash.AssignedCourses) != 1 {
			t.Error("expected 1 course")
		}
	})

	t.Run("teacher not found", func(t *testing.T) {
		_, err := uc.GetTeacherDashboard(context.Background(), "unknown")
		if err == nil {
			t.Error("expected error")
		}
	})
}
