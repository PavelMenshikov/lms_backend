package repository

import (
	"context"
	"database/sql"
	"time"

	"lms_backend/internal/domain"
)

type ScheduleRepository interface {
	GetStudentLessonsInRange(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error)
}

type ScheduleRepoImpl struct {
	db *sql.DB
}

var _ ScheduleRepository = (*ScheduleRepoImpl)(nil)

func NewScheduleRepository(db *sql.DB) *ScheduleRepoImpl {
	return &ScheduleRepoImpl{db: db}
}

func (r *ScheduleRepoImpl) GetStudentLessonsInRange(ctx context.Context, userID string, start, end time.Time) ([]domain.ScheduleLesson, error) {
	query := `
		SELECT 
			l.id, l.title, c.title as course_name, l.lesson_time, l.duration_min,
			u.first_name || ' ' || u.last_name as teacher_name, u.email as teacher_email,
			l.online_url as discord_url, COALESCE(ula.comment_teacher, '') as teacher_comment,
			COALESCE(uas.status, 'not_submitted') as homework_status
		FROM lessons l
		JOIN modules m ON l.module_id = m.id
		JOIN courses c ON m.course_id = c.id
		JOIN user_courses uc ON c.id = uc.course_id
		JOIN users u ON l.teacher_id = u.id
		LEFT JOIN user_lesson_attendance ula ON l.id = ula.lesson_id AND ula.user_id = $1
		LEFT JOIN assignments a ON l.id = a.lesson_id
		LEFT JOIN user_assignments_submission uas ON a.id = uas.assignment_id AND uas.user_id = $1
		WHERE uc.user_id = $1 AND l.lesson_time BETWEEN $2 AND $3
		ORDER BY l.lesson_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []domain.ScheduleLesson
	for rows.Next() {
		var l domain.ScheduleLesson
		err := rows.Scan(
			&l.ID, &l.Title, &l.CourseName, &l.StartTime, &l.DurationMin,
			&l.TeacherName, &l.TeacherEmail, &l.DiscordURL, &l.TeacherComment,
			&l.HomeworkStatus,
		)
		if err != nil {
			return nil, err
		}
		l.EndTime = l.StartTime.Add(time.Duration(l.DurationMin) * time.Minute)
		l.Color = "#4F46E5" // Логика выбора цвета может быть расширена на уровне курсов
		lessons = append(lessons, l)
	}
	return lessons, nil
}
