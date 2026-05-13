package domain

import "time"

type StudentStatistics struct {
	ID                  string     `json:"id" db:"id"`
	StudentID           string     `json:"student_id" db:"student_id"`
	TotalLessons        int        `json:"total_lessons" db:"total_lessons"`
	AttendedLessons     int        `json:"attended_lessons" db:"attended_lessons"`
	AbsentExcused       int        `json:"absent_excused" db:"absent_excused"`
	AbsentUnexcused     int        `json:"absent_unexcused" db:"absent_unexcused"`
	FreezeDays          int        `json:"freeze_days" db:"freeze_days"`
	RemainingLessons    int        `json:"remaining_lessons" db:"remaining_lessons"`
	RemainingExcused    int        `json:"remaining_excused" db:"remaining_excused"`
	LastAttendanceDate  *time.Time `json:"last_attendance_date,omitempty" db:"last_attendance_date"`
	CurrentFreezeEndDate *time.Time `json:"current_freeze_end_date,omitempty" db:"current_freeze_end_date"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
}
