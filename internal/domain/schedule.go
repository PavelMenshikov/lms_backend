package domain

import "time"

type ScheduleLesson struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	CourseName     string    `json:"course_name"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	DurationMin    int       `json:"duration_min"`
	TeacherName    string    `json:"teacher_name"`
	TeacherEmail   string    `json:"teacher_email"`
	DiscordURL     string    `json:"discord_url"`
	TeacherComment string    `json:"teacher_comment"`
	HomeworkStatus string    `json:"homework_status"`
	Color          string    `json:"color"`
}

type WeeklySchedule struct {
	StartDate time.Time                   `json:"start_date"`
	EndDate   time.Time                   `json:"end_date"`
	Days      map[string][]ScheduleLesson `json:"days"`
}

type MonthlySchedule struct {
	Month int                      `json:"month"`
	Year  int                      `json:"year"`
	Days  map[int][]ScheduleLesson `json:"days"`
}
