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
	Reviews         []*TeacherReview `json:"reviews"`
}
