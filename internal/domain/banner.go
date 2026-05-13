package domain

import "time"

type BannerType string

const (
	BannerTypeInfo        BannerType = "INFO"
	BannerTypeWarning     BannerType = "WARNING"
	BannerTypeAnnouncement BannerType = "ANNOUNCEMENT"
	BannerTypePromotion   BannerType = "PROMOTION"
)

type Banner struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Content     string     `json:"content" db:"content"`
	Type        BannerType `json:"type" db:"type"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	Priority    int        `json:"priority" db:"priority"`
	StartDate   *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	TargetRoles []string   `json:"target_roles,omitempty" db:"target_roles"`
	CreatedBy   string     `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
