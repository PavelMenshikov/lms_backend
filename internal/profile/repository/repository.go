package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lms_backend/internal/domain"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) error
	UpdateTeacherSchedule(ctx context.Context, userID string, scheduleJSON []byte) error
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
		SELECT 
			u.id, u.first_name, u.last_name, u.first_name || ' ' || u.last_name as full_name,
			u.email, u.role, u.created_at,
			COALESCE(u.phone, ''), COALESCE(u.city, ''), COALESCE(u.language, 'ru'), COALESCE(u.gender, ''), 
			COALESCE(u.birth_date, '0001-01-01 00:00:00'::timestamp), COALESCE(u.school_name, ''),
			COALESCE(u.experience_years, 0), COALESCE(u.whatsapp_link, ''), COALESCE(u.telegram_link, ''), 
			COALESCE(u.avatar_url, ''),
			COALESCE(ROUND(tr.avg_rating, 1), 0.0) as rating
		FROM users u
		LEFT JOIN (SELECT teacher_id, AVG(rating) as avg_rating FROM teacher_reviews GROUP BY teacher_id) tr ON tr.teacher_id = u.id
		WHERE u.id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.FullName,
		&u.Email, &u.Role, &u.CreatedAt,
		&u.Phone, &u.City, &u.Language, &u.Gender,
		&u.BirthDate, &u.SchoolName,
		&u.ExperienceYears, &u.Whatsapp, &u.Telegram, &u.AvatarURL,
		&u.Rating,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}

	if u.Role == domain.RoleStudent {
		// Fetch parents
		parentRows, err := r.db.QueryContext(ctx, `
			SELECT u.first_name || ' ' || u.last_name, COALESCE(u.phone, ''), u.email
			FROM child_parent_link cpl
			JOIN users u ON cpl.parent_id = u.id
			WHERE cpl.child_id = $1
		`, userID)
		if err == nil {
			defer parentRows.Close()
			for parentRows.Next() {
				var p domain.ParentInfo
				if parentRows.Scan(&p.FullName, &p.Phone, &p.Email) == nil {
					u.Parents = append(u.Parents, p)
				}
			}
		}
		if u.Parents == nil {
			u.Parents = []domain.ParentInfo{}
		}

		// Fetch courses completed count
		r.db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM user_courses WHERE user_id = $1 AND progress_percent >= 100
		`, userID).Scan(&u.CoursesCompleted)

		// Fetch groups count
		r.db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM user_courses uc
			JOIN groups g ON uc.group_id = g.id
			WHERE uc.user_id = $1
		`, userID).Scan(&u.GroupsCount)
	}

	return u, nil
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
func (r *ProfileRepoImpl) UpdateTeacherSchedule(ctx context.Context, userID string, scheduleJSON []byte) error {
	query := `
		INSERT INTO teachers (id, working_hours) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE SET working_hours = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, scheduleJSON)
	return err
}
