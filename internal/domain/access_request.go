package domain

import "time"

type AccessRequestStatus string

const (
	AccessRequestStatusPending  AccessRequestStatus = "PENDING"
	AccessRequestStatusApproved AccessRequestStatus = "APPROVED"
	AccessRequestStatusRejected AccessRequestStatus = "REJECTED"
)

type AccessRequest struct {
	ID            string              `json:"id" db:"id"`
	UserID        string              `json:"user_id" db:"user_id"`
	ResourceType  string              `json:"resource_type" db:"resource_type"`
	ResourceID    string              `json:"resource_id" db:"resource_id"`
	Reason        string              `json:"reason" db:"reason"`
	Status        AccessRequestStatus `json:"status" db:"status"`
	ReviewedBy    *string             `json:"reviewed_by,omitempty" db:"reviewed_by"`
	ReviewedAt    *time.Time          `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewComment *string             `json:"review_comment,omitempty" db:"review_comment"`
	CreatedAt     time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" db:"updated_at"`
}
