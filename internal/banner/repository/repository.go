package repository

import (
	"context"
	"database/sql"
	"lms_backend/internal/domain"

	"github.com/lib/pq"
)

type BannerRepository interface {
	Create(ctx context.Context, banner *domain.Banner) error
	GetByID(ctx context.Context, id string) (*domain.Banner, error)
	GetActive(ctx context.Context, role *string) ([]*domain.Banner, error)
	GetAll(ctx context.Context) ([]*domain.Banner, error)
	Update(ctx context.Context, banner *domain.Banner) error
	Delete(ctx context.Context, id string) error
}

type bannerRepository struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) BannerRepository {
	return &bannerRepository{db: db}
}

func (r *bannerRepository) Create(ctx context.Context, banner *domain.Banner) error {
	query := `
		INSERT INTO banners (id, title, content, type, is_active, priority, start_date, end_date, target_roles, created_by)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		banner.Title, banner.Content, banner.Type, banner.IsActive, banner.Priority,
		banner.StartDate, banner.EndDate, pq.Array(banner.TargetRoles), banner.CreatedBy,
	).Scan(&banner.ID, &banner.CreatedAt, &banner.UpdatedAt)
}

func (r *bannerRepository) GetByID(ctx context.Context, id string) (*domain.Banner, error) {
	var banner domain.Banner
	query := `
		SELECT id, title, content, type, is_active, priority, start_date, end_date, target_roles, created_by, created_at, updated_at
		FROM banners
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&banner.ID, &banner.Title, &banner.Content, &banner.Type, &banner.IsActive,
		&banner.Priority, &banner.StartDate, &banner.EndDate, pq.Array(&banner.TargetRoles),
		&banner.CreatedBy, &banner.CreatedAt, &banner.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

func (r *bannerRepository) GetActive(ctx context.Context, role *string) ([]*domain.Banner, error) {
	query := `
		SELECT id, title, content, type, is_active, priority, start_date, end_date, target_roles, created_by, created_at, updated_at
		FROM banners
		WHERE is_active = true
		  AND (start_date IS NULL OR start_date <= CURRENT_TIMESTAMP)
		  AND (end_date IS NULL OR end_date >= CURRENT_TIMESTAMP)
		  AND ($1::text IS NULL OR target_roles IS NULL OR $1 = ANY(target_roles))
		ORDER BY priority DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*domain.Banner
	for rows.Next() {
		var banner domain.Banner
		err := rows.Scan(
			&banner.ID, &banner.Title, &banner.Content, &banner.Type, &banner.IsActive,
			&banner.Priority, &banner.StartDate, &banner.EndDate, pq.Array(&banner.TargetRoles),
			&banner.CreatedBy, &banner.CreatedAt, &banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}
	return banners, nil
}

func (r *bannerRepository) GetAll(ctx context.Context) ([]*domain.Banner, error) {
	query := `
		SELECT id, title, content, type, is_active, priority, start_date, end_date, target_roles, created_by, created_at, updated_at
		FROM banners
		ORDER BY priority DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*domain.Banner
	for rows.Next() {
		var banner domain.Banner
		err := rows.Scan(
			&banner.ID, &banner.Title, &banner.Content, &banner.Type, &banner.IsActive,
			&banner.Priority, &banner.StartDate, &banner.EndDate, pq.Array(&banner.TargetRoles),
			&banner.CreatedBy, &banner.CreatedAt, &banner.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}
	return banners, nil
}

func (r *bannerRepository) Update(ctx context.Context, banner *domain.Banner) error {
	query := `
		UPDATE banners
		SET title = $1, content = $2, type = $3, is_active = $4, priority = $5,
		    start_date = $6, end_date = $7, target_roles = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		banner.Title, banner.Content, banner.Type, banner.IsActive, banner.Priority,
		banner.StartDate, banner.EndDate, pq.Array(banner.TargetRoles), banner.ID,
	).Scan(&banner.UpdatedAt)
}

func (r *bannerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM banners WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
