-- +goose Up
-- +goose StatementBegin
ALTER TYPE assignment_status ADD VALUE IF NOT EXISTS 'on_revision';

ALTER TABLE user_assignments_submission ADD COLUMN IF NOT EXISTS submission_files JSONB DEFAULT '[]'::jsonb;
-- +goose StatementEnd