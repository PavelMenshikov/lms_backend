-- Sprint 3: Teacher Substitutions and Cancelled Lessons (FIX-013, FIX-015, FIX-026, FIX-028)
-- +goose Up
-- +goose StatementBegin

ALTER TABLE lessons
    ADD COLUMN is_cancelled          BOOLEAN     NOT NULL DEFAULT FALSE,
    ADD COLUMN cancelled_at          TIMESTAMPTZ,
    ADD COLUMN cancellation_reason   TEXT,
    ADD COLUMN substituted_teacher_id UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_lessons_substituted_teacher ON lessons(substituted_teacher_id) WHERE substituted_teacher_id IS NOT NULL;
CREATE INDEX idx_lessons_cancelled ON lessons(is_cancelled) WHERE is_cancelled = TRUE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_lessons_substituted_teacher;
DROP INDEX IF EXISTS idx_lessons_cancelled;

ALTER TABLE lessons
    DROP COLUMN substituted_teacher_id,
    DROP COLUMN cancellation_reason,
    DROP COLUMN cancelled_at,
    DROP COLUMN is_cancelled;

-- +goose StatementEnd
