package domain

import "time"

type Role string

const (
	RoleStudent   Role = "student"
	RoleParent    Role = "parent"
	RoleTeacher   Role = "teacher"
	RoleModerator Role = "moderator"
	RoleCurator   Role = "curator"
	RoleAdmin     Role = "admin"
)

type User struct {
	ID              string    `json:"id" db:"id"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"-" db:"password_hash"`
	Role            Role      `json:"role" db:"role"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	Phone           string    `json:"phone" db:"phone"`
	City            string    `json:"city" db:"city"`
	Language        string    `json:"language" db:"language"`
	Gender          string    `json:"gender" db:"gender"`
	BirthDate       time.Time `json:"birth_date" db:"birth_date"`
	SchoolName      string    `json:"school_name" db:"school_name"`
	ExperienceYears int       `json:"experience_years" db:"experience_years"`
	Whatsapp        string    `json:"whatsapp" db:"whatsapp_link"`
	Telegram        string    `json:"telegram" db:"telegram_link"`
	AvatarURL       string    `json:"avatar_url" db:"avatar_url"`
}

type UserFilter struct {
	Role     Role   `json:"role"`
	CourseID string `json:"course_id"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

type StudentTableItem struct {
	ID              string    `json:"id"`
	Photo           string    `json:"photo"`
	FullName        string    `json:"full_name"`
	CreatedAt       time.Time `json:"created_at"`
	Gender          string    `json:"gender"`
	Age             int       `json:"age"`
	Status          string    `json:"status"`
	Course          string    `json:"course"`
	Group           string    `json:"group"`
	Curator         string    `json:"curator"`
	Teacher         string    `json:"teacher"`
	Stream          string    `json:"stream"`
	Performance     int       `json:"performance"`
	ParentPhone     string    `json:"parent_phone"`
	City            string    `json:"city"`
	School          string    `json:"school"`
	Language        string    `json:"language"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
}

type TeacherTableItem struct {
	ID              string    `json:"id"`
	Photo           string    `json:"photo"`
	FullName        string    `json:"full_name"`
	CreatedAt       time.Time `json:"created_at"`
	Gender          string    `json:"gender"`
	Phone           string    `json:"phone"`
	Groups          string    `json:"groups"`
	City            string    `json:"city"`
	Email           string    `json:"email"`
	ExperienceYears int       `json:"experience_years"`
	Language        string    `json:"language"`
}

type CuratorTableItem struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	Groups    string    `json:"groups"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
}

type ModeratorTableItem struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
}

type AllUsersTableItem struct {
	ID        string    `json:"id"`
	Photo     string    `json:"photo"`
	FullName  string    `json:"full_name"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}