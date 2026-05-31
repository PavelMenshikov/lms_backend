-- +goose Up
ALTER TABLE users
    ADD COLUMN intro_broadcast_url TEXT,
    ADD COLUMN graduation_broadcast_url TEXT,
    ADD COLUMN subscription_end_date TIMESTAMP,
    ADD COLUMN balance NUMERIC(12, 2) NOT NULL DEFAULT 0,
    ADD COLUMN loss_reason TEXT;

-- +goose Down
ALTER TABLE users
    DROP COLUMN IF EXISTS loss_reason,
    DROP COLUMN IF EXISTS balance,
    DROP COLUMN IF EXISTS subscription_end_date,
    DROP COLUMN IF EXISTS graduation_broadcast_url,
    DROP COLUMN IF EXISTS intro_broadcast_url;
