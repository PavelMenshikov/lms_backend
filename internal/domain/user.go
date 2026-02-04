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
	Role   Role
	Limit  int
	Offset int
}
