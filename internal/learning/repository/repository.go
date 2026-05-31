package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"lms_backend/internal/domain"

	"golang.org/x/sync/errgroup"
)

type LearningRepository interface {
	GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error)
	GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error)
	GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error)
	GetAssignmentIDByLesson(ctx context.Context, lessonID string) (string, error)
	SaveSubmission(ctx context.Context, userID, assignmentID, text string, files []string) error
	SetLessonAttendance(ctx context.Context, userID, lessonID, status, recordingURL, teacherComment string) error

	GetTeachersList(ctx context.Context) ([]*domain.TeacherPublicInfo, error)
	GetTeacherByID(ctx context.Context, id string) (*domain.TeacherPublicInfo, error)
	AddTeacherReview(ctx context.Context, review *domain.TeacherReview) error
	GetTeacherReviews(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error)
	GetTeacherCourses(ctx context.Context, teacherID string) ([]*domain.StudentCoursePreview, error)

	GetTestByID(ctx context.Context, testID string) (*domain.Test, error)
	GetProjectByID(ctx context.Context, projectID string) (*domain.Project, error)
	GetTeacherSubstitutions(ctx context.Context, teacherID string) ([]*domain.Lesson, error)
	GetTeacherUpcomingLessons(ctx context.Context, teacherID string) ([]*domain.Lesson, error)
	GetTeacherCancelledLessons(ctx context.Context, teacherID string) ([]*domain.Lesson, error)
}

type LearningRepoImpl struct {
	db *sql.DB
}

var _ LearningRepository = (*LearningRepoImpl)(nil)

func NewLearningRepository(db *sql.DB) *LearningRepoImpl {
	return &LearningRepoImpl{db: db}
}

func (r *LearningRepoImpl) GetMyCourses(ctx context.Context, userID string) ([]*domain.StudentCoursePreview, error) {
	query := `
		SELECT c.id, c.title, c.description, c.image_url, c.is_main, uc.progress_percent
		FROM courses c
		JOIN user_courses uc ON c.id = uc.course_id
		WHERE uc.user_id = $1 AND c.status = 'active'
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []*domain.StudentCoursePreview
	for rows.Next() {
		c := &domain.StudentCoursePreview{}
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.ImageURL, &c.IsMain, &c.ProgressPercent); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *LearningRepoImpl) GetCourseContent(ctx context.Context, courseID, userID string) (*domain.StudentCourseView, error) {
	view := &domain.StudentCourseView{
		Modules:      []*domain.StudentModuleView{},
		RootLessons:  []*domain.StudentLessonRef{},
		RootTests:    []domain.Test{},
		RootProjects: []domain.Project{},
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		view.Course = &domain.Course{}
		return r.db.QueryRowContext(egCtx, "SELECT id, title, description, image_url, is_main FROM courses WHERE id = $1", courseID).
			Scan(&view.Course.ID, &view.Course.Title, &view.Course.Description, &view.Course.ImageURL, &view.Course.IsMain)
	})

	eg.Go(func() error {
		rowsM, err := r.db.QueryContext(egCtx, "SELECT id, title, description, order_num FROM modules WHERE course_id = $1 ORDER BY order_num ASC", courseID)
		if err != nil {
			return err
		}
		defer rowsM.Close()
		for rowsM.Next() {
			m := &domain.StudentModuleView{Lessons: []*domain.StudentLessonRef{}}
			if err := rowsM.Scan(&m.ID, &m.Title, &m.Description, &m.OrderNum); err != nil {
				return err
			}
			view.Modules = append(view.Modules, m)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	queryL := `
		SELECT 
			l.id, l.module_id, l.title, l.order_num, l.duration_min,
			CASE WHEN ula.status IN ('visited', 'trial') OR uas.status = 'accepted' THEN true ELSE false END as is_completed
		FROM lessons l
		LEFT JOIN user_lesson_attendance ula ON l.id = ula.lesson_id AND ula.user_id = $2
		LEFT JOIN assignments a ON l.id = a.lesson_id
		LEFT JOIN user_assignments_submission uas ON a.id = uas.assignment_id AND uas.user_id = $2
		WHERE l.course_id = $1 AND l.is_published = true
		ORDER BY l.order_num ASC`

	rowsL, err := r.db.QueryContext(ctx, queryL, courseID, userID)
	if err != nil {
		return nil, err
	}
	defer rowsL.Close()

	allLessonsMap := make(map[string]*domain.StudentLessonRef)

	for rowsL.Next() {
		var l domain.StudentLessonRef
		var mid sql.NullString
		if err := rowsL.Scan(&l.ID, &mid, &l.Title, &l.OrderNum, &l.DurationMin, &l.IsCompleted); err != nil {
			return nil, err
		}

		l.Tests = []domain.Test{}
		l.Projects = []domain.Project{}
		allLessonsMap[l.ID] = &l

		if mid.Valid {
			for _, m := range view.Modules {
				if m.ID == mid.String {
					m.Lessons = append(m.Lessons, allLessonsMap[l.ID])
				}
			}
		} else {
			view.RootLessons = append(view.RootLessons, allLessonsMap[l.ID])
		}
	}

	eg2, egCtx2 := errgroup.WithContext(ctx)

	eg2.Go(func() error {
		testsQuery := `SELECT id, lesson_id, title, description, passing_score FROM tests WHERE lesson_id IN (SELECT id FROM lessons WHERE course_id = $1)`
		rowsT, err := r.db.QueryContext(egCtx2, testsQuery, courseID)
		if err != nil {
			return err
		}
		defer rowsT.Close()
		for rowsT.Next() {
			var t domain.Test
			var lid sql.NullString
			if err := rowsT.Scan(&t.ID, &lid, &t.Title, &t.Description, &t.PassingScore); err != nil {
				return err
			}
			if lid.Valid {
				s := lid.String
				t.LessonID = &s
				if less, ok := allLessonsMap[s]; ok {
					less.Tests = append(less.Tests, t)
				}
			} else {
				view.RootTests = append(view.RootTests, t)
			}
		}
		return nil
	})

	eg2.Go(func() error {
		projQuery := `SELECT id, lesson_id, title, description, max_score FROM projects WHERE lesson_id IN (SELECT id FROM lessons WHERE course_id = $1)`
		rowsP, err := r.db.QueryContext(egCtx2, projQuery, courseID)
		if err != nil {
			return err
		}
		defer rowsP.Close()
		for rowsP.Next() {
			var p domain.Project
			var lid sql.NullString
			if err := rowsP.Scan(&p.ID, &lid, &p.Title, &p.Description, &p.MaxScore); err != nil {
				return err
			}
			if lid.Valid {
				s := lid.String
				p.LessonID = &s
				if less, ok := allLessonsMap[s]; ok {
					less.Projects = append(less.Projects, p)
				}
			} else {
				view.RootProjects = append(view.RootProjects, p)
			}
		}
		return nil
	})

	if err := eg2.Wait(); err != nil {
		return nil, err
	}

	return view, nil
}

func (r *LearningRepoImpl) GetLessonDetail(ctx context.Context, lessonID, userID string) (*domain.StudentLessonDetail, error) {
	lesson := &domain.Lesson{}
	var contentRaw []byte

	query := `SELECT id, title, video_url, presentation_url, content_text, content, duration_min 
              FROM lessons WHERE id = $1 AND is_published = true`

	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(
		&lesson.ID, &lesson.Title, &lesson.VideoURL, &lesson.PresentationURL,
		&lesson.ContentText, &contentRaw, &lesson.DurationMin,
	)
	if err != nil {
		return nil, fmt.Errorf("lesson not found: %w", err)
	}

	if len(contentRaw) > 0 {
		_ = json.Unmarshal(contentRaw, &lesson.Content)
	}
	if lesson.Content == nil {
		lesson.Content = []domain.ContentBlock{}
	}

	res := &domain.StudentLessonDetail{Lesson: lesson}

	attendanceQuery := `
		SELECT COALESCE(status::text, ''), COALESCE(recording_url, ''), COALESCE(comment_teacher, '')
		FROM user_lesson_attendance
		WHERE lesson_id = $1 AND user_id = $2`

	var attStatus, recURL, attComment string
	if err := r.db.QueryRowContext(ctx, attendanceQuery, lessonID, userID).Scan(&attStatus, &recURL, &attComment); err == nil {
		res.AttendanceStatus = attStatus
		res.RecordingURL = recURL
		res.IsCompleted = (attStatus == "visited" || attStatus == "trial")
	}

	homeworkQuery := `
		SELECT uas.status, COALESCE(uas.grade, 0), COALESCE(uas.teacher_comment, '')
		FROM assignments a
		JOIN user_assignments_submission uas ON a.id = uas.assignment_id
		WHERE a.lesson_id = $1 AND uas.user_id = $2
		LIMIT 1`

	var hwStatus string
	var grade int
	var hwComment string
	err = r.db.QueryRowContext(ctx, homeworkQuery, lessonID, userID).Scan(&hwStatus, &grade, &hwComment)

	if err == nil {
		res.AssignmentStatus = hwStatus
		if hwComment != "" {
			res.TeacherComment = hwComment
		}
		res.Grade = grade
		if !res.IsCompleted && hwStatus == "accepted" {
			res.IsCompleted = true
		}
	}

	return res, nil
}

func (r *LearningRepoImpl) GetAssignmentIDByLesson(ctx context.Context, lessonID string) (string, error) {
	var id string
	query := `SELECT id FROM assignments WHERE lesson_id = $1 LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, lessonID).Scan(&id)
	return id, err
}

func (r *LearningRepoImpl) SaveSubmission(ctx context.Context, userID, assignmentID, text string, files []string) error {
	filesJSON, _ := json.Marshal(files)
	query := `
		INSERT INTO user_assignments_submission (user_id, assignment_id, submission_text, submission_files, status, submitted_at)
		VALUES ($1, $2, $3, $4, 'pending_check', NOW())
		ON CONFLICT (user_id, assignment_id) 
		DO UPDATE SET 
			submission_text = EXCLUDED.submission_text, 
			submission_files = EXCLUDED.submission_files, 
			status = 'pending_check', 
			submitted_at = NOW()
		WHERE user_assignments_submission.status != 'accepted'
	`
	_, err := r.db.ExecContext(ctx, query, userID, assignmentID, text, filesJSON)
	return err
}

func (r *LearningRepoImpl) SetLessonAttendance(ctx context.Context, userID, lessonID, status, recordingURL, teacherComment string) error {
	query := `
		INSERT INTO user_lesson_attendance (user_id, lesson_id, status, recording_url, comment_teacher)
		VALUES ($1, $2, $3::attendance_status, NULLIF($4, ''), NULLIF($5, ''))
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			status = $3::attendance_status,
			recording_url = COALESCE(NULLIF($4, ''), user_lesson_attendance.recording_url),
			comment_teacher = COALESCE(NULLIF($5, ''), user_lesson_attendance.comment_teacher)
	`
	_, err := r.db.ExecContext(ctx, query, userID, lessonID, status, recordingURL, teacherComment)
	return err
}

func (r *LearningRepoImpl) GetTeachersList(ctx context.Context) ([]*domain.TeacherPublicInfo, error) {
	query := `
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar_url, ''), 
		       COALESCE(ROUND(tr.avg_rating, 1), 0.0) as rating, 
		       COALESCE(u.experience_years, 0), u.email, COALESCE(u.city, ''), COALESCE(u.phone, '')
		FROM users u
		LEFT JOIN (SELECT teacher_id, AVG(rating) as avg_rating FROM teacher_reviews GROUP BY teacher_id) tr ON tr.teacher_id = u.id
		WHERE u.role = 'teacher'
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teachers []*domain.TeacherPublicInfo
	for rows.Next() {
		t := &domain.TeacherPublicInfo{}
		if err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.AvatarURL, &t.Rating, &t.ExperienceYears, &t.Email, &t.City, &t.Phone); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, nil
}

func (r *LearningRepoImpl) GetTeacherByID(ctx context.Context, id string) (*domain.TeacherPublicInfo, error) {
	t := &domain.TeacherPublicInfo{}

	query := `
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar_url, ''), 
		       COALESCE(ROUND(tr.avg_rating, 1), 0.0) as rating, 
		       COALESCE(u.experience_years, 0), COALESCE(t.bio, ''), u.email, COALESCE(u.city, ''), COALESCE(u.phone, ''),
		       COALESCE(t.working_hours, '{}'::jsonb)
		FROM users u
		LEFT JOIN teachers t ON u.id = t.id
		LEFT JOIN (SELECT teacher_id, AVG(rating) as avg_rating FROM teacher_reviews GROUP BY teacher_id) tr ON tr.teacher_id = u.id
		WHERE u.id = $1 AND u.role = 'teacher'
	`

	var scheduleRaw []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.FirstName, &t.LastName, &t.AvatarURL,
		&t.Rating, &t.ExperienceYears, &t.Bio, &t.Email, &t.City, &t.Phone,
		&scheduleRaw,
	)
	if err != nil {
		return nil, err
	}

	if len(scheduleRaw) > 0 {
		var schedule map[string]interface{}
		if err := json.Unmarshal(scheduleRaw, &schedule); err == nil {
			t.Schedule = schedule
		}
	}
	if t.Schedule == nil {
		t.Schedule = map[string]interface{}{}
	}

	return t, nil
}

func (r *LearningRepoImpl) AddTeacherReview(ctx context.Context, rev *domain.TeacherReview) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO teacher_reviews (teacher_id, student_id, rating, comment)
		VALUES ($1, $2, $3, $4)
	`
	if _, err := tx.ExecContext(ctx, query, rev.TeacherID, rev.StudentID, rev.Rating, rev.Comment); err != nil {
		tx.Rollback()
		return err
	}
	updateRating := `
		UPDATE teachers SET rating = (
			SELECT AVG(rating) FROM teacher_reviews WHERE teacher_id = $1
		) WHERE id = $1
	`
	if _, err := tx.ExecContext(ctx, updateRating, rev.TeacherID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *LearningRepoImpl) GetTeacherReviews(ctx context.Context, teacherID string) ([]*domain.TeacherReview, error) {
	query := `
		SELECT tr.id, tr.teacher_id, tr.student_id, u.first_name || ' ' || u.last_name, tr.rating, tr.comment, tr.created_at
		FROM teacher_reviews tr
		JOIN users u ON tr.student_id = u.id
		WHERE tr.teacher_id = $1
		ORDER BY tr.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []*domain.TeacherReview
	for rows.Next() {
		rev := &domain.TeacherReview{}
		if err := rows.Scan(&rev.ID, &rev.TeacherID, &rev.StudentID, &rev.StudentName, &rev.Rating, &rev.Comment, &rev.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}
func (r *LearningRepoImpl) GetTestByID(ctx context.Context, testID string) (*domain.Test, error) {
	t := &domain.Test{}
	err := r.db.QueryRowContext(ctx, "SELECT id, lesson_id, title, description, passing_score, created_at FROM tests WHERE id = $1", testID).
		Scan(&t.ID, &t.LessonID, &t.Title, &t.Description, &t.PassingScore, &t.CreatedAt)
	return t, err
}

func (r *LearningRepoImpl) GetProjectByID(ctx context.Context, projectID string) (*domain.Project, error) {
	p := &domain.Project{}
	err := r.db.QueryRowContext(ctx, "SELECT id, lesson_id, title, description, max_score, created_at FROM projects WHERE id = $1", projectID).
		Scan(&p.ID, &p.LessonID, &p.Title, &p.Description, &p.MaxScore, &p.CreatedAt)
	return p, err
}
func (r *LearningRepoImpl) GetTeacherCourses(ctx context.Context, teacherID string) ([]*domain.StudentCoursePreview, error) {
	query := `
		SELECT DISTINCT c.id, c.title, c.description, c.image_url, c.is_main, 0 as progress_percent
		FROM courses c
		JOIN streams s ON c.id = s.course_id
		JOIN groups g ON s.id = g.stream_id
		WHERE g.teacher_id = $1 AND c.status != 'archived'
		ORDER BY c.title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*domain.StudentCoursePreview
	for rows.Next() {
		c := &domain.StudentCoursePreview{}
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.ImageURL, &c.IsMain, &c.ProgressPercent); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func scanLessons(rows *sql.Rows) ([]*domain.Lesson, error) {
	var lessons []*domain.Lesson
	for rows.Next() {
		l := &domain.Lesson{}
		var cancelledAt sql.NullTime
		var cancellationReason sql.NullString
		var substitutedTeacherID sql.NullString
		err := rows.Scan(&l.ID, &l.CourseID, &l.ModuleID, &l.TeacherID, &l.Title, &l.LessonTime, &l.DurationMin, &l.OrderNum, &l.IsPublished, &l.VideoURL, &l.PresentationURL, &l.ContentText, &l.HasHomework, &l.IsInteractive, &l.IsCancelled, &cancelledAt, &cancellationReason, &substitutedTeacherID)
		if err != nil {
			return nil, err
		}
		if cancelledAt.Valid {
			l.CancelledAt = &cancelledAt.Time
		}
		if cancellationReason.Valid {
			l.CancellationReason = cancellationReason.String
		}
		if substitutedTeacherID.Valid {
			l.SubstitutedTeacherID = &substitutedTeacherID.String
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *LearningRepoImpl) GetTeacherSubstitutions(ctx context.Context, teacherID string) ([]*domain.Lesson, error) {
	query := `
		SELECT id, course_id, module_id, teacher_id, title, lesson_time, duration_min, order_num,
		       is_published, COALESCE(video_url, ''), COALESCE(presentation_url, ''), COALESCE(content_text, ''),
		       COALESCE(has_homework, FALSE), COALESCE(is_interactive, FALSE),
		       is_cancelled, cancelled_at, cancellation_reason, substituted_teacher_id
		FROM lessons
		WHERE substituted_teacher_id = $1 AND is_cancelled = FALSE
		ORDER BY lesson_time ASC
		LIMIT 50
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLessons(rows)
}

func (r *LearningRepoImpl) GetTeacherUpcomingLessons(ctx context.Context, teacherID string) ([]*domain.Lesson, error) {
	query := `
		SELECT id, course_id, module_id, teacher_id, title, lesson_time, duration_min, order_num,
		       is_published, COALESCE(video_url, ''), COALESCE(presentation_url, ''), COALESCE(content_text, ''),
		       COALESCE(has_homework, FALSE), COALESCE(is_interactive, FALSE),
		       is_cancelled, cancelled_at, cancellation_reason, substituted_teacher_id
		FROM lessons
		WHERE (teacher_id = $1 OR substituted_teacher_id = $1)
		  AND lesson_time > NOW() AND is_cancelled = FALSE
		ORDER BY lesson_time ASC
		LIMIT 20
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLessons(rows)
}

func (r *LearningRepoImpl) GetTeacherCancelledLessons(ctx context.Context, teacherID string) ([]*domain.Lesson, error) {
	query := `
		SELECT id, course_id, module_id, teacher_id, title, lesson_time, duration_min, order_num,
		       is_published, COALESCE(video_url, ''), COALESCE(presentation_url, ''), COALESCE(content_text, ''),
		       COALESCE(has_homework, FALSE), COALESCE(is_interactive, FALSE),
		       is_cancelled, cancelled_at, cancellation_reason, substituted_teacher_id
		FROM lessons
		WHERE teacher_id = $1 AND is_cancelled = TRUE
		ORDER BY cancelled_at DESC
		LIMIT 20
	`
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLessons(rows)
}
