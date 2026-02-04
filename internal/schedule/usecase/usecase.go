package usecase

import (
	"context"
	"time"

	"lms_backend/internal/domain"
	"lms_backend/internal/schedule/repository"
)

type ScheduleUseCase struct {
	repo repository.ScheduleRepository
}

func NewScheduleUseCase(repo repository.ScheduleRepository) *ScheduleUseCase {
	return &ScheduleUseCase{repo: repo}
}

func (uc *ScheduleUseCase) GetWeeklySchedule(ctx context.Context, userID string, date time.Time) (*domain.WeeklySchedule, error) {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := date.AddDate(0, 0, -(weekday - 1))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end := start.AddDate(0, 0, 7).Add(-time.Second)

	lessons, err := uc.repo.GetStudentLessonsInRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	days := make(map[string][]domain.ScheduleLesson)
	for _, l := range lessons {
		dayKey := l.StartTime.Format("2006-01-02")
		days[dayKey] = append(days[dayKey], l)
	}

	return &domain.WeeklySchedule{
		StartDate: start,
		EndDate:   end,
		Days:      days,
	}, nil
}

func (uc *ScheduleUseCase) GetMonthlySchedule(ctx context.Context, userID string, year, month int) (*domain.MonthlySchedule, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	lessons, err := uc.repo.GetStudentLessonsInRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}

	days := make(map[int][]domain.ScheduleLesson)
	for _, l := range lessons {
		day := l.StartTime.Day()
		days[day] = append(days[day], l)
	}

	return &domain.MonthlySchedule{
		Month: month,
		Year:  year,
		Days:  days,
	}, nil
}
