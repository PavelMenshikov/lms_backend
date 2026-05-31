-- +goose Up
-- +goose StatementBegin

-- Users: role-based listings (admin pages) and sorting
CREATE INDEX IF NOT EXISTS idx_users_role_created ON users(role, created_at DESC);

-- user_courses: student enrollment lookups and course aggregations
CREATE INDEX IF NOT EXISTS idx_user_courses_user_id ON user_courses(user_id, course_id);
CREATE INDEX IF NOT EXISTS idx_user_courses_course_id ON user_courses(course_id, user_id);

-- attendance: per-user stats and per-lesson lookups
CREATE INDEX IF NOT EXISTS idx_ula_user_status ON user_lesson_attendance(user_id, status);
CREATE INDEX IF NOT EXISTS idx_ula_lesson_user ON user_lesson_attendance(lesson_id, user_id);

-- submissions: per-user homework stats and per-assignment lookups
CREATE INDEX IF NOT EXISTS idx_submissions_user_status ON user_assignments_submission(user_id, status);
CREATE INDEX IF NOT EXISTS idx_submissions_assignment_user ON user_assignments_submission(assignment_id, user_id);

-- teacher_reviews: rating aggregations
CREATE INDEX IF NOT EXISTS idx_reviews_teacher_id ON teacher_reviews(teacher_id);

-- assignments: lesson lookups
CREATE INDEX IF NOT EXISTS idx_assignments_lesson_id ON assignments(lesson_id);

-- lessons: course_id + ordering (frequent sort)
CREATE INDEX IF NOT EXISTS idx_lessons_course_order ON lessons(course_id, order_num);

-- user_courses: group aggregations
CREATE INDEX IF NOT EXISTS idx_user_courses_group_course ON user_courses(group_id, course_id);

-- courses: status filters
CREATE INDEX IF NOT EXISTS idx_courses_status_created ON courses(status, created_at DESC);

-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS idx_users_role_created;
DROP INDEX IF EXISTS idx_user_courses_user_id;
DROP INDEX IF EXISTS idx_user_courses_course_id;
DROP INDEX IF EXISTS idx_ula_user_status;
DROP INDEX IF EXISTS idx_ula_lesson_user;
DROP INDEX IF EXISTS idx_submissions_user_status;
DROP INDEX IF EXISTS idx_submissions_assignment_user;
DROP INDEX IF EXISTS idx_reviews_teacher_id;
DROP INDEX IF EXISTS idx_assignments_lesson_id;
DROP INDEX IF EXISTS idx_lessons_course_order;
DROP INDEX IF EXISTS idx_user_courses_group_course;
DROP INDEX IF EXISTS idx_courses_status_created;
