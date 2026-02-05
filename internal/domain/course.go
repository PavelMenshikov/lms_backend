package domain

import "time"

type CourseStatus string

const (
	CourseStatusDraft    CourseStatus = "draft"
	CourseStatusActive   CourseStatus = "active"
	CourseStatusArchived CourseStatus = "archived"
)

type Course struct {
	ID                  string       `json:"id" db:"id"`
	Title               string       `json:"title" db:"title"`
	Description         string       `json:"description" db:"description"`
	IsMain              bool         `json:"is_main" db:"is_main"`
	ImageURL            string       `json:"image_url" db:"image_url"`
	Status              CourseStatus `json:"status" db:"status"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
	HasHomework         bool         `json:"has_homework" db:"has_homework"`
	IsHomeworkMandatory bool         `json:"is_homework_mandatory" db:"is_homework_mandatory"`
	IsTestMandatory     bool         `json:"is_test_mandatory" db:"is_test_mandatory"`
	IsProjectMandatory  bool         `json:"is_project_mandatory" db:"is_project_mandatory"`
	IsDiscordMandatory  bool         `json:"is_discord_mandatory" db:"is_discord_mandatory"`
	IsAntiCopyEnabled   bool         `json:"is_anti_copy_enabled" db:"is_anti_copy_enabled"`
}

type Module struct {
	ID          string `json:"id" db:"id"`
	CourseID    string `json:"course_id" db:"course_id"`
	Title       string `json:"title" db:"title"`
	OrderNum    int    `json:"order_num" db:"order_num"`
	Description string `json:"description" db:"description"`
}

type Lesson struct {
	ID              string    `json:"id" db:"id"`
	ModuleID        string    `json:"module_id" db:"module_id"`
	TeacherID       string    `json:"teacher_id" db:"teacher_id"`
	Title           string    `json:"title" db:"title"`
	LessonTime      time.Time `json:"lesson_time" db:"lesson_time"`
	DurationMin     int       `json:"duration_min" db:"duration_min"`
	OrderNum        int       `json:"order_num" db:"order_num"`
	IsPublished     bool      `json:"is_published" db:"is_published"`
	VideoURL        string    `json:"video_url" db:"video_url"`
	PresentationURL string    `json:"presentation_url" db:"presentation_url"`
	ContentText     string    `json:"content_text" db:"content_text"`
	HasHomework     bool      `json:"has_homework" db:"has_homework"`
	IsInteractive   bool      `json:"is_interactive" db:"is_interactive"`
}

type Stream struct {
	ID        string    `json:"id" db:"id"`
	CourseID  string    `json:"course_id" db:"course_id"`
	Title     string    `json:"title" db:"title"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Group struct {
	ID        string    `json:"id" db:"id"`
	StreamID  string    `json:"stream_id" db:"stream_id"`
	CuratorID string    `json:"curator_id" db:"curator_id"`
	TeacherID string    `json:"teacher_id" db:"teacher_id"`
	Title     string    `json:"title" db:"title"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Test struct {
	ID           string    `json:"id" db:"id"`
	LessonID     string    `json:"lesson_id" db:"lesson_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	PassingScore int       `json:"passing_score" db:"passing_score"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Project struct {
	ID          string    `json:"id" db:"id"`
	LessonID    string    `json:"lesson_id" db:"lesson_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	MaxScore    int       `json:"max_score" db:"max_score"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CourseStructure struct {
	Course  *Course            `json:"course"`
	Modules []*ModuleStructure `json:"modules"`
}

type ModuleStructure struct {
	Module  *Module   `json:"module"`
	Lessons []*Lesson `json:"lessons"`
}
