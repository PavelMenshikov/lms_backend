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
	GetAdminCounters(ctx context.Context) (totalStudents, newStudents int, studentsDelta float64, totalTeachers, activeCourses int, err error)
	GetAllPerformanceStats(ctx context.Context) (*domain.AllPerformanceStats, error)
	GetPerformanceStats(ctx context.Context) (domain.PerformanceZones, error)
	GetHwPerformanceStats(ctx context.Context) (domain.PerformanceZones, error)
	GetAttendancePerformanceStats(ctx context.Context) (domain.PerformanceZones, error)
	GetLessonActivity(ctx context.Context) ([]domain.DailyLessonActivity, error)
	GetCuratorGroups(ctx context.Context, curatorID string) ([]domain.Group, error)
	GetCuratorAttendanceStats(ctx context.Context, curatorID string) ([]domain.CuratorGroupAttendance, error)
	GetCuratorHomeworkStats(ctx context.Context, curatorID string) ([]domain.CuratorHomeworkStats, error)
	GetCuratorPerformanceZones(ctx context.Context, curatorID string) (domain.PerformanceZones, error)
}
