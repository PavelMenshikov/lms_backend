-- +goose Up
-- +goose StatementBegin
ALTER TABLE lessons ALTER COLUMN teacher_id DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_user_courses_user_id ON user_courses(user_id);
CREATE INDEX IF NOT EXISTS idx_user_courses_course_id ON user_courses(course_id);
CREATE INDEX IF NOT EXISTS idx_groups_teacher_id ON groups(teacher_id);
CREATE INDEX IF NOT EXISTS idx_groups_curator_id ON groups(curator_id);
-- +goose StatementEnd