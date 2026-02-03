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
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}