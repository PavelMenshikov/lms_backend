package domain

import "time"

type StudentCoursePreview struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ImageURL        string `json:"image_url"`
	ProgressPercent int    `json:"progress_percent"`
	IsMain          bool   `json:"is_main"`
}

type StudentCourseView struct {
	Course  *Course              `json:"course"`
	Modules []*StudentModuleView `json:"modules"`
}

type StudentModuleView struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	OrderNum    int                 `json:"order_num"`
	Lessons     []*StudentLessonRef `json:"lessons"`
}

type StudentLessonRef struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	OrderNum    int    `json:"order_num"`
	DurationMin int    `json:"duration_min"`
	IsCompleted bool   `json:"is_completed"`
	IsLocked    bool   `json:"is_locked"`
}

type StudentLessonDetail struct {
	Lesson           *Lesson           `json:"lesson"`
	Materials        []*LessonMaterial `json:"materials"`
	PreviousLessonID string            `json:"previous_lesson_id,omitempty"`
	NextLessonID     string            `json:"next_lesson_id,omitempty"`
	IsCompleted      bool              `json:"is_completed"`
	AssignmentStatus string            `json:"assignment_status,omitempty"`
	TeacherComment   string            `json:"teacher_comment,omitempty"`
}

type LessonMaterial struct {
	ID        string `json:"id" db:"id"`
	LessonID  string `json:"lesson_id" db:"lesson_id"`
	Title     string `json:"title" db:"title"`
	S3Path    string `json:"s3_path" db:"s3_path"`
	PublicURL string `json:"public_url"`
}

type SubmissionRecord struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StudentName string    `json:"student_name"`
	LessonTitle string    `json:"lesson_title"`
	CourseTitle string    `json:"course_title"`
	Text        string    `json:"submission_text"`
	Link        string    `json:"submission_link"`
	SubmittedAt time.Time `json:"submitted_at"`
}

type StudentSubmissionInput struct {
	LessonID   string
	UserID     string
	TextAnswer string
}
