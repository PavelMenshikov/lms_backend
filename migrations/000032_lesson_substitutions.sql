-- Sprint 3 gap: lesson_substitutions table per TOR
-- +goose Up
-- +goose StatementBegin

CREATE TABLE lesson_substitutions (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id            UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    original_teacher_id  UUID NOT NULL REFERENCES users(id),
    substitute_teacher_id UUID NOT NULL REFERENCES users(id),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lesson_substitutions_lesson ON lesson_substitutions(lesson_id);
CREATE INDEX idx_lesson_substitutions_original ON lesson_substitutions(original_teacher_id);
CREATE INDEX idx_lesson_substitutions_substitute ON lesson_substitutions(substitute_teacher_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS lesson_substitutions;

-- +goose StatementEnd
