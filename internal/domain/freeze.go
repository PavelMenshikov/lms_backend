package domain

import "time"

type FreezeStatus string

const (
	FreezeStatusPending  FreezeStatus = "PENDING"
	FreezeStatusApproved FreezeStatus = "APPROVED"
	FreezeStatusRejected FreezeStatus = "REJECTED"
)

type FreezeRequest struct {
	ID            string       `json:"id" db:"id"`
	StudentID     string       `json:"student_id" db:"student_id"`
	RequestedBy   string       `json:"requested_by" db:"requested_by"`
	StartDate     time.Time    `json:"start_date" db:"start_date"`
	EndDate       time.Time    `json:"end_date" db:"end_date"`
	Reason        string       `json:"reason" db:"reason"`
	Status        FreezeStatus `json:"status" db:"status"`
	ReviewedBy    *string      `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt    *time.Time   `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewComment *string      `json:"review_comment,omitempty" db:"review_comment"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

type FreezePeriod struct {
	ID              string    `json:"id" db:"id"`
	StudentID       string    `json:"student_id" db:"student_id"`
	FreezeRequestID *string   `json:"freeze_request_id,omitempty" db:"freeze_request_id"`
	StartDate       time.Time `json:"start_date" db:"start_date"`
	EndDate         time.Time `json:"end_date" db:"end_date"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	UsedDays        int       `json:"used_days"`
	RemainingDays   int       `json:"remaining_days"`
}
