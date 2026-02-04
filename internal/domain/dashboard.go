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

type AdminHomeDashboard struct {
	TotalStudents    int                   `json:"total_students"`
	StudentsDelta    float64               `json:"students_delta"`
	NewStudentsMonth int                   `json:"new_students_month"`
	TotalTeachers    int                   `json:"total_teachers"`
	ActiveCourses    int                   `json:"active_courses"`
	Performance      PerformanceZones      `json:"performance"`
	LessonActivity   []DailyLessonActivity `json:"lesson_activity"`
}

type PerformanceZones struct {
	Green  int `json:"green"`
	Yellow int `json:"yellow"`
	Red    int `json:"red"`
}

type DailyLessonActivity struct {
	Date       string `json:"date"`
	Group      int    `json:"group"`
	Trial      int    `json:"trial"`
	Individual int    `json:"individual"`
}
