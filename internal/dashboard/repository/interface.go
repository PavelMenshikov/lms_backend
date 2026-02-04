package repository

import (
	"context"
	"lms_backend/internal/domain"
)

type UserDataRepository interface {
	GetLastLessonData(ctx context.Context, userID string) (*domain.LastLesson, error)
	GetActiveCoursesCount(ctx context.Context, userID string) (int, error)
	GetAttendancePercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error)
	GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error)
	GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error)
	GetAdminCounters(ctx context.Context) (totalStudents, newStudents, totalTeachers, activeCourses int, err error)
	GetPerformanceStats(ctx context.Context) (domain.PerformanceZones, error)
	GetLessonActivity(ctx context.Context) ([]domain.DailyLessonActivity, error)
}
