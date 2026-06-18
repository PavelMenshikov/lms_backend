package repository

import (
	"context"
	"fmt"

	"lms_backend/internal/domain"
)

func (r *ContentAdminRepoImpl) CreateUser(ctx context.Context, u *domain.User) (string, error) {
	var newID string
	query := `INSERT INTO users (first_name, last_name, email, password_hash, role, phone, city, language, gender, birth_date, school_name, experience_years, whatsapp_link, telegram_link, avatar_url, intro_broadcast_url, graduation_broadcast_url, balance)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Password, u.Role, u.Phone, u.City, u.Language, u.Gender, u.BirthDate, u.SchoolName, u.ExperienceYears, u.Whatsapp, u.Telegram, u.AvatarURL, u.IntroBroadcastURL, u.GraduationBroadcastURL, u.Balance).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) GetUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, error) {
	query := `
        SELECT id, first_name, last_name, first_name || ' ' || last_name as full_name, email, role, created_at, 
               COALESCE(phone, ''), COALESCE(city, ''), COALESCE(school_name, ''), COALESCE(avatar_url, '')
        FROM users WHERE 1=1 
    `
	args := []interface{}{}
	argID := 1

	if filter.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argID)
		args = append(args, filter.Role)
		argID++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, filter.Limit)
		argID++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argID)
		args = append(args, filter.Offset)
		argID++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.FullName, &u.Email, &u.Role, &u.CreatedAt,
			&u.Phone, &u.City, &u.SchoolName, &u.AvatarURL); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *ContentAdminRepoImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	query := `
		SELECT 
			id, first_name, last_name, first_name || ' ' || last_name as full_name, 
			email, role, created_at, COALESCE(phone, ''), COALESCE(city, ''), 
			COALESCE(school_name, ''), COALESCE(language, ''), COALESCE(gender, ''), 
			COALESCE(birth_date, NOW()), COALESCE(experience_years, 0), 
			COALESCE(whatsapp_link, ''), COALESCE(telegram_link, ''), COALESCE(avatar_url, ''),
			COALESCE(intro_broadcast_url, ''), COALESCE(graduation_broadcast_url, ''),
			subscription_end_date, COALESCE(balance, 0), COALESCE(loss_reason, '')
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.FullName, &u.Email, &u.Role, &u.CreatedAt,
		&u.Phone, &u.City, &u.SchoolName, &u.Language, &u.Gender, &u.BirthDate,
		&u.ExperienceYears, &u.Whatsapp, &u.Telegram, &u.AvatarURL,
		&u.IntroBroadcastURL, &u.GraduationBroadcastURL,
		&u.SubscriptionEndDate, &u.Balance, &u.LossReason,
	)
	return u, err
}

func (r *ContentAdminRepoImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, first_name, last_name, email, role, COALESCE(phone, '') FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.Phone)
	return u, err
}

func (r *ContentAdminRepoImpl) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, first_name, last_name, email, COALESCE(phone, '') FROM users WHERE phone = $1`
	err := r.db.QueryRowContext(ctx, query, phone).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Phone)
	return u, err
}

func (r *ContentAdminRepoImpl) UpdateUser(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE users SET
			first_name = $1, last_name = $2, email = $3, role = $4,
			phone = $5, city = $6, school_name = $7, experience_years = $8, 
			whatsapp_link = $9, telegram_link = $10,
			gender = $11, language = $12, birth_date = $13,
			intro_broadcast_url = $14, graduation_broadcast_url = $15,
			balance = $16, loss_reason = $17, subscription_end_date = $18
		WHERE id = $19
	`
	_, err := r.db.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Role,
		u.Phone, u.City, u.SchoolName, u.ExperienceYears,
		u.Whatsapp, u.Telegram,
		u.Gender, u.Language, u.BirthDate,
		u.IntroBroadcastURL, u.GraduationBroadcastURL,
		u.Balance, u.LossReason, u.SubscriptionEndDate,
		u.ID,
	)
	return err
}

func (r *ContentAdminRepoImpl) DeleteUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}

func (r *ContentAdminRepoImpl) GetParentsByStudentID(ctx context.Context, studentID string) ([]domain.User, error) {
	query := `
		SELECT 
			u.id, u.first_name, u.last_name, u.first_name || ' ' || u.last_name as full_name, 
			u.email, u.phone 
		FROM users u 
		JOIN child_parent_link cpl ON u.id = cpl.parent_id 
		WHERE cpl.child_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parents []domain.User
	for rows.Next() {
		u := domain.User{}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.FullName, &u.Email, &u.Phone); err != nil {
			return nil, err
		}
		parents = append(parents, u)
	}
	return parents, nil
}

func (r *ContentAdminRepoImpl) LinkParentToStudent(ctx context.Context, studentID, parentID string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO child_parent_link (child_id, parent_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", studentID, parentID)
	return err
}

func (r *ContentAdminRepoImpl) UnlinkAllParents(ctx context.Context, studentID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM child_parent_link WHERE child_id = $1", studentID)
	return err
}

func (r *ContentAdminRepoImpl) GetDetailedStudentList(ctx context.Context, filter domain.UserFilter) ([]*domain.StudentTableItem, error) {
	query := `
		SELECT 
			u.id, 
			COALESCE(u.avatar_url, '') as photo, 
			COALESCE(u.first_name || ' ' || u.last_name, '') as full_name, 
			u.created_at, 
			COALESCE(u.gender, ''),
			COALESCE(EXTRACT(YEAR FROM AGE(u.birth_date)), 0) as age,
			COALESCE(STRING_AGG(DISTINCT uc.status, ', '), 'inactive') as status, 
			COALESCE(STRING_AGG(DISTINCT c.title, ', '), '') as course, 
			COALESCE(STRING_AGG(DISTINCT g.title, ', '), '') as "group", 
			COALESCE(STRING_AGG(DISTINCT cur.first_name || ' ' || cur.last_name, ', '), '') as curator,
			COALESCE(STRING_AGG(DISTINCT teach.first_name || ' ' || teach.last_name, ', '), '') as teacher,
			COALESCE(STRING_AGG(DISTINCT s.title, ', '), '') as stream,
			COALESCE(AVG(uc.progress_percent)::INT, 0) as performance,
			COALESCE(MIN(par.phone), '') as parent_phone,
			COALESCE(u.city, '') as city, 
			COALESCE(u.school_name, '') as school, 
			COALESCE(u.language, '') as language,
			COALESCE(u.phone, '') as phone, 
			u.email,
			COALESCE(u.intro_broadcast_url, ''),
			COALESCE(u.graduation_broadcast_url, ''),
			COALESCE(u.balance, 0)
		FROM users u
		LEFT JOIN user_courses uc ON u.id = uc.user_id
		LEFT JOIN courses c ON uc.course_id = c.id
		LEFT JOIN groups g ON uc.group_id = g.id
		LEFT JOIN users cur ON g.curator_id = cur.id
		LEFT JOIN users teach ON g.teacher_id = teach.id
		LEFT JOIN streams s ON uc.stream_id = s.id
		LEFT JOIN child_parent_link cpl ON cpl.child_id = u.id
		LEFT JOIN users par ON par.id = cpl.parent_id
		WHERE u.role = 'student'`

	args := []interface{}{}
	if filter.CourseID != "" {
		query += " AND uc.course_id = $1"
		args = append(args, filter.CourseID)
	}

	query += ` GROUP BY u.id ORDER BY u.created_at DESC LIMIT 500`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch detailed students: %w", err)
	}
	defer rows.Close()

	var list []*domain.StudentTableItem
	for rows.Next() {
		item := &domain.StudentTableItem{}
		err := rows.Scan(
			&item.ID, &item.Photo, &item.FullName, &item.CreatedAt, &item.Gender, &item.Age,
			&item.Status, &item.Course, &item.Group, &item.Curator, &item.Teacher, &item.Stream,
			&item.Performance, &item.ParentPhone, &item.City, &item.School, &item.Language,
			&item.Phone, &item.Email,
			&item.IntroBroadcastURL, &item.GraduationBroadcastURL, &item.Balance,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedTeacherList(ctx context.Context) ([]*domain.TeacherTableItem, error) {
	query := `
		SELECT u.id, COALESCE(u.avatar_url, ''), COALESCE(u.first_name || ' ' || u.last_name, ''),
			u.created_at, COALESCE(u.gender, ''), COALESCE(u.phone, ''),
			COALESCE(STRING_AGG(g.title, ', '), '') as groups,
			COALESCE(u.city, ''), u.email, COALESCE(u.experience_years, 0),
			COALESCE(u.language, '')
		FROM users u
		LEFT JOIN groups g ON g.teacher_id = u.id
		WHERE u.role = 'teacher'
		GROUP BY u.id
		ORDER BY u.created_at DESC
		LIMIT 200
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.TeacherTableItem
	for rows.Next() {
		item := &domain.TeacherTableItem{}
		if err := rows.Scan(&item.ID, &item.Photo, &item.FullName, &item.CreatedAt, &item.Gender, &item.Phone, &item.Groups, &item.City, &item.Email, &item.ExperienceYears, &item.Language); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedCuratorList(ctx context.Context) ([]*domain.CuratorTableItem, error) {
	query := `
		SELECT u.id, COALESCE(u.first_name || ' ' || u.last_name, ''), u.created_at,
			COALESCE(STRING_AGG(g.title, ', '), '') as groups, COALESCE(u.phone, ''), u.email
		FROM users u
		LEFT JOIN groups g ON g.curator_id = u.id
		WHERE u.role = 'curator'
		GROUP BY u.id
		ORDER BY u.created_at DESC
		LIMIT 200
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.CuratorTableItem
	for rows.Next() {
		item := &domain.CuratorTableItem{}
		if err := rows.Scan(&item.ID, &item.FullName, &item.CreatedAt, &item.Groups, &item.Phone, &item.Email); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetDetailedModeratorList(ctx context.Context) ([]*domain.ModeratorTableItem, error) {
	query := `SELECT id, COALESCE(first_name || ' ' || last_name, ''), created_at, COALESCE(phone, ''), email FROM users WHERE role = 'moderator' ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.ModeratorTableItem
	for rows.Next() {
		item := &domain.ModeratorTableItem{}
		if err := rows.Scan(&item.ID, &item.FullName, &item.CreatedAt, &item.Phone, &item.Email); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}

func (r *ContentAdminRepoImpl) GetAllUsersList(ctx context.Context) ([]*domain.AllUsersTableItem, error) {
	query := `SELECT id, COALESCE(avatar_url, ''), first_name || ' ' || last_name, role, created_at FROM users ORDER BY created_at DESC LIMIT 200`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*domain.AllUsersTableItem
	for rows.Next() {
		item := &domain.AllUsersTableItem{}
		if err := rows.Scan(&item.ID, &item.Photo, &item.FullName, &item.Role, &item.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}
