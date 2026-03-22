-- +goose Up
-- +goose StatementBegin
ALTER TABLE teachers ADD COLUMN IF NOT EXISTS working_hours JSONB DEFAULT '{}'::jsonb;
-- +goose StatementEnd

-- +goose Down
ALTER TABLE teachers DROP COLUMN IF EXISTS working_hours;