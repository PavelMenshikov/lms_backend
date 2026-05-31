package reports

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ReportsService interface {
	GenerateLessonsReport(ctx context.Context, startDate, endDate time.Time) (*excelize.File, error)
}

type reportsService struct {
	db *sql.DB
}

func NewReportsService(db *sql.DB) ReportsService {
	return &reportsService{db: db}
}

type LessonReportRow struct {
	Date             string
	Time             string
	Course           string
	Group            string
	Student          string
	Teacher          string
	LessonStatus     string
	AttendanceStatus string
	Reason           string
	Comment          string
	Freeze           string
	ChangedBy        string
	ChangedAt        string
}

func (s *reportsService) GenerateLessonsReport(ctx context.Context, startDate, endDate time.Time) (*excelize.File, error) {
	query := `
		SELECT
			l.scheduled_at::date as lesson_date,
			l.scheduled_at::time as lesson_time,
			c.title as course_name,
			COALESCE(g.name, 'Без группы') as group_name,
			CONCAT(u.first_name, ' ', u.last_name) as student_name,
			CONCAT(t.first_name, ' ', t.last_name) as teacher_name,
			l.status as lesson_status,
			COALESCE(ar.status::text, 'Не отмечено') as attendance_status,
			COALESCE(ar.reason, '') as reason,
			COALESCE(ar.comment, '') as comment,
			CASE
				WHEN fp.is_active = true AND fp.start_date <= l.scheduled_at::date AND fp.end_date >= l.scheduled_at::date
				THEN 'Да'
				ELSE 'Нет'
			END as is_frozen,
			COALESCE(CONCAT(ub.first_name, ' ', ub.last_name), '') as changed_by,
			COALESCE(ar.updated_at::text, '') as changed_at
		FROM lessons l
		JOIN courses c ON l.course_id = c.id
		LEFT JOIN groups g ON l.group_id = g.id
		LEFT JOIN user_enrollments ue ON ue.course_id = c.id
		LEFT JOIN users u ON ue.user_id = u.id AND u.role = 'student'
		LEFT JOIN users t ON l.teacher_id = t.id
		LEFT JOIN attendance_records ar ON ar.lesson_id = l.id AND ar.student_id = u.id
		LEFT JOIN users ub ON ar.updated_by = ub.id
		LEFT JOIN freeze_periods fp ON fp.student_id = u.id
		WHERE l.scheduled_at BETWEEN $1 AND $2
		ORDER BY l.scheduled_at DESC, student_name
	`

	rows, err := s.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheetName := "Отчёт по занятиям"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Заголовки
	headers := []string{
		"Дата", "Время", "Курс", "Группа", "Ученик", "Учитель",
		"Статус занятия", "Статус посещения", "Причина", "Комментарий",
		"Заморозка", "Кем изменено", "Дата изменения",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Стиль заголовков
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle(sheetName, "A1", "M1", headerStyle)

	// Данные
	rowNum := 2
	for rows.Next() {
		var row LessonReportRow
		err := rows.Scan(
			&row.Date, &row.Time, &row.Course, &row.Group, &row.Student,
			&row.Teacher, &row.LessonStatus, &row.AttendanceStatus, &row.Reason,
			&row.Comment, &row.Freeze, &row.ChangedBy, &row.ChangedAt,
		)
		if err != nil {
			return nil, err
		}

		values := []interface{}{
			row.Date, row.Time, row.Course, row.Group, row.Student,
			row.Teacher, row.LessonStatus, row.AttendanceStatus, row.Reason,
			row.Comment, row.Freeze, row.ChangedBy, row.ChangedAt,
		}

		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, rowNum)
			f.SetCellValue(sheetName, cell, value)
		}
		rowNum++
	}

	// Автоширина колонок
	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheetName, col, col, 15)
	}

	return f, nil
}
