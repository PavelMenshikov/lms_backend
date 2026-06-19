package domain

import "time"

type TeacherProfile struct {
	User
	Bio             string   `json:"bio"`
	Rating          float64  `json:"rating"`
	ExperienceYears int      `json:"experience_years"`
	Courses         []string `json:"courses"`
}

type TeacherReview struct {
	ID          string    `json:"id"`
	TeacherID   string    `json:"teacher_id"`
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	CourseTitle string    `json:"course_title"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type TeacherPublicInfo struct {
	ID              string           `json:"id"`
	FirstName       string           `json:"first_name"`
	LastName        string           `json:"last_name"`
	AvatarURL       string           `json:"avatar_url"`
	Rating          float64          `json:"rating"`
	ExperienceYears int              `json:"experience_years"`
	Bio             string           `json:"bio"`
	Email           string           `json:"email"`
	City            string           `json:"city"`
	Phone           string           `json:"phone"`
	Schedule        interface{}      `json:"schedule"`
	Reviews         []*TeacherReview `json:"reviews"`
}
type TeacherMonthlyReport struct {
	TeacherID              string  `json:"teacher_id"`
	Year                   int     `json:"year"`
	Month                  int     `json:"month"`
	TotalLessons           int     `json:"total_lessons"`
	LessonsPerWeek         float64 `json:"lessons_per_week"`
	SubstitutionsCount     int     `json:"substitutions_count"`
	ReplacedCount          int     `json:"replaced_count"`
	SubstitutionsDelta     int     `json:"substitutions_delta"`
	ReplacedDelta          int     `json:"replaced_delta"`
	AvgRating              float64 `json:"avg_rating"`
	TotalStudents          int     `json:"total_students"`
	AttendanceAvg          float64 `json:"attendance_avg_percent"`
	AverageHomeworkScore   float64 `json:"average_homework_score"`
	TotalCancelled         int     `json:"total_cancelled"`
	CancelledDelta         int     `json:"cancelled_delta"`
}

type TeacherCertificate struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	StudentName string    `json:"student_name"`
	IssuedAt    time.Time `json:"issued_at"`
	CourseName  string    `json:"course_name"`
	CertificateURL string `json:"certificate_url"`
}

type TeacherDashboardData struct {
	Profile          *TeacherPublicInfo      `json:"profile"`
	AssignedCourses  []*StudentCoursePreview `json:"assigned_courses"`
	MyReviews        []*TeacherReview        `json:"my_reviews"`
	Substitutions    []*Lesson               `json:"substitutions,omitempty"`
	CancelledLessons []*Lesson               `json:"cancelled_lessons,omitempty"`
	UpcomingLessons  []*Lesson               `json:"upcoming_lessons,omitempty"`
}
