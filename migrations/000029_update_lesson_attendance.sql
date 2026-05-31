-- +goose Up
CREATE TYPE attendance_status AS ENUM ('visited', 'missing_valid', 'missing_invalid', 'frozen', 'trial');

ALTER TABLE user_lesson_attendance
    DROP COLUMN is_attended,
    ADD COLUMN status attendance_status NOT NULL DEFAULT 'visited',
    ADD COLUMN recording_url VARCHAR(500);

-- +goose Down
ALTER TABLE user_lesson_attendance
    DROP COLUMN recording_url,
    DROP COLUMN status,
    ADD COLUMN is_attended BOOLEAN NOT NULL DEFAULT false;

DROP TYPE attendance_status;
