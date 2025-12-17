package domain

type ChildParentLink struct {
	ID       string `json:"id" db:"id"`
	ChildID  string `json:"child_id" db:"child_id"`
	ParentID string `json:"parent_id" db:"parent_id"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type StudentDetails struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
