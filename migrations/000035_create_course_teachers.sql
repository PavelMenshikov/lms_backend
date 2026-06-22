-- +goose Up
CREATE TABLE IF NOT EXISTS course_teachers (
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (course_id, teacher_id)
);

CREATE INDEX idx_course_teachers_teacher_id ON course_teachers(teacher_id);

-- +goose Down
DROP TABLE IF EXISTS course_teachers;
