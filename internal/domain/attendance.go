package domain

import "time"

type AttendanceStatus string

const (
	AttendanceStatusAttended        AttendanceStatus = "ATTENDED"
	AttendanceStatusAbsentExcused   AttendanceStatus = "ABSENT_EXCUSED"
	AttendanceStatusAbsentUnexcused AttendanceStatus = "ABSENT_UNEXCUSED"
	AttendanceStatusFreeze          AttendanceStatus = "FREEZE"
)

type AttendanceRecord struct {
	ID        string           `json:"id" db:"id"`
	LessonID  string           `json:"lesson_id" db:"lesson_id"`
	StudentID string           `json:"student_id" db:"student_id"`
	Status    AttendanceStatus `json:"status" db:"status"`
	Reason    *string          `json:"reason,omitempty" db:"reason"`
	Comment   *string          `json:"comment,omitempty" db:"comment"`
	MarkedBy  *string          `json:"marked_by,omitempty" db:"marked_by"`
	MarkedAt  time.Time        `json:"marked_at" db:"marked_at"`
	UpdatedBy *string          `json:"updated_by,omitempty" db:"updated_by"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}
