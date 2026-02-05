package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		       COALESCE(phone, ''), COALESCE(city, ''), COALESCE(language, 'ru'), COALESCE(gender, ''), 
               COALESCE(birth_date, '0001-01-01 00:00:00Z'), COALESCE(school_name, ''),
		       COALESCE(experience_years, 0), COALESCE(whatsapp_link, ''), COALESCE(telegram_link, ''), COALESCE(avatar_url, '')
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.CreatedAt,
		&u.Phone, &u.City, &u.Language, &u.Gender, &u.BirthDate, &u.SchoolName,
		&u.ExperienceYears, &u.Whatsapp, &u.Telegram, &u.AvatarURL,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("profile not found")
	}
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
	res, err := r.db.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Phone, u.City,
		u.Language, u.SchoolName, u.Whatsapp,
		u.Telegram, u.AvatarURL, u.ID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("profile with ID %s not found", u.ID)
	}
	return nil
}
