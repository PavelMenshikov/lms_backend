package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
}

type ProfileRepoImpl struct {
	db *sql.DB
}

var _ ProfileRepository = (*ProfileRepoImpl)(nil)

func NewProfileRepository(db *sql.DB) *ProfileRepoImpl {
	return &ProfileRepoImpl{db: db}
}

func (r *ProfileRepoImpl) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT id, first_name, last_name, email, role, created_at,
		       phone, city, language, gender, birth_date, school_name,
		       experience_years, whatsapp_link, telegram_link, avatar_url
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.CreatedAt,
		&u.Phone, &u.City, &u.Language, &u.Gender, &u.BirthDate, &u.SchoolName,
		&u.ExperienceYears, &u.Whatsapp, &u.Telegram, &u.AvatarURL,
	)
	return u, err
}

func (r *ProfileRepoImpl) UpdateProfile(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET
			first_name = $1, last_name = $2, phone = $3, city = $4,
			language = $5, school_name = $6, whatsapp_link = $7,
			telegram_link = $8, avatar_url = $9
		WHERE id = $10
	`
	_, err := r.db.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Phone, u.City,
		u.Language, u.SchoolName, u.Whatsapp,
		u.Telegram, u.AvatarURL, u.ID,
	)
	return err
}
