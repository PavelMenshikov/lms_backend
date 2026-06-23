-- +goose Up
CREATE TABLE IF NOT EXISTS test_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    options JSONB NOT NULL DEFAULT '[]'::jsonb,
    correct_answer VARCHAR(255) NOT NULL DEFAULT '',
    points INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_test_questions_test_id ON test_questions(test_id);

-- +goose Down
DROP TABLE IF EXISTS test_questions;
