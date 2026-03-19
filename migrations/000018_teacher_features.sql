-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS course_teachers (
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, teacher_id)
);

ALTER TABLE teachers ADD COLUMN IF NOT EXISTS working_hours JSONB DEFAULT '{}'::jsonb;

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS course_teachers;
ALTER TABLE teachers DROP COLUMN IF EXISTS working_hours;