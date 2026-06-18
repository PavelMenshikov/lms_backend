package repository

import (
	"context"
	"database/sql"
	"strings"

	"lms_backend/internal/domain"
)

func (r *ContentAdminRepoImpl) CreateStream(ctx context.Context, s *domain.Stream) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO streams (course_id, title, start_date) VALUES ($1, $2, $3) RETURNING id", s.CourseID, s.Title, s.StartDate).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) GetStreamsByCourse(ctx context.Context, courseID string) ([]*domain.Stream, error) {
	query := `SELECT id, course_id, title, start_date, created_at FROM streams`
	args := []interface{}{}
	if courseID != "" {
		query += " WHERE course_id = $1"
		args = append(args, courseID)
	}
	query += " ORDER BY created_at DESC LIMIT 200"
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var streams []*domain.Stream
	for rows.Next() {
		s := &domain.Stream{}
		if err := rows.Scan(&s.ID, &s.CourseID, &s.Title, &s.StartDate, &s.CreatedAt); err != nil {
			return nil, err
		}
		streams = append(streams, s)
	}
	return streams, nil
}

func (r *ContentAdminRepoImpl) CreateGroup(ctx context.Context, g *domain.Group) (string, error) {
	var newID string
	err := r.db.QueryRowContext(ctx, "INSERT INTO groups (stream_id, curator_id, teacher_id, title) VALUES ($1, $2, $3, $4) RETURNING id", g.StreamID, g.CuratorID, g.TeacherID, g.Title).Scan(&newID)
	return newID, err
}

func (r *ContentAdminRepoImpl) GetGroupsByStream(ctx context.Context, streamID string) ([]*domain.Group, error) {
	query := `SELECT id, stream_id, curator_id, teacher_id, title, created_at FROM groups`
	args := []interface{}{}
	if streamID != "" {
		query += " WHERE stream_id = $1"
		args = append(args, streamID)
	}
	query += " ORDER BY created_at DESC LIMIT 200"
	query = `SELECT id, stream_id, COALESCE(curator_id, ''), COALESCE(teacher_id, ''), title, created_at FROM groups` + strings.TrimPrefix(query, "SELECT id, stream_id, curator_id, teacher_id, title, created_at FROM groups")
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*domain.Group
	for rows.Next() {
		g := &domain.Group{}
		if err := rows.Scan(&g.ID, &g.StreamID, &g.CuratorID, &g.TeacherID, &g.Title, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *ContentAdminRepoImpl) GetStudentEnrollment(ctx context.Context, userID string) (map[string]string, error) {
	var cid, sid, gid sql.NullString
	query := `SELECT course_id, stream_id, group_id FROM user_courses WHERE user_id = $1 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&cid, &sid, &gid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	res := make(map[string]string)
	if cid.Valid {
		res["course_id"] = cid.String
	}
	if sid.Valid {
		res["stream_id"] = sid.String
	}
	if gid.Valid {
		res["group_id"] = gid.String
	}

	return res, nil
}
