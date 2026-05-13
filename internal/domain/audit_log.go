package domain

import "time"

type AuditLog struct {
	ID         string     `json:"id" db:"id"`
	UserID     *string    `json:"user_id,omitempty" db:"user_id"`
	Action     string     `json:"action" db:"action"`
	EntityType string     `json:"entity_type" db:"entity_type"`
	EntityID   string     `json:"entity_id" db:"entity_id"`
	OldValues  *string    `json:"old_values,omitempty" db:"old_values"`
	NewValues  *string    `json:"new_values,omitempty" db:"new_values"`
	IPAddress  *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent  *string    `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
