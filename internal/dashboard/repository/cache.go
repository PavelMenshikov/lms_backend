package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"lms_backend/internal/domain"
)

const dashboardCacheTTL = 5 * time.Minute

type CachedDashboardRepo struct {
	next UserDataRepository
	rdb  *redis.Client
}

func NewCachedDashboardRepo(next UserDataRepository, rdb *redis.Client) *CachedDashboardRepo {
	return &CachedDashboardRepo{next: next, rdb: rdb}
}

func cacheKey(format string, args ...interface{}) string {
	return fmt.Sprintf("dash:"+format, args...)
}

func cacheGet(ctx context.Context, rdb *redis.Client, key string, dest interface{}) (bool, error) {
	data, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

func cacheSet(ctx context.Context, rdb *redis.Client, key string, val interface{}) {
	data, err := json.Marshal(val)
	if err != nil {
		return
	}
	rdb.Set(ctx, key, data, dashboardCacheTTL)
}

func (c *CachedDashboardRepo) GetLastLessonData(ctx context.Context, userID string) (*domain.LastLesson, error) {
	key := cacheKey("last_lesson:%s", userID)
	var v domain.LastLesson
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return &v, nil
	}
	vp, err := c.next.GetLastLessonData(ctx, userID)
	if err == nil && vp != nil {
		cacheSet(ctx, c.rdb, key, vp)
	}
	return vp, err
}

func (c *CachedDashboardRepo) GetActiveCoursesCount(ctx context.Context, userID string) (int, error) {
	key := cacheKey("user:%s:active_courses", userID)
	var v int
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetActiveCoursesCount(ctx, userID)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetAttendancePercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	key := cacheKey("user:%s:attendance", userID)
	var v domain.StatisticSummary
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return &v, nil
	}
	vp, err := c.next.GetAttendancePercentage(ctx, userID)
	if err == nil && vp != nil {
		cacheSet(ctx, c.rdb, key, vp)
	}
	return vp, err
}

func (c *CachedDashboardRepo) GetAssignmentsCompletionPercentage(ctx context.Context, userID string) (*domain.StatisticSummary, error) {
	key := cacheKey("user:%s:assignments", userID)
	var v domain.StatisticSummary
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return &v, nil
	}
	vp, err := c.next.GetAssignmentsCompletionPercentage(ctx, userID)
	if err == nil && vp != nil {
		cacheSet(ctx, c.rdb, key, vp)
	}
	return vp, err
}

func (c *CachedDashboardRepo) GetUpcomingLessons(ctx context.Context, userID string) ([]domain.UpcomingLesson, error) {
	key := cacheKey("user:%s:upcoming", userID)
	var v []domain.UpcomingLesson
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetUpcomingLessons(ctx, userID)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetAdminCounters(ctx context.Context) (totalStudents, newStudents int, studentsDelta float64, totalTeachers, activeCourses int, err error) {
	type counters struct {
		TotalStudents int
		NewStudents   int
		StudentsDelta float64
		TotalTeachers int
		ActiveCourses int
	}

	key := cacheKey("admin:counters")
	var ct counters
	ok, cacheErr := cacheGet(ctx, c.rdb, key, &ct)
	if cacheErr == nil && ok {
		return ct.TotalStudents, ct.NewStudents, ct.StudentsDelta, ct.TotalTeachers, ct.ActiveCourses, nil
	}

	ts, ns, sd, tt, ac, err := c.next.GetAdminCounters(ctx)
	if err == nil {
		cacheSet(ctx, c.rdb, key, counters{ts, ns, sd, tt, ac})
	}
	return ts, ns, sd, tt, ac, err
}

func (c *CachedDashboardRepo) GetAllPerformanceStats(ctx context.Context) (*domain.AllPerformanceStats, error) {
	key := cacheKey("admin:all_perf")
	var v domain.AllPerformanceStats
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return &v, nil
	}
	vp, err := c.next.GetAllPerformanceStats(ctx)
	if err == nil && vp != nil {
		cacheSet(ctx, c.rdb, key, vp)
	}
	return vp, err
}

func (c *CachedDashboardRepo) GetPerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	key := cacheKey("admin:perf_zones")
	var v domain.PerformanceZones
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetPerformanceStats(ctx)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetHwPerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	key := cacheKey("admin:hw_zones")
	var v domain.PerformanceZones
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetHwPerformanceStats(ctx)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetAttendancePerformanceStats(ctx context.Context) (domain.PerformanceZones, error) {
	key := cacheKey("admin:att_zones")
	var v domain.PerformanceZones
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetAttendancePerformanceStats(ctx)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetLessonActivity(ctx context.Context) ([]domain.DailyLessonActivity, error) {
	key := cacheKey("admin:lesson_activity")
	var v []domain.DailyLessonActivity
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetLessonActivity(ctx)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetCuratorGroups(ctx context.Context, curatorID string) ([]domain.Group, error) {
	key := cacheKey("curator:%s:groups", curatorID)
	var v []domain.Group
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetCuratorGroups(ctx, curatorID)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetCuratorAttendanceStats(ctx context.Context, curatorID string) ([]domain.CuratorGroupAttendance, error) {
	key := cacheKey("curator:%s:attendance", curatorID)
	var v []domain.CuratorGroupAttendance
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetCuratorAttendanceStats(ctx, curatorID)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetCuratorHomeworkStats(ctx context.Context, curatorID string) ([]domain.CuratorHomeworkStats, error) {
	key := cacheKey("curator:%s:homework", curatorID)
	var v []domain.CuratorHomeworkStats
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetCuratorHomeworkStats(ctx, curatorID)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}

func (c *CachedDashboardRepo) GetCuratorPerformanceZones(ctx context.Context, curatorID string) (domain.PerformanceZones, error) {
	key := cacheKey("curator:%s:zones", curatorID)
	var v domain.PerformanceZones
	ok, err := cacheGet(ctx, c.rdb, key, &v)
	if err == nil && ok {
		return v, nil
	}
	v, err = c.next.GetCuratorPerformanceZones(ctx, curatorID)
	if err == nil {
		cacheSet(ctx, c.rdb, key, v)
	}
	return v, err
}
