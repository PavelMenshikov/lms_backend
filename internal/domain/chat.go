package domain

import "time"

type ChatMessage struct {
	ID          string    `json:"id" db:"id"`
	ModuleID    string    `json:"module_id" db:"module_id"`
	StudentID   string    `json:"student_id" db:"student_id"`
	SenderID    string    `json:"sender_id" db:"sender_id"`
	SenderName  string    `json:"sender_name" db:"sender_name"`
	SenderRole  Role      `json:"sender_role" db:"sender_role"`
	MessageText string    `json:"message_text" db:"message_text"`
	FileURL     string    `json:"file_url" db:"file_url"`
	IsRead      bool      `json:"is_read" db:"is_read"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ChatHistoryRequest struct {
	ModuleID  string `json:"module_id"`
	StudentID string `json:"student_id"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}
