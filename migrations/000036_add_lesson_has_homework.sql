-- +goose Up
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS has_homework BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE lessons DROP COLUMN IF EXISTS has_homework;
