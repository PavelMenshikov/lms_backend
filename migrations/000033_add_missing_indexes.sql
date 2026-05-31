-- +goose Up
-- +goose StatementBegin

-- lessons: teacher lookups for monthly report and teacher dashboard
CREATE INDEX IF NOT EXISTS idx_lessons_teacher_time ON lessons(teacher_id, lesson_time);
CREATE INDEX IF NOT EXISTS idx_lessons_sub_teacher_time ON lessons(substituted_teacher_id, lesson_time);
CREATE INDEX IF NOT EXISTS idx_lessons_lesson_time ON lessons(lesson_time);
CREATE INDEX IF NOT EXISTS idx_lessons_cancelled_teacher ON lessons(is_cancelled, teacher_id);

-- attendance: primary lookup and aggregation
CREATE INDEX IF NOT EXISTS idx_ula_user_lesson ON user_lesson_attendance(user_id, lesson_id);
CREATE INDEX IF NOT EXISTS idx_ula_lesson_status ON user_lesson_attendance(lesson_id, status);

-- substitutions: monthly report lookups
CREATE INDEX IF NOT EXISTS idx_lesson_subs_substitute_time ON lesson_substitutions(substitute_teacher_id, created_at);
CREATE INDEX IF NOT EXISTS idx_lesson_subs_original_time ON lesson_substitutions(original_teacher_id, created_at);

-- groups: curator and teacher dashboard
CREATE INDEX IF NOT EXISTS idx_groups_curator_id ON groups(curator_id);
CREATE INDEX IF NOT EXISTS idx_groups_teacher_id ON groups(teacher_id);

-- ordering and filtering
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);

-- user_courses: group lookups
CREATE INDEX IF NOT EXISTS idx_user_courses_group_id ON user_courses(group_id);

-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS idx_lessons_teacher_time;
DROP INDEX IF EXISTS idx_lessons_sub_teacher_time;
DROP INDEX IF EXISTS idx_lessons_lesson_time;
DROP INDEX IF EXISTS idx_lessons_cancelled_teacher;
DROP INDEX IF EXISTS idx_ula_user_lesson;
DROP INDEX IF EXISTS idx_ula_lesson_status;
DROP INDEX IF EXISTS idx_lesson_subs_substitute_time;
DROP INDEX IF EXISTS idx_lesson_subs_original_time;
DROP INDEX IF EXISTS idx_groups_curator_id;
DROP INDEX IF EXISTS idx_groups_teacher_id;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_courses_status;
DROP INDEX IF EXISTS idx_user_courses_group_id;
