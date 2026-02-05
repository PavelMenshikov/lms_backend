-- +goose Up
-- +goose StatementBegin
DO $$ 
BEGIN
    ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'moderator';
    ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'curator';
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +goose StatementEnd