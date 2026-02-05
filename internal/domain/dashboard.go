package domain

import "time"

type LastLesson struct {
	CourseTitle      string `json:"course_title"`
	ModuleName       string `json:"module_name"`
	LessonTitle      string `json:"lesson_title"`
	AssignmentStatus string `json:"assignment_status"`
	LessonID         string `json:"lesson_id"`
	HomeworkID       string `json:"homework_id"`
}

type StatisticSummary struct {
	Percentage float64        `json:"percentage"`
	Delta      float64        `json:"delta_vs_last_month"`
	Breakdown  map[string]int `json:"breakdown"`
}

type UpcomingLesson struct {
	Date          time.Time `json:"date"`
	TimeRange     string    `json:"time_range"`
	CourseTitle   string    `json:"course_title"`
	TeacherName   string    `json:"teacher_name"`
	IsHomeworkDue bool      `json:"is_homework_due"`
}

type HomeDashboard struct {
	UserRole           Role              `json:"role"`
	LastLessonData     *LastLesson       `json:"last_lesson"`
	ActiveCoursesCount int               `json:"active_courses_count"`
	AttendanceStats    *StatisticSummary `json:"attendance_stats"`
	AssignmentStats    *StatisticSummary `json:"assignment_stats"`
	UpcomingLessons    []UpcomingLesson  `json:"upcoming_lessons"`
	User               *User             `json:"user_data"`
}

type PerformanceZones struct {
	Green      int     `json:"green_count"`
	GreenPct   float64 `json:"green_percent"`
	Yellow     int     `json:"yellow_count"`
	YellowPct  float64 `json:"yellow_percent"`
	Red        int     `json:"red_count"`
	RedPct     float64 `json:"red_percent"`
	TotalRated int     `json:"total_rated"`
}

type DailyLessonActivity struct {
	Date       string `json:"date"`
	Group      int    `json:"group"`
	Trial      int    `json:"trial"`
	Individual int    `json:"individual"`
}

type AdminHomeDashboard struct {
	TotalStudents     int                   `json:"total_students"`
	StudentsDelta     float64               `json:"students_delta_percent"`
	NewStudentsMonth  int                   `json:"new_students_month"`
	NewStudentsDelta  float64               `json:"new_students_delta_percent"`
	TotalTeachers     int                   `json:"total_teachers"`
	ActiveCourses     int                   `json:"active_courses"`
	Performance       PerformanceZones      `json:"performance_zones"`
	LessonActivity    []DailyLessonActivity `json:"lesson_activity"`
	UpdatePeriodMonth string                `json:"update_period_month"`
}

type AdminCourseStats struct {
	TotalStudents        int            `json:"total_students"`
	NewStudentsMonth     int            `json:"new_students_month"`
	FrozenStudents       int            `json:"frozen_students"`
	GraduatedStudents    int            `json:"graduated_students"`
	AverageScore         float64        `json:"average_score"`
	AverageLagLessons    float64        `json:"average_lag_lessons"`
	AverageWatchTimeMin  float64        `json:"average_watch_time_minutes"`
	SuccessRateBreakdown map[string]int `json:"success_rate_breakdown"`
}

type AdminStudentProgress struct {
	UserID          string    `json:"user_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PhotoURL        string    `json:"photo_url"`
	Status          string    `json:"status"`
	StartDate       time.Time `json:"start_date"`
	ProgressPercent int       `json:"progress_percent"`
	LessonsAttended int       `json:"lessons_attended"`
	HomeworksDone   int       `json:"homeworks_done"`
}
